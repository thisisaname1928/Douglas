package main

import (
	"fmt"

	"github.com/thisisaname1928/goParsingDocx/docx"
)

func main() {
	v, _ := docx.Parse2Fluid("./test.docx")
	docx.DecompressDocxMedia("./test.docx", "./media/")
	var index uint64 = 0
	i, q := docx.ParseFluid2Question(&index, v)
	for i == 0 {
		fmt.Println(i)
		fmt.Println(q.Content)
		fmt.Println(index)
		i, q = docx.ParseFluid2Question(&index, v)
	}

	// for _, val := range v {
	// 	val.Text += " "
	// 	for i, c := range val.Text {
	// 		for _, p := range val.Properties {
	// 			if p.Start == i {
	// 				for _, pr := range p.Property {
	// 					if pr.Type == docx.Bold && pr.Value != "false" {
	// 						fmt.Print("\x1b[1m")
	// 					}
	// 					if pr.Type == docx.Italic && pr.Value != "false" {
	// 						fmt.Print("\x1b[3m")
	// 					}
	// 					if pr.Type == docx.Underline && pr.Value != "false" {
	// 						fmt.Print("\x1b[4m")
	// 					}
	// 					if pr.Type == docx.Marked && pr.Value != "auto" {
	// 						fmt.Print("\x1b[43m")
	// 					}
	// 				}
	// 			}

	// 			if p.End == i {
	// 				for _, pr := range p.Property {
	// 					if pr.Type == docx.Bold {
	// 						fmt.Print("\x1b[22m")
	// 					}
	// 					if pr.Type == docx.Italic {
	// 						fmt.Print("\x1b[23m")
	// 					}
	// 					if pr.Type == docx.Underline {
	// 						fmt.Print("\x1b[24m")
	// 					}
	// 					if pr.Type == docx.Marked {
	// 						fmt.Print("\x1b[49m")
	// 					}
	// 				}
	// 			}
	// 		}
	// 		fmt.Printf("%c", c)
	// 	}
	// 	fmt.Println()
	// }
}
