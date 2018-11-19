package main

import (
	"github.com/Financial-Times/service-status-go/gtg"
)

type gtgService struct {
	dfc  diskFreeChecker
	mc   memoryChecker
	lac  loadAverageChecker
	ntpc ntpChecker
	apiSc apiServerChecker
}

func newGtgService(rootDiskThresholdPercent int, mountsThresholdPercent int, memoryThresholdPercent float64, apiServerURL string) *gtgService {
	return &gtgService{
		dfc:  diskFreeCheckerImpl{rootDiskThresholdPercent, mountsThresholdPercent},
		mc:   memoryCheckerImpl{memoryThresholdPercent},
		lac:  loadAverageCheckerImpl{},
		ntpc: &ntpCheckerImpl{},
		apiSc: &apiServerCheckerImpl{apiServerURL},
	}
}

func (service *gtgService) Check() gtg.Status {
	mountedDiskSpaceCheck := func() gtg.Status {
		return gtgCheck(service.mountedDiskSpaceChecker)
	}
	rootDiskSpaceCheck := func() gtg.Status {
		return gtgCheck(service.rootDiskSpaceChecker)
	}
	memoryUsageCheck := func() gtg.Status {
		return gtgCheck(service.memoryUsageChecker)
	}
	loadAverageCheck := func() gtg.Status {
		return gtgCheck(service.loadAverageChecker)
	}
	clockSyncCheck := func() gtg.Status {
		return gtgCheck(service.clockSyncChecker)
	}

	apiServerCertCheck := func() gtg.Status {
		return gtgCheck(service.apiServerCertChecker)
	}

	return gtg.FailFastParallelCheck(
		[]gtg.StatusChecker{mountedDiskSpaceCheck, rootDiskSpaceCheck, memoryUsageCheck, loadAverageCheck, clockSyncCheck, apiServerCertCheck})()
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
}

func (service *gtgService) mountedDiskSpaceChecker() (string, error) {
	if _, err := service.dfc.MountedDiskSpaceCheck(); err != nil {
		return err.Error(), err
	}
	return "Mounted disk space check OK.", nil
}

func (service *gtgService) rootDiskSpaceChecker() (string, error) {
	if _, err := service.dfc.RootDiskSpaceCheck(); err != nil {
		return err.Error(), err
	}
	return "Root disk space check OK.", nil
}

func (service *gtgService) memoryUsageChecker() (string, error) {
	if _, err := service.mc.AvMemoryCheck(); err != nil {
		return err.Error(), err
	}
	return "Memory usage check OK.", nil
}

func (service *gtgService) loadAverageChecker() (string, error) {
	if _, err := service.lac.Check(); err != nil {
		return err.Error(), err
	}
	return "Load Average check OK.", nil
}

func (service *gtgService) clockSyncChecker() (string, error) {
	if _, err := service.ntpc.Check(); err != nil {
		return err.Error(), err
	}
	return "Clock sync check OK.", nil
}

func (service *gtgService) apiServerCertChecker() (string, error) {
	if _, err := service.apiSc.CheckCertificate(); err != nil {
		return err.Error(), err
	}
	return "Api server certificate check OK.", nil
}
