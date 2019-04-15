package web

import (
	"database/sql"
	"fmt"
	"github.com/parvmor/docshare/web/controllers"
	"net/http"
	_ "github.com/lib/pq"
)

const hashCost = 8
var db *sql.DB

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/signin", app.SigninHandler)
	http.HandleFunc("/signup", app.SignupHandler)
	http.HandleFunc("/home", app.HomeHandler)
	http.HandleFunc("/request", app.RequestHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
	})

	initDB()

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}

func initDB(){
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("postgres", "dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
}
