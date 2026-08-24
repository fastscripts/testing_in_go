package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name                    string
		url                     string
		expectedStatusCode      int
		expectedURL             string
		expectedFirstStatusCode int
	}{
		{"home", "/", http.StatusOK, "/", http.StatusOK},
		//{"static", "/static/css/main.css", http.StatusOK},
		{"404", "/notfound", http.StatusNotFound, "/notfound", http.StatusNotFound},
		{"profile", "/user/profile", http.StatusOK, "/", http.StatusTemporaryRedirect},
	}

	routes := app.routes()
	// create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, e := range theTests {
		fmt.Println(e.name)
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("%s: expected %d, got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

		if resp.Request.URL.Path != e.expectedURL {
			t.Errorf("%s: expected to be redirected to %s, but got %s", e.name, e.expectedURL, resp.Request.URL.Path)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("%s: expected first status code %d, got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

		resp2, _ := client.Get(ts.URL + e.url)
		if resp2.StatusCode != e.expectedFirstStatusCode {
			t.Errorf("%s: expected first returned status code to be %d, but got %d", e.name, e.expectedStatusCode, resp2.StatusCode)
		}
	}
}

func TestAppHome_Old(t *testing.T) {
	// create request
	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(rr, req)

	// status code should be 200
	if rr.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), "<small>Your session:") {
		t.Errorf("expected to find %q in response body", "<small>Your session:")
	}
}

func TestAppHome(t *testing.T) {
	var theTests = []struct {
		name         string
		method       string
		url          string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "GET", "/", "", "<small>Your session:"},
		{"second visit", "GET", "/", "hit the home page", "<small>Your session: hit the home page"},
	}
	for _, e := range theTests {
		req, _ := http.NewRequest(e.method, e.url, nil)
		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("%s: expected %d, got %d", e.name, http.StatusOK, rr.Code)
		}

		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s: expected to find %s in response body", e.name, e.expectedHTML)
		}

	}

}

func TestApp_renderWithBadTemplate(t *testing.T) {
	// set templatepath to a location that does not exist
	pathToTemplates = "./testdata/"
	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("expected error but did not get one")
	}

	pathToTemplates = "./../../templates/"

}
func getCtx(r *http.Request) context.Context {
	ctx := context.WithValue(r.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(ctx)
}

func Test_app_Login(t *testing.T) {
	var theTests = []struct {
		name               string
		postetData         url.Values
		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name: "valid credentials",
			postetData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/user/profile",
		},
		{
			name: "missing form data",
			postetData: url.Values{
				"email":    {""},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
		},
		{
			name: "wrong password",
			postetData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"123"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
		},
		{
			name: "user not found",
			postetData: url.Values{
				"email":    {"nonexistent@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLocation:   "/",
		},
	}

	for _, e := range theTests {
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(e.postetData.Encode()))
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status code %d, got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if rr.Header().Get("Location") != e.expectedLocation {
			t.Errorf("%s: expected location %s, got %s", e.name, e.expectedLocation, rr.Header().Get("Location"))
		}
	}
}
