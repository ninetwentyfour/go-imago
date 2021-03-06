// go run *.go --race
package main

import (
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/yvasiyarov/gorelic"
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/zenazn/goji/graceful"
	"net/http"
	"time"
)

var agent = new(gorelic.Agent)

func init() {
	if ConNewRelicKey != "" {
		agent = gorelic.NewAgent()
		agent.NewrelicLicense = ConNewRelicKey
		agent.NewrelicName = "Imago Go"
		agent.Run()
	}

	createRedisPool()
}

// url structure http://imago.in/width/height/url/format
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", agent.WrapHTTPHandlerFunc(homeHandler)).Methods(ConMethod)
	r.HandleFunc("/get_image", agent.WrapHTTPHandlerFunc(oldHandler)).Methods(ConMethod)
	r.HandleFunc(ConUrl+"json", agent.WrapHTTPHandlerFunc(jsonHandler)).Methods(ConMethod)
	r.HandleFunc(ConUrl+"html", agent.WrapHTTPHandlerFunc(htmlHandler)).Methods(ConMethod)
	r.HandleFunc(ConUrl+"image", agent.WrapHTTPHandlerFunc(imageHandler)).Methods(ConMethod)

	http.Handle(ConRootUrl, r)

	LogInfo("Listening...")
	handler := http.TimeoutHandler(SetHeaders(Logger(rateLimit(r))), time.Second*25, "Timeout!")
	LogFatal(graceful.ListenAndServe(ConPort, handler))
}
