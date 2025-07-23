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

func calStart(len int) int {
	if len == 0 {
		return 0
	}

	return len - 1
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
					currentLine.Text += " "
					currentProperty.Start = len(currentLine.Text)
					currentProperty.End = len("[image]") + currentProperty.Start
					currentLine.Text += "[image]"
					currentLine.Text += " "
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
					currentProperty.Start = len(currentLine.Text)
					currentProperty.End = len(*r.Text.Text) + currentProperty.Start

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
		if len(currentLine.Text) > 0 {
			currentLine.Text += " "
			str = append(str, currentLine)
		}
	}

	return str, nil
}
