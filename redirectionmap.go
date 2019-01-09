package lingo

import (
	"log"
	"math"
	"sync"
	"time"
)

func NewRedirectionMap(timeout time.Duration, capacity int) *RedirectionMap {
	return &RedirectionMap{
		values: make(map[string]*Redirection),
		timeout: timeout,
		capacity: capacity,
	}
}

type RedirectionMap struct {
	lock sync.RWMutex
	values map[string]*Redirection
	timeout time.Duration
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

	if len(r.values) >= r.capacity {
		r.deleteOldest()
	}

	// Cancel previous auto delete timer if present.
	if previous, contain := r.values[key]; contain {
		previous.autoDelete.Stop()
	}

	// Schedule auto delete.
	autoDelete := time.AfterFunc(r.timeout, func() {
		r.delete(key)
	})

	// Add redirection to map.
	r.values[key] = &Redirection{link, autoDelete, time.Now().UnixNano()}

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