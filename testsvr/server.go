package testsvr

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/dou"
)

const (
	ERROR_FIR_NOT_CREATED = "ERROR_FIR_NOT_CREATED"
)

type DouglasFir struct {
	ServerPort string
	Douglas    dou.DouFile // test file
	Created    bool        // check for if init success
	HttpServer *http.Server
}

func copyFile(dest string, src string) {
	f, e := os.ReadFile(src)

	if e != nil {
		fmt.Println("internal error: " + fmt.Sprintf("%v", e))
	}

	e = os.WriteFile(dest, f, 0755)

	if e != nil {
		fmt.Println("internal error: " + fmt.Sprintf("%v", e))
	}
}

// create new test server
func NewDouglasFir(serverPort string, path string, key string) (DouglasFir, error) {
	var fir DouglasFir
	fir.ServerPort = serverPort
	fir.Created = false

	var df dou.DouFile
	df, e := dou.Open(path, key)
	if e != nil {
		return fir, e
	}

	fir.Douglas = df

	fir.Created = true

	// create new data folder
	uuid := uuid.New().String()
	e = os.Mkdir("./testsvr/testdata/"+uuid, 0755)
	if e != nil {
		fmt.Println("internal error: " + fmt.Sprintf("%v", e))
	}

	// copy a backup .dou file into it
	copyFile("./testsvr/testdata/"+uuid+"/test.dou", path)

	return fir, nil
}

func (fir DouglasFir) OpenServer() error {
	if !fir.Created {
		return errors.New(ERROR_FIR_NOT_CREATED)
	}
	server := mux.NewRouter()

	server.HandleFunc("/", route)

	fir.HttpServer = &http.Server{Addr: "0.0.0.0:" + fir.ServerPort, Handler: server}

	return fir.HttpServer.ListenAndServe()
}

func (fir DouglasFir) CloseServer() {
	fir.HttpServer.Close()
}

func route(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sus"))
}
