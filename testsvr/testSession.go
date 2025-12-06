package testsvr

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// this for temporary saving sessions data
type TestSessionData struct {
	AnswerSheet [][]string `json:"AnswerSheet"` // Ans[0]=["1", "1", "1", "1"] example
	IP          string     `json:"IP"`
	StartTime   time.Time  `json:"StartTime"`
	IsLocked    bool
	UUID        string
}

// access Sessions data by uuid
// Session Data should be saved and destroyed after the test is done
// Resume session by lookup this TestSessions
type TestSessions struct {
	SessionsData map[string]TestSessionData
	mutex        sync.Mutex
}

func (session *TestSessions) Init() {
	session.SessionsData = make(map[string]TestSessionData)
}

func (session *TestSessions) NewSession(UUID string, IP string, startTime time.Time, numberOfQuestions int) {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	session.SessionsData[UUID] = TestSessionData{make([][]string, numberOfQuestions), IP, startTime, false, UUID}

	// assign 4 element to Answer sheet
	for i := 0; i < numberOfQuestions; i++ {
		session.SessionsData[UUID].AnswerSheet[i] = make([]string, 4)
	}

	fmt.Println("create Session for ", UUID)
}

func (session *TestSessions) UpdateAnswerSheet(i int, UUID string, Answer [4]string) {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	if session.SessionsData[UUID].IsLocked {
		fmt.Println("MEET A LOCK")
		return
	}

	if i < len(session.SessionsData[UUID].AnswerSheet) {
		for k := 0; k < 4; k++ {
			session.SessionsData[UUID].AnswerSheet[i][k] = Answer[k] // do a copy
		}
	}
}

// update single answer
func (session *TestSessions) UpdateAnswer(UUID string, index int, answerIndex int, data string) {
	if !session.CheckSession(UUID) {
		return
	}

	session.mutex.Lock()
	defer session.mutex.Unlock()

	if session.SessionsData[UUID].IsLocked {
		return
	}

	session.SessionsData[UUID].AnswerSheet[index][answerIndex] = data
}

func (session *TestSessions) DoneSession(testUUID string, UUID string, endTime time.Time) {
	session.mutex.Lock()
	defer session.mutex.Unlock()
	// check if session available
	if _, there := session.SessionsData[UUID]; !there {
		return
	}

	if session.SessionsData[UUID].IsLocked {
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

func (session *TestSessions) CheckSession(UUID string) bool {
	session.mutex.Lock()
	defer session.mutex.Unlock()
	_, ok := session.SessionsData[UUID]
	return ok
}

// convert test sessions answer sheet into json
func (session *TestSessions) CopyAnsSheet(UUID string) ([][]string, error) {
	if !session.CheckSession(UUID) {
		return [][]string{}, errors.New("ACCESS_DENIED")
	}

	session.mutex.Lock()
	arr := make([][]string, len(session.SessionsData[UUID].AnswerSheet))

	// do a deep & manual copy
	for i := range arr {
		arr[i] = make([]string, len(session.SessionsData[UUID].AnswerSheet[i]))
		copy(arr[i], session.SessionsData[UUID].AnswerSheet[i])
	}

	session.mutex.Unlock()
	return arr, nil
}

func (session *TestSessions) LockSession(UUID string) {
	if !session.CheckSession(UUID) {
		return
	}

	session.mutex.Lock()
	defer session.mutex.Unlock()

	tmp := session.SessionsData[UUID]
	tmp.IsLocked = true
	session.SessionsData[UUID] = tmp
}

func (session *TestSessions) CheckSessionLock(UUID string) bool {
	// if not exist, better return true
	if !session.CheckSession(UUID) {
		return true
	}

	session.mutex.Lock()
	defer session.mutex.Unlock()

	return session.SessionsData[UUID].IsLocked
}

func (session *TestSessions) GetSessionStartTime(UUID string) (time.Time, error) {
	if !session.CheckSession(UUID) {
		return time.Now(), errors.New("ACCESS_DENIED")
	}

	session.mutex.Lock()
	defer session.mutex.Unlock()

	return session.SessionsData[UUID].StartTime, nil
}

func (session *TestSessions) CloseAllTestSessions(testUUID string) {
	// fetch for all UUIDs
	var UUIDsList []string

	for v := range session.SessionsData {
		UUIDsList = append(UUIDsList, v)
	}

	for i := 0; i < len(UUIDsList); i++ {
		session.DoneSession(testUUID, UUIDsList[i], time.Now())
	}
}
