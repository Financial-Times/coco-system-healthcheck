package main

import (
	"flag"
	"fmt"
	healthchecks "github.com/Financial-Times/coco-system-healthcheck/checks"
	"github.com/Financial-Times/go-fthealth"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	checks   []fthealth.Check
	hostPath = flags.String("hostPath", "/host_dir", "The path where the host fs is mounted to the container")
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s.\n", r.URL.Path[1:])
}

func main() {
	flag.Parse()
	mux := mux.NewRouter()
	mux.HandleFunc("/", handler)

	healthchecks.RegisterChecks(hostDir, &checks)

	mux.HandleFunc("/__health", fthealth.Handler("myserver", "a server", checks...))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
