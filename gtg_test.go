package main

import (
	"errors"
	"testing"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckHappyFlow(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedApisc := &mockedAPIServerChecker{}

	gtg := gtgService{
		dfc:   mockedDfc,
		lac:   mockedLac,
		mc:    mockedMc,
		ntpc:  mockedNtpc,
		apiSc: mockedApisc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedApisc.On("CheckCertificate").Return("Exp date: 3.10.3000", nil)

	status := gtg.Check()

	assert.Equal(t, true, status.GoodToGo)
}

func TestCheckInsufficientRootDiskSpace(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedApisc := &mockedAPIServerChecker{}

	gtg := gtgService{
		dfc:   mockedDfc,
		lac:   mockedLac,
		mc:    mockedMc,
		ntpc:  mockedNtpc,
		apiSc: mockedApisc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", errors.New("Insuficient root disk space"))
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedApisc.On("CheckCertificate").Return("Exp date: 3.10.3000", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckInsufficientMountedDiskSpace(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedApisc := &mockedAPIServerChecker{}

	gtg := gtgService{
		dfc:   mockedDfc,
		lac:   mockedLac,
		mc:    mockedMc,
		ntpc:  mockedNtpc,
		apiSc: mockedApisc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", errors.New("Insuficient mounted disk space"))
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedApisc.On("CheckCertificate").Return("Exp date: 3.10.3000", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckHighAverageCPULoad(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedApisc := &mockedAPIServerChecker{}

	gtg := gtgService{
		dfc:   mockedDfc,
		lac:   mockedLac,
		mc:    mockedMc,
		ntpc:  mockedNtpc,
		apiSc: mockedApisc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", errors.New("The average load is above recommended average"))
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedApisc.On("CheckCertificate").Return("Exp date: 3.10.3000", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckHighAverageMemoryLoad(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedApisc := &mockedAPIServerChecker{}

	gtg := gtgService{
		dfc:   mockedDfc,
		lac:   mockedLac,
		mc:    mockedMc,
		ntpc:  mockedNtpc,
		apiSc: mockedApisc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", errors.New("The average memory load is above recommended average"))
	mockedNtpc.On("Check").Return("", nil)
	mockedApisc.On("CheckCertificate").Return("Exp date: 3.10.3000", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckNtpOutOfSync(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedApisc := &mockedAPIServerChecker{}

	gtg := gtgService{
		dfc:   mockedDfc,
		lac:   mockedLac,
		mc:    mockedMc,
		ntpc:  mockedNtpc,
		apiSc: mockedApisc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", errors.New("The ntp is out of sync"))
	mockedApisc.On("CheckCertificate").Return("Exp date: 3.10.3000", nil)

	status := gtg.Check()

	assert.Equal(t, false, status.GoodToGo)
}

func TestCheckApiServerNotValidCert(t *testing.T) {
	mockedDfc := &mockedDiskFreeChecker{}
	mockedLac := &mockedLoadAverageChecker{}
	mockedMc := &mockedMemoryChecker{}
	mockedNtpc := &mockedNtpChecker{}
	mockedApisc := &mockedAPIServerChecker{}

	gtg := gtgService{
		dfc:   mockedDfc,
		lac:   mockedLac,
		mc:    mockedMc,
		ntpc:  mockedNtpc,
		apiSc: mockedApisc,
	}

	mockedDfc.On("RootDiskSpaceCheck").Return("", nil)
	mockedDfc.On("MountedDiskSpaceCheck").Return("", nil)
	mockedLac.On("Check").Return("", nil)
	mockedMc.On("AvMemoryCheck").Return("", nil)
	mockedNtpc.On("Check").Return("", nil)
	mockedApisc.On("CheckCertificate").Return("Exp date: 3.10.2018", errors.New("the API server certificate expires in less than one month"))

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

type mockedAPIServerChecker struct {
	mock.Mock
}

func (m *mockedAPIServerChecker) Checks() []fthealth.Check {
	return nil
}

func (m *mockedAPIServerChecker) CheckCertificate() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
