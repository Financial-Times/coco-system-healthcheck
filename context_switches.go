package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"time"
)

type contextSwitchChecker struct {
	threshold uint64
}

func (csc contextSwitchChecker) Checks() []fthealth.Check {

	go csc.updateCsCount()

	check := fthealth.Check{
		BusinessImpact:   "System may become unresponsive",
		Name:             "Context switches check",
		PanicGuide:       "Check the system with vmstat and investigate cause",
		Severity:         2,
		TechnicalSummary: "Number of context switches as reported by /proc/stat is unusually high",
		Checker:          csc.ctxCheck,
	}

	return []fthealth.Check{check}
}

func (csc contextSwitchChecker) count() uint64 {
	d, err := linuxproc.ReadStat(*hostPath + "/proc/stat")
	if err != nil {
		panic(fmt.Sprintf("Cannot read disk info of %s file system.", *hostPath+"/proc/stat"))
	}
	return d.ContextSwitches
}

func (csc contextSwitchChecker) ctxCheck() error {
	perSec := <-latestIntPerSec
	if perSec > csc.threshold {
		return fmt.Errorf("%d context switches per second. (>%d)", perSec, csc.threshold)
	}
	return nil
}

var latestCsPerSec chan uint64 = make(chan uint64)

func (csc contextSwitchChecker) updateCsCount() {
	ticker := time.NewTicker(1 * time.Second)
	latestPerSec := uint64(0)
	prevInt := uint64(0)
	for {
		select {
		case latestCsPerSec <- latestPerSec:
		case <-ticker.C:
			newInt := csc.count()
			if prevInt != 0 {
				latestPerSec = newInt - prevInt
			}
			prevInt = newInt
		}

	}
}
