package app

import (
	"encoding/json"
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
	v := mux.Vars(r)
	_, e := os.Stat("./app/tests/" + v["UUID"] + ".dat")

	if e != nil {
		w.Write([]byte("NO SUCH FILE OR DIR"))
		return
	}

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

type GetConfigRequest struct {
	UUID string `json:"UUID"`
}

func exportAPI(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	// client should upload file with uuid then using ExportRequest with uuid

	switch v["NAME"] {
	case "getConfig":
		var request GetConfigRequest
		var response ExportConfigResponse
		decoder := json.NewDecoder(r.Body)
		e := decoder.Decode(&request)

		if e != nil {
			response.Status = false
			response.Msg = "invalid request"
		} else {
			response, _ = getExportConfig(request.UUID)
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(&response)
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

type StypeN struct {
	Stype string `json:"stype"`
	N     uint64 `json:"N"`
}

type ExportConfigResponse struct {
	Status            bool     `json:"status"`
	Msg               string   `json:"msg"`
	NumberOfQuestions uint64   `json:"numberOfQuestions"`
	Stype             []StypeN `json:"stype"`
}

func search4Stype4Cfg(ques []StypeN, name string) int {
	for i, v := range ques {
		if v.Stype == name {
			return i
		}
	}

	return -1
}

// function to get some info about test
func getExportConfig(UUID string) (ExportConfigResponse, error) {
	var response ExportConfigResponse
	res, e := GenQues("./app/tests/" + UUID + ".dat")
	if e != nil {
		response.Status = false
		response.Msg = fmt.Sprintf("%v", e)

		return response, e
	}

	response.Status = true
	response.NumberOfQuestions = uint64(len(res))
	for _, v := range res {
		curIdx := search4Stype4Cfg(response.Stype, v.Stype)
		if curIdx != -1 { // there are some questions already have that stype in response
			response.Stype[curIdx].N++
		} else {
			response.Stype = append(response.Stype, StypeN{v.Stype, 1})
		}
	}

	return response, nil
}
