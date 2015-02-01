// go run *.go --race
// should i resize image after capture like imago does? or let wkhtml screenshot at the sizr requested?
package main

import (
	"github.com/gorilla/mux"
	"github.com/zenazn/goji/graceful"
	"net/http"
	"os"
	"os/exec"
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
	r.HandleFunc("/user/{name:[a-z]+}/profile", profile).Methods("GET")

	http.Handle(ConRootUrl, r)

	LogInfo("Listening...")
	handler := http.TimeoutHandler(SetHeaders(Logger(rateLimit(r))), time.Second*25, "Timeout!")
	LogFatal(graceful.ListenAndServe(ConPort, handler))
}

func profile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	cmd := exec.Command("echo", "woooohoooo")

	output, _ := cmd.CombinedOutput()

	w.Write([]byte("Hello " + name + " ENV VAR IS: " + os.Getenv("DEIS_TEST_VAR") + " EXEC: " + string(output)))
}
