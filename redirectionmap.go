package lingo

import (
	"log"
	"sync"
	"time"
)

var (
	maxTime = time.Unix(1<<63-62135596801, 0)
)

func NewRedirectionMap(timeout time.Duration, capacity int) *RedirectionMap {
	rm := new(RedirectionMap)
	rm.values = make(map[string]*Redirection)
	rm.timeout = timeout
	rm.capacity = capacity

	return rm
}

type RedirectionMap struct {
	lock     sync.RWMutex
	values   map[string]*Redirection
	timeout  time.Duration
	capacity int
}

func (r *RedirectionMap) Get(key string) (value *Redirection, ok bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	value, ok = r.values[key]
	log.Printf("'%s' fetched.", key)
	return
}

func (r *RedirectionMap) Add(key, link string) {
	// Lock map mutex.
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.capacity > 0 && len(r.values) >= r.capacity {
		r.deleteOldest()
	}

	// Cancel previous auto delete timer if present.
	if previous, contain := r.values[key]; contain && r.timeout > 0 {
		previous.autoDelete.Stop()
	}

	// Schedule auto delete.
	var autoDelete *time.Timer = nil
	if r.timeout > 0 {
		autoDelete = time.AfterFunc(r.timeout, func() {
			r.delete(key)
		})
	}

	// Add redirection to map.
	r.values[key] = &Redirection{link, autoDelete, time.Now()}

	log.Printf("'%s' added.", key)
}

func (r *RedirectionMap) delete(key string) {
	// Lock map mutex.
	r.lock.Lock()
	defer r.lock.Unlock()

	log.Printf("'%s' deleted.", key)
	delete(r.values, key)
}

func (r *RedirectionMap) deleteOldest() {
	// Find oldest redirection.
	min := maxTime
	var oldest string

	for key, value := range r.values {
		if value.creation.Before(min) {
			min = value.creation
			oldest = key
		}
	}

	// Cancel oldest timer.
	if r.timeout > 0 {
		r.values[oldest].autoDelete.Stop()
	}

	// Remove entry from map.
	delete(r.values, oldest)

	log.Printf("'%s' deleted (map full).", oldest)
}
