package main

import (
	"flag"
	"github.com/scotow/lingo"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	DefaultExpireDuration 	= time.Hour
	DefaultCapacity 		= 100
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
		w.WriteHeader(http.StatusInternalServerError)
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

func listeningAddress() string {
	port, set := os.LookupEnv("PORT")
	if !set {
		port = "8080"
	}

	return ":" + port
}

func main() {
	durationFlag := flag.String("d", DefaultExpireDuration.String(), "expiration duration of links in seconds or using golang duration format. 0 for no timeout (memory leak)")
	capacityFlag := flag.Int("n", DefaultCapacity, "maximum capacity of the redirection map, 0 for infinite capacity")

	flag.Parse()

	var duration time.Duration
	durationSec, err := strconv.Atoi(*durationFlag)
	if err == nil {
		duration = time.Duration(durationSec) * time.Second
	} else {
		duration, err = time.ParseDuration(*durationFlag)
		if err != nil {
			log.Fatalln("invalid duration", duration)
			return
		}
	}

	redirectionMap = *lingo.NewRedirectionMap(duration, *capacityFlag)

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(listeningAddress(), nil))
}