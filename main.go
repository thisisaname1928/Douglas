package main

import (
	"fmt"

	"github.com/thisisaname1928/goParsingDocx/docx"
)

func main() {
	fmt.Println(docx.Parse2Html())
	docx.DecompressDocxMedia("./test.docx", "./testmedia/media/")
}
