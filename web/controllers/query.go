package controllers

import (
	"net/http"
	b64 "encoding/base64"
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
		Success       bool
		Embed					string
	}{
		Success:       false,
		Embed:				 ""
	}

	if r.FormValue("submitted") == "true" {
		fileBytes, err := app.Fabric.QueryGetFile(r.FormValue("filename"), user)
		if err != nil {
			http.Error(w, "Unable to query Blockchain", 500)
		}

		fileString := b64.StdEncoding.DecodeString(string(fileBytes))

		data.Success = true
		data.Embed = "<embed src=data:application/pdf;base64," + fileString + ">"
	}

	renderTemplate(w, r, "getfile.html", data)
}
