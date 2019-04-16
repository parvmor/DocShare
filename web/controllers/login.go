package controllers

import (
	"crypto/aes"
	"crypto/rsa"
	"net/http"

	"github.com/gorilla/sessions"
	ipfs "github.com/ipfs/go-ipfs-api"
	"golang.org/x/crypto/bcrypt"
)

var (
	userpass = make(map[string][]byte)
	key      = []byte("super-secret-key")
	store    = sessions.NewCookieStore(key)
	// RSAKeySize is 2048 bits
	RSAKeySize = 2048
	// AESKeySize is 16 bytes
	AESKeySize = 16
	// BlockSize for AES
	BlockSize = aes.BlockSize
	keypair   = make(map[string]rsa.PrivateKey)
	aeskey    = make(map[string][]byte)
	shell     = ipfs.NewShell("0.0.0.0:5001")
)

// SignupHandler function
func (app *Application) SignupHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	if r.FormValue("submitted") == "submit" {
		user, pass := r.FormValue("username"), r.FormValue("password")
		passhash, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, ok := userpass[user]; ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		priv, err := GenerateRSAKey()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		keypair[user] = *priv
		aeskey[user] = Argon2Key([]byte(userpass[user]), key, uint32(AESKeySize))
		userpass[user] = passhash
		session.Values["authenticated"] = true
		session.Values["user"] = user
		session.Save(r, w)
		data := &struct {
			Success	bool
			Embed		string
		}{
			Success:	false,
			Embed:		"",
		}
		renderTemplate(w, r, "getfile.html", data)
	} else {
		renderTemplate(w, r, "signup.html", nil)
	}
}

// SigninHandler function
func (app *Application) SigninHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	if r.FormValue("submitted") == "submit" {
		user, pass := r.FormValue("username"), r.FormValue("password")
		if _, ok := userpass[user]; !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		passhash := userpass[user]
		// Compare the stored hashed password, with the hashed version of the password that was received
		if err := bcrypt.CompareHashAndPassword(passhash, []byte(pass)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		session.Values["authenticated"] = true
		session.Values["user"] = user
		session.Save(r, w)
		
		data := &struct {
			Success	bool
			Embed		string
		}{
			Success:	false,
			Embed:		"",
		}
		renderTemplate(w, r, "getfile.html", data)
	} else {
		renderTemplate(w, r, "default.html", nil)
	}
}

// SignoutHandler function
func (app *Application) SignoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Values["user"] = ""
	session.Save(r, w)
	renderTemplate(w, r, "default.html", nil)
}
