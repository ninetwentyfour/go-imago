// go run *.go --race
// should i resize image after capture like imago does? or let wkhtml screenshot at the sizr requested?
package main

import (
	"github.com/gorilla/mux"
	"github.com/zenazn/goji/graceful"
	"net/http"
	"time"
)

func init() {
	createRedisPool()
}

// url structure http://imago.in/width/height/url/format
func main() {
	r := mux.NewRouter()

	r.HandleFunc(ConUrl+"json", jsonHandler).Methods(ConMethod)
	r.HandleFunc(ConUrl+"html", htmlHandler).Methods(ConMethod)
	r.HandleFunc(ConUrl+"image", imageHandler).Methods(ConMethod)

	http.Handle(ConRootUrl, r)

	LogInfo("Listening...")
	handler := http.TimeoutHandler(SetHeaders(Logger(rateLimit(r))), time.Second*25, "Timeout!")
	LogFatal(graceful.ListenAndServe(ConPort, handler))
}
