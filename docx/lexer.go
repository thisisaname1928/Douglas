package docx

type Token struct {
	Value FluidString
	Type  int
}

const (
	TOKEN_NONE            = 0
	TOKEN_QUES            = 1
	TOKEN_ANSWER_A        = 2
	TOKEN_ANSWER_B        = 3
	TOKEN_ANSWER_C        = 4
	TOKEN_ANSWER_D        = 5
	TOKEN_NEW_LINE        = 6
	TOKEN_TEXT_CONTENT    = 7
	TOKEN_EOF             = 8
	TOKEN_TN_ANSWER_KEY   = 9
	TOKEN_TNDS_ANSWER_KEY = 10
	TOKEN_TLN_ANSWER_KEY  = 11
)

const (
	TOKEN_QUES_VALUE     = "Câu"
	TOKEN_ANSWER_A_VALUE = "A."
	TOKEN_ANSWER_B_VALUE = "B."
	TOKEN_ANSWER_C_VALUE = "C."
	TOKEN_ANSWER_D_VALUE = "D."
)

func isTNAnswerKey(src []rune, index int) bool {
	if index+2 >= len(src) {
		return false
	}

	if src[index+1] != '.' {
		return false
	}

	if src[index] == 'A' || src[index] == 'B' || src[index] == 'C' || src[index] == 'D' {
		return true
	}

	return false
}

func HasPrefix(src []rune, index int, pref string) bool {

	aRune := []rune(pref)
	i := 0

	for ; i < len(aRune); i++ {
		if index+i >= len(src) {
			break
		}

		if src[index+i] != aRune[i] {
			break
		}
	}

	return i == len(aRune)
}

func Lex(src []FluidString) []Token {
	var tokens []Token

	for i := range src {
		aRune := []rune(src[i].Text)
		currentTokenBeginIndex := 0
		currentTokenType := TOKEN_NONE

		for k := 0; k < len(aRune); {
			// identify tokens
			if HasPrefix(aRune, k, TOKEN_QUES_VALUE) {
				if k != 0 {
					// finish last token
					var curTok = Token{CopyFluid(src[i], currentTokenBeginIndex, k-1), currentTokenType}
					tokens = append(tokens, curTok)
				}
				currentTokenBeginIndex = k
				currentTokenType = TOKEN_QUES

				k += len(TOKEN_QUES_VALUE) - 1
				continue
			} else if isTNAnswerKey(aRune, k) {
				// finish last token
				if currentTokenType != TOKEN_NONE {
					var curTok = Token{CopyFluid(src[i], currentTokenBeginIndex, k-1), currentTokenType}
					tokens = append(tokens, curTok)
				}

				currentTokenBeginIndex = k
				currentTokenType = TOKEN_TN_ANSWER_KEY

				k += 2
				continue
			} else { // TOKEN_TEXT_CONTENT
				// finish last token
				if currentTokenType != TOKEN_NONE && currentTokenType != TOKEN_TEXT_CONTENT {
					var curTok = Token{CopyFluid(src[i], currentTokenBeginIndex, k-1), currentTokenType}
					tokens = append(tokens, curTok)
				}

				if currentTokenType == TOKEN_TEXT_CONTENT {
					k++
				} else { // create new text content token
					currentTokenBeginIndex = k
					currentTokenType = TOKEN_TEXT_CONTENT
					k++
				}
			}

		}

		// finish last token again
		var curTok = Token{CopyFluid(src[i], currentTokenBeginIndex, len(aRune)-1), currentTokenType}
		tokens = append(tokens, curTok)
		// append NEW_LINE TOKEN

		tokens = append(tokens, Token{FluidString{"<br>", []FluidProperty{}}, TOKEN_NEW_LINE})
	}

	tokens = append(tokens, Token{FluidString{}, TOKEN_EOF})

	return tokens
}
