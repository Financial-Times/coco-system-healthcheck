package checks

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func meminfo() linuxproc.MemInfo {
	meminfo, err := linuxproc.ReadMemInfo(baseDir + "/proc/meminfo")
	if err != nil {
		fmt.Errorf("meminfo read fail")
	}
	return *meminfo
}

func avMemoryCheck() error {
	meminfo := meminfo()
	available := meminfo.MemAvailable
	total := meminfo.MemTotal
	availablePercent := float64(available) / float64(total) * 100
	if availablePercent < 20 {
		fmt.Errorf("Low available memory: %2.1f %%", availablePercent)
	}
	return nil
}

func MemInfo(checks *[]fthealth.Check) {
	memAvCheck := fthealth.Check{
		BusinessImpact:   "No newspaper",
		Name:             "Memory load check",
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "Spin up more boxes",
		Checker:          avMemoryCheck,
	}
	*checks = append(*checks, memAvCheck)
}
