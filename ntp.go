package main

import (
	"errors"
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/bt51/ntpclient"
	"time"
)

type ntpChecker struct{}

type offsetResult struct {
	val time.Duration
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
		return offset.val.String(), offset.err
	}

	if offset.val > 1*time.Minute || offset.val < -1*time.Minute {
		return offset.val.String(), fmt.Errorf("offset is greater then limit of 1 minute: %s", offset.val.String())
	}
	return offset.val.String(), nil
}

func ntpLoop() {
	update := make(chan offsetResult)
	go func() {
		for {
			update <- ntpOffset()
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

func ntpOffset() offsetResult {
	t, err := ntpclient.GetNetworkTime("0.pool.ntp.org", 123)
	if err != nil {
		return offsetResult{0, fmt.Errorf("Could not get time form 0.pool.ntp.org")}
	}
	return offsetResult{time.Since(*t), nil}
}
