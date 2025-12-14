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

func checkIfTokenMarked(tok Token) bool {
	for _, v := range tok.Value.Properties {
		for _, v2 := range v.Property {
			if v2.Type == Marked && v2.Value != "auto" {
				return true
			}
		}
	}

	return false
}

func getTNAnswerIndex(ch rune) int {
	switch ch {
	case 'A':
		return 0
	case 'B':
		return 1
	case 'C':
		return 2
	case 'D':
		return 3
	}

	return 0
}

func getTNDSAnswerIndex(ch rune) int {
	switch ch {
	case 'a':
		return 0
	case 'b':
		return 1
	case 'c':
		return 2
	case 'd':
		return 3
	}

	return 0
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

		// parse first text content
		if tokens[i].Type == TOKEN_TEXT_CONTENT {
			textContent := MakeFluidStringInstance(tokens[i].Value)
			anaRune := []rune(textContent.Text)

			j := 1
			for ; j < len(anaRune); j++ {
				if anaRune[j-1] == ':' {
					break
				}
			}

			i++
			currentQuesContent = CopyFluid(textContent, j, len(anaRune)-1)
		}

		// get whole question content
	getQuestionContentLoop:
		for ; i < len(tokens); i++ {
			switch tokens[i].Type {
			case TOKEN_TN_ANSWER_KEY:
				if tokens[i].Value.Text != "A." {
					break
				}
				CurrentQuestion.Type = TN
				break getQuestionContentLoop
			case TOKEN_TNDS_ANSWER_KEY:
				if tokens[i].Value.Text != "a)" {
					break
				}
				CurrentQuestion.Type = TNDS
				break getQuestionContentLoop
			case TOKEN_TLN_ANSWER_KEY:
				CurrentQuestion.Type = TLN
				break getQuestionContentLoop
			}

			currentQuesContent = ConcatFluid(currentQuesContent, tokens[i].Value)
		}

		// parse TN answers
		if CurrentQuestion.Type == TN {
			for tokens[i].Type != TOKEN_EOF {
				// parse content

				if tokens[i].Type == TOKEN_TN_ANSWER_KEY {
					ansIndx := getTNAnswerIndex([]rune(tokens[i].Value.Text)[0])
					if checkIfTokenMarked(tokens[i]) {
						for l := range CurrentQuestion.TrueAnswer {
							if l == ansIndx {
								CurrentQuestion.TrueAnswer[l] = true
							} else {
								CurrentQuestion.TrueAnswer[l] = false
							}
						}
					}

					var ansContent FluidString

					i++

					for ; (tokens[i].Type != TOKEN_TN_ANSWER_KEY || ansIndx >= getTNAnswerIndex([]rune(tokens[i].Value.Text)[0])) && tokens[i].Type != TOKEN_EOF && tokens[i].Type != TOKEN_QUES; i++ {
						ansContent = ConcatFluid(ansContent, tokens[i].Value)
					}

					CurrentQuestion.Answer[ansIndx] = ParseFluid2HtmlNonMark(ansContent)
				} else {
					break
				}
			}
		} else if CurrentQuestion.Type == TNDS { // parse TNDS answers
			for tokens[i].Type != TOKEN_EOF {
				// parse content
				if tokens[i].Type == TOKEN_TNDS_ANSWER_KEY {
					ansIndx := getTNDSAnswerIndex([]rune(tokens[i].Value.Text)[0])
					if checkIfTokenMarked(tokens[i]) {
						for l := range CurrentQuestion.TrueAnswer {
							if l == ansIndx {
								CurrentQuestion.TrueAnswer[l] = true
							}
						}
					}

					var ansContent FluidString

					i++

					for ; tokens[i].Type != TOKEN_TNDS_ANSWER_KEY && tokens[i].Type != TOKEN_EOF && tokens[i].Type != TOKEN_QUES; i++ {
						ansContent = ConcatFluid(ansContent, tokens[i].Value)
					}

					CurrentQuestion.Answer[ansIndx] = ParseFluid2HtmlNonMark(ansContent)
				} else {
					break
				}
			}
		} else if CurrentQuestion.Type == TLN { // parse TLN answers
			for tokens[i].Type != TOKEN_EOF {
				// parse content
				if tokens[i].Type == TOKEN_TLN_ANSWER_KEY {
					// get answer content
					i++
					var ansContent FluidString = tokens[i].Value

					// parse content
					ansContentRune := []rune(ansContent.Text)

					x := 0
					for x = range ansContentRune {
						if ansContentRune[x] != ' ' { // remove space
							break
						}
					}

					for y := 0; y < len(ansContentRune) && y < 4; y++ {
						if x+y >= len(ansContentRune) {
							break
						}

						CurrentQuestion.TLNA[y] = string(ansContentRune[x+y])
					}
					i++
				} else {
					break
				}
			}
		}

		CurrentQuestion.Content = ParseFluid2HtmlNonMark(currentQuesContent)

		res = append(res, CurrentQuestion)
		if i >= len(tokens) {
			break
		}
		if tokens[i].Type != TOKEN_QUES {
			i++
		}
	}

	return res
}
