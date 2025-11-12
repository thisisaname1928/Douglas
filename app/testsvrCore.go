package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func checkIsTestFolder(path string) bool {
	r, e := os.Stat(path)

	if e != nil {
		return false
	}

	// check if there is a test.dou in it
	if r.IsDir() {
		r1, e := os.Stat(path + "/test.dou")
		if e != nil {
			return false
		}

		if !r1.IsDir() {
			return true
		}
	}

	return false
}

func ListTestAPI(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	res, e := listTest()

	var response struct {
		Status bool       `json:"status"`
		Msg    string     `json:"msg"`
		List   []testInfo `json:"list"`
	}

	if e != nil {
		response.Status = false
		response.Msg = fmt.Sprint(e)
		encoder.Encode(response)
		return
	}

	response.Status = true
	response.Msg = "ok"
	response.List = res
	encoder.Encode(response)
}

func getNumberOfSubFolder(path string) int {
	f, e := os.ReadDir(path)

	if e != nil {
		return 0
	}

	return len(f)
}

type testInfo struct {
	TestUUID          string `json:"uuid"`
	NumberOfCandinate int    `json:"numberOfCandinate"`
	Name              string `json:"name"`
}

type testInfoJson struct {
	Name string `json:"name"`
}

func getTestName(path string) string {
	var testInf testInfoJson

	b, e := os.ReadFile(path + "/info.json")
	if e != nil {
		return ""
	}

	json.Unmarshal(b, &testInf)
	return testInf.Name
}

func getTestInfo(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var request struct {
		UUID string `json:"uuid"`
	}

	decoder.Decode(&request)

	var path = request.UUID
	var testInf = testInfo{path, getNumberOfSubFolder("./testsvr/testdata/" + path + "/testdat"), getTestName("./testsvr/testdata/" + path)}

	encoder.Encode(testInf)
}

func listTest() ([]testInfo, error) {
	f, e := os.ReadDir("./testsvr/testdata")

	if e != nil {
		return []testInfo{}, e
	}

	res := []testInfo{}

	// loop through and check
	for _, v := range f {
		if checkIsTestFolder("./testsvr/testdata/" + v.Name()) {
			var tmp = testInfo{v.Name(), getNumberOfSubFolder("./testsvr/testdata/" + v.Name() + "/testdat"), getTestName("./testsvr/testdata/" + v.Name())}
			res = append(res, tmp)
		}
	}

	return res, nil
}
