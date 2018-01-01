package observe

import (
	"context"
	"net/http"
	"time"
)

func NewObservation(r *http.Request) (*Observation, error) {
	referrer := r.Referer()
	// if referrer is passed in url string, prefer that
	q := r.URL.Query()
	if r := q.Get("referrer"); r != "" {
		referrer = r
	}

	return &Observation{
		URL:        referrer,
		RemoteAddr: r.RemoteAddr,
		Observed:   time.Now(),
		Header:     r.Header,
	}, nil
}

type Observation struct {
	URL        string
	RemoteAddr string
	Observed   time.Time
	Header     map[string][]string
}

type Storer interface {
	Store(context.Context, *Observation) error
}
