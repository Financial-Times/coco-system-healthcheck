package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/bt51/ntpclient"
	"github.com/kr/pretty"
)

type result struct {
	val string
	err error
}

var earliestTime = time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)

var pools = []string{
	"169.254.169.123", // AWS NTP server
	"0.coreos.pool.ntp.org",
	"1.coreos.pool.ntp.org",
	"2.coreos.pool.ntp.org",
	"3.coreos.pool.ntp.org",
	"0.pool.ntp.org",
	"1.pool.ntp.org",
	"2.pool.ntp.org",
	"3.pool.ntp.org",
}

type ntpChecker interface {
	Checks() []fthealth.Check
	Check() (string, error)
}

type ntpCheckerImpl struct {
	sync.RWMutex
	timeDrift     time.Duration
	pollingPeriod time.Duration
	result        result
}

func (ntpc *ntpCheckerImpl) Checks() []fthealth.Check {
	ntpc.result = result{
		err: errors.New("No value yet"),
	}

	go func() {
		for {
			ntpc.Lock()
			ntpTime, err := ntpc.callNtpWithPools(pools)
			ntpc.result = result{
				val: ntpTime.String(),
				err: err,
			}
			ntpc.Unlock()
			time.Sleep(ntpc.pollingPeriod)
		}
	}()

	ntpCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "NTP sync check",
		PanicGuide:       "https://dewey.ft.com/upp-system-healthcheck.html",
		Severity:         2,
		TechnicalSummary: "System time has drifted out of sync of the box, investigate `timedatectl` and `systemd-timesyncd.service`",
		Checker:          ntpc.Check,
	}
	return []fthealth.Check{ntpCheck}
}

func (ntpc *ntpCheckerImpl) Check() (string, error) {
	ntpc.RLock()
	defer ntpc.RUnlock()
	return ntpc.result.val, ntpc.result.err
}

func (ntpc *ntpCheckerImpl) callNtpWithPools(pools []string) (*time.Time, error) {
	var errs []error
	for _, pool := range pools {
		ntpTime, err := ntpclient.GetNetworkTime(pool, 123)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if ntpTime.Before(earliestTime) {
			errs = append(errs, fmt.Errorf("Time from pool %s was way too old %s", pool, ntpTime.String()))
			continue
		}
		drift := time.Since(*ntpTime)
		if drift > ntpc.timeDrift || drift < -ntpc.timeDrift {
			return nil, fmt.Errorf("offset is greater then limit of %s: %s", ntpc.timeDrift.String(), drift)
		}
		return ntpTime, nil
	}
	return nil, fmt.Errorf("None of the pools are reachable: %# v", pretty.Formatter(errs))
}
