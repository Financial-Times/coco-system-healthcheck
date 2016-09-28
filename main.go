package main

import (
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"log"
	"net/http"
)

var (
	checks   []fthealth.Check
	hostPath *string
)

func main() {
	app := cli.App("System-healthcheck", "A service that report on current VM status at __health")

	hostPath = app.String(cli.StringOpt{
		Name:   "hostPath",
		Value:  "",
		Desc:   "The dir path of the mounted host fs (in the container)",
		EnvVar: "SYS_HC_HOST_PATH",
	})

	checks = append(checks, diskFreeChecker{20}.Checks()...)
	checks = append(checks, memoryChecker{15}.Checks()...)
	checks = append(checks, loadAverageChecker{}.Checks()...)
	checks = append(checks, ntpChecker{}.Checks()...)
	checks = append(checks, tcpChecker{}.Checks()...)

	mux := mux.NewRouter()
	mux.HandleFunc("/__health", fthealth.Handler("myserver", "a server", checks...))

	log.Printf("Starting http server on 8080\n")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
