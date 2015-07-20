package main

import (
	healthchecks "./checks"
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	checks []fthealth.Check
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s.\n", r.URL.Path[1:])
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", handler)

	healthchecks.DiskChecks(&checks)

	mux.HandleFunc("/__health", fthealth.Handler("myserver", "a server", checks...))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
