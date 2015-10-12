package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func loadAvgCheck() error {
	l, err := linuxproc.ReadLoadAvg(baseDir + "/proc/loadavg")

	if err != nil {
		return fmt.Errorf("Couldn't read loadavg data")
	}

	cpuInfo, err := linuxproc.ReadCPUInfo(baseDir + "/proc/cpuinfo")

	if err != nil {
		return fmt.Errorf("Couldn't read cpuinfo data")
	}

	fiveMinLimit := (1.5 * float64(cpuInfo.NumCPU()))
	fifteenMinLimit := (0.9 * float64(cpuInfo.NumCPU()))

	if l.Last5Min > fiveMinLimit || l.Last15Min > fifteenMinLimit {
		return fmt.Errorf("Load avg is above the recommended threshold: Last5Min: %2.2f, Last15Min: %2.2f", l.Last5Min, l.Last15Min)
	}

	return nil
}

func LoadAvg(checks *[]fthealth.Check) {
	loadAvgCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "CPU load average check",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "CPU is quite busy lately. This might not be a problem if it happens intermittently, however if it persists consider upgrading or adding new boxes.",
		Checker:          loadAvgCheck,
	}

	*checks = append(*checks, loadAvgCheck)
}
