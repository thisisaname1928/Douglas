package dou

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/thisisaname1928/goParsingDocx/docx"
	"github.com/thisisaname1928/goParsingDocx/security"
)

// idk how 2 create a error so I create a constant instead
const (
	ERROR_KEY_NOT_MATCH = "ERROR_KEY_NOT_MATCH"
)

type MediaData struct {
	Name string
	Data []byte
}

type DouFile struct {
	Info  DouInfo
	Data  DouData
	Media []MediaData
}

// only work with path from tag <relationship> of docx file
func (file DouFile) OpenMedia(path string) []byte {
	ConvertPath(&path)

	for i := range file.Media {
		if file.Media[i].Name == ("/" + path) {
			return file.Media[i].Data
		}
	}

	return []byte{}
}

func douCheckMedia(path string) bool {
	ConvertPath(&path)
	arr := strings.Split(path, "/")

	if len(arr) < 2 {
		return false
	}

	return arr[len(arr)-2] == "media"
}

func Open(path string, key string) (DouFile, error) {
	var result DouFile

	// read info
	var info DouInfo
	dat, e := docx.DecompressFile(path, "info.json")
	if e != nil {
		return result, e
	}

	e = json.Unmarshal(dat, &info)
	if e != nil {
		return result, e
	}

	result.Info = info

	// check encryption
	var needDecrypting = result.Info.Encrypted
	if needDecrypting {
		if info.Key != security.EncryptKey(key) {
			return result, errors.New(ERROR_KEY_NOT_MATCH)
		}
	}

	// read test data
	var data DouData
	dat, e = docx.DecompressFile(path, "data.json")
	if e != nil {
		return result, e
	}
	if needDecrypting {
		dat, e = security.Decrypt(dat, key)
		if e != nil {
			return result, e
		}
	}

	e = json.Unmarshal(dat, &data)
	if e != nil {
		return result, e
	}

	result.Data = data

	// read media
	// we have to manual load all of them:(

	f, e := os.ReadFile(path)

	if e != nil {
		return result, e
	}

	a, e := zip.NewReader(bytes.NewReader(f), int64(len(f)))

	if e != nil {
		return result, e
	}

	for _, f := range a.File {
		if douCheckMedia(f.Name) {
			dat, e := f.Open()

			if e != nil {
				return result, e
			}

			var medDat MediaData
			medDat.Name = f.Name
			ConvertPath(&medDat.Name)
			medDat.Data, _ = io.ReadAll(dat) // Im too lazy to add a error handler

			if needDecrypting {
				medDat.Data, _ = security.Decrypt(medDat.Data, key)
			}

			result.Media = append(result.Media, medDat)
			dat.Close()
		}
	}

	return result, nil
}
