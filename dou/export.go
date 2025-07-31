package dou

import (
	"encoding/json"
	"fmt"

	"github.com/thisisaname1928/goParsingDocx/app"
	"github.com/thisisaname1928/goParsingDocx/docx"
)

const (
	DOU_REVISION_1 = 0xdebe
)

type DouQuestion struct {
	Stype    string          `json:"stype"`
	Question []docx.Question `json:"questions"`
}

type DouData struct {
	Revision     uint64        `json:"revision"`     // version data
	TestDuration uint64        `json:"testDuration"` // in second, if this field is zero then no time limit
	Questions    []DouQuestion `json:"questions"`    // store info about dou file version
}

func search4Stype(qs []DouQuestion, stype string) int {
	for i := range qs {
		if qs[i].Stype == stype {
			return i
		}
	}

	return -1
}

func Export(input string, output string, testDuration uint64, useEncryption bool, key string) error {
	// make sure...
	app.ConvertPath(&input)
	app.ConvertPath(&output)

	// EXPORT QUESTIONS DATA HERE
	fluid, e := docx.Parse2Fluid(input)
	if e != nil {
		return e
	}

	var ques []docx.Question
	var index uint64 = 0
	i, q := docx.ParseFluid2Question(&index, fluid)
	for i == 0 {
		ques = append(ques, q)
		i, q = docx.ParseFluid2Question(&index, fluid)
	}

	// parse questions into better type
	var douQues []DouQuestion
	for _, v := range ques {
		index := search4Stype(douQues, v.Stype)
		if index == -1 { // there is no stype like that
			var currentDouQues DouQuestion
			currentDouQues.Stype = v.Stype
			currentDouQues.Question = append(currentDouQues.Question, v)
			douQues = append(douQues, currentDouQues)
		} else {
			douQues[index].Question = append(douQues[index].Question, v)
		}
	}

	// convert to json
	var dou DouData
	dou.Questions = douQues
	dou.Revision = DOU_REVISION_1
	dou.TestDuration = testDuration
	jsonQuestionData, e := json.Marshal(&dou)

	fmt.Println(string(jsonQuestionData), e)

	return nil
}
