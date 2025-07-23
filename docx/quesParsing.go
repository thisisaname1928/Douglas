package docx

import (
	"strings"
)

const (
	TN  = 0x12
	TLN = 0x13
)

type Question struct {
	Type       int
	Stype      string
	Content    string
	Answer     [4]string
	TrueAnswer [4]bool
	TLNA       [4]string // true answer for TLN question type
}

const (
	NOT_A_QUESTION = 1
	OUT_OF_RANGE   = 2
)

// GOLANG CAN'T PERFORM A DEEPCOPY BY ITS SELF:(
func copyFluidString(input []FluidString, index uint64) FluidString {
	var cpyInput FluidString
	cpyInput.Text = input[index].Text
	cpyInput.Properties = make([]FluidProperty, len(input[index].Properties))
	for i := range cpyInput.Properties {
		cpyInput.Properties[i].Start = input[index].Properties[i].Start
		cpyInput.Properties[i].End = input[index].Properties[i].End
		cpyInput.Properties[i].Property = make([]Prop, len(input[index].Properties[i].Property))
		copy(cpyInput.Properties[i].Property, input[index].Properties[i].Property)
	}

	return cpyInput
}

func ParseFluid2Question(index *uint64, input []FluidString) (int, Question) {
	var res Question

	currentStr := copyFluidString(input, *index)

	analyseStrRune := []rune(strings.ToLower(currentStr.Text))
	// pass space
	for analyseStrRune[0] == ' ' {
		DelFirstCharacterRune(&analyseStrRune)
		DelFirstCharacter(&currentStr)
	}

	if !strings.HasPrefix(string(analyseStrRune), "câu") {
		return NOT_A_QUESTION, res
	}

	// pass "câu"
	DelNCharacter(&currentStr, 3)
	DelNCharacterRune(&analyseStrRune, 3)

	// next find for : or [
	for analyseStrRune[0] != ':' && analyseStrRune[0] != '[' {
		DelFirstCharacterRune(&analyseStrRune)
		DelFirstCharacter(&currentStr)
	}

	res.Stype = "NONE"

	// parse stype
	if analyseStrRune[0] == '[' {
		res.Stype = ""
		DelFirstCharacterRune(&analyseStrRune)
		DelFirstCharacter(&currentStr)

		for analyseStrRune[0] != ']' {
			res.Stype += string(currentStr.Text[0])

			DelFirstCharacterRune(&analyseStrRune)
			DelFirstCharacter(&currentStr)
		}

		DelFirstCharacterRune(&analyseStrRune)
		DelFirstCharacter(&currentStr)
	}

	// find for :
	for analyseStrRune[0] != ':' {
		DelFirstCharacterRune(&analyseStrRune)
		DelFirstCharacter(&currentStr)
	}
	DelFirstCharacterRune(&analyseStrRune)
	DelFirstCharacter(&currentStr)

	// parse question content
	res.Content = ParseFluid2Html(currentStr) + "<br>"

	*index++
	currentStr = copyFluidString(input, *index)
	analyseStrRune = []rune(strings.ToLower(currentStr.Text))

	for !strings.HasPrefix(currentStr.Text, "A.") && !strings.HasPrefix(string(analyseStrRune), "đáp án:") {
		res.Content += ParseFluid2Html(currentStr) + "<br>"

		if *index >= uint64(len(input))-1 {
			return OUT_OF_RANGE, res
		}
		*index++
		currentStr = copyFluidString(input, *index)
		analyseStrRune = []rune(strings.ToLower(currentStr.Text))
	}

	// parse TN question

	return 0, res
}
