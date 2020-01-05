package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fardog/observe"
)

func getProjectID() (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	req.Header.Set("Metadata-Flavor", "Google")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	projectID, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(projectID), nil
}

func main() {
	projectID, err := getProjectID()
	if err != nil {
		log.Fatal("Could not read project id", err)
	}

	opts := &observe.BigQueryOptions{}
	store, err := observe.NewBigQuery(string(projectID), "observe.observations", opts)
	if err != nil {
		log.Fatalf("unable to instantiate bigquery client: %s", err.Error())
	}

	handler := observe.NewHandler(store, &observe.HandlerOptions{
		ServerHeader: "observe/1.0 (+https://github.com/fardog/observe)",
	})

	http.HandleFunc("/observe.gif", handler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
