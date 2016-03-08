package main

import (
	"errors"
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/bt51/ntpclient"
	"time"
)

var offsetCh chan offsetResult
var pools = [4]string{"0.pool.ntp.org", "1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org"}

type ntpChecker struct{}

type offsetResult struct {
	val time.Duration
	err error
}

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

	if offset.val > 2*time.Second || offset.val < -2*time.Second {
		return offset.val.String(), fmt.Errorf("offset is greater then limit of 2 seconds: %s", offset.val.String())
	}
	return offset.val.String(), nil
}

func ntpLoop() {
	update := make(chan offsetResult)
	go func() {
		for {
			update <- ntpOffset()
			time.Sleep(1 * time.Minute)
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

func callNtp() (*time.Time, error) {
	var err error
	for _, pool := range pools {
		t, err := ntpclient.GetNetworkTime(pool, 123)
		if err == nil {
			return t, err
		}
	}
	return nil, err
}

func ntpOffset() offsetResult {
	t, err := callNtp()
	if err != nil {
		return offsetResult{0, fmt.Errorf("Could not get time %v", err)}
	}
	return offsetResult{time.Since(*t), nil}
}
