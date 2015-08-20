package main

import (
	healthchecks "./checks"
	"flag"
	"github.com/Financial-Times/go-fthealth"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	checks   []fthealth.Check
	hostPath = flag.String("hostPath", "", "The dir path of the mounted host fs (in the container)")
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
