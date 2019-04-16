package controllers

import (
	"net/http"
)

// RequestHandler function
func (app *Application) GetFileHandler(w http.ResponseWriter, r *http.Request) {
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
		Embed					string
	}{
		TransactionID: "",
		Success:       false,
		Response:      false,
		Embed:				 ""
	}
	if r.FormValue("submitted") == "true" {
		txnid, err := app.Fabric.QueryGetFile(r.FormValue("filename"), user)
		if err != nil {
			http.Error(w, "Unable to query Blockchain", 500)
		}
		data.TransactionId = txnid
		data.Success = true
		data.Response = true
		data.Embed = 
	}
	renderTemplate(w, r, "request.html", data)
}
