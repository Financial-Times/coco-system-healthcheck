package main

import (
	fthealth "github.com/Financial-Times/go-fthealth"
)

var baseDir string

func RegisterChecks(path string, checks *[]fthealth.Check) {
	baseDir = path
	DiskFreeChecks(checks)
	MemInfo(checks)
	LoadAvg(checks)
	DiskInodes(checks)
}
