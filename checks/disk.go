package checks

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"os"
)

func spaceAvailablePercent(disk *linuxproc.Disk) float64 {
	return (float64(disk.Free) / float64(disk.All) * 100)
}

func diskSpaceCheck(path string) error {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		return fmt.Errorf("Cannot read disk info of %s file system.", path)
	}
	if spaceAvailablePercent(d) < 20 {
		return fmt.Errorf("Low free space on %s. Free disk space: %2.1f %%", path, spaceAvailablePercent)
	}
	return nil
}

func rootDiskSpaceCheck() error {
	return diskSpaceCheck(baseDir + "/")
}

func mountedDiskSpaceCheck() error {
	path := baseDir + "/vol"
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil
	}
	return diskSpaceCheck(path)
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

	mountedDiskSpaceCheck := fthealth.Check{
		BusinessImpact:   "No newspaper",
		Name:             "Mounted disk space check (/vol)",
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "rm -rf some shit",
		Checker:          mountedDiskSpaceCheck,
	}

	*checks = append(*checks, rootDiskSpaceCheck)
	*checks = append(*checks, mountedDiskSpaceCheck)
}
