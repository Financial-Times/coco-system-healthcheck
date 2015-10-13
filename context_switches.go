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

func (csc contextSwitchChecker) switchCount() uint64 {
	d, err := linuxproc.ReadStat(*hostPath + "/proc/stat")
	if err != nil {
		panic(fmt.Sprintf("Cannot read disk info of %s file system.", *hostPath+"/proc/stat"))
	}
	return d.ContextSwitches
}

func (csc contextSwitchChecker) ctxCheck() error {
	first := csc.switchCount()
	time.Sleep(200 * time.Millisecond)
	count := csc.switchCount() - first
	perSec := count * 5
	if perSec > csc.threshold {
		return fmt.Errorf("%d context switches per second. (>%d)", perSec, csc.threshold)
	}
	return nil
}
