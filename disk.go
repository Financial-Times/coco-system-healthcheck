package main

import (
	"fmt"
	"os"

	"errors"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"io/ioutil"
)

type diskFreeChecker interface {
	Checks() []fthealth.Check
	RootDiskSpaceCheck() (string, error)
	MountedDiskSpaceCheck() (string, error)
}

type diskFreeCheckerImpl struct {
	rootThresholdPercent   int
	mountsThresholdPercent int
}

func (dff diskFreeCheckerImpl) Checks() []fthealth.Check {
	rootCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "Root disk space check",
		PanicGuide:       "https://dewey.ft.com/upp-system-healthcheck.html",
		Severity:         2,
		TechnicalSummary: fmt.Sprintf("Free space on root volume is under %d%%", dff.rootThresholdPercent),
		Checker:          dff.RootDiskSpaceCheck,
	}

	mountedCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing and delivery workflow might be effected",
		Name:             "AWS EBS volumes mounted under '" + *awsEbsMountPath + "'",
		PanicGuide:       "https://dewey.ft.com/upp-system-healthcheck.html",
		Severity:         2,
		TechnicalSummary: fmt.Sprintf("Free space on mounted volumes under '%s' is under %d%%", *awsEbsMountPath, dff.mountsThresholdPercent),
		Checker:          dff.MountedDiskSpaceCheck,
	}

	return []fthealth.Check{rootCheck, mountedCheck}
}

func (dff diskFreeCheckerImpl) diskSpaceCheck(path string, thresholdPercent int) (string, error) {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		return "", fmt.Errorf("Cannot read disk info of %s file system.", path)
	}
	pctAvail := float64(d.Free) / float64(d.All) * 100
	if pctAvail < float64(thresholdPercent) {
		return fmt.Sprintf("Free space on %s: %2.1f%%", path, pctAvail), fmt.Errorf("Low free space on %s. Free disk space: %2.1f%%", path, pctAvail)
	}
	return fmt.Sprintf("Free space on %s: %2.1f%%", path, pctAvail), nil
}

func (dff diskFreeCheckerImpl) RootDiskSpaceCheck() (string, error) {
	return dff.diskSpaceCheck(*hostPath+"/", dff.rootThresholdPercent)
}

func (dff diskFreeCheckerImpl) MountedDiskSpaceCheck() (string, error) {
	path := *hostPath + *awsEbsMountPath

	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return "", nil
	}

	// look for the zone folder in the path
	zonesDir, err := ioutil.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("Cannot list zones folder in path '%s'. Cause: %s ", path, err.Error())
	}

	if len(zonesDir) == 0 { // if no zone folder is present then return ok
		return "", nil
	}

	zoneDirPath := path + "/" + zonesDir[0].Name()
	vols, err := ioutil.ReadDir(zoneDirPath)
	if err != nil {
		return "", fmt.Errorf("Cannot list volumes in path '%s'. Cause: %s ", zoneDirPath, err.Error())
	}

	aggStatus := ""
	aggError := ""
	for _, vol := range vols {
		volStatus, err := dff.diskSpaceCheck(zoneDirPath+"/"+vol.Name(), dff.mountsThresholdPercent)
		if err != nil {
			aggError = aggError + err.Error() + "; "
		} else {
			aggStatus = aggStatus + volStatus + "; "
		}
	}

	if len(aggError) == 0 {
		return aggStatus, nil
	} else {
		return aggStatus, errors.New(aggError)
	}
}
