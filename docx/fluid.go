package docx

const (
	Italic    = 6
	Bold      = 2
	Marked    = 3
	Underline = 4
	ImgSource = 5
)

type Prop struct {
	Type  int
	Value string
}

type FluidProperty struct {
	Start    int
	End      int
	Property []Prop
}

type FluidString struct {
	Text       string
	Properties []FluidProperty
}

func calLen(inp string) int {
	return len([]rune(inp))
}

func Parse2Fluid(path string) ([]FluidString, error) {
	var str []FluidString

	doc, e := Parse(path)

	if e != nil {
		return nil, e
	}

	paras := doc.Body.Paragraphs
	rIDTab := GetRID(path)

	for _, v := range paras { // parse per line
		var currentLine FluidString
		for _, r := range v.Runs {
			if r.Drawing != nil { // parse drawing
				for _, drawing := range *r.Drawing {
					var currentProperty FluidProperty
					currentLine.Text += ""
					currentProperty.Start = calLen(currentLine.Text)
					currentProperty.End = currentProperty.Start
					var p Prop
					p.Type = ImgSource
					p.Value = rIDTab[ParseDrawing(&drawing)]

					currentProperty.Property = append(currentProperty.Property, p)
					currentLine.Properties = append(currentLine.Properties, currentProperty)
				}
			}
			if r.Text.Text != nil {
				if r.RunProperties.Properties != nil { // parse properties
					var currentProperty FluidProperty
					currentProperty.Start = calLen(currentLine.Text)
					currentProperty.End = calLen(*r.Text.Text) + currentProperty.Start

					for _, pr := range *r.RunProperties.Properties {
						var p Prop
						switch pr.XMLName.Local {
						case "b":
							p.Type = Bold
							if pr.Val != nil {
								p.Value = *pr.Val
							}
						case "i":
							p.Type = Italic
							if pr.Val != nil {
								p.Value = *pr.Val
							}
						case "u":
							p.Type = Underline
							if pr.Val != nil {
								p.Value = *pr.Val
							}
						case "shd":
							p.Type = Marked
							if pr.Fill != nil {
								p.Value = *pr.Fill
							}
						}
						currentProperty.Property = append(currentProperty.Property, p)

					}
					currentLine.Properties = append(currentLine.Properties, currentProperty)
				}
				currentLine.Text += *r.Text.Text

			}
		}
		currentLine.Text += " "
		str = append(str, currentLine)

	}

	return str, nil
}

func DelFirstCharacterStr(str *string) {
	r := []rune(*str)
	*str = string(r[1:])
}

func DelNCharacterStr(str *string, n int) {
	for i := 0; i < n; i++ {
		DelFirstCharacterStr(str)
	}
}

func DelFirstCharacterRune(str *[]rune) {
	tmp := *str
	*str = tmp[1:]
}

func DelNCharacterRune(str *[]rune, n int) {
	for i := 0; i < n; i++ {
		DelFirstCharacterRune(str)
	}
}

func DelFirstCharacter(str *FluidString) {
	tmpText := ""
	r := []rune(str.Text) // go use utf8 by default so this is the safe way to delete
	if len(r) > 0 {
		tmpText = string(r[1:])
	}

	for i := range str.Properties {
		if str.Properties[i].Start > 0 {
			str.Properties[i].Start--
		}

		if str.Properties[i].End > 0 {
			str.Properties[i].End--
		}
	}

	str.Text = tmpText
}

func DelNCharacter(str *FluidString, n int) {
	for i := 0; i < n; i++ {
		DelFirstCharacter(str)
	}
}

func ParseFluid2Html(str FluidString) string {
	output := ""
	addLabel := false

	if str.Text == "" {
		str.Text += " "
	}

	text := []rune(str.Text)
	for _, v := range text {
		if v != ' ' {
			addLabel = true
			break
		}
	}

	for i, c := range text {
		for _, prop := range str.Properties {
			if prop.Start == i {
				for _, pr := range prop.Property {
					if pr.Type == ImgSource {
						output += "<img src=\"" + pr.Value + "\">"
					}
					if pr.Type == Bold && pr.Value != "false" {
						output += "<b>"
					}
					if pr.Type == Italic && pr.Value != "false" {
						output += "<i>"
					}
					if pr.Type == Underline && pr.Value != "false" {
						output += "<u>"
					}
					if pr.Type == Marked && pr.Value != "auto" {
						output += "<mark>"
					}
				}
			}

			if prop.End == i {
				for _, pr := range prop.Property {
					if pr.Type == Bold && pr.Value != "false" {
						output += "</b>"
					}
					if pr.Type == Italic && pr.Value != "false" {
						output += "</i>"
					}
					if pr.Type == Underline && pr.Value != "false" {
						output += "</u>"
					}
					if pr.Type == Marked && pr.Value != "auto" {
						output += "</mark>"
					}
				}
			}

		}
		output += string(c)
	}

	if addLabel {
		output = "<label class=\"ques_content\">" + output + "</label>"
	}
	return output
}

func ParseFluid2HtmlNonMark(str FluidString) string {
	output := ""
	addLabel := false

	if str.Text == "" {
		str.Text += " "
	}

	text := []rune(str.Text)
	for _, v := range text {
		if v != ' ' {
			addLabel = true
			break
		}
	}

	for i, c := range text {
		for _, prop := range str.Properties {
			if prop.Start == i {
				for _, pr := range prop.Property {
					if pr.Type == ImgSource {
						output += "<img src=\"" + pr.Value + "\">"
					}
					if pr.Type == Bold && pr.Value != "false" {
						output += "<b>"
					}
					if pr.Type == Italic && pr.Value != "false" {
						output += "<i>"
					}
					if pr.Type == Underline && pr.Value != "false" {
						output += "<u>"
					}
				}
			}

			if prop.End == i {
				for _, pr := range prop.Property {
					if pr.Type == Bold && pr.Value != "false" {
						output += "</b>"
					}
					if pr.Type == Italic && pr.Value != "false" {
						output += "</i>"
					}
					if pr.Type == Underline && pr.Value != "false" {
						output += "</u>"
					}
				}
			}

		}
		output += string(c)
	}

	if addLabel {
		output = "<label class=\"ques_content\">" + output + "</label>"
	}
	return output
}

// a1, a2 as source index, b1 b2 as chop index
func chopRange(a1 int, a2 int, b1 int, b2 int) (int, int) {
	c1, c2 := -1, -1

	if b1 <= a1 && b2 >= a2 {

		c1, c2 = a1, a2
	} else if b1 <= a1 && b2 <= a2 && b2 >= a1 { // b2 in range [a1,a2] and b1 not in
		c1, c2 = a1, b2

	} else if b1 >= a1 && b1 <= a2 && b2 >= a2 { // b1 in range [a1,a2] and b2 not in
		c1, c2 = b1, a2

	} else if b1 > a1 && b2 < a2 {
		c1, c2 = b1, b2

	}

	return c1, c2
}

func CopyFluid(fluid FluidString, beginIndex int, endIndex int) FluidString {
	var res FluidString
	aRune := []rune(fluid.Text)

	for i := beginIndex; i <= endIndex; i++ {
		if i >= len(aRune) {
			break
		}

		res.Text += string(aRune[i])
	}

	//check 4 property
	for i := range fluid.Properties {
		b, e := chopRange(beginIndex, endIndex, fluid.Properties[i].Start, fluid.Properties[i].End)

		if min(b, e) != -1 {
			var prop FluidProperty
			prop.Property = fluid.Properties[i].Property
			prop.Start, prop.End = b, e
			res.Properties = append(res.Properties, prop)
		}
	}

	return res
}
