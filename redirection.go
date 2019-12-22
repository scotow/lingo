package lingo

import (
	"regexp"
	"time"
)

type Redirection struct {
	Payload string

	autoDelete *time.Timer
	creation   time.Time
}

func (r *Redirection) GetValidUrl() string {
	// Default to HTTP link.
	if prefixed, _ := regexp.MatchString("https?://.+", r.Payload); !prefixed {
		return "http://" + r.Payload
	} else {
		return r.Payload
	}
}
