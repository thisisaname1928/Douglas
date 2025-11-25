package docx

func FindStype(fluid FluidString) string {
	aRune := []rune(fluid.Text)
	i := 0
	stype := ""

	for ; i < len(aRune); i++ {
		if aRune[i] == ':' {
			break
		}
		if aRune[i] == '[' {
			i++
			break
		}
	}

	if i >= len(aRune) || aRune[i] == ':' {
		return "NONE"
	}

	for ; i < len(aRune); i++ {
		if aRune[i] == ':' {
			return "NONE"
		}
		if aRune[i] == ']' {
			break
		}
		stype += string(aRune[i])
	}

	return stype
}

func BetterParse(tokens []Token) []Question {
	var res []Question

	for i := 0; i < len(tokens); {
		if tokens[i].Type == TOKEN_EOF {
			break
		}

		if tokens[i].Type != TOKEN_QUES {
			i++
			continue
		}

		var CurrentQuestion Question
		var currentQuesContent FluidString

		i++
		CurrentQuestion.Stype = FindStype(tokens[i].Value) // find next stype

		// get whole question content
	getQuestionContentLoop:
		for ; i < len(tokens); i++ {
			switch tokens[i].Type {
			case TOKEN_TN_ANSWER_KEY:
				CurrentQuestion.Type = TN
				break getQuestionContentLoop
			case TOKEN_TNDS_ANSWER_KEY:
				CurrentQuestion.Type = TNDS
				break getQuestionContentLoop
			case TOKEN_TLN_ANSWER_KEY:
				CurrentQuestion.Type = TLN
				break getQuestionContentLoop
			}

			currentQuesContent = ConcatFluid(currentQuesContent, tokens[i].Value)
		}

		CurrentQuestion.Content = ParseFluid2HtmlNonMark(currentQuesContent)

		res = append(res, CurrentQuestion)
		break
	}

	return res
}
