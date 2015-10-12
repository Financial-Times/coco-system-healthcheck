package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"os"
)

func inodeCheck(path string) error {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		return fmt.Errorf("Cannot read disk info of %s file system.", path)
	}
	if d.FreeInodes < 1024 {
		return fmt.Errorf("Lack of free inodes on %s : %d (< 1024)", path, d.FreeInodes)
	}
	return nil
}

func rootInodesCheck() error {
	return inodeCheck(*hostPath + "/")
}

func mountedInodesCheck() error {
	path := *hostPath + "/vol"
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil
	}
	return inodeCheck(path)
}

func DiskInodes(checks *[]fthealth.Check) {
	impact := "Filesystem may appear full. Services that require the filesystem may stop"
	panicGuide := "Check the filesystem with df -i <path> and investigate"

	rootCheck := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "Root disk inode check",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "Please free some inodes on the 'root' mount",
		Checker:          rootInodesCheck,
	}

	mountedCheck := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "Persistent disk inode check on '/vol' (always true for stateless nodes)",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "Please clear some inodes on the 'vol' mount",
		Checker:          mountedInodesCheck,
	}

	*checks = append(*checks, rootCheck)
	*checks = append(*checks, mountedCheck)
}
