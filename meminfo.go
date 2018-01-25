package main

import (
	"fmt"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

type memoryChecker interface {
	Checks() []fthealth.Check
	AvMemoryCheck() (string, error)
}

type memoryCheckerImpl struct {
	thresholdPercent float64
}

func (mc memoryCheckerImpl) Checks() []fthealth.Check {

	check := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "Memory load check",
		PanicGuide:       "https://dewey.ft.com/upp-system-healthcheck.html",
		Severity:         2,
		TechnicalSummary: "Check the memory usage of services/containers on this host, please confirm these values.",
		Checker:          mc.AvMemoryCheck,
	}

	return []fthealth.Check{check}
}

func (mc memoryCheckerImpl) AvMemoryCheck() (string, error) {
	meminfo, err := linuxproc.ReadMemInfo(*hostPath + "/proc/meminfo")
	if err != nil {
		return "", err
	}
	available := meminfo.MemAvailable
	total := meminfo.MemTotal
	availablePercent := float64(available) / float64(total) * 100
	if availablePercent < mc.thresholdPercent {
		return fmt.Sprintf("%2.1f%%", availablePercent), fmt.Errorf("Low available memory: %2.1f%%", availablePercent)
	}
	return fmt.Sprintf("%2.1f%%", availablePercent), nil
}
