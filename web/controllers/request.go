package controllers

import (
	"net/http"
)

// RequestHandler function
func (app *Application) RequestHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	renderTemplate(w, r, "request.html", nil)
}
