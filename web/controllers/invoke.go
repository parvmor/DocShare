package controllers

import (
	"io/ioutil"
	"net/http"
)

// PutFileHandler function
func (app *Application) PutFileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user := session.Values["user"].(string)
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
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, "Unable to upload the file", 500)
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to upload the file", 500)
		}

		app.Fabric.InvokePutFile(fileBytes, handler.Filename, user)
	}
	renderTemplate(w, r, "home.html", data)
}

// ShareFileHandler function
func (app *Application) ShareFileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user := session.Values["user"].(string)
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
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, "Unable to upload the file", 500)
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to upload the file", 500)
		}

		receiver := r.FormValue("receiver")
		app.Fabric.InvokeShareFile(fileBytes, handler.Filename, user, receiver)
	}
	renderTemplate(w, r, "home.html", data)
}
