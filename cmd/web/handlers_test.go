package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		//{"static", "/static/css/main.css", http.StatusOK},
		{"404", "/notfound", http.StatusNotFound},
	}

	routes := app.routes()
	// create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("%s: expected %d, got %d", e.name, e.expectedStatusCode, resp.StatusCode)
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
