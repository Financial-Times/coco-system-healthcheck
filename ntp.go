package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	"os/exec"
	"strings"
)

type ntpChecker struct{}

func (ntpc ntpChecker) Checks() []fthealth.Check {
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
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command("/usr/bin/timedatectl").Output(); err != nil {
		return fmt.Errorf("Could not run `timedatectl`")
	}
	out := string(cmdOut)
	anchor := "NTP synchronized: "
	index := strings.Index(out, anchor)
	if index == -1 {
		return fmt.Errorf("Failed to get `NTP synchronized` from `timedatectl` command")
	}
	answer := out[index+len(anchor) : index+len(anchor)+3]
	if answer != "yes" {
		return fmt.Errorf("NTP is not synchronized")
	}
	return nil
}
