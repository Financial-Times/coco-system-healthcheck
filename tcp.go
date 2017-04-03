package main

import (
	"fmt"
	"net"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
)

type tcpChecker interface {
	Checks() []fthealth.Check
	Check() (string, error)
}

type tcpCheckerImpl struct{}

func (tcpc tcpCheckerImpl) Checks() []fthealth.Check {
	check := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "TCP connection to port 8080 is available",
		PanicGuide:       "https://dewey.ft.com/upp-system-healthcheck.html",
		Severity:         2,
		TechnicalSummary: "HTTP connections to port 8080 will not be successful to any of the services deployed on this machine if this falls.",
		Checker:          tcpc.Check,
	}

	return []fthealth.Check{check}
}

func (tcpc tcpCheckerImpl) Check() (string, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return fmt.Sprintf("Connection error: [%v]", err), fmt.Errorf("Connecting to port 8080 was unsuccessful: [%v]", err)
	}
	defer conn.Close()
	return "Connected successfully.", nil
}
