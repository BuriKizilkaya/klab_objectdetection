package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("&E(H+MbQeThWmZq4t7w!z%C*F-J@NcRf")
	store = sessions.NewCookieStore(key)
)

func render_index_template(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("templates/index.gohtml"))
	template.Execute(w, nil)
}

func handle_upload_file(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file-uploaded")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isFileAllowed(fileHeader.Filename) {
		http.Error(w, strings.TrimLeft(fileHeader.Filename, "."), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	saveFilePath := "./_output/" + fileHeader.Filename
	saveFile, err := os.OpenFile(saveFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer saveFile.Close()
	io.Copy(saveFile, file)

	session, err := store.Get(r, "GoLang-Object-Detection-App")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["uploaded_img_file_path"] = saveFilePath
	session.Save(r, w)
	render_index_template(w, r)
}

func root(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		render_index_template(w, r)
	case "POST":
		handle_upload_file(w, r)
	}
}

func isFileAllowed(fileName string) bool {
	allowedExtensions := [...]string{"jpeg", "jpg", "png"}

	for _, allowedExtension := range allowedExtensions {
		if strings.HasSuffix(fileName, allowedExtension) {
			return true
		}
	}

	return false
}

func main() {
	store.Options.MaxAge = 60 * 60 * 1 // 1 hour

	mux := http.NewServeMux()

	// setup handles
	mux.HandleFunc("/", root)

	log.Fatal(http.ListenAndServe(":9000", mux))
}
