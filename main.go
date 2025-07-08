package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/emirpasic/gods/stacks/arraystack"
)

type RunProperty struct {
	XMLName xml.Name `xml:"rPr"`
}

type ParagraphProperty struct {
	XMLName xml.Name `xml:"pPr"`
}

type Text struct {
	XMLName xml.Name `xml:"t"`
	Text    string   `xml:",chardata"`
}

type Run struct {
	XMLName xml.Name `xml:"r"`
	Texts   []Text   `xml:"t"`
}

type Paragraph struct {
	XMLName xml.Name `xml:"p"`
	Runs    []Run    `xml:"r"`
}

type DocumentBody struct {
	XMLName    xml.Name    `xml:"body"`
	Paragraphs []Paragraph `xml:"p"`
}

type Document struct {
	XMLName xml.Name     `xml:"document"`
	Body    DocumentBody `xml:"body"`
}

type Properties struct {
	IsBold        bool
	IsItalic      bool
	IsUnderline   bool
	IsStrike      bool
	IsSubscript   bool
	IsSuperscript bool
	IsMarked      bool
}

type DocxRun struct {
	RPr  Properties
	Text string
}

type DocxParagraph struct {
	PPr  Properties
	Runs []DocxRun
}

func ParseProperties(htmlOutput *string, st *arraystack.Stack, decoder *xml.Decoder) {
	needToStop := false
	for {
		if needToStop {
			break
		}
		tok, e := decoder.Token()
		if tok == nil || e == io.EOF {
			break
		}

		switch token := tok.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "b":
				st.Push("b")
				*htmlOutput += "<b>"
			case "i":
				st.Push("i")
				*htmlOutput += "<i>"
			case "u":
				st.Push("u")
				*htmlOutput += "<u>"
			case "shd":
				st.Push("mark")
				*htmlOutput += "<mark>"
			}
		case xml.EndElement: // if reach the properties end tag
			if (token.Name.Local == "pPr") || (token.Name.Local == "rPr") {
				needToStop = true
				break
			}
		}
	}
}

func main() {
	f, e := os.Open("./test.xml")
	if e != nil {
		panic(e)
	}

	s, e := io.ReadAll(f)

	if e != nil {
		panic(e)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(s)))

	var htmlOutput string

	// read until meet <body>
	shouldEnd := false
	for {
		if shouldEnd {
			fmt.Println("I meet body!!")
			break
		}

		tok, e := decoder.Token()
		if tok == nil || e == io.EOF {
			fmt.Println("END1")
			break
		}
		// meet body

		for {
			if tok == nil || e == io.EOF {
				fmt.Println("END2")
				break
			}

			// check is a body tag
			switch token := tok.(type) {
			case xml.StartElement:
				if token.Name.Local == "body" {
					shouldEnd = true
					break
				}
			}

			if shouldEnd {
				break
			}
			// fetch next token
			tok, e = decoder.Token()
		}

	}

	htmlOutput += "<body>"

	for {
		stack := arraystack.New()
		tok, e := decoder.Token()
		if tok == nil || e == io.EOF {
			break
		}

		switch token := tok.(type) {
		case xml.StartElement:
			if token.Name.Local == "p" { // parse paragraph
				stack.Push("div")
				htmlOutput += "<div><br>"

				// start parse tag inside paragraph
				needToStop := false
				for {
					if needToStop {
						break
					}

					paragraphTok, e := decoder.Token()
					if paragraphTok == nil || e == io.EOF {
						break
					}
					switch paragraphToken := paragraphTok.(type) {
					case xml.StartElement:
						if paragraphToken.Name.Local == "rPr" { // parse paragraph properties
							ParseProperties(&htmlOutput, stack, decoder)
						} // then parse runs and their text

						nextTok, e := decoder.Token()
						if nextTok == nil || e == io.EOF {
							break
						}

					case xml.EndElement:
						if paragraphToken.Name.Local == "p" {
							needToStop = true
							break
						}
					}
				}

				for stack.Size() > 0 {
					v, _ := stack.Pop()
					htmlOutput += "</" + v.(string) + ">"
				}
			} else if token.Name.Local == "t" {
				nexttok, _ := decoder.Token()
				switch t := nexttok.(type) {
				case xml.CharData:
					htmlOutput += string(t)
				}

				// pass </t>
				decoder.Token()
			}
		case xml.EndElement:
			if token.Name.Local == "body" {
				break
			}
		}
	}

	// done
	htmlOutput += "</body>"
	fmt.Println(htmlOutput)
}
