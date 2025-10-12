package testsvr

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/dou"
)

const (
	ERROR_FIR_NOT_CREATED = "ERROR_FIR_NOT_CREATED"
)

type DouglasFir struct {
	ServerPort   string
	Douglas      dou.DouFile // test file
	Created      bool        // check for if init success
	UUID         string
	HttpServer   *http.Server
	TestSessions TestSessions
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
func NewDouglasFir(serverPort string, path string, key string) (*DouglasFir, error) {
	var fir DouglasFir
	fir.ServerPort = serverPort
	fir.Created = false

	var df dou.DouFile
	df, e := dou.Open(path, key)
	if e != nil {
		return &fir, e
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

	return &fir, nil
}

func isExist(path string) bool {
	_, e := os.Stat(path)

	return !os.IsNotExist(e)
}

func OpenOldTest(uuid string, key string) (*DouglasFir, error) {
	if !isExist(fmt.Sprintf("./testsvr/testdata/%v", uuid)) {
		return nil, errors.New(ERROR_FIR_NOT_CREATED)
	}
	if !isExist(fmt.Sprintf("./testsvr/testdata/%v/test.dou", uuid)) {
		return nil, errors.New(ERROR_FIR_NOT_CREATED)
	}

	var fir DouglasFir
	fir.UUID = uuid

	var df dou.DouFile
	df, e := dou.Open(fmt.Sprintf("./testsvr/testdata/%v/test.dou", uuid), key)
	if e != nil {
		return nil, e
	}
	fir.Douglas = df

	fir.Created = true

	if !isExist(fmt.Sprintf("./testsvr/testdata/%v/testdat", uuid)) {
		os.Mkdir(fmt.Sprintf("./testsvr/testdata/%v/testdat", uuid), 0755)
	}

	return &fir, nil
}

func GetIp() (string, error) {
	it, e := net.Interfaces()
	fmt.Println(e)

	for _, v := range it {
		if v.Flags&net.FlagUp == 0 || v.Flags&net.FlagLoopback != 0 || v.Flags&net.FlagMulticast == 0 {
			continue
		}

		// detect wireless lan network interface
		if runtime.GOOS == "linux" {
			if !strings.HasPrefix(v.Name, "wl") {
				continue
			}
		} else if runtime.GOOS == "darwin" {
			if !strings.HasPrefix(v.Name, "en0") {
				continue
			}
		} else if runtime.GOOS == "windows" {
			if !strings.HasPrefix(v.Name, "Wireless LAN adapter	Wi-Fi") {
				continue
			}
		}

		addrs, e := v.Addrs()
		if e != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// not what we are finding
			if ip == nil || ip.IsLoopback() {
				continue
			}

			if ip.IsPrivate() {
				return ip.String(), nil
			}
		}
	}

	return "", errors.New("cant detect wlan ip")
}

func (fir *DouglasFir) OpenServer(port string) error {
	// init testSession
	fir.TestSessions.Init()

	if !fir.Created {
		return errors.New(ERROR_FIR_NOT_CREATED)
	}
	server := mux.NewRouter()

	// ROUTING
	server.HandleFunc("/", route)
	server.HandleFunc("/rsrc/{FILE}", res)
	server.HandleFunc("/favicon.ico", favicon)
	server.HandleFunc("/api/{NAME}", fir.testsvrAPI)
	server.HandleFunc("/taketest/media/{FILE}", fir.mediaRoute)
	server.HandleFunc("/taketest/{UUID}", fir.takeTestRoute)

	fir.HttpServer = &http.Server{Addr: "0.0.0.0:" + port, Handler: server}

	return fir.HttpServer.ListenAndServe()
}

func (fir *DouglasFir) CloseServer() {
	fir.HttpServer.Close()
}

func detectFileExt(path string) string {
	fileExtSpl := strings.Split(path, ".")
	ext := fileExtSpl[len(fileExtSpl)-1]

	switch ext {
	case "js":
		return "text/javascript"
	case "css":
		return "text/css"
	case "ico":
		return "image/x-icon"
	default:
		return ""
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

func route(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./testsvr/frontend/taketest/index.html")
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if e == nil {
		w.Write(f)
	}
}

func (fir *DouglasFir) takeTestRoute(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./testsvr/frontend/realtaketest/index.html")
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if e == nil {
		w.Write(f)
	}
}

func res(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./testsvr/frontend/resources/")
}

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
