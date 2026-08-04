package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_form_Has(t *testing.T) {

	form := NewForm(nil)

	has := form.Has("whatever")
	if has {
		t.Error("form should not have field 'whatever', but Has() returned true")
	}

	postData := url.Values{}
	postData.Add("a", "a")

	form = NewForm(postData)

	has = form.Has("a")
	if !has {
		t.Errorf("form.Has() = %v, want %v", has, true)
	}
}

func Test_form_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form should be invalid when required fields are missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form should be valid when required fields are present")
	}
}

func Test_form_check(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "Password field is required")
	if form.Valid() {
		t.Error("Valid() returns false, and it should by true when calling Check()")
	}
}

func TestForm_ErrorGet(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "Password field is required")
	if form.Errors.Get("password") == "" {
		t.Error("form.Errors.Get() should return an error message for the 'password' field")
	}

	if form.Errors.Get("nonexistent") != "" {
		t.Error("form.Errors.Get() should return an empty string for a nonexistent field")
	}
}
