package main

import (
	"bytes"
	"errors"
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type vethChecker struct{}

type vethResult struct {
	val string
	err error
}

var vethCh chan vethResult

func (vethc vethChecker) Checks() []fthealth.Check {
	vethCh = make(chan vethResult)
	go vethCollector()

	vethCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "veth interface check",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "Box is likely experiencing docker veth bug, please restart: `sudo locksmithctl restart`",
		Checker:          vethc.Check,
	}
	return []fthealth.Check{vethCheck}
}

func (vethc vethChecker) Check() (string, error) {
	veth := <-vethCh
	if veth.err != nil {
		return veth.val, veth.err
	}

	iveth, err := strconv.Atoi(veth.val)
	if err != nil {
		return "", err
	}

	if iveth > 100 {
		return veth.val, fmt.Errorf("veth is greater then limit of 100: %f", veth)
	}
	return veth.val, nil
}

func vethCollector() {
	update := make(chan vethResult)
	go func() {
		for {
			update <- vethCmd()
			time.Sleep(10 * time.Minute)
		}
	}()

	veth := vethResult{err: errors.New("Veth count not initialised")}
	for {
		select {
		case vethCh <- veth:
		case veth = <-update:
		}
	}
}

func vethCmd() vethResult {
	out := new(bytes.Buffer)

	ipC := exec.Command("ip", "a")
	grepC := exec.Command("grep", "veth")
	wcC := exec.Command("wc", "-l")
	grepC.Stdin, _ = ipC.StdoutPipe()
	wcC.Stdin, _ = grepC.StdoutPipe()
	wcC.Stdout = out

	_ = grepC.Start()
	_ = wcC.Start()
	_ = ipC.Run()
	_ = grepC.Wait()
	_ = wcC.Wait()

	// Strip newline
	strippedOut := strings.Replace(out.String(), "\n", "", -1)
	return vethResult{strippedOut, nil}
}
