package dou

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/thisisaname1928/goParsingDocx/docx"
	"github.com/thisisaname1928/goParsingDocx/security"
)

const (
	DOU_REVISION_1 = 0xdebe
)

type DouQuestion struct {
	Stype    string          `json:"stype"`
	Question []docx.Question `json:"questions"`
}

type TestStructure struct {
	Stype  string  `json:"stype"`
	N      uint64  `json:"number"`
	Points float64 `json:"point"` // point per ques
}

type DouData struct {
	Revision         uint64          `json:"revision"`     // version data
	TestDuration     uint64          `json:"testDuration"` // in second, if this field is zero then no time limit
	Questions        []DouQuestion   `json:"questions"`    // store info about dou file version
	UseTestStructure bool            `json:"useTestStructure"`
	TestStruct       []TestStructure `json:"testStructure"` // store info about how 2 display the test, and point per type
}

type DouInfo struct {
	Revision  uint64 `json:"revision"`
	Author    string `json:"author"`
	Encrypted bool   `json:"encrypted"`
	Key       string `json:"key"` // store as sha256
}

func ConvertPath(path *string) {
	r := []rune(*path)
	for i := range r {
		if r[i] == '\\' {
			r[i] = '/'
		}
	}

	*path = string(r)
}

func search4Stype(qs []DouQuestion, stype string) int {
	for i := range qs {
		if qs[i].Stype == stype {
			return i
		}
	}

	return -1
}

// func recreateFolder(path string) error {
// 	e := os.RemoveAll(path)
// 	if e != nil && e != os.ErrNotExist {
// 		return e
// 	}

// 	e = os.Mkdir(path, 0755) // 0755 is permission on unix-like os, idk if windows have it?
// 	if e != nil {
// 		return e
// 	}
// 	return nil
// }

func Export(input string, output string, author string, testDuration uint64, useTestStructure bool, testStructure []TestStructure, useEncryption bool, key string) error {
	// make sure...
	ConvertPath(&input)
	ConvertPath(&output)

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

	dou.UseTestStructure = useTestStructure
	if useTestStructure {
		dou.TestStruct = testStructure
	}

	jsonQuestionData, e := json.Marshal(&dou)
	if e != nil {
		return e
	}

	// pack into a zip
	archive, e := os.Create(output)
	if e != nil {
		return e
	}
	defer archive.Close()
	writer := zip.NewWriter(archive)

	// copy jsonQuestionData
	w, e := writer.Create("data.json")
	if e != nil {
		return e
	}

	if useEncryption {
		jsonQuestionData, e = security.Encrypt(jsonQuestionData, key)
		if e != nil {
			return e
		}
	}

	io.Copy(w, bytes.NewReader(jsonQuestionData))

	// add word media

	f, e := os.ReadFile(input)

	if e != nil {
		return e
	}

	a, e := zip.NewReader(bytes.NewReader(f), int64(len(f)))

	if e != nil {
		return e
	}

	for _, f := range a.File {
		if docx.CheckIsMediaFile(f.Name) {
			dat, e := f.Open()

			if e != nil {
				return e
			}
			defer dat.Close()

			tmp := strings.Split(f.Name, "/")
			fname := tmp[len(tmp)-1]

			w, e := writer.Create("/media/" + fname)
			if e != nil {
				return e
			}

			byteDat, e := io.ReadAll(dat)
			if e != nil {
				return e
			}
			if useEncryption {
				byteDat, e = security.Encrypt(byteDat, key)
				if e != nil {
					return e
				}
			}

			io.Copy(w, bytes.NewReader(byteDat))

		}
	}

	// write dou info
	var info DouInfo
	info.Revision = DOU_REVISION_1
	info.Author = author
	info.Encrypted = useEncryption

	if useEncryption {
		info.Key = security.EncryptKey(key)
	}

	jsonDouInfo, e := json.Marshal(&info)
	if e != nil {
		return e
	}
	w, e = writer.Create("info.json")
	if e != nil {
		return e
	}
	io.Copy(w, bytes.NewReader(jsonDouInfo))

	// end
	writer.Close()

	return nil
}
