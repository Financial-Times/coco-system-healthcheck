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
	checks                   []fthealth.Check
	hostPath                 *string
	awsEbsMountPath          *string
	ntpTimeDrift             *string
	ntpPollingPeriod         *string
	rootDiskThresholdPercent *int
	mountsThresholdPercent   *int
)

const (
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

	rootDiskThresholdPercent = app.Int(cli.IntOpt{
		Name:   "rootDiskThresholdPercent",
		Value:  20,
		Desc:   "For monitoring the root disk of the instances: when the free space goes bellow this percentage, the health check will fail",
		EnvVar: "ROOT_DISK_THRESHOLD",
	})

	awsEbsMountPath = app.String(cli.StringOpt{
		Name:   "awsEbsMountPath",
		Value:  "",
		Desc:   "The folder path where the AWS EBSs are mounted by Kubernetes",
		EnvVar: "AWS_EBS_MOUNTS_PATH",
	})

	mountsThresholdPercent = app.Int(cli.IntOpt{
		Name:   "mountsThresholdPercent",
		Value:  10,
		Desc:   "For monitoring the AWS EBSs that are mounted by Kubernetes: when the free space goes bellow this percentage, the health check will fail",
		EnvVar: "MOUNTS_THRESHOLD",
	})

	ntpTimeDrift = app.String(cli.StringOpt{
		Name:   "ntpTimeDrift",
		Value:  "2s",
		Desc:   "Time drift to allow for in NTP check, either in past or future",
		EnvVar: "NTP_TIME_DRIFT",
	})

	ntpTimeDriftDuration, err := time.ParseDuration(*ntpTimeDrift)
	if err != nil {
		ntpTimeDriftDuration = time.Second * 2
		log.Printf("Invalid time drift, using default 2s")
	}

	ntpPollingPeriod = app.String(cli.StringOpt{
		Name:   "ntpPollingPeriod",
		Value:  "1m",
		Desc:   "Polling period for NTP check",
		EnvVar: "NTP_POLLING_PERIOD",
	})

	ntpPollingPeriodDuration, err := time.ParseDuration(*ntpPollingPeriod)
	if err != nil {
		ntpPollingPeriodDuration = time.Minute
		log.Printf("Invalid polling period drift, using default 1m")
	}

	ntpChecker := &ntpCheckerImpl{
		timeDrift:     ntpTimeDriftDuration,
		pollingPeriod: ntpPollingPeriodDuration,
	}

	checks = append(checks, diskFreeCheckerImpl{*rootDiskThresholdPercent, *mountsThresholdPercent}.Checks()...)
	checks = append(checks, memoryCheckerImpl{memoryThresholdPercent}.Checks()...)
	checks = append(checks, loadAverageCheckerImpl{}.Checks()...)
	checks = append(checks, ntpChecker.Checks()...)

	r := mux.NewRouter()
	timedHC := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{
			SystemCode:  "upp-system-healthcheck",
			Name:        "System Healthcheck",
			Description: "Monitors system parameters.",
			Checks:      checks,
		},
		Timeout: 10 * time.Second,
	}
	r.HandleFunc("/__health", fthealth.Handler(timedHC))
	gtgService := newGtgService(*rootDiskThresholdPercent, *mountsThresholdPercent, memoryThresholdPercent)
	r.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(gtgService.Check))

	log.Print("Starting http server on 8080\n")
	panic(http.ListenAndServe(":8080", r))
}
