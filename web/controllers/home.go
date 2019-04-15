package controllers

import (
	"net/http"
	"io/ioutil"
)

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// helloValue, err := app.Fabric.QueryHello()
	// if err != nil {
	// 	http.Error(w, "Unable to query the blockchain", 500)
	// }

	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
	}{
		TransactionId: "",
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
