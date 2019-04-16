package controllers

import (
	"io/ioutil"
	"net/http"
	b64 "encoding/base64"
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

		fileBytes = []byte(b64.StdEncoding.EncodeToString(fileBytes))

		txnid, err := app.Fabric.InvokePutFile(fileBytes, handler.Filename, user)
		if err != nil {
			http.Error(w, "Unable to query Blockchain", 500)
		}

		data.TransactionId = txnid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "putfile.html", data)
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

		fileBytes = []byte(b64.StdEncoding.EncodeToString(fileBytes))

		receiver := r.FormValue("receiver")
		txnid, err := app.Fabric.InvokeShareFile(fileBytes, handler.Filename, user, receiver)
		if err != nil {
			http.Error(w, "Unable to query Blockchain", 500)
		}
		data.TransactionId = txnid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "sharefile.html", data)
}
