package main

import (
	fthealth "github.com/Financial-Times/go-fthealth"
)

func RegisterChecks(checks *[]fthealth.Check) {
	DiskFreeChecks(checks)
	MemInfo(checks)
	LoadAvg(checks)
	DiskInodes(checks)
	ContextSwitches(checks)
	Interrupts(checks)
	Iops(checks)
}
