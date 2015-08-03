package checks

import (
	fthealth "github.com/Financial-Times/go-fthealth"
)

var baseDir string

func RegisterChecks(path string, checks *[]fthealth.Check) {
	baseDir = path
	DiskChecks(checks)
	MemInfo(checks)
	LoadAvg(checks)
}
