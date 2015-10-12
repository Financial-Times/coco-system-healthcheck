package main

import (
	fthealth "github.com/Financial-Times/go-fthealth"
)

func RegisterChecks(checks *[]fthealth.Check) {
	DiskFreeChecks(checks)
	MemInfo(checks)
	LoadAvg(checks)
	DiskInodes(checks)
	*checks = append(*checks, contextSwitchChecker{120000}.Checks()...)
	Interrupts(checks)
	Iops(checks)
}
