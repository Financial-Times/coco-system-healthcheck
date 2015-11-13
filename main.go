package main

import (
	"flag"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	checks   []fthealth.Check
	hostPath = flag.String("hostPath", "", "The dir path of the mounted host fs (in the container)")
)

func main() {
	flag.Parse()

	checks = append(checks, diskFreeChecker{20}.Checks()...)
	checks = append(checks, memoryChecker{20}.Checks()...)
	checks = append(checks, loadAverageChecker{}.Checks()...)
	//checks = append(checks, inodeChecker{1024}.Checks()...)
	//checks = append(checks, contextSwitchChecker{threshold: 120000}.Checks()...)
	//checks = append(checks, interruptsChecker{threshold: 10000}.Checks()...)
	//checks = append(checks, iopsChecker{1000}.Checks()...)
	checks = append(checks, ntpChecker{}.Checks()...)

	mux := mux.NewRouter()
	mux.HandleFunc("/__health", fthealth.Handler("myserver", "a server", checks...))

	log.Printf("Starting http server on 8080\n")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
