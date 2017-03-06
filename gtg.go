package main

import "net/http"

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

func (service *gtgService) Check(writer http.ResponseWriter, req *http.Request) {
	if _, err := service.dfc.MountedDiskSpaceCheck(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if _, err := service.dfc.RootDiskSpaceCheck(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if _, err := service.mc.AvMemoryCheck(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if _, err := service.lac.Check(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if _, err := service.ntpc.Check(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if _, err := service.tcpc.Check(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}
