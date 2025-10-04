package testsvr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (fir DouglasFir) mediaRoute(w http.ResponseWriter, r *http.Request) {
	// NOTE: should add a check on user uuid
	v := mux.Vars(r)
	requestMedia := "media/" + v["FILE"] // remove /
	mediaResource, e := fir.Douglas.OpenMedia(requestMedia)

	fmt.Println(requestMedia)

	if e != nil {
		w.WriteHeader(404)
		return
	}

	fileExt := strings.Split(r.RequestURI, ".")

	w.Header().Add("Content-Type", "image/"+fileExt[len(fileExt)-1])
	w.Write(mediaResource)
}
