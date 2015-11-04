package main

import (
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

type memoryChecker struct {
	thresholdPercent float64
}

func (mc memoryChecker) Checks() []fthealth.Check {

	check := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "Memory load check",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "Check the memory usage of services/containers on this host, please proceed conform these values.",
		Checker:          mc.avMemoryCheck,
	}

	return []fthealth.Check{check}
}

func (mc memoryChecker) avMemoryCheck() (string, error) {
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
