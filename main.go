package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

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

func main() {
	f, e := os.Open("./tmp/word/document.xml")
	if e != nil {
		panic(e)
	}

	s, e := io.ReadAll(f)

	if e != nil {
		panic(e)
	}

	var doc Document
	e = xml.Unmarshal(s, &doc)

	if e != nil {
		panic(e)
	}

	for i, p := range doc.Body.Paragraphs {
		fmt.Println("p: ", i)
		for j, r := range p.Runs {
			fmt.Println("r: ", j)
			for k, t := range r.Texts {
				fmt.Println("t: ", k)
				fmt.Println(t.Text)
			}
		}
	}

}
