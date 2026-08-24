package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_routes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/static/*", "GET"},
		{"/login", "POST"},
		{"/user/profile", "GET"},
	}

	var app application
	mux := app.routes()

	chiRouts := mux.(chi.Routes)
	for _, route := range registered {
		if !routeExists(route.route, route.method, chiRouts) {
			t.Errorf("Route %s %s not registered", route.method, route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRouts chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRouts, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})
	return found
}
