package main

import (
	"fmt"

	"errors"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

type loadAverageChecker interface {
	Checks() []fthealth.Check
	Check() (string, error)
}

type loadAverageCheckerImpl struct{}

func (lac loadAverageCheckerImpl) Checks() []fthealth.Check {
	check := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "CPU load average check",
		PanicGuide:       "https://dewey.ft.com/upp-system-healthcheck.html",
		Severity:         2,
		TechnicalSummary: "CPU is quite busy lately. This might not be a problem if it happens intermittently, however if it persists consider upgrading or adding new boxes.",
		Checker:          lac.Check,
	}

	return []fthealth.Check{check}
}

func (lac loadAverageCheckerImpl) Check() (string, error) {
	l, err := linuxproc.ReadLoadAvg(*hostPath + "/proc/loadavg")

	if err != nil {
		return "", errors.New("Couldn't read loadavg data")
	}

	cpuInfo, err := linuxproc.ReadCPUInfo(*hostPath + "/proc/cpuinfo")

	if err != nil {
		return "", errors.New("Couldn't read cpuinfo data")
	}

	fiveMinLimit := (1.8 * float64(cpuInfo.NumCPU()))
	fifteenMinLimit := (1.1 * float64(cpuInfo.NumCPU()))

	if l.Last5Min > fiveMinLimit || l.Last15Min > fifteenMinLimit {
		return fmt.Sprintf("Last5Min: %2.2f, Last15Min: %2.2f", l.Last5Min, l.Last15Min), fmt.Errorf("Load avg is above the recommended threshold: Last5Min: %2.2f, Last15Min: %2.2f", l.Last5Min, l.Last15Min)
	}

	return fmt.Sprintf("Last5Min: %2.2f, Last15Min: %2.2f", l.Last5Min, l.Last15Min), nil
}
