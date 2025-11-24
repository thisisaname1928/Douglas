package docx

import (
	"fmt"
	"testing"
)

func TestDocx(t *testing.T) {

	fluid, e := Parse2Fluid("/home/ngqt/projs/goParsingDocx/test.docx")

	if e != nil {
		panic(e)
	}

	//sf := CopyFluid(fluid[0], 0, 5)

	for _, v := range Lex(fluid) {
		fmt.Print(v.Value.Text, " ")
	}
}
