package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
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

func (mc memoryChecker) avMemoryCheck() error {

	meminfo, err := linuxproc.ReadMemInfo(*hostPath + "/proc/meminfo")
	if err != nil {
		return err
	}
	available := meminfo.MemAvailable
	total := meminfo.MemTotal
	availablePercent := float64(available) / float64(total) * 100
	if availablePercent < mc.thresholdPercent {
		return fmt.Errorf("Low available memory: %2.1f %%", availablePercent)
	}
	return nil
}
