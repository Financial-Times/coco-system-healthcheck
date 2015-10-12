package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"time"
)

func count() uint64 {
	d, err := linuxproc.ReadStat(*hostPath + "/proc/stat")
	if err != nil {
		panic(fmt.Sprintf("Cannot read disk info of %s file system.", *hostPath+"/proc/stat"))
	}
	return d.Interrupts
}

func intCheck() error {
	perSec := <-latestIntPerSec
	threshold := uint64(3000)
	if perSec > threshold {
		return fmt.Errorf("%d interrupts per second. (>%d)", perSec, threshold)
	}
	return nil
}

var latestIntPerSec chan uint64 = make(chan uint64)

func updateIntCount() {
	ticker := time.NewTicker(1 * time.Second)
	latestPerSec := uint64(0)
	prevInt := uint64(0)
	for {
		select {
		case latestIntPerSec <- latestPerSec:
		case <-ticker.C:
			newInt := count()
			if prevInt != 0 {
				latestPerSec = newInt - prevInt
			}
			prevInt = newInt
		}

	}
}

func Interrupts(checks *[]fthealth.Check) {

	go updateIntCount()

	impact := "System may become unresponsive"
	panicGuide := "Check the system with vmstat and investigate cause"

	check := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "Interrupts check",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "Number of interrupts as reported by /proc/stat is unusually high",
		Checker:          intCheck,
	}

	*checks = append(*checks, check)
}
