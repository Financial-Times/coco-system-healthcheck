package main

import (
	"errors"
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ntpChecker struct{}

type offsetResult struct {
	val string
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

func (ntpc ntpChecker) Check() (string, error) {
	offset := <-offsetCh
	if offset.err != nil {
		return offset.val, offset.err
	}

	offsetFloat, err := strconv.ParseFloat(offset.val, 64)
	if err != nil {
		return offset.val, fmt.Errorf("Could not parse offset: %v", err)
	}

	if offsetFloat > 100 || offsetFloat < -100 {
		return offset.val, fmt.Errorf("offset is greater then limit of 100: %f", offsetFloat)
	}
	return offset.val, nil
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
		log.Printf("Didn't get an offset, ntpd out: %v", out)
		return offsetResult{"", fmt.Errorf("ntpd did not return an offset value")}
	}

	return offsetResult{offset, nil}
}
