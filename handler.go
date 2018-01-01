package observe

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

type HandlerOptions struct {
	ServerHeader string
}

func NewHandler(store Storer, opts *HandlerOptions) *Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}

	return &Handler{
		store:   store,
		options: opts,
	}
}

type Handler struct {
	store   Storer
	options *HandlerOptions
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if h.options.ServerHeader != "" {
		w.Header().Set("server", h.options.ServerHeader)
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	referrer := r.Referer()
	// if referrer is passed in url string, prefer that
	q := r.URL.Query()
	if r := q.Get("referrer"); r != "" {
		referrer = r
	}

	o := &Observation{
		URL:        referrer,
		RemoteAddr: r.RemoteAddr,
		Observed:   time.Now(),
		Header:     r.Header,
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "image/gif")
	w.Write([]byte(TransparentGIF))

	// if the user does not wish to be tracked, abort
	if r.Header.Get("DNT") == "1" {
		log.Info("not observed; user does not wish to be tracked")
	}

	go func() {
		if err := h.store.Store(o); err != nil {
			log.Errorf("unable to store observation: %v", err)
		}
	}()

	log.WithFields(log.Fields{
		"RemoteAddr": o.RemoteAddr,
		"URL":        o.URL,
	}).Info("observed")
}
