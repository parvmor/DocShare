package controllers

import (
	"bytes"
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

		// Encrypt the file bytes
		iv := RandomBytes(BlockSize)
		ciphertext := make([]byte, len(fileBytes))
		cipher := CFBEncrypter(aeskey[user], iv)
		cipher.XORKeyStream(ciphertext, fileBytes)
		value := append(iv, ciphertext...)

		// Put them in the blockchain
		txnid, err := app.fabric.InvokePutFile(value, handler.Filename, user)

		if err != nil {
			http.Error(w, "Unable to query Blockchain", 500)
		}

		data.TransactionID = txnid
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

		// Generate a random aeskey
		ek := RandomBytes(AESKeySize)
		iv := RandomBytes(BlockSize)
		// encrpyt the cipher text using it
		ciphertext := make([]byte, len(fileBytes))
		cipher := CFBEncrypter(ek, iv)
		cipher.XORKeyStream(ciphertext, fileBytes)
		value := append(iv, ciphertext...)
		// Put it in IPFS
		cid, err := shell.Add(bytes.NewReader(value))
		if err != nil {
			http.Error(w, "Unable to upload the file", 500)
		}
		// Create new value
		value = append(ek, cid...)
		// Encrpyt it using public key of receiver
		receiver := r.FormValue("receiver")
		pubkey := keypair[receiver].PublicKey
		sharingdata, err := RSAEncrypt(&pubkey, value, []byte("sharing"))
		if err != nil {
			http.Error(w, "Unable to upload the file", 500)
		}

		app.fabric.InvokeShareFile(sharingdata, handler.Filename, user, receiver)

		receiver := r.FormValue("receiver")
		txnid, err := app.Fabric.InvokeShareFile(fileBytes, handler.Filename, user, receiver)
		if err != nil {
			http.Error(w, "Unable to query Blockchain", 500)
		}
		data.TransactionID = txnid
		data.Success = true
		data.Response = true

	}
	renderTemplate(w, r, "sharefile.html", data)
}
