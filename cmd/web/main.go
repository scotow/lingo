package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/scotow/lingo"
)

var (
	port     = flag.Int("p", 8080, "listening port")
	duration = flag.Duration("d", time.Hour, "expiration duration of links. 0 for no timeout (memory leak)")
	capacity = flag.Int("n", 100, "maximum capacity of the redirection map, 0 for infinite capacity (memory leak)")
)

var (
	redirectionMap lingo.RedirectionMap
)

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		redirect(w, r)
	case "POST":
		store(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	redirection, ok := redirectionMap.Get(r.URL.Path[1:])

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	http.Redirect(w, r, redirection.GetValidUrl(), http.StatusFound)
	_, _ = w.Write([]byte(redirection.Payload))
}

func store(_ http.ResponseWriter, r *http.Request) {
	redirectionMap.Add(r.URL.Path[1:], r.FormValue("url"))
}

func main() {
	flag.Parse()

	if *capacity < 0 {
		log.Fatalln("invalid capacity")
	}

	redirectionMap = *lingo.NewRedirectionMap(*duration, *capacity)

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
