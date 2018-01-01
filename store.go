package observe

import (
	"time"
)

type Observation struct {
	URL        string
	RemoteAddr string
	Observed   time.Time
	Header     map[string][]string
}

type Storer interface {
	Store(*Observation) error
}
