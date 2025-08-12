package testsvr

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/thisisaname1928/goParsingDocx/docx"
)

func (fir DouglasFir) getTestDataPath() string {
	return "./testsvr/testdata/" + fir.UUID + "/testdat/"
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

	for _, v := range info.Questions {
		switch v.Type {
		case docx.TN:
			// TODO
		case docx.TNDS:
			//TODO
		case docx.TLN:
		}
	}

	return output, nil
}
