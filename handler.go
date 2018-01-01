package observe

import (
	"context"
	"log"
	"net/http"
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
	ctx := context.Background()

	h.HandleWithContext(ctx, w, r)
}

func (h *Handler) HandleWithContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.options.ServerHeader != "" {
		w.Header().Set("server", h.options.ServerHeader)
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	o, err := NewObservation(r)
	if err != nil {
		log.Printf("unable to create observation from request: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "image/gif")
	w.Write([]byte(TransparentGIF))

	// if the user does not wish to be tracked, abort
	if r.Header.Get("DNT") == "1" {
		log.Println("not observed; user does not wish to be tracked")
		return
	}

	if err := h.store.Store(ctx, o); err != nil {
		log.Printf("unable to store observation: %v", err)
	}

	log.Printf("observed %v, %v", o.RemoteAddr, o.URL)
}
