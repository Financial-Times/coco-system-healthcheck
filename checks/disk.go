package checks

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func spaceAvailablePercent(disk *linuxproc.Disk) float64 {
	return (float64(disk.Free) / float64(disk.All) * 100)
}

func diskSpaceCheck(path string) error {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		return fmt.Errorf("Cannot read disk info of volume mounted under %s.", path)
	}
	if spaceAvailablePercent(d) < 20 {
		return fmt.Errorf("No space on root")
	}
	return nil
}

func rootDiskSpaceCheck() error {
	return diskSpaceCheck("/")
}

func bootDiskSpaceCheck() error {
	return diskSpaceCheck("/boot")
}

func DiskChecks(checks *[]fthealth.Check) {
	rootDiskSpaceCheck := fthealth.Check{
		BusinessImpact:   "No newspaper",
		Name:             "Root disk space check.",
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "rm -rf some shit",
		Checker:          rootDiskSpaceCheck,
	}

	bootDiskSpaceCheck := fthealth.Check{
		BusinessImpact:   "No newspaper",
		Name:             "Boot disk space check",
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "rm -rf some shit",
		Checker:          bootDiskSpaceCheck,
	}

	*checks = append(*checks, rootDiskSpaceCheck)
	*checks = append(*checks, bootDiskSpaceCheck)
}
