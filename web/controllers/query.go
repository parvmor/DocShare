package controllers

import (
  "io"
  "io/ioutil"
	"net/http"
  "bytes"
  "fmt"
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
		Embed:				 "",
	}

	if r.FormValue("submitted") == "true" {
		fileBytes, err := app.Fabric.QueryGetFile(r.FormValue("filename"), user)
		if err != nil {
			http.Error(w, "Unable to query Blockchain 1", 500)
		}

		iv := fileBytes[:BlockSize]
		fileBytes = fileBytes[BlockSize:]
		stream := CFBDecrypter(aeskey[user], iv);
		stream.XORKeyStream(fileBytes, fileBytes);
    w.Header().Set("Content-type", "application/pdf")
    io.Copy(w, bytes.NewReader(fileBytes))
    return
  }

	renderTemplate(w, r, "getfile.html", data)
}

// ReceiveFileHandler function
func (app *Application) ReceiveFileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user := session.Values["user"].(string)
	data := &struct {
		Success	bool
		Embed		string
	}{
		Success:	false,
		Embed:		"",
	}

	if r.FormValue("submitted") == "true" {
		fileBytes, err := app.Fabric.QueryGetFile(r.FormValue("sharer") + "_" + r.FormValue("filename"), user)

		if err != nil {
			http.Error(w, "Unable to query Blockchain 1", 500)
            return
		}

    priv := keypair[user]
		value, err := RSADecrypt(&priv, fileBytes, []byte("sharing"))
		if err != nil {
			http.Error(w, "Unable to query Blockchain 2", 500)
            return
		}

		ek := value[:AESKeySize]
		cid := value[AESKeySize:]
		reader, err := shell.Cat(string(cid))
		if err != nil {
			http.Error(w, "Unable to query Blockchain 3", 500)
			return
		}
		fileBytes, err = ioutil.ReadAll(reader)
		if err != nil {
			http.Error(w, "Unable to query Blockchain 4", 500)
			return
		}

		iv := fileBytes[:BlockSize]
		fileBytes = fileBytes[BlockSize:]
		stream := CFBDecrypter(ek, iv);
		stream.XORKeyStream(fileBytes, fileBytes);
    w.Header().Set("Content-type", "application/pdf")
    io.Copy(w, bytes.NewReader(fileBytes))
    return

	}
	renderTemplate(w, r, "receivefile.html", data)
}
