package docx

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Property struct {
	XMLName xml.Name
	Val     *string `xml:"val,attr"`
}

type Text struct {
	XMLName xml.Name `xml:"t"`
	Text    *string  `xml:",chardata"`
}

type RunProperties struct {
	XMLName    xml.Name    `xml:"rPr"`
	Properties *[]Property `xml:",any"`
}

type Run struct {
	XMLName xml.Name `xml:"r"`
	Text    Text     `xml:"t"`

	RunProperties RunProperties `xml:"rPr"`
}

type Paragraph struct {
	XMLName xml.Name `xml:"p"`
	Runs    []Run    `xml:"r"`
}

type Body struct {
	XMLName    xml.Name    `xml:"body"`
	Paragraphs []Paragraph `xml:"p"`
}

type Document struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

func parseSubPropery(p *Property) bool {
	if p != nil {

		if p.Val != nil { // when property tag have w:val
			if *p.Val != "false" {
				return true
			} else {
				return false
			}
		} else { // when it doesn't have
			return true
		}

	}
	// if property isn't exist
	return false
}

func ParseProperties(p *[]Property, text *string) {
	if p != nil {
		for _, prop := range *p {
			if parseSubPropery(&prop) {
				switch prop.XMLName.Local {
				case "b":
					*text = "\033[1m" + *text
				case "i":
					*text = "\033[3m" + *text
				case "u":
					*text = "\033[4m" + *text
				case "shd":
					*text = "\033[7m" + *text
				}
			}
		}
	}

	*text += "\033[0m"
}

func Parse2Html() string {
	s, _ := os.ReadFile("./tmp/word/document.xml")
	var doc Document
	xml.Unmarshal(s, &doc)

	p := doc.Body.Paragraphs

	output := ""
	for _, v := range p {
		r := v.Runs
		hasAnyText := false
		for _, i := range r {
			if i.Text.Text != nil {
				hasAnyText = true
				text := i.Text.Text
				ParseProperties(i.RunProperties.Properties, text)
				fmt.Print(*text)
			}

		}
		if hasAnyText {
			fmt.Println()
		}
	}

	return output
}
