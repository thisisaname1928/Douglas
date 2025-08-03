package app

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/dou"
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

func exportConfigRouteRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/export/config/")
}

func exportConfigRoute(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/frontend/export/config/index.html")
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if e == nil {
		w.Write(f)
	}
}

type ExportRequest struct {
	Author           string              `json:"author"`
	UseEncryption    bool                `json:"useEncryption"`
	Key              string              `json:"key"`
	UseTestStructure string              `json:"useTestStructure"`
	TestStruct       []dou.TestStructure `json:"testStructure"`
}

// gen a v4 uuid
func genUUID() string {
	id := uuid.New()
	return id.String()
}

func exportAPI(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	// client should upload file with uuid then using ExportRequest with uuid

	switch v["NAME"] {
	case "export":
	case "upload":
		f, e := os.Create("./app/tests/" + r.Header.Get("uuid") + ".dat")
		if e != nil {
			fmt.Println(e)
		}
		defer f.Close()

		dat, e := io.ReadAll(r.Body)
		if e != nil {
			fmt.Println(e)
		}
		f.Write(dat)
	case "genUUID":
		w.Write([]byte(genUUID()))
		w.Header().Add("Content-Type", "text/plain")
	}
}
