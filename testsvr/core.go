package testsvr

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/thisisaname1928/goParsingDocx/docx"
	"github.com/thisisaname1928/goParsingDocx/dou"
)

func shuffleQuesAr(ques *[]docx.Question) {
	for i := range *ques {
		// shuffle question answer
		// pass TLN type
		if (*ques)[i].Type == docx.TLN {
			continue
		}

		rand.Shuffle(4, func(r int, l int) {
			(*ques)[i].Answer[r], (*ques)[i].Answer[l] = (*ques)[i].Answer[l], (*ques)[i].Answer[r]
			(*ques)[i].TrueAnswer[r], (*ques)[i].TrueAnswer[l] = (*ques)[i].TrueAnswer[l], (*ques)[i].TrueAnswer[r]
		})
	}
}

type DouglasTest struct {
	Test []docx.Question `json:"test"`
}

func (fir *DouglasFir) GetNumberOfQuestions() int {
	c := 0

	for _, v := range fir.Douglas.Data.TestStruct {
		c += int(v.N)
	}

	return c
}

func (fir *DouglasFir) ShuffleNewTest() DouglasTest {
	var test []docx.Question
	// load teststruct into a map
	testStructure := make(map[string]dou.TestStructure)
	for _, v := range fir.Douglas.Data.TestStruct {
		testStructure[v.Stype] = v
	}

	// make a copy of current test
	var curQuestionList []dou.DouQuestion = make([]dou.DouQuestion, len(fir.Douglas.Data.Questions))
	copy(curQuestionList, fir.Douglas.Data.Questions)

	for i := range curQuestionList {
		shuffleQuesAr(&curQuestionList[i].Question) // shuffle sub question

		// get first N questions
		for i, v := range curQuestionList[i].Question {
			if i >= int(testStructure[v.Stype].N) {
				break
			}

			v.Point = testStructure[v.Stype].Points
			test = append(test, v)
		}

	}

	// shuffle questions list
	rand.Shuffle(len(test), func(a, b int) {
		test[a], test[b] = test[b], test[a]
	})

	var testdg DouglasTest
	testdg.Test = make([]docx.Question, len(test))
	copy(testdg.Test, test)
	return testdg
}

// get file that store the data of test sessions uuid
func (fir *DouglasFir) Route2UUID(UUID string) string {
	return "./testsvr/testdata/" + fir.UUID + "/testdat/" + UUID + ".json"
}

func (fir *DouglasFir) CheckIfTestDone(UUID string) bool {
	// check if test haven't available in test sessions
	buf, e := os.ReadFile(fir.Route2UUID(UUID))

	if e != nil {
		return false
	}

	var testInfo testsvrInfo
	e = json.Unmarshal(buf, &testInfo)
	if e != nil {
		return false
	}

	return testInfo.Done
}

func calcTNQuestion(point float64, trueAns [4]bool, userAns string) float64 {
	index := -1
	switch userAns {
	case "A":
		index = 0
	case "B":
		index = 1
	case "C":
		index = 2
	case "D":
		index = 3
	}

	if index == -1 {
		return 0
	}

	if trueAns[index] {
		return point
	}

	return 0
}

func calcTLNQuestion(point float64, trueAns [4]string, userAns [4]string) float64 {
	if trueAns == userAns {
		return point
	}

	return 0
}

func calcTNDSQuestion(point float64, trueAns [4]bool, userAns [4]string) float64 {
	// convert
	var uAns [4]bool
	for i := range userAns {
		uAns[i] = (userAns[i] == "T")
	}

	if trueAns == uAns {
		return point
	}

	return 0
}

func (fir *DouglasFir) CalculateMark(UUID string) (int, float64, error) {
	var mark float64 = 0
	var trueQuesCount = 0
	// load
	buf, e := os.ReadFile(fir.Route2UUID(UUID))

	if e != nil {
		return 0, 0, errors.New("TEST_NOT_FOUND")
	}

	var testInfo testsvrInfo
	e = json.Unmarshal(buf, &testInfo)
	if e != nil {
		return 0, 0, errors.New("BAD_TEST")
	}

	// check if test done
	if !testInfo.Done {
		return 0, 0, errors.New("TEST_IS_NOT_DONE")
	}

	for i := range testInfo.AnswerSheet {
		switch testInfo.Questions[i].Type {
		case docx.TN:
			point := calcTNQuestion(testInfo.Questions[i].Point, testInfo.Questions[i].TrueAnswer, testInfo.AnswerSheet[i][0])

			mark += point
			if point == testInfo.Questions[i].Point {
				trueQuesCount++
			}
		case docx.TLN:
			point := calcTLNQuestion(testInfo.Questions[i].Point, testInfo.Questions[i].TLNA, [4]string(testInfo.AnswerSheet[i]))

			mark += point
			if point == testInfo.Questions[i].Point {
				trueQuesCount++
			}
		case docx.TNDS:
			point := calcTNDSQuestion(testInfo.Questions[i].Point, testInfo.Questions[i].TrueAnswer, [4]string(testInfo.AnswerSheet[i]))

			mark += point
			if point == testInfo.Questions[i].Point {
				trueQuesCount++
			}
		}
	}

	return trueQuesCount, mark, nil
}

func currentServerTime() time.Time {
	return time.Now()
}

type getCurrentServerTimeRequest struct {
	UUID string `json:"uuid"`
}

func (fir *DouglasFir) getCurrentServerTime(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request getCurrentServerTimeRequest
	decoder.Decode(&request)

	// done test
	if !fir.checkTestTime(request.UUID) {
		st, _ := fir.TestSessions.GetSessionStartTime(request.UUID)
		requestIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		if fir.verifyIP(request.UUID, requestIP) {
			fir.TestSessions.DoneSession(fir.UUID, request.UUID, st.Add(time.Minute*time.Duration(fir.Douglas.Data.TestDuration)))
		}

	}

	w.Write([]byte(currentServerTime().Format(time.RFC3339)))
}

// return false if test end of time
func (fir *DouglasFir) checkTestTime(UUID string) bool {
	st, e := fir.TestSessions.GetSessionStartTime(UUID)

	if e != nil {
		return false
	}

	duration := currentServerTime().Sub(st)

	return duration.Seconds() <= float64(fir.Douglas.Data.TestDuration*60)
}
