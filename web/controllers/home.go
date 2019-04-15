package controllers

import (
	"net/http"
)

// HomeHandler function
func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	renderTemplate(w, r, "home.html", nil)
}
