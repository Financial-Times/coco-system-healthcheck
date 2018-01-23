package main

import (
	"log"
	"net/http"

	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

var (
	checks   []fthealth.Check
	hostPath *string
)

const (
	diskThresholdPercent   = 20
	memoryThresholdPercent = 15
)

func main() {
	app := cli.App("System-healthcheck", "A service that report on current VM status at __health")

	hostPath = app.String(cli.StringOpt{
		Name:   "hostPath",
		Value:  "",
		Desc:   "The dir path of the mounted host fs (in the container)",
		EnvVar: "SYS_HC_HOST_PATH",
	})

	checks = append(checks, diskFreeCheckerImpl{diskThresholdPercent}.Checks()...)
	checks = append(checks, memoryCheckerImpl{memoryThresholdPercent}.Checks()...)
	checks = append(checks, loadAverageCheckerImpl{}.Checks()...)
	checks = append(checks, ntpCheckerImpl{}.Checks()...)

	r := mux.NewRouter()
	timedHC := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{
			SystemCode:  "system-healthcheck",
			Name:        "System Healthcheck",
			Description: "Monitors system parameters.",
			Checks:      checks,
		},
		Timeout: 10 * time.Second,
	}
	r.HandleFunc("/__health", fthealth.Handler(timedHC))
	gtgService := newGtgService(diskThresholdPercent, memoryThresholdPercent)
	r.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(gtgService.Check))

	log.Print("Starting http server on 8080\n")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
