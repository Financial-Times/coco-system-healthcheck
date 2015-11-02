package main

import (
	"errors"
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ntpChecker struct{}

type offsetResult struct {
	val float64
	err error
}

var offsetCh chan offsetResult

func (ntpc ntpChecker) Checks() []fthealth.Check {
	offsetCh = make(chan offsetResult)
	go ntpLoop()

	ntpCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "NTP sync check",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "System time has drifted out of sync of the box, investigate `timedatectl` and `systemd-timesyncd.service`",
		Checker:          ntpc.Check,
	}
	return []fthealth.Check{ntpCheck}
}

func (ntpc ntpChecker) Check() error {
	offset := <-offsetCh
	if offset.err != nil {
		return offset.err
	}
	if offset.val > 100 || offset.val < -100 {
		return fmt.Errorf("offset is greater then limit of 100: %f", offset.val)
	}
	return nil
}

func ntpLoop() {
	update := make(chan offsetResult)
	go func() {
		for {
			update <- ntpOffset(ntpCmd)
			time.Sleep(10 * time.Minute)
		}
	}()

	offset := offsetResult{err: errors.New("Ntp offset not initialised")}
	for {
		select {
		case offsetCh <- offset:
		case o := <-update:
			offset = o
		}
	}
}

func ntpCmd() string {
	// Throw away the error, it never returns 0, and we can catch an actual error later
	cmdOut, _ := exec.Command("ntpd", "-q", "-n", "-w", "-p", "pool.ntp.org").CombinedOutput()
	return string(cmdOut)
}

func ntpOffset(ntpCmd func() string) offsetResult {
	var offset string
	out := strings.Split(ntpCmd(), " ")
	for _, str := range out {
		if strings.Index(str, "offset") != -1 {
			offset = strings.Split(str, ":")[1]
			break
		}
	}
	if offset == "" {
		return offsetResult{0, fmt.Errorf("ntpd did not return an offset value")}
	}
	offsetFloat, err := strconv.ParseFloat(offset, 64)
	if err != nil {
		return offsetResult{0, fmt.Errorf("Could not parse offset: %v", err)}
	}

	return offsetResult{offsetFloat, nil}
}
