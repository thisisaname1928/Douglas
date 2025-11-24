package docx

import (
	"slices"
)

func searchIsXmlTag(s []rune, index int) (bool, []rune) {
	var res []rune
	if s[index] == '<' {

		for i := index + 1; i < len(s); i++ {
			if s[i] == ' ' {
				break
			}
			if s[i] == '>' {
				break
			}
			if s[i] == '/' && s[i+1] == '>' {
				break
			}
			res = append(res, s[i])
		}

		return true, res
	} else {
		return false, res
	}
}

var keepableTag = []string{"?xml", "pic:nvPicPr", "pic:blipFill", "a:blip", "w:i", "w:iCs", "w:shd",
	"w:document", "w:body", "w:p", "w:r", "w:pPr", "a:graphicData", "pic:pic", "w:inline",
	"w:rPr", "w:t", "w:drawing", "wp:anchor", "a:graphic", "w:highlight", "w:b"}

func isKeepable(tag string) bool {
	for i := range keepableTag {
		if tag == keepableTag[i] {
			return true
		}
	}

	return false
}

func isEndTag(tag []rune) bool {
	return tag[0] == '/'
}

func StripTag(s string) string {
	var currentTag []string
	res := ""
	anaRune := []rune(s)

	for i := 0; i < len(anaRune); {
		isTag, tag := searchIsXmlTag(anaRune, i)

		if !isTag {
			// end single tag
			if anaRune[i] == '/' && anaRune[i+1] == '>' {
				if isKeepable(currentTag[len(currentTag)-1]) {
					res += "/>"
				}
				currentTag = slices.Delete(currentTag, len(currentTag)-1, len(currentTag))
				i += 2
				continue
			}

			if isKeepable(currentTag[len(currentTag)-1]) { // keep
				res += string(anaRune[i])
				i++
				continue
			} else {
				i++
				continue
			}
		} else { // found a tag
			if isEndTag(tag) {
				if isKeepable(currentTag[len(currentTag)-1]) {
					res += "<" + string(tag) + ">"
				}

				currentTag = slices.Delete(currentTag, len(currentTag)-1, len(currentTag))
				i += len(tag) + 2
				continue
			}

			// found new tag
			strTag := string(tag)
			if isKeepable(string(strTag)) {
				res += "<" + strTag
			}

			currentTag = append(currentTag, strTag)

			i += len(tag) + 1
			continue
		}
	}

	return res
}
