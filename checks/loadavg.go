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
		return fmt.Errorf("Load avg is above the recommended threshold: %#v", l)
	}

	return nil
}

func LoadAvg(checks *[]fthealth.Check) {
	loadAvgCheck := fthealth.Check{
		BusinessImpact:   "No newspaper soon...",
		Name:             "Load avg check",
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "CPU is quite busy, has too much work. Profile the services or spin up more boxes",
		Checker:          loadAvgCheck,
	}

	*checks = append(*checks, loadAvgCheck)
}
