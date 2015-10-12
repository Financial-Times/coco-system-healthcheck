package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"time"
)

func switchCount() uint64 {
	d, err := linuxproc.ReadStat(*hostPath + "/proc/stat")
	if err != nil {
		panic(fmt.Sprintf("Cannot read disk info of %s file system.", *hostPath+"/proc/stat"))
	}
	return d.ContextSwitches
}

func ctxCheck() error {
	first := switchCount()
	time.Sleep(200 * time.Millisecond)
	count := switchCount() - first
	perSec := count * 5
	threshold := uint64(120000)
	if perSec > threshold {
		return fmt.Errorf("%d context switches per second. (>%d)", perSec, threshold)
	}
	return nil
}

func ContextSwitches(checks *[]fthealth.Check) {
	impact := "System may become unresponsive"
	panicGuide := "Check the system with vmstat and investigate cause"

	check := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "Context switches check",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "Number of context switches as reported by /proc/stat is unusually high",
		Checker:          ctxCheck,
	}

	*checks = append(*checks, check)
}
