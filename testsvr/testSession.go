package testsvr

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// this for temporary saving sessions data
type TestSessionData struct {
	AnswerSheet [][]string `json:"AnswerSheet"` // Ans[0]=["1", "1", "1", "1"] example
	IP          string     `json:"IP"`
	StartTime   time.Time  `json:"StartTime"`
}

// access Sessions data by uuid
// Session Data should be saved and destroyed after the test is done
// Resume session by lookup this TestSessions
type TestSessions struct {
	SessionsData map[string]TestSessionData
}

func (session *TestSessions) Init() {
	session.SessionsData = make(map[string]TestSessionData)
}

func (session *TestSessions) NewSession(UUID string, IP string, startTime time.Time, numberOfQuestions int) {
	session.SessionsData[UUID] = TestSessionData{make([][]string, numberOfQuestions), IP, startTime}

	// assign 4 element to Answer sheet
	for i := 0; i < numberOfQuestions; i++ {
		session.SessionsData[UUID].AnswerSheet[i] = make([]string, 4)
	}

	fmt.Println("create Session for ", UUID)
}

func (session *TestSessions) UpdateAnswerSheet(i int, UUID string, Answer [4]string) {
	if i < len(session.SessionsData[UUID].AnswerSheet) {
		for k := 0; k < 4; k++ {
			session.SessionsData[UUID].AnswerSheet[i][k] = Answer[k] // do a copy
		}
	}
}

func (session *TestSessions) DoneSession(testUUID string, UUID string, endTime time.Time) {
	// check if session available
	if _, there := session.SessionsData[UUID]; !there {
		return
	}

	f, e := os.ReadFile("./testsvr/testdata/" + testUUID + "/testdat/" + UUID + ".json")

	if e != nil {
		return
	}

	var info testsvrInfo

	e = json.Unmarshal(f, &info)

	if e != nil {
		return
	}

	info.Done = true
	info.EndTime = endTime
	info.AnswerSheet = session.SessionsData[UUID].AnswerSheet // just a shallow copy, might cause a bug!

	value, e := json.Marshal(info)

	if e != nil {
		return
	}

	// save
	e = os.WriteFile("./testsvr/testdata/"+testUUID+"/testdat/"+UUID+".json", value, 0664)

	if e != nil {
		return
	}
	// destroy
	delete(session.SessionsData, UUID)
}
