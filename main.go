package main

import (
	"fmt"
	"time"

	"github.com/thisisaname1928/goParsingDocx/testsvr"
)

func main() {

	s, e := testsvr.OpenOldTest("4c32f096-7ab8-447b-81bb-87258994da49", "abc")
	fmt.Println(e)

	var data testsvr.TestSessions
	data.NewSession("31569b3a-b643-4e5e-b51c-eccaa86be19a", "??", time.Now(), 3)
	data.UpdateAnswerSheet(0, "31569b3a-b643-4e5e-b51c-eccaa86be19a", [4]string{"a", "ok", "l", "m"})
	data.DoneSession(s.UUID, "31569b3a-b643-4e5e-b51c-eccaa86be19a", time.Now())

	fmt.Println(testsvr.GetIp())

	s.OpenServer("5999")

	//app.StartApp()
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
