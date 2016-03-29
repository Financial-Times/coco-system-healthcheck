package main

import (
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"os"
	"strings"
)

type diskFreeChecker struct {
	thresholdPercent float64
}

func (dff diskFreeChecker) Checks() []fthealth.Check {
	mtabPath := *hostPath + "/etc/mtab"
	mounts, err := linuxproc.ReadMounts(mtabPath)
	if err != nil {
		return []fthealth.Check{fthealth.Check{
			BusinessImpact:   "A part of the publishing workflow might be affected",
			Name:             "Disk space check",
			PanicGuide:       "Please refer to technical summary",
			Severity:         2,
			TechnicalSummary: fmt.Sprintf("Please check that service can read %q mount", mtabPath),
			Checker:          func() (string, error) { return "", fmt.Errorf("Cannot read disk info of %s mtab. %s", mtabPath, err) },
		},
		}
	}

	var fthealthChecks = []fthealth.Check{}
	for _, mount := range mounts.Mounts {
		if strings.HasPrefix(mount.Device, "/dev") {
			fthealthChecks = append(fthealthChecks, fthealth.Check{
				BusinessImpact:   "A part of the publishing workflow might be affected",
				Name:             fmt.Sprintf("Persistent disk space check mounted on %q (always true for stateless nodes)", mount.MountPoint),
				PanicGuide:       "Please refer to technical summary",
				Severity:         2,
				TechnicalSummary: fmt.Sprintf("Please clear some disk space on the %q mount", mount.MountPoint),
				Checker: func() (string, error) {
					return dff.mountedDiskSpaceCheck(mount.MountPoint)
				},
			})
		}
	}

	return fthealthChecks
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

func (dff diskFreeChecker) mountedDiskSpaceCheck(mountPoint string) (string, error) {
	path := *hostPath + mountPoint
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return "", nil
	}
	return dff.diskSpaceCheck(path)
}
