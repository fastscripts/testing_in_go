package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/fastscripts/testing_in_go/pkg/data"
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

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]interface{})

	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "hit the home page "+time.Now().Format(time.RFC1123))
	}
	_ = app.render(w, r, "profile.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP    string
	Data  map[string]interface{}
	Error string
	Flash string
	User  data.User
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad Request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(r.Context())

	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	// execute the template

	err = parsedTemplate.Execute(w, td)
	if err != nil {
		http.Error(w, "bad Request", http.StatusBadRequest)
		return err
	}
	return nil
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// validate data
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		//fmt.Fprint(w, "failed validation")
		// redirect to login Page with error message
		app.Session.Put(r.Context(), "error", "valid Email and Password are required")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "invalid Email or Password")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//log.Println("From database:", user.FirstName)
	//log.Println(password, user.FirstName)

	// authenticate user
	if !app.authenticate(r, user, password) {
		// if not authenticated then redirect to login page with error message
		app.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// privent fiaxation attack
	_ = app.Session.RenewToken(r.Context())

	// store success message in session

	// redirect to some other page
	app.Session.Put(r.Context(), "flash", "Successful login")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)

	fmt.Fprint(w, email)
}

func (app *application) authenticate(r *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}
	app.Session.Put(r.Context(), "user", user)
	return true
}
