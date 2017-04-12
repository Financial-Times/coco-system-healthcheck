package main

import (
	"log"
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
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

	checks = append(checks, diskFreeCheckerImpl{20}.Checks()...)
	checks = append(checks, memoryCheckerImpl{15}.Checks()...)
	checks = append(checks, loadAverageCheckerImpl{}.Checks()...)
	checks = append(checks, ntpCheckerImpl{}.Checks()...)
	checks = append(checks, tcpCheckerImpl{}.Checks()...)

	router := mux.NewRouter()
	router.HandleFunc("/__health", fthealth.Handler("myserver", "a server", checks...))

	gtgService := newGtgService(20, 15)
	gtgHandler := httphandlers.NewGoodToGoHandler(gtgService.Check)
	router.HandleFunc(httphandlers.GTGPath, gtgHandler)

	log.Print("Starting http server on 8080\n")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
