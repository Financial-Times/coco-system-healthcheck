package checks

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

	if l.Last1Min > 5.0 || l.Last5Min > 1.0 || l.Last15Min > 0.7 {
		return fmt.Errorf("Load avg is above the recommended threshold: Last1Min: %2.2f , Last5Min: %2.2f, Last15Min: %2.2f", l.Last1Min, l.Last5Min, l.Last15Min)
	}

	return nil
}

func LoadAvg(checks *[]fthealth.Check) {
	loadAvgCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be effected",
		Name:             "CPU load average check",
		PanicGuide:       "Please refer to the technical summary",
		Severity:         2,
		TechnicalSummary: "CPU is quite busy lately. This might not be a problem if it happens intermittently, however if it persists consider upgrading or adding new boxes.",
		Checker:          loadAvgCheck,
	}

	*checks = append(*checks, loadAvgCheck)
}
