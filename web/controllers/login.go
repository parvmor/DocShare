package controllers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	userpass = make(map[string][]byte)
	key      = []byte("super-secret-key")
	store    = sessions.NewCookieStore(key)
)

// SignupHandler function
func (app *Application) SignupHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	if r.FormValue("submitted") == "submit" {
		user, pass := r.FormValue("username"), r.FormValue("password")
		passhash, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, ok := userpass[user]; ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userpass[user] = passhash
		session.Values["authenticated"] = true
		session.Save(r, w)

		renderTemplate(w, r, "home.html", nil)
	} else {
		renderTemplate(w, r, "signup.html", nil)
	}
}

// SigninHandler function
func (app *Application) SigninHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	if r.FormValue("submitted") == "submit" {
		user, pass := r.FormValue("username"), r.FormValue("password")
		if _, ok := userpass[user]; !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		passhash = userpass[user]
		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword(passhash, []byte(pass)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		session.Values["authenticated"] = true
		session.Save(r, w)

		renderTemplate(w, r, "home.html", nil)
	} else {
		renderTemplate(w, r, "signup.html", nil)
	}
}

// SignoutHandler function
func (app *Application) SignoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
}
