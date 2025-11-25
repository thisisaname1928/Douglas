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
	tokens := Lex(fluid)

	for _, v := range tokens {
		switch v.Type {
		case TOKEN_QUES:
			fmt.Print("TOKEN_QUES ")
		case TOKEN_TEXT_CONTENT:
			fmt.Print("TOKEN_TEXT_CONTENT ")
		case TOKEN_ANSWER_A:
			fmt.Print("TOKEN_ANSWER_A ")
		case TOKEN_ANSWER_B:
			fmt.Print("TOKEN_ANSWER_B ")
		case TOKEN_ANSWER_C:
			fmt.Print("TOKEN_ANSWER_C ")
		case TOKEN_ANSWER_D:
			fmt.Print("TOKEN_ANSWER_D ")
		case TOKEN_EOF:
			fmt.Print("TOKEN_EOF ")
		case TOKEN_NEW_LINE:
			fmt.Println("TOKEN_NEW_LINE ")
		case TOKEN_TN_ANSWER_KEY:
			fmt.Print("TOKEN_TN_ANSWER_KEY ")
		}
	}

	fmt.Println("\n", BetterParse(tokens)[0].Content)
}
