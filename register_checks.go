package main

import (
	fthealth "github.com/Financial-Times/go-fthealth"
)

func RegisterChecks(checks *[]fthealth.Check) {
	DiskFreeChecks(checks)
	MemInfo(checks)
	*checks = append(*checks, loadAverageChecker{}.Checks()...)
	*checks = append(*checks, inodeChecker{1024}.Checks()...)
	*checks = append(*checks, contextSwitchChecker{120000}.Checks()...)
	*checks = append(*checks, interruptsChecker{3000}.Checks()...)
	*checks = append(*checks, iopsChecker{100}.Checks()...)
}
