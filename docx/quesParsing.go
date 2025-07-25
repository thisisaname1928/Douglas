package docx

import (
	"strings"
)

const (
	TN  = 0x12
	TLN = 0x13
)

type Question struct {
	Type       int       `json:"type"`
	Stype      string    `json:"stype"`
	Content    string    `json:"content"`
	Answer     [4]string `json:"answers"`
	TrueAnswer [4]bool   `json:"TNAnswers"`
	TLNA       [4]string `json:"TLNAnswers"` // true answer for TLN question type
}

const (
	NOT_A_QUESTION = 1
	OUT_OF_RANGE   = 2
)

// GOLANG CAN'T PERFORM A DEEPCOPY BY ITS SELF:(
func copyFluidString(input []FluidString, index uint64) FluidString {
	var cpyInput FluidString

	defer func() {
		if r := recover(); r != nil {
		}
	}()
	if index > uint64(len(input)-1) {
		return cpyInput
	}
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

	if *index >= uint64(len(input)-1) {
		return OUT_OF_RANGE, res
	}

	currentStr := copyFluidString(input, *index)

	if len(currentStr.Text) <= 0 {
		return OUT_OF_RANGE, res
	}

	analyseStrRune := []rune(strings.ToLower(currentStr.Text))

	// pass space
	for analyseStrRune[0] == ' ' {
		DelFirstCharacterRune(&analyseStrRune)
		DelFirstCharacter(&currentStr)

		if len(analyseStrRune) <= 0 {
			return OUT_OF_RANGE, res
		}
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
	if strings.HasPrefix(currentStr.Text, "A.") {
		res.Type = TN
		t, ans := parseAnswer(index, input)
		copy(res.Answer[:], t[:])
		copy(res.TrueAnswer[:], ans)
	} else if strings.HasPrefix(string(currentStr.Text), "đáp án:") || strings.HasPrefix(string(currentStr.Text), "Đáp án:") {
		res.Type = TLN
		DelNCharacter(&currentStr, 7)

		analyseStrRune = []rune(currentStr.Text)

		if len(analyseStrRune) > 0 {
			for analyseStrRune[0] == ' ' {
				DelFirstCharacter(&currentStr)
				DelFirstCharacterRune(&analyseStrRune)
				if len(analyseStrRune) <= 0 {
					break
				}
			}
		}

		for i := 0; i < 4; i++ {
			if len(analyseStrRune) > 0 {
				res.TLNA[i] = string(analyseStrRune[0])
				DelFirstCharacter(&currentStr)
				DelFirstCharacterRune(&analyseStrRune)
			}
		}
		*index++
	}

	return 0, res
}

func parseAnswer(index *uint64, input []FluidString) ([]string, []bool) {
	var output [4]string
	var ans [4]bool

	output[0], ans[0] = parseSubAnswer(index, input, "A.", "B.")
	output[1], ans[1] = parseSubAnswer(index, input, "B.", "C.")
	output[2], ans[2] = parseSubAnswer(index, input, "C.", "D.")
	output[3], ans[3] = parseSubAnswer(index, input, "D.", "Câu")

	return output[:], ans[:]
}

func parseSubAnswer(index *uint64, input []FluidString, curPrefix string, nextPrefix string) (string, bool) {
	var output string
	res := false

	currentStr := copyFluidString(input, *index)
	// check is answer
	for _, p := range currentStr.Properties {
		for _, rp := range p.Property {
			if rp.Type == Marked && rp.Value != "auto" {
				res = true
			}
		}
	}

	for !strings.HasPrefix(currentStr.Text, curPrefix) {
		DelFirstCharacter(&currentStr)

		if len(currentStr.Text) <= 0 {
			*index++
			if *index > uint64(len(input)-1) {
				return output, res
			}
			currentStr = copyFluidString(input, *index)
		}
	}

	DelNCharacter(&currentStr, 2)

	output = ParseFluid2HtmlNonMark(currentStr)

	*index++
	if *index > uint64(len(input)-1) {
		return output, res
	}
	currentStr = copyFluidString(input, *index)

	for !strings.HasPrefix(currentStr.Text, nextPrefix) {
		output += "<br>" + ParseFluid2HtmlNonMark(currentStr)
		*index++
		if *index > uint64(len(input)-1) {
			return output, res
		}
		currentStr = copyFluidString(input, *index)
	}

	return output, res
}
