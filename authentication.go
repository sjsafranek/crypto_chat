package main

import (
	"errors"
	// "fmt"
	"net/http"
	"html/template"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var users_file = "users.json"
var db Users

func init() {
	secret := uuid.New().String()
	store = sessions.NewCookieStore([]byte(secret))
	db = Users{}
	db.Fetch(users_file)
	user := User{Username: "admin", Email:"admin@email.com"}
	user.SetPassword("dev")
	db.Add(&user)
	db.Save(users_file)
}

func HasSession(r *http.Request) bool {
	session, err := store.Get(r, "chat-session")

	if nil != err {
		return false
	}
	if nil == session.Values["loggedin"] {
		return false
	}
	return true
}

func GetUserFromSession(r *http.Request) *User {
	session, _ := store.Get(r, "chat-session")
	username := session.Values["username"].(string)
	user, err := db.Get(username)
	if nil != err {
		logger.Error(err)
	}
	return user
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	template_file := "./templates/login.html"
	tmpl, _ := template.ParseFiles(template_file)

	if "GET" == r.Method {
		if HasSession(r) {
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
			return
		}
		err := tmpl.Execute(w, "")
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := r.ParseForm()
	if nil != err {
		logger.Error(err)
		http.Error(w, "Unable to parse form", http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	user, err := db.Get(username)
	if nil != err {
		logger.Warn(err)
		err := tmpl.Execute(w, err)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	password := r.Form["password"][0]
	if !user.IsPassword(password) {
		err = errors.New("Incorrect password")
		logger.Warn(err)
		err := tmpl.Execute(w, err)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// create session
	session, _ := store.Get(r, "chat-session")
	session.Values["loggedin"] = true
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// delete session
	session, _ := store.Get(r, "chat-session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
