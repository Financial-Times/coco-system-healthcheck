package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var reachablePools = []string{
	"0.coreos.pool.ntp.org",
	"1.coreos.pool.ntp.org",
	"2.coreos.pool.ntp.org",
	"3.coreos.pool.ntp.org",
	"0.pool.ntp.org",
	"2.pool.ntp.org",
}
var someReachablePools = []string{
	"foobar8237.ntp.org",
	"0.coreos.pool.ntp.org",
	"1.pool.ntp.org",
	"foobarlllx237.ntp.org",
	"3.pool.ntp.org",
}
var unreachablePools = []string{
	"foobar8237.ntp.org",
	"foobarlllx237.ntp.org",
}

func TestCallNtpWithPoolsAllPoolsReachable(t *testing.T) {
	ntpc := &ntpCheckerImpl{
		timeDrift: time.Second,
	}
	result, err := ntpc.callNtpWithPools(reachablePools)

	assert.Nil(t, err, "Error should be nil when all pools are reachable")
	assert.NotNil(t, result, "Result should not be nil when all pools are reachable")
}

func TestCallNtpWithPoolsSomePoolsReachable(t *testing.T) {
	ntpc := &ntpCheckerImpl{
		timeDrift: time.Second,
	}
	result, err := ntpc.callNtpWithPools(someReachablePools)

	assert.Nil(t, err, "Error should be nil when some pools are reachable")
	assert.NotNil(t, result, "Result should not be nil when some pools are reachable")
}

func TestCallNtpWithPoolsAllPoolsUnreachable(t *testing.T) {
	ntpc := &ntpCheckerImpl{
		timeDrift: time.Second,
	}
	result, err := ntpc.callNtpWithPools(unreachablePools)

	assert.NotNil(t, err, "Error should not be nil when all pools are unreachable")
	assert.Nil(t, result, "Result should be nil when all pools are unreachable")
}
