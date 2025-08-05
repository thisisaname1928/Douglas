package testsvr

import (
	"errors"
	"net/http"

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
