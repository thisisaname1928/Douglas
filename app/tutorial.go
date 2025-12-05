package app

import (
	"net/http"
	"os"
)

func tutorialPageRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/tutorial/")
}

func tutorialPage(w http.ResponseWriter, r *http.Request) {
	b, e := os.ReadFile("./app/frontend/tutorial/index.html")

	if e != nil {
		w.WriteHeader(404)
		return
	}

	w.Write(b)
}
