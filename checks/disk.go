package checks

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"log"
)

func spaceAvailablePercent(disk *linuxproc.Disk) float64 {
	return (float64(disk.Free) / float64(disk.All) * 100)
}

func diskSpaceCheck(path string) error {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		log.Printf("Cannot read disk info of %s file system.", path)
		return nil
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
	return diskSpaceCheck(baseDir + "/vol")
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
