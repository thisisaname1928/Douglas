package app

import (
	"io"
	"net/http"
	"os"
)

func exportRouteRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/export/")
}

func exportRoute(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/frontend/export/index.html")
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if e == nil {
		w.Write(f)
	}
}
