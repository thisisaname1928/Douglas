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
	UUID       string
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

	fir.UUID = uuid
	os.Mkdir("./testsvr/testdata/"+uuid+"/testdat", 0755)

	return fir, nil
}

func isExist(path string) bool {
	_, e := os.Stat(path)

	return !os.IsNotExist(e)
}

func OpenOldTest(uuid string, key string) (DouglasFir, error) {
	if !isExist(fmt.Sprintf("./testsvr/testdata/%v", uuid)) {
		return DouglasFir{}, errors.New(ERROR_FIR_NOT_CREATED)
	}
	if !isExist(fmt.Sprintf("./testsvr/testdata/%v/test.dou", uuid)) {
		return DouglasFir{}, errors.New(ERROR_FIR_NOT_CREATED)
	}

	var fir DouglasFir
	fir.UUID = uuid

	var df dou.DouFile
	df, e := dou.Open(fmt.Sprintf("./testsvr/testdata/%v/test.dou", uuid), key)
	if e != nil {
		return DouglasFir{}, e
	}
	fir.Douglas = df

	fir.Created = true

	if !isExist(fmt.Sprintf("./testsvr/testdata/%v/testdat", uuid)) {
		os.Mkdir(fmt.Sprintf("./testsvr/testdata/%v/testdat", uuid), 0755)
	}

	return fir, nil
}

func (fir DouglasFir) OpenServer(port string) error {
	if !fir.Created {
		return errors.New(ERROR_FIR_NOT_CREATED)
	}
	server := mux.NewRouter()

	server.HandleFunc("/", route)

	fir.HttpServer = &http.Server{Addr: "0.0.0.0:" + port, Handler: server}

	return fir.HttpServer.ListenAndServe()
}

func (fir DouglasFir) CloseServer() {
	fir.HttpServer.Close()
}

func route(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sus"))
}
