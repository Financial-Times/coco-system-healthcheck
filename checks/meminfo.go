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
		return fmt.Errorf("Low available memory: %2.1f %%", availablePercent)
	}
	return nil
}

func MemInfo(checks *[]fthealth.Check) {
	memAvCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be effected",
		Name:             "Memory load check",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "Check the memory usage of services/containers on this host, please proceed conform these values.",
		Checker:          avMemoryCheck,
	}
	*checks = append(*checks, memAvCheck)
}
