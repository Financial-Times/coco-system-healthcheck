package main

import (
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"net"
)

type tcpChecker struct{}

func (tcpc tcpChecker) Checks() []fthealth.Check {
	check := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "TCP connection to port 8080 is available",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "HTTP connections to port 8080 will not be successful to any of the services deployed on this machine if this falls.",
		Checker:          tcpc.doCheck,
	}

	return []fthealth.Check{check}
}

func (tcpc tcpChecker) doCheck() (string, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return fmt.Sprintf("Connection error: [%v]", err), fmt.Errorf("Connecting to port 8080 was unsuccessful: [%v]", err)
	}
	defer conn.Close()
	return "Connected successfully.", nil
}
