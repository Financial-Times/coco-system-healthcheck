package main

import (
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"os"
)

type diskFreeChecker struct {
	thresholdPercent float64
}

func (dff diskFreeChecker) Checks() []fthealth.Check {
	rootCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "Root disk space check",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "Please clear some disk space on the 'root' mount",
		Checker:          dff.rootDiskSpaceCheck,
	}

	mountedCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be effected",
		Name:             "Persistent disk space check mounted on '/vol' (always true for stateless nodes)",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "Please clear some disk space on the 'vol' mount",
		Checker:          dff.mountedDiskSpaceCheck,
	}

	return []fthealth.Check{rootCheck, mountedCheck}
}

func (dff diskFreeChecker) diskSpaceCheck(path string) (string, error) {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		return "", fmt.Errorf("Cannot read disk info of %s file system.", path)
	}
	pctAvail := (float64(d.Free) / float64(d.All) * 100)
	if pctAvail < dff.thresholdPercent {
		return fmt.Sprintf("%2.1f%%", pctAvail), fmt.Errorf("Low free space on %s. Free disk space: %2.1f%%", path, pctAvail)
	}
	return fmt.Sprintf("%2.1f%%", pctAvail), nil
}

func (dff diskFreeChecker) rootDiskSpaceCheck() (string, error) {
	return dff.diskSpaceCheck(*hostPath + "/")
}

func (dff diskFreeChecker) mountedDiskSpaceCheck() (string, error) {
	path := *hostPath + "/vol"
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return "", nil
	}
	return dff.diskSpaceCheck(path)
}
