package main

import (
	"fmt"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/bt51/ntpclient"
)

var offsetCh chan result
var pools = [4]string{"0.pool.ntp.org", "1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org"}

type ntpChecker struct{}

func (ntpc ntpChecker) Checks() []fthealth.Check {
	offsetCh = make(chan result)
	go loop(ntpOffset, 60, offsetCh)

	ntpCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "NTP sync check",
		PanicGuide:       "Please refer to the technical summary section below",
		Severity:         2,
		TechnicalSummary: "System time has drifted out of sync of the box, investigate `timedatectl` and `systemd-timesyncd.service`",
		Checker:          ntpc.Check,
	}
	return []fthealth.Check{ntpCheck}
}

func (ntpc ntpChecker) Check() (string, error) {
	offset := <-offsetCh
	return offset.val, offset.err
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

func ntpOffset() result {
	t, err := callNtp()
	if err != nil {
		return result{err: fmt.Errorf("Could not get time %v", err)}
	}
	tsn := time.Since(*t)
	if tsn > 2*time.Second || tsn < -2*time.Second {
		return result{err: fmt.Errorf("offset is greater then limit of 2 seconds: %s", tsn)}
	}

	return result{val: tsn.String()}
}
