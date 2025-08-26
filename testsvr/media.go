package testsvr

import (
	"net/http"
	"strings"
)

func (fir DouglasFir) mediaRoute(w http.ResponseWriter, r *http.Request) {
	// NOTE: should add a check on user uuid
	requestMedia := r.RequestURI[1:]
	mediaResource, e := fir.Douglas.OpenMedia(requestMedia)

	if e != nil {
		w.WriteHeader(404)
		return
	}

	fileExt := strings.Split(r.RequestURI, ".")

	w.Header().Add("Content-Type", "image/"+fileExt[len(fileExt)-1])
	w.Write(mediaResource)
}
