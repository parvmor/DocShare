package controllers

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"net/http"
	_ "github.com/lib/pq"
)

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

func (app *Application) SignupHandler(w http.ResponseWriter, r *http.Request){

	if r.FormValue("submitted") == "submit" {
		// Parse and decode the request body into a new `Credentials` instance
		creds := Credentials{r.FormValue("password"), r.FormValue("username")}

		// Salt and hash the password using the bcrypt algorithm
		// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

		// Next, insert the username, along with the hashed password into the database
		if _, err = db.Query("insert into users values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
			// If there is any issue with inserting into the database, return a 500 error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := &struct {
			Hello string
		}{
			Hello: "Hi There",
		}

		renderTemplate(w, r, "home.html", data)
	} else {
			renderTemplate(w, r, "signup.html", struct{}{})
	}

}

func (app *Application) SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("submitted") == "submit" {
		// Parse and decode the request body into a new `Credentials` instance
		creds := Credentials{r.FormValue("password"), r.FormValue("username")}
		// Get the existing entry present in the database for the given username
		result := db.QueryRow("select password from users where username=$1", creds.Username)
		if err != nil {
			// If there is an issue with the database, return a 500 error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// We create another instance of `Credentials` to store the credentials we get from the database
		storedCreds := &Credentials{}
		// Store the obtained password in `storedCreds`
		err = result.Scan(&storedCreds.Password)
		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// If the error is of any other type, send a 500 status
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
			// If the two passwords don't match, return a 401 status
			w.WriteHeader(http.StatusUnauthorized)
		}

		data := &struct {
			Hello string
		}{
			Hello: "Hi There",
		}

		renderTemplate(w, r, "home.html", data)
	} else {
		renderTemplate(w, r, "signup.html", nil)
	}

}
