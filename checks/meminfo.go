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

func avMemeryCheck() error {
	meminfo := meminfo()
	available := meminfo.MemAvailable
	total := meminfo.MemTotal
	availablePercent := float64(available) / float64(total) * 100
	if availablePercent < 20 {
		fmt.Errorf("Low available memery: %2.f", availablePercent)
	}
	return nil
}

func MemInfo(checks *[]fthealth.Check) {
	meminfo := meminfo()
	available := meminfo.MemAvailable
	total := meminfo.MemTotal
	availablePercent := float64(available) / float64(total) * 100

	memAvCheck := fthealth.Check{
		BusinessImpact:   "No newspaper",
		Name:             fmt.Sprintf("Memory Available %d KB (%2.f%%)", available, availablePercent),
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "Spin up more boxes",
		Checker:          avMemeryCheck,
	}
	*checks = append(*checks, memAvCheck)
}
