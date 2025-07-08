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
	f, e := os.Open("./tmp/word/document.xml")
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
				stack.Push("div><br")
				htmlOutput += "<div>"

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
						switch paragraphToken.Name.Local {
						case "pPr": // parse paragraph properties
							ParseProperties(&htmlOutput, stack, decoder)
						case "rPr": // parse run properties and text
							localstack := arraystack.New()
							ParseProperties(&htmlOutput, localstack, decoder)

							parseTextNeedToStop := false
							for {
								if parseTextNeedToStop {
									break
								}
								textTagToken, e := decoder.Token()
								if e == io.EOF || textTagToken == nil {
									needToStop = true
									break
								}
								switch textTag := textTagToken.(type) {
								case xml.StartElement:
									if textTag.Name.Local == "t" {
										textToken, e := decoder.Token()
										if e == io.EOF || textToken == nil {
											needToStop = true
											break
										}
										// check if contain text then get it
										switch text := textToken.(type) {
										case xml.CharData:
											htmlOutput += string(text)
											decoder.Token()
											parseTextNeedToStop = true
										case xml.EndElement:
											parseTextNeedToStop = true
										}
									}
								}
							}

							// and end tag
							for localstack.Size() > 0 {
								v, _ := localstack.Pop()
								htmlOutput += "</" + v.(string) + ">"
							}

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
