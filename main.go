package main

import (
	"flag"
	healthchecks "github.com/Financial-Times/coco-system-healthcheck/checks"
	"github.com/Financial-Times/go-fthealth"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	checks   []fthealth.Check
	hostPath = flag.String("hostPath", "/host_dir", "The dir path where the host fs is mounted in the container")
)

func main() {
	flag.Parse()

	healthchecks.RegisterChecks(*hostPath, &checks)

	mux := mux.NewRouter()
	mux.HandleFunc("/__health", fthealth.Handler("myserver", "a server", checks...))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
