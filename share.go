package main

import (
	"log"
	"net/http"
	"sync"
	"time"
	"regexp"
	"math"
)

const (
	timeout = time.Hour
	max = 100
)

type redirection struct {
	url string
	autoDelete *time.Timer
	creation int64
}

type redirectionMap struct {
	sync.RWMutex
	values map[string]*redirection
}

func newRedirectionMap() *redirectionMap {
	return &redirectionMap{
		values: make(map[string]*redirection),
	}
}

func (r *redirectionMap) get(key string) (value *redirection, ok bool) {
	r.RLock()
	defer r.RUnlock()

	value, ok = r.values[key]
	log.Printf("'%s' fetched.", key)
	return
}

func (r *redirectionMap) add(key, link string) {
	// Default to HTTP link.
	if prefixed, _ := regexp.MatchString("https?://.+", link); !prefixed {
		link = "http://" + link
	}

	// Lock map muttex.
	r.Lock()
	defer r.Unlock()

	if len(r.values) >= max {
		r.deleteOldest()
	}

	// Cancel previous auto delete timer if present.
	if previous, contain := r.values[key]; contain {
		previous.autoDelete.Stop()
	}

	// Schedule auto delete.
	autoDelete := time.AfterFunc(timeout, func() {
		r.delete(key)
	})

	// Add redirection to map.
	r.values[key] = &redirection{link, autoDelete, time.Now().UnixNano()}

	log.Printf("'%s' added.", key)
}

func (r *redirectionMap) delete(key string) {
	// Lock map muttex.
	r.Lock()
	defer r.Unlock()

	log.Printf("'%s' deleted.", key)
	delete(r.values, key)
}

func (r *redirectionMap) deleteOldest() {
	// Find oldest redirection.
	min := int64(math.MaxInt64)
	var oldest string

	for key, value := range r.values {
		if value.creation < min {
			min = value.creation
			oldest = key
		}
	}

	// Cancel oldest timer.
	r.values[oldest].autoDelete.Stop()

	// Remove entry from map.
	delete(r.values, oldest)

	log.Printf("'%s' deleted (map full).", oldest)
}

type shareHandler struct {
	*redirectionMap
}

func (h *shareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.Redirect(w, r)
	case "POST":
		h.Store(w, r)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *shareHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	redirection, ok := h.get(r.URL.Path[1:])

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, redirection.url, http.StatusFound)
}

func (h *shareHandler) Store(w http.ResponseWriter, r *http.Request) {
	h.add(r.URL.Path[1:], r.FormValue("url"))
}

func main() {
	r := newRedirectionMap()
	http.Handle("/", &shareHandler{r})

	log.Fatal(http.ListenAndServe("localhost:6000", nil))
}