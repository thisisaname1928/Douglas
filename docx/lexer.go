package docx

import "fmt"

type Token struct {
	Value FluidString
	Type  int
}

const (
	TOKEN_NONE         = 0
	TOKEN_QUES         = 1
	TOKEN_ANSWER_A     = 2
	TOKEN_ANSWER_B     = 3
	TOKEN_ANSWER_C     = 4
	TOKEN_ANSWER_D     = 5
	TOKEN_NEW_LINE     = 6
	TOKEN_TEXT_CONTENT = 7
)

const (
	TOKEN_QUES_VALUE = "Câu"
)

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

	currentTokenBeginIndex := 0
	currentTokenType := TOKEN_TEXT_CONTENT
	for i := range src {
		aRune := []rune(src[i].Text)

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

				k += len(TOKEN_QUES_VALUE)
				continue
			} else { // TOKEN_TEXT_CONTENT
				if currentTokenType == TOKEN_TEXT_CONTENT {
					k++
				} else { // finish last token
					var curTok = Token{CopyFluid(src[i], currentTokenBeginIndex, k-1), currentTokenType}
					tokens = append(tokens, curTok)

					currentTokenBeginIndex = k
					currentTokenType = TOKEN_TEXT_CONTENT
					k++
				}
			}

		}

		// finish last token again
		fmt.Println("FC", currentTokenBeginIndex)
		var curTok = Token{CopyFluid(src[i], currentTokenBeginIndex, len(aRune)-1), currentTokenType}
		tokens = append(tokens, curTok)
	}

	return tokens
}
