package controllers

import (
	"io/ioutil"
	"net/http"
)

// HomeHandler function
func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	data := &struct {
		TransactionID string
		Success       bool
		Response      bool
	}{
		TransactionID: "",
		Success:       false,
		Response:      false,
	}

	if r.FormValue("submitted") == "true" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err1 := r.FormFile("uploadfile")
		fileBytes, err2 := ioutil.ReadAll(file)
		if err1 != nil || err2 != nil {
			http.Error(w, "Unable to upload the file", 500)
		}

		app.fabric.Query(handler.Filename)
	}
	renderTemplate(w, r, "home.html", data)
}
