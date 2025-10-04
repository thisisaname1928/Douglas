package testsvr

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/docx"
)

// copy to prevent circular import:)
func genUUID() string {
	id := uuid.New()
	return id.String()
}

type startTestRequest struct {
	Name  string `json:"name"`
	Class string `json:"className"`
}

type startTestResponse struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	UUID   string `json:"uuid"`
}

type testsvrInfo struct {
	Name      string          `json:"name"`
	Class     string          `json:"class"`
	IP        string          `json:"IP"`
	StartTime time.Time       `json:"startTime"`
	Questions []docx.Question `json:"questions"`
}

func (fir DouglasFir) handleStartTest(w http.ResponseWriter, r *http.Request) {
	var request startTestRequest
	var response startTestResponse

	encoder := json.NewEncoder(w)

	decoder := json.NewDecoder(r.Body)
	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = fmt.Sprintf("%v", e)
		encoder.Encode(&response)
		return
	}

	fmt.Printf("Student %v class: %v register\n", request.Name, request.Class)

	uuid := genUUID()
	response.Msg = "ok"

	// create a sub test
	var info testsvrInfo
	info.Class = request.Class
	info.Name = request.Name
	info.Questions = fir.ShuffleNewTest().Test
	info.StartTime = time.Now()

	// save IP
	IP, _, _ := net.SplitHostPort(r.RemoteAddr)
	info.IP = IP

	f, e := os.Create("./testsvr/testdata/" + fir.UUID + "/testdat/" + uuid + ".json")
	if e != nil {
		response.Msg = fmt.Sprintf("internal error: %v", e)
		encoder.Encode(&response)
		return
	}
	b, _ := json.Marshal(&info)
	f.Write(b)
	f.Close()

	response.Status = true
	response.UUID = uuid

	encoder.Encode(&response)
}

func (fir DouglasFir) testsvrAPI(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	switch v["NAME"] {
	case "startTest":
		fir.handleStartTest(w, r)
	case "getTest":
		fir.handleGetTest(w, r)
	}
}
