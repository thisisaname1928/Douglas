package docx

import (
	"fmt"
	"testing"
)

func TestCopyFluid(t *testing.T) {
	var f = FluidString{"abc", []FluidProperty{{0, 2, []Prop{{1, "ABC"}}}}}

	fmt.Println(CopyFluid(f, 1, 2))
}

func TestDocx(t *testing.T) {
	//fluid, _ := Parse2Fluid("C:\\Users\\nguye\\OneDrive\\Documents\\Câu_1.docx")
	fluid := String2Fluid("Câu * 1: 1        23. *a) Dap          an A=A.*A. b) Dap An B *c) Dap an C d) Dap an D")
	//return
	//sf := CopyFluid(fluid[0], 0, 5)
	tokens := Lex(fluid)
	fmt.Println(tokens)

	for _, v := range tokens {
		switch v.Type {
		case TOKEN_QUES:
			fmt.Print("TOKEN_QUES ")
		case TOKEN_TEXT_CONTENT:
			fmt.Print("TOKEN_TEXT_CONTENT:", v.Value.Text, " ")
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
		case TOKEN_TLN_ANSWER_KEY:
			fmt.Print("TOKEN_TLN_ANSWER_KEY ")
		case TOKEN_TNDS_ANSWER_KEY:
			fmt.Print("TOKEN_TNDS_ANSWER_KEY ")
		case TOKEN_ANSWER_MARK:
			fmt.Print("TOKEN_ANSWER_MARK ")
		}
	}

	ques := BetterParse(tokens)
	ex := ques[0]

	fmt.Println("\n", ex.Content)
	fmt.Println("\n", ex.TrueAnswer)
	fmt.Println("\nA.", ex.Answer[0])
	fmt.Println("\nB.", ex.Answer[1])
	fmt.Println("\nC.", ex.Answer[2])
	fmt.Println("\nD.", ex.Answer[3])
}
