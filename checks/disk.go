package checks

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	"syscall"
)

var stat syscall.Statfs_t

func spaceAvailablePercent(path string) float64 {
	syscall.Statfs(path, &stat)
	avail := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)
	return (float64(avail) / float64(total) * 100)
}

func diskSpaceCheck() error {
	spaceAv := spaceAvailablePercent("/")
	if spaceAv < 20 {
		return fmt.Errorf("No space on root")
	}
	return nil
}

func bootDiskSpaceCheck() error {
	spaceAv := spaceAvailablePercent("/boot")
	if spaceAv < 20 {
		return fmt.Errorf("No space on root")
	}
	return nil
}

func DiskChecks(checks *[]fthealth.Check) {
	spaceAv := spaceAvailablePercent("/")
	// health checks
	rootDiskSpaceCheck := fthealth.Check{
		BusinessImpact:   "No newspaper",
		Name:             fmt.Sprintf("Disk Space \"/\" %.2f%%", spaceAv),
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "rm -rf some shit",
		Checker:          diskSpaceCheck,
	}

	spaceAv = spaceAvailablePercent("/boot")
	bootDiskSpaceCheck := fthealth.Check{
		BusinessImpact:   "No newspaper",
		Name:             fmt.Sprintf("Disk Space \"/boot\" %.2f%%", spaceAv),
		PanicGuide:       "Keep calm and carry on",
		Severity:         2,
		TechnicalSummary: "rm -rf some shit",
		Checker:          bootDiskSpaceCheck,
	}

	*checks = append(*checks, rootDiskSpaceCheck)
	*checks = append(*checks, bootDiskSpaceCheck)
}
