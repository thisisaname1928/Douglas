package testsvr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func (fir DouglasFir) getTest(uuid string) ([]byte, error) {
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

func (fir DouglasFir) handleGetTest(w http.ResponseWriter, r *http.Request) {

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
