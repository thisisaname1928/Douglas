package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/go-getter"
)

func updatePageRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/update/")
}

func updatePage(w http.ResponseWriter, r *http.Request) {
	b, e := os.ReadFile("./app/frontend/update/index.html")

	if e != nil {
		w.WriteHeader(404)
		return
	}

	w.Write(b)
}

func downloadUpdate(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Msg    string `json:"msg"`
		Status bool   `json:"status"`
	}

	os.WriteFile("./tmp.json", []byte{}, os.FileMode(0777))

	encoder := json.NewEncoder(w)

	client := getter.Client{Src: "https://raw.githubusercontent.com/thisisaname1928/Douglas/refs/heads/master/appVersion.json", Dst: "./tmp.json", Mode: getter.ClientModeFile}

	e := client.Get()
	if e != nil {
		response.Msg = "Tải file lỗi!"
		response.Status = false
		encoder.Encode(response)
		return
	}

	f, e := os.ReadFile("./tmp.json")
	if e != nil {
		response.Msg = "Tải file lỗi!: " + fmt.Sprint(e)
		response.Status = false
		encoder.Encode(response)
		return
	}

	var config appVersion
	e = json.Unmarshal(f, &config)
	if e != nil {
		response.Msg = "Tải file lỗi!" + fmt.Sprint(e)
		response.Status = false
		encoder.Encode(response)
		return
	}

	// download real update file
	client = getter.Client{Src: config.UpdatePathWin64, Dst: "./update.zip", Mode: getter.ClientModeFile}

	e = client.Get()
	if e != nil {
		response.Msg = "Tải file lỗi!" + fmt.Sprint(e)
		response.Status = false
		encoder.Encode(response)
		return
	}

	// renew appVersion
	config.ShouldUpdate = true
	b, e := json.Marshal(config)
	if e != nil {
		response.Msg = "Tải file lỗi!" + fmt.Sprint(e)
		response.Status = false
		encoder.Encode(response)
		return
	}

	e = os.WriteFile("./appVersion.json", b, os.FileMode(0777))
	if e != nil {
		response.Msg = "Lỗi nội bộ!"
		response.Status = false
		encoder.Encode(response)
		return
	}

	response.Msg = "ok"
	response.Status = true
	encoder.Encode(response)
}
