package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"log"
	"time"
)

type iopsChecker struct {
	threshold uint64
}

func (ic iopsChecker) Checks() []fthealth.Check {

	go ic.updateIopsCount()

	check := fthealth.Check{
		BusinessImpact:   "System may become unresponsive",
		Name:             "Iops check",
		PanicGuide:       "Check the system with iostat and investigate cause",
		Severity:         2,
		TechnicalSummary: "Number of iops as reported by /proc/diskstat is unusually high",
		Checker:          ic.iopsCheck,
	}

	return []fthealth.Check{check}
}

func (ic iopsChecker) iops(name string) (uint64, error) {
	stats, err := linuxproc.ReadDiskStats(*hostPath + "/proc/diskstats")
	if err != nil {
		log.Fatalf("Cannot read disk stat info for %v\n", err)
	}
	for _, stat := range stats {
		if stat.Name == name {
			return stat.WriteIOs + stat.ReadIOs, nil
		}
	}
	return 0, fmt.Errorf("disk not found %v", name)
}

func (ic iopsChecker) iopsCheck() error {
	perSec := <-latestPerSec
	if perSec > ic.threshold {
		return fmt.Errorf("%d iops per second. (>%d)", perSec, ic.threshold)
	}
	return nil
}

var latestPerSec chan uint64 = make(chan uint64)

func (ic iopsChecker) updateIopsCount() {
	ticker := time.NewTicker(1 * time.Second)
	latest := uint64(0)
	prevInt := uint64(0)
	for {
		select {
		case latestPerSec <- latest:
		case <-ticker.C:
			newInt, err := ic.iops("xvda")
			if err != nil {
				log.Printf("failed to read IOPS : %v\n", err.Error())
				continue
			}
			if prevInt != 0 {
				latest = newInt - prevInt
			}
			prevInt = newInt
		}

	}
}
