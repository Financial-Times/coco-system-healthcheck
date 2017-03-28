package main

import (
	"github.com/Financial-Times/service-status-go/gtg"
)

type gtgService struct {
	dfc  diskFreeChecker
	mc   memoryChecker
	lac  loadAverageChecker
	ntpc ntpChecker
	tcpc tcpChecker
}

func newGtgService(diskThresholdPercent, memoryThresholdPercent float64) *gtgService {
	return &gtgService{
		dfc:  diskFreeCheckerImpl{diskThresholdPercent},
		mc:   memoryCheckerImpl{memoryThresholdPercent},
		lac:  loadAverageCheckerImpl{},
		ntpc: ntpCheckerImpl{},
		tcpc: tcpCheckerImpl{},
	}
}

func (service *gtgService) Check() gtg.Status {
	if _, err := service.dfc.MountedDiskSpaceCheck(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	if _, err := service.dfc.RootDiskSpaceCheck(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	if _, err := service.mc.AvMemoryCheck(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	if _, err := service.lac.Check(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	if _, err := service.ntpc.Check(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	if _, err := service.tcpc.Check(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}

	return gtg.Status{GoodToGo: true}
}
