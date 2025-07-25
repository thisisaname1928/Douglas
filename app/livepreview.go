package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/docx"
)

func livePreview(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/frontend/livePreview/index.html")
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if e == nil {
		w.Write(f)
	}
}

func detectFileExt(path string) string {
	fileExtSpl := strings.Split(path, ".")
	ext := fileExtSpl[len(fileExtSpl)-1]

	switch ext {
	case "js":
		return "text/javascript"
	case "css":
		return "text/css"
	case "ico":
		return "image/x-icon"
	default:
		return ""
	}
}

func livePreviewRes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file, e := os.Open("./app/frontend/livePreview/" + vars["FILE"])
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if str := detectFileExt(vars["FILE"]); str != "" {
		w.Header().Add("Content-Type", str)
	} else {
		contentType := http.DetectContentType(f)
		w.Header().Add("Content-Type", contentType)
	}

	if e == nil {
		w.Write(f)
	}
}

type GenJsonAPIRequest struct {
	Path string `json:"path"`
}

type GenJsonResponse struct {
	Status    bool            `json:"status"`
	Error     string          `json:"error"`
	Questions []docx.Question `json:"questions"`
}

func livePreviewAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["NAME"] == "genJson" {
		decoder := json.NewDecoder(r.Body)
		var request GenJsonAPIRequest
		decoder.Decode(&request)

		res, e := GenQues(request.Path)

		var response GenJsonResponse
		encoder := json.NewEncoder(w)

		if e != nil {
			response.Status = false
			response.Error = fmt.Sprintf("%v", e)
		} else {
			response.Status = true
			response.Questions = res
		}
		encoder.Encode(&response)
	}
}

func convertPath(path *string) {
	r := []rune(*path)
	for i := range r {
		if r[i] == '\\' {
			r[i] = '/'
		}
	}

	*path = string(r)
}

func GenQues(path string) ([]docx.Question, error) {
	convertPath(&path)
	dest := strings.Split(path, "/")
	ddest := dest[len(dest)-1]

	_, e := os.Stat("./app/tests/" + ddest)
	if os.IsNotExist(e) {
		os.Mkdir("./app/tests/"+ddest, 0775)
	}

	if e != nil {
		return []docx.Question{}, e
	}

	docx.DecompressDocxMedia(path, "./app/media/"+ddest+"/")
	fluid, e := docx.Parse2Fluid(path)
	if e != nil {
		return []docx.Question{}, e
	}

	var index uint64 = 0
	var ques []docx.Question
	i, q := docx.ParseFluid2Question(&index, fluid)

	for i == 0 {
		ques = append(ques, q)
		i, q = docx.ParseFluid2Question(&index, fluid)
	}
	return ques, nil
}

func mediaRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file, e := os.Open("./media/" + vars["FILE"])
	if e != nil {
		w.Write([]byte{})
	}
	f, e := io.ReadAll(file)

	if str := detectFileExt("./media/" + vars["FILE"]); str != "" {
		w.Header().Add("Content-Type", str)
	} else {
		contentType := http.DetectContentType(f)
		w.Header().Add("Content-Type", contentType)
	}

	if e == nil {
		w.Write(f)
	}
}
