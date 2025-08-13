package testsvr

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/thisisaname1928/goParsingDocx/docx"
)

func (fir DouglasFir) getTestDataPath() string {
	return "./testsvr/testdata/" + fir.UUID + "/testdat/"
}

func openTemplate(path string) (string, error) {
	b, e := os.ReadFile("./testsvr/" + path)

	if e != nil {
		return "", e
	}

	return string(b), nil
}

type templateElement struct {
	Name  string
	Value string
}

func loadTemplate(es []templateElement, path string) (string, error) {
	content, e := openTemplate(path)
	if e != nil {
		return "", e
	}

	// replace
	for _, v := range es {
		content = strings.ReplaceAll(content, fmt.Sprintf("{{.%v}}", v.Name), v.Value)
	}

	return content, nil
}

func (fir DouglasFir) RenderTest(uuid string) (string, error) {
	output := ""

	b, e := os.ReadFile(fir.getTestDataPath() + uuid + ".json")

	if e != nil {
		return "", fmt.Errorf("internal error: %v", e)
	}

	var info testsvrInfo
	e = json.Unmarshal(b, &info)
	if e != nil {
		return "", fmt.Errorf("internal error: %v", e)
	}

	for i, v := range info.Questions {
		switch v.Type {
		case docx.TN:
			s, e := loadTemplate([]templateElement{{"quesIndex", fmt.Sprint(i)}}, "TNTemplate.html")
			if e == nil {
				output += s
			}
		case docx.TNDS:
			s, e := loadTemplate([]templateElement{{"quesIndex", fmt.Sprint(i)}}, "TNDSTemplate.html")
			if e == nil {
				output += s
			}
		case docx.TLN:
			s, e := loadTemplate([]templateElement{{"quesIndex", fmt.Sprint(i)}}, "TLNTemplate.html")
			if e == nil {
				output += s
			}
		}
	}

	return output, nil
}
