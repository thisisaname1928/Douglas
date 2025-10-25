package testsvr

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func getIP(r *http.Request) string {
	requestIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	return requestIP
}

func (fir *DouglasFir) verifyIP(uuid string, IP string) bool {
	fir.TestSessions.mutex.Lock()
	defer fir.TestSessions.mutex.Unlock()
	test, ok := fir.TestSessions.SessionsData[uuid]
	if !ok {
		return false
	}

	return test.IP == IP
}

func (fir *DouglasFir) getTest(uuid string) ([]byte, error) {
	b, e := os.ReadFile(fir.getTestDataPath() + uuid + ".json")
	if e != nil {
		return []byte{}, e
	}

	return b, nil
}

type getTestRequest struct {
	UUID string `json:"uuid"`
}

type getTestResponse struct {
	UUID     string      `json:"uuid"`
	Status   bool        `json:"status"`
	Msg      string      `json:"msg"`
	TestInfo testsvrInfo `json:"test"`
}

func (fir *DouglasFir) handleGetTest(w http.ResponseWriter, r *http.Request) {
	var response getTestResponse
	encoder := json.NewEncoder(w)
	response.Status = false

	decoder := json.NewDecoder(r.Body)

	var req getTestRequest
	decoder.Decode(&req)

	// get testdat
	s, e := fir.getTest(req.UUID)
	if e != nil {
		response.Msg = fmt.Sprint(e)
		encoder.Encode(&response)
		return
	}

	// unmarshal into testsvrInfo
	var test testsvrInfo
	e = json.Unmarshal(s, &test)
	if e != nil {
		response.Msg = fmt.Sprint(e)
		encoder.Encode(&response)
		return
	}

	// create new test session
	if !fir.TestSessions.CheckSession(req.UUID) {
		fir.TestSessions.NewSession(req.UUID, test.IP, test.StartTime, fir.GetNumberOfQuestions())

		// reload answer from testdats
		for i := range test.AnswerSheet {
			fir.TestSessions.UpdateAnswerSheet(i, req.UUID, [4]string(test.AnswerSheet[i]))
		}

		// if test already done
		if test.Done {
			fir.TestSessions.LockSession(req.UUID)
		}
	}

	// remove test answer
	for i := range test.Questions {
		for j := range test.Questions[i].TrueAnswer {
			test.Questions[i].TrueAnswer[j] = false
		}

		for j := range test.Questions[i].TLNA {
			test.Questions[i].TLNA[j] = ""
		}
	}

	response.Status = true
	response.TestInfo = test
	response.UUID = req.UUID
	encoder.Encode(&response)
}

// update answer sheet api

type updateAnswerSheetRequest struct {
	UUID        string   `json:"UUID"`
	Index       int      `json:"index"`
	AnswerIndex int      `json:"answerIndex"`
	AnswerData  string   `json:"data"`
	ShouldClear bool     `json:"shouldClear"`
	AnswerSheet []string `json:"answerSheet"`
}

type updateAnswerSheetResponse struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
}

func (fir *DouglasFir) handleUpdateAnswerSheet(w http.ResponseWriter, r *http.Request) {

	var request updateAnswerSheetRequest
	var response updateAnswerSheetResponse
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	// check if test sessions is locked
	if fir.TestSessions.CheckSession(request.UUID) {
		response.Status = false
		response.Msg = "TEST_SESSION_LOCKED"
		encoder.Encode(response)
		return
	}

	// check if AnswerSheet valid
	// if len(request.AnswerSheet) != 4 {
	// 	response.Status = false
	// 	response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
	// 	encoder.Encode(response)
	// 	return
	// }

	// check if answer data isn't bad
	if len(request.AnswerData) != 1 || request.AnswerIndex > 3 {
		response.Status = false
		response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	// COMPARING IP
	requestIP, _, _ := net.SplitHostPort(r.RemoteAddr)

	if !fir.verifyIP(request.UUID, requestIP) {
		response.Status = false
		response.Msg = "TEST_ACCESS_DENIED"
		encoder.Encode(response)
		return
	}

	if request.Index > fir.GetNumberOfQuestions() {
		response.Status = false
		response.Msg = "OUT_OF_RANGE"
		encoder.Encode(response)
		return
	}

	// update
	if request.ShouldClear { // update and clear all of the
		fir.TestSessions.UpdateAnswerSheet(request.Index, request.UUID, [4]string(request.AnswerSheet))
	} else {
		fir.TestSessions.UpdateAnswer(request.UUID, request.Index, request.AnswerIndex, request.AnswerData)
	}

	response.Status = true
	response.Msg = "ok"
	encoder.Encode(response)
}

// done test

type doneTestRequest struct {
	UUID string `json:"uuid"`
}

type doneTestResponse struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
}

func (fir *DouglasFir) handleDoneTest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	var request doneTestRequest
	var response doneTestResponse
	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	requestIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	if !fir.verifyIP(request.UUID, requestIP) {
		response.Status = false
		response.Msg = "TEST_ACCESS_DENIED"
		encoder.Encode(response)
		return
	}

	// verifyIP just done checking if testSession availble for us, so we dont need to check
	fir.TestSessions.DoneSession(fir.UUID, request.UUID, time.Now())
	response.Status = true
	response.Msg = "ok"
	encoder.Encode(response)

	fmt.Print("PPP: ")
	fmt.Println(fir.CalculateMark(request.UUID))
}

func (fir *DouglasFir) getTestStatus(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	// they are the same so
	var request doneTestRequest
	var response doneTestResponse
	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	response.Status = fir.CheckIfTestDone(request.UUID)
	response.Msg = "ok"
	encoder.Encode(response)
}

type getPointResponse struct {
	Status        bool    `json:"status"`
	Msg           string  `json:"msg"`
	Point         float64 `json:"point"`
	TrueQuesCount int     `json:"trueQuesCount"`
}

func (fir *DouglasFir) getTestPoint(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	// they are the same so
	var request doneTestRequest
	var response getPointResponse
	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	if !fir.CheckIfTestDone(request.UUID) {
		response.Status = false
		response.Msg = "TEST_NOT_DONE"
		encoder.Encode(response)
		return
	}

	c, p, e := fir.CalculateMark(request.UUID)
	if e != nil {
		response.Status = false
		response.Msg = "TEST_NOT_DONE"
		encoder.Encode(response)
		return
	}

	response.Status = true
	response.TrueQuesCount = c
	response.Msg = "ok"
	response.Point = p
	encoder.Encode(response)
}

type getCurrentAnsSheetResponse struct {
	Status   bool       `json:"status"`
	Msg      string     `json:"msg"`
	AnsSheet [][]string `json:"ansSheet"`
}

func (fir *DouglasFir) getCurrentAnsSheet(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var request doneTestRequest
	var response getCurrentAnsSheetResponse
	e := decoder.Decode(&request)
	if e != nil {
		response.Status = false
		response.Msg = "CLIENT_MAKE_A_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	if !fir.verifyIP(request.UUID, getIP(r)) {
		response.Status = false
		response.Msg = "ACCESS_DENIED"
		encoder.Encode(response)
		return
	}

	response.Status = true
	response.Msg = "ok"
	response.AnsSheet, e = fir.TestSessions.CopyAnsSheet(request.UUID)

	if e != nil {
		response.Status = false
		response.Msg = "ACCESS_DENIED"
	}

	encoder.Encode(response)
}

// func (fir *DouglasFir) startSession(w http.ResponseWriter, r *http.Request) {
// 	encoder := json.NewEncoder(w)
// 	decoder := json.NewDecoder(r.Body)

// 	var request doneTestRequest
// 	var response doneTestResponse

// }
