package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"
)

var pathToTemplates = "./templates/"

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]interface{})

	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "hit the home page "+time.Now().Format(time.RFC1123))
	}
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
}

type TemplateData struct {
	IP   string
	Data map[string]interface{}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad Request", http.StatusBadRequest)
		return err
	}

	data.IP = app.ipFromContext(r.Context())
	// execute the template

	err = parsedTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "bad Request", http.StatusBadRequest)
		return err
	}
	return nil
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "bad Request", http.StatusBadRequest)
		return
	}

	// validate the form data
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		fmt.Fprint(w, "failed validation")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	w.Write([]byte("Email: " + email + "\n"))
	w.Write([]byte("Password: " + password + "\n"))
}
