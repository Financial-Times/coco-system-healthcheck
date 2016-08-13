package main

import (
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"io/ioutil"
	"strconv"
	"strings"
)

type fdChecker struct{}

type fileDescriptors struct {
	current float64
	max     float64
}

func (fdc fdChecker) Checks() []fthealth.Check {
	impact := "Running out of file descriptors will prevent most services from functioning"
	panicGuide := "Restarted offending procs, limit services that are using up a lot of files."

	fdCheck := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "File descriptors check",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "To give an indication of which proc is using up fds, us this command: `ls /proc/*/fd 2>/dev/null`",
		Checker:          fdc.Check,
	}

	return []fthealth.Check{fdCheck}
}

func readDisk(path string) (fileDescriptors, error) {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		return fileDescriptors{}, err
	}
	ar := strings.Split(string(r), "\t")
	ar[2] = strings.Trim(ar[2], " ")
	ar[2] = strings.Replace(ar[2], "\n", "", -1)
	cur, _ := strconv.ParseFloat(ar[0], 64)
	max, _ := strconv.ParseFloat(ar[2], 64)

	return fileDescriptors{cur, max}, nil
}

func (fdc fdChecker) Check() (string, error) {
	fd, err := readDisk("/proc/sys/fs/file-nr")
	if err != nil {
		return "", fmt.Errorf("Cannot read proc for file-nr.")
	}

	currentPercent := fd.current / fd.max * 100

	if currentPercent > 80 {
		return fmt.Sprintf("%f", currentPercent), fmt.Errorf("Lack of file descriptors: %.2f", currentPercent)
	}
	return fmt.Sprintf("%.2f%%", currentPercent), nil
}
