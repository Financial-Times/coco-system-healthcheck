package main

import (
	"errors"
	"testing"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckHappyFlow(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedTcpc := &mockedTcpChecker{}

	gtg := gtgService{
		dfc:  mockedDfc,
		lac:  mockedLac,
		mc:   mockedMc,
		ntpc: mockedNtpc,
		tcpc: mockedTcpc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedTcpc.On("Check").Return("", nil)

	status := gtg.Check()

	assert.Equal(t, true, status.GoodToGo)
}

func TestCheckInsufficientRootDiskSpace(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedTcpc := &mockedTcpChecker{}

	gtg := gtgService{
		dfc:  mockedDfc,
		lac:  mockedLac,
		mc:   mockedMc,
		ntpc: mockedNtpc,
		tcpc: mockedTcpc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", errors.New("Insuficient root disk space"))
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedTcpc.On("Check").Return("", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckInsufficientMountedDiskSpace(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedTcpc := &mockedTcpChecker{}

	gtg := gtgService{
		dfc:  mockedDfc,
		lac:  mockedLac,
		mc:   mockedMc,
		ntpc: mockedNtpc,
		tcpc: mockedTcpc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", errors.New("Insuficient mounted disk space"))
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedTcpc.On("Check").Return("", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckHighAverageCPULoad(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedTcpc := &mockedTcpChecker{}

	gtg := gtgService{
		dfc:  mockedDfc,
		lac:  mockedLac,
		mc:   mockedMc,
		ntpc: mockedNtpc,
		tcpc: mockedTcpc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", errors.New("The average load is above recommended average"))
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedTcpc.On("Check").Return("", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckHighAverageMemoryLoad(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedTcpc := &mockedTcpChecker{}

	gtg := gtgService{
		dfc:  mockedDfc,
		lac:  mockedLac,
		mc:   mockedMc,
		ntpc: mockedNtpc,
		tcpc: mockedTcpc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", errors.New("The average memory load is above recommended average"))
	mockedNtpc.On("Check").Return("", nil)
	mockedTcpc.On("Check").Return("", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckNtpOutOfSync(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedTcpc := &mockedTcpChecker{}

	gtg := gtgService{
		dfc:  mockedDfc,
		lac:  mockedLac,
		mc:   mockedMc,
		ntpc: mockedNtpc,
		tcpc: mockedTcpc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", errors.New("The ntp is out of sync"))
	mockedTcpc.On("Check").Return("", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckUnsuccessfulTcpConnection(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedTcpc := &mockedTcpChecker{}

	gtg := gtgService{
		dfc:  mockedDfc,
		lac:  mockedLac,
		mc:   mockedMc,
		ntpc: mockedNtpc,
		tcpc: mockedTcpc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedTcpc.On("Check").Return("", errors.New("Unsuccessful connection to TCP port"))

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

type mockedDiskFreeChecker struct {
	mock.Mock
}

func (m *mockedDiskFreeChecker) Checks() []fthealth.Check {
	return nil
}

func (m *mockedDiskFreeChecker) RootDiskSpaceCheck() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockedDiskFreeChecker) MountedDiskSpaceCheck() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type mockedLoadAverageChecker struct {
	mock.Mock
}

func (m *mockedLoadAverageChecker) Checks() []fthealth.Check {
	return nil
}

func (m *mockedLoadAverageChecker) Check() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type mockedMemoryChecker struct {
	mock.Mock
}

func (m *mockedMemoryChecker) Checks() []fthealth.Check {
	return nil
}

func (m *mockedMemoryChecker) AvMemoryCheck() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type mockedNtpChecker struct {
	mock.Mock
}

func (m *mockedNtpChecker) Checks() []fthealth.Check {
	return nil
}

func (m *mockedNtpChecker) Check() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type mockedTcpChecker struct {
	mock.Mock
}

func (m *mockedTcpChecker) Checks() []fthealth.Check {
	return nil
}

func (m *mockedTcpChecker) Check() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
