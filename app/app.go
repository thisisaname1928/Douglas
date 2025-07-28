package app

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func favicon(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/icon.ico")
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if str := detectFileExt("./app/icon.ico"); str != "" {
		w.Header().Add("Content-Type", str)
	} else {
		contentType := http.DetectContentType(f)
		w.Header().Add("Content-Type", contentType)
	}

	if e == nil {
		w.Write(f)
	}
}

func addResource(w http.ResponseWriter, r *http.Request, path string) {
	vars := mux.Vars(r)
	file, e := os.Open(path + vars["FILE"])
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if str := detectFileExt(vars["FILE"]); str != "" {
		w.Header().Add("Content-Type", str)
	} else {
		contentType := http.DetectContentType(f)
		w.Header().Add("Content-Type", contentType)
	}

	if e == nil {
		w.Write(f)
	}
}

func StartApp() {
	server := mux.NewRouter()
	server.HandleFunc("/LivePreview", livePreview)
	server.HandleFunc("/LivePreview/{FILE}", livePreviewRes)
	server.HandleFunc("/LivePreview/API/{NAME}", livePreviewAPI)
	server.HandleFunc("/favicon.ico", favicon)
	server.HandleFunc("/media/{FILE}", mediaRoute)
	server.HandleFunc("/Home", homePage)
	server.HandleFunc("/Home/{FILE}", homePageRes)
	fmt.Println("dia chi web app: http://localhost:8080")
	http.ListenAndServe("localhost:8080", server)
}
