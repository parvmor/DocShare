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
		Embed:				 "",
	}

	if r.FormValue("submitted") == "true" {
		fileBytes, err := app.Fabric.QueryGetFile(r.FormValue("filename"), user)
		if err != nil {
			http.Error(w, "Unable to query Blockchain", 500)
		}

		iv := fileBytes[:BlockSize]
		fileBytes = fileBytes[BlockSize:]
		stream := CFBDecrypter(aeskey[user], iv);
		stream.XORKeyStream(fileBytes, fileBytes);
		fileString := b64.StdEncoding.EncodeString(fileBytes)

		data.Success = true
		data.Embed = "<embed src=data:application/pdf;base64," + string(fileString) + ">"
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
			http.Error(w, "Unable to query Blockchain", 500)
		}

		value := RSADecrypt(keypair[user], fileBytes, []byte("sharing"))
		ek := value[:AESKeySize]
		cid := value[AESKeySize:]

		reader, err := shell.Cat(string(cid))
		if err != nil {
			return nil, err
		}
		fileBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		iv := fileBytes[:BlockSize]
		fileBytes = fileBytes[BlockSize:]
		stream := CFBDecrypter(ek, iv);
		stream.XORKeyStream(fileBytes, fileBytes);
		fileString := b64.StdEncoding.EncodeString(fileBytes)

		data.Success = true
		data.Embed = "<embed src=data:application/pdf;base64," + fileString + ">"

	}
	renderTemplate(w, r, "receivefile.html", data)
}
