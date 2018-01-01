package appengine

import (
	"log"
	"net/http"

	"github.com/fardog/observe"
	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/observe.gif", observeHandler)
}

func observeHandler(w http.ResponseWriter, r *http.Request) {
	projectID := appengine.AppID(appengine.NewContext(r))
	table := "observe.observations"

	opts := &observe.BigQueryOptions{
		Context: appengine.NewContext(r),
	}

	store, err := observe.NewBigQuery(projectID, table, opts)
	if err != nil {
		log.Fatalf("unable to instantiate bigquery client: %s", err.Error())
	}

	handler := observe.NewHandler(store, &observe.HandlerOptions{
		ServerHeader: "observe/1.0 (+https://github.com/fardog/observe)",
	})

	handler.HandleWithContext(appengine.NewContext(r), w, r)
}
