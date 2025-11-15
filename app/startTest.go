package app

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/dou"
	"github.com/thisisaname1928/goParsingDocx/testsvr"
)

func startTestRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/startTest/")
}

func startTest(w http.ResponseWriter, r *http.Request) {
	b, e := os.ReadFile("./app/frontend/startTest/index.html")

	if e != nil {
		w.WriteHeader(404)
		return
	}

	w.Write(b)
}

func testInfoRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/startTest/testInfo/")
}

func testIn4(w http.ResponseWriter, r *http.Request) {
	b, e := os.ReadFile("./app/frontend/startTest/testInfo/index.html")

	if e != nil {
		w.WriteHeader(404)
		return
	}

	w.Write(b)
}

func uploadTestAPI(_ http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)

	os.WriteFile("./uploadedTest.dou", b, os.FileMode.Perm(0777))
}

func loadTestAPI(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var request struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}

	var response struct {
		Status bool   `json:"status"`
		Msg    string `json:"msg"`
	}

	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	t, e := testsvr.NewDouglasFir("8000", "./uploadedTest.dou", request.Key)
	if e != nil {
		response.Status = false

		if e.Error() == dou.ERROR_KEY_NOT_MATCH {
			response.Msg = "AUTH_FAILED"
		} else {
			response.Msg = "BAD_TEST_UPLOAD"
		}
		encoder.Encode(response)
		return
	}

	var info testInfoJson
	info.Name = request.Name
	info.Key = request.Key

	b, _ := json.Marshal(info)
	os.WriteFile("./testsvr/testdata/"+t.UUID+"/info.json", b, 0666)

	response.Status = true
	response.Msg = "ok"
	encoder.Encode(response)
}

func startTestAPI(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	switch v["NAME"] {
	case "getTestList":
		ListTestAPI(w, r)
	case "upload":
		uploadTestAPI(w, r)
	case "load":
		loadTestAPI(w, r)
	case "getTestInfo":
		getTestInfo(w, r)
	case "startATest":
		startATest(w, r)
	case "stopATest":
		stopATest(w, r)
	case "getTestIp":
		getTestIP(w, r)
	case "getCandinateList":
		getCandinateListAPI(w, r)
	}
}

func getCandinateListAPI(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var request struct {
		TestUUID string `json:"uuid"`
	}

	var response struct {
		Status     bool        `json:"status"`
		Msg        string      `json:"string"`
		Candinates []candinate `json:"candinates"`
	}

	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = "ERR_CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	can, e := getCandinateList(request.TestUUID)
	if e != nil {
		response.Status = false
		response.Msg = e.Error()
		encoder.Encode(response)
		return
	}

	response.Candinates = can
	response.Status = true
	response.Msg = "ok"
	encoder.Encode(response)

}

func startATest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	var request struct {
		TestUUID string `json:"uuid"`
	}

	var response struct {
		Status bool   `json:"status"`
		Msg    string `json:"msg"`
	}

	decoder.Decode(&request)

	e := testPool.OpenNewTest(request.TestUUID, getTestKey(request.TestUUID))

	if e != nil {
		response.Msg = e.Error()
		response.Status = false
		encoder.Encode(response)
		return
	}

	response.Msg = ""
	response.Status = true
	encoder.Encode(response)
}

func stopATest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	var request struct {
		TestUUID string `json:"uuid"`
	}

	var response struct {
		Status bool   `json:"status"`
		Msg    string `json:"msg"`
	}

	decoder.Decode(&request)

	testPool.CloseTest(request.TestUUID)

	response.Msg = ""
	response.Status = true
	encoder.Encode(response)
}

func getTestIP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var request struct {
		TestUUID string `json:"uuid"`
	}

	e := decoder.Decode(&request)

	if e != nil {
		w.Write([]byte("NaN"))
	} else {
		w.Write([]byte(testPool.GetServerIP(request.TestUUID)))
	}
}
