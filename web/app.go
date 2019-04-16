package web

import (
	"fmt"
	"net/http"

	"github.com/parvmor/docshare/web/controllers"
)

// Serve an HTTP server
func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/signin", app.SigninHandler)
	http.HandleFunc("/signup", app.SignupHandler)
	http.HandleFunc("/putfile", app.PutFileHandler)
	http.HandleFunc("/getfile", app.GetFileHandler)
	http.HandleFunc("/sharefile", app.ShareFileHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
