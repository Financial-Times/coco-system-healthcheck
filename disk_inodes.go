package main

import (
	"fmt"
	"os"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

type inodeChecker struct {
	threshold uint64
}

func (ic inodeChecker) Checks() []fthealth.Check {
	impact := "Filesystem may appear full. Services that require the filesystem may stop"
	panicGuide := "Check the filesystem with df -i <path> and investigate"

	rootCheck := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "Root disk inode check",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "Please free some inodes on the 'root' mount",
		Checker:          ic.rootInodesCheck,
	}

	mountedCheck := fthealth.Check{
		BusinessImpact:   impact,
		Name:             "Persistent disk inode check on '/vol' (always true for stateless nodes)",
		PanicGuide:       panicGuide,
		Severity:         2,
		TechnicalSummary: "Please clear some inodes on the 'vol' mount",
		Checker:          ic.mountedInodesCheck,
	}

	return []fthealth.Check{rootCheck, mountedCheck}
}

func (ic inodeChecker) inodeCheck(path string) (string, error) {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		return "", fmt.Errorf("Cannot read disk info of %s file system.", path)
	}
	if d.FreeInodes < ic.threshold {
		return fmt.Sprintf("%d", d.FreeInodes), fmt.Errorf("Lack of free inodes on %s : %d (< %d)", path, d.FreeInodes, ic.threshold)
	}
	return fmt.Sprintf("%d", d.FreeInodes), nil
}

func (ic inodeChecker) rootInodesCheck() (string, error) {
	return ic.inodeCheck(*hostPath + "/")
}

func (ic inodeChecker) mountedInodesCheck() (string, error) {
	path := *hostPath + "/vol"
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return "", nil
	}
	return ic.inodeCheck(path)
}
