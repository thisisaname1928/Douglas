package testsvr

import (
	"math/rand"

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

func (fir DouglasFir) ShuffleNewTest() DouglasTest {
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
			if i > int(testStructure[v.Stype].N) {
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
