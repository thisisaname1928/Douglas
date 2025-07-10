package docx

import "encoding/xml"

type Drawing struct {
	XMLName xml.Name        `xml:"drawing"`
	Anchor  *ImageContainer `xml:"anchor"` // anchor mean image that isn't at the same line with text
	Inline  *ImageContainer `xml:"inline"`
}

type ImageContainer struct { // this can be an anchor or an inline
	XMLName xml.Name
	Graphic Graphic `xml:"graphic"`
}

type Blip struct {
	XMLName xml.Name `xml:"blip"`
	Embed   *string  `xml:"embed,attr"`
}

type BlipFill struct {
	XMLName xml.Name `xml:"blipFill"`
	Blip    Blip     `xml:"blip"`
}

type Picture struct {
	XMLName  xml.Name `xml:"pic"`
	BlipFill BlipFill `xml:"blipFill"`
}

type GraphicData struct {
	XMLName xml.Name `xml:"graphicData"`
	Picture *Picture `xml:"pic"`
}

type Graphic struct {
	XMLName     xml.Name    `xml:"graphic"`
	GraphicData GraphicData `xml:"graphicData"`
}

func ParseDrawing(drawing *Drawing) string {
	if drawing == nil { // check if there is a <w:drawing>
		return ""
	}

	res := ""

	r, v := parseImageContainer(drawing.Anchor)
	if r {
		res = v

		return res
	}

	r, v = parseImageContainer(drawing.Inline)
	if r {
		res = v

		return res
	}

	return res
}

func parseImageContainer(container *ImageContainer) (bool, string) {
	if container == nil {
		return false, ""
	}

	picture := container.Graphic.GraphicData.Picture

	if picture == nil {
		return false, ""
	}

	res := ""

	blip := picture.BlipFill.Blip

	if blip.Embed != nil {
		res = *blip.Embed
		if res != "" {
			return true, res
		}
	}

	return false, res
}

// get rId from word/_rels/document.xml.rels

type Relationship struct {
	XMLName xml.Name `xml:"Relationship"`
	ID      string   `xml:"Id,attr"`
	Target  string   `xml:"Target,attr"`
}

type Relationships struct {
	XMLName       xml.Name       `xml:"Relationships"`
	Relationships []Relationship `xml:"Relationship"`
}

func GetRID(path string) map[string]string {
	res := map[string]string{}

	content, e := DecompressFile(path, "word/_rels/document.xml.rels")

	if e != nil {
		return res
	}

	var rIDs Relationships
	e = xml.Unmarshal(content, &rIDs)

	if e != nil { // if the xml is invalid
		return res
	}

	for _, v := range rIDs.Relationships {
		res[v.ID] = v.Target
	}

	return res
}
