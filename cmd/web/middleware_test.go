package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_addIPToConttext(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-FORWARDED-FOR", "192.168.1.1", "", false},
		{"", "", "hello:world", false},
	}

	var app application

	// create a dummy handler that we`ll use to test the context value`

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure that the value exists in the context
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Error(contextUserKey, "value not found in context")
		}

		// make sure we got a string back
		ip, ok := val.(string)
		if !ok {
			t.Error(contextUserKey, "value is not a string")
		}
		t.Logf("ip: %s", ip)
	})

	for _, e := range tests {
		//crete the handler to test
		handlerToTest := app.addIPToContext(nextHandler)

		// create a dummy request to pass to the handler
		req := httptest.NewRequest("GET", "/", nil)

		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headerName) > 0 {
			req.Header.Set(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		// create a dummy response recorder
		rr := httptest.NewRecorder()
		// call the handler
		handlerToTest.ServeHTTP(rr, req)
	}
}

func Test_application_ipFRomContext(t *testing.T) {
	var app application

	ctx := context.Background()

	ctx = context.WithValue(ctx, contextUserKey, "192.168.1.1")

	ip := app.ipFromContext(ctx)

	if ip != "192.168.1.1" {
		t.Errorf("expected ip to be %s, got %s", "192.168.1.1", ip)
	}
}
