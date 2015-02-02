// go run *.go --race
package main

import (
	"github.com/gorilla/mux"
	"github.com/yvasiyarov/gorelic"
	"github.com/zenazn/goji/graceful"
	"net/http"
	"time"
)

// var agent = new(gorelic.Agent)

func init() {
	createRedisPool()
}

// url structure http://imago.in/width/height/url/format
func main() {
	agent := new(gorelic.Agent)
	if ConNewRelicKey != "" {
		agent = gorelic.NewAgent()
		agent.NewrelicLicense = ConNewRelicKey
		agent.NewrelicName = "Imago Go"
		agent.Run()
	}

	r := mux.NewRouter()

	r.HandleFunc("/", agent.WrapHTTPHandlerFunc(homeHandler)).Methods(ConMethod)
	r.HandleFunc(ConUrl+"json", agent.WrapHTTPHandlerFunc(jsonHandler)).Methods(ConMethod)
	r.HandleFunc(ConUrl+"html", agent.WrapHTTPHandlerFunc(htmlHandler)).Methods(ConMethod)
	r.HandleFunc(ConUrl+"image", agent.WrapHTTPHandlerFunc(imageHandler)).Methods(ConMethod)

	http.Handle(ConRootUrl, r)

	LogInfo("Listening...")
	handler := http.TimeoutHandler(SetHeaders(Logger(rateLimit(r))), time.Second*25, "Timeout!")
	LogFatal(graceful.ListenAndServe(ConPort, handler))
}
