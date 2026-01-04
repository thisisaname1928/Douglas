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
	"github.com/thisisaname1928/goParsingDocx/dou"
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
	Name                 string                      `json:"name"`
	Class                string                      `json:"class"`
	IP                   string                      `json:"IP"`
	StartTime            time.Time                   `json:"startTime"`
	EndTime              time.Time                   `json:"endTime"`
	Done                 bool                        `json:"done"`
	MaxTestDuration      uint64                      `json:"testDuration"`
	Questions            []docx.Question             `json:"questions"`
	AnswerSheet          [][]string                  `json:"answerSheet"`
	WarnTimes            int                         `json:"warnTimes"`
	AdditionalExportData dou.DouAdditionalExportData `json:"additionalExportData"`
}

// I forgot to export it:)
type TestsvrInfo testsvrInfo

func (fir *DouglasFir) handleStartTest(w http.ResponseWriter, r *http.Request) {
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
	info.Done = false
	info.MaxTestDuration = fir.Douglas.Data.TestDuration
	info.AdditionalExportData = fir.Douglas.Data.AdditionalExportData

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

	fir.TestSessions.NewSession(uuid, info.IP, info.StartTime, len(info.Questions))

	response.Status = true
	response.UUID = uuid

	encoder.Encode(&response)
}

func (fir *DouglasFir) getSchoolName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fir.ExtraInfo.SchoolName))
}

func (fir *DouglasFir) testsvrAPI(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	switch v["NAME"] {
	case "startTest":
		fir.handleStartTest(w, r)
	case "getTest":
		fir.handleGetTest(w, r)
	case "updateAnswer":
		fir.handleUpdateAnswerSheet(w, r)
	case "handleDoneTest":
		fir.handleDoneTest(w, r)
	case "getTestStatus":
		fir.getTestStatus(w, r)
	case "getPoint":
		fir.getTestPoint(w, r)
	case "getCurrentAnsSheet":
		fir.getCurrentAnsSheet(w, r)
	case "getCurrentServerTime":
		fir.getCurrentServerTime(w, r)
	case "isAdmin":
		fir.isAdmin(w, r)
	case "getTestName":
		fir.getTestName(w, r)
	case "getSchoolName":
		fir.getSchoolName(w, r)
	case "warn":
		fir.warnAPI(w, r)
	}
}
