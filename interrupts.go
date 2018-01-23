package main

import (
	"fmt"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

type interruptsChecker struct {
	threshold       uint64
	latestIntPerSec chan uint64
}

func (ic interruptsChecker) Checks() []fthealth.Check {
	ic.latestIntPerSec = make(chan uint64)

	go ic.updateIntCount()

	check := fthealth.Check{
		BusinessImpact:   "System may become unresponsive",
		Name:             "Interrupts check",
		PanicGuide:       "Check the system with vmstat and investigate cause",
		Severity:         2,
		TechnicalSummary: "Number of interrupts as reported by /proc/stat is unusually high",
		Checker:          ic.intCheck,
	}

	return []fthealth.Check{check}
}

func (ic interruptsChecker) count() uint64 {
	d, err := linuxproc.ReadStat(*hostPath + "/proc/stat")
	if err != nil {
		panic(fmt.Sprintf("Cannot read disk info of %s file system.", *hostPath+"/proc/stat"))
	}
	return d.Interrupts
}

func (ic interruptsChecker) intCheck() (string, error) {
	perSec := <-ic.latestIntPerSec
	threshold := uint64(ic.threshold)
	if perSec > threshold {
		return fmt.Sprintf("%d", perSec), fmt.Errorf("%d interrupts per second. (>%d)", perSec, threshold)
	}
	return fmt.Sprintf("%d", perSec), nil
}

func (ic interruptsChecker) updateIntCount() {
	ticker := time.NewTicker(1 * time.Second)
	latestPerSec := uint64(0)
	prevInt := uint64(0)
	for {
		select {
		case ic.latestIntPerSec <- latestPerSec:
		case <-ticker.C:
			newInt := ic.count()
			if prevInt != 0 {
				latestPerSec = newInt - prevInt
			}
			prevInt = newInt
		}

	}
}
