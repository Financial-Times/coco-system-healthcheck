package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"log"
	"time"
)

func iops(name string) (uint64, error) {
	stats, err := linuxproc.ReadDiskStats(*hostPath + "/proc/diskstats")
	if err != nil {
		log.Fatalf("Cannot read disk stat info for %v\n", err)
	}
	for _, stat := range stats {
		if stat.Name == name {
			return stat.WriteIOs + stat.ReadIOs, nil
		}
		println(stat.Name)
	}
	return 0, fmt.Errorf("disk not found %v", name)
}

func iopsCheck() error {
	perSec := <-latestPerSec
	threshold := uint64(100)
	if perSec > threshold {
		return fmt.Errorf("%d iops per second. (>%d)", perSec, threshold)
	}
	return nil
}

var latestPerSec chan uint64 = make(chan uint64)

func updateIopsCount() {
	ticker := time.NewTicker(1 * time.Second)
	latest := uint64(0)
	prevInt := uint64(0)
	for {
		select {
		case latestPerSec <- latest:
		case <-ticker.C:
			newInt, err := iops("sda")
			if err != nil {
				log.Print("failed to read IOPS : %v\n", err.Error())
				continue
			}
			if prevInt != 0 {
				latest = newInt - prevInt
			}
			prevInt = newInt
		}

	}
}

func Iops(checks *[]fthealth.Check) {

	go updateIopsCount()

	impact := "System may become unresponsive"
	panicGuide := "Check the system with iostat and investigate cause"

	check := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "Iops check",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "Number of iops as reported by /proc/diskstat is unusually high",
		Checker:          iopsCheck,
	}

	*checks = append(*checks, check)
}
