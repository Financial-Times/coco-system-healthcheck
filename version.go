package main

import (
	"errors"
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	versionUri string = "http://%s.release.core-os.net/amd64-usr/current/version.txt"
	channelF   string = "/etc/coreos/update.conf"
	releaseF   string = "/usr/share/coreos/release"
)

var resultCh chan result

type versionChecker struct{}
type result struct {
	val string
	err error
}

func (v versionChecker) Checks() []fthealth.Check {
	resultCh = make(chan result)
	go loop()
	check := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be affected",
		Name:             "CoreOS version",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "CoreOS version is out of date or cannot be determined",
		Checker:          v.Check,
	}
	return []fthealth.Check{check}
}

func (v versionChecker) Check() (string, error) {
	result := <-resultCh
	return result.val, result.err
}

func loop() {
	updateCh := make(chan result)
	go func() {
		for {
			updateCh <- latest()
			time.Sleep(5 * time.Minute)
		}
	}()

	result := result{err: errors.New("No version yet")}
	for {
		select {
		case resultCh <- result:
		case result = <-updateCh:
		}
	}
}

func latest() result {
	channel, err := valFromFile("GROUP=", channelF)
	if err != nil {
		return result{err: err}
	}
	release, err := valFromFile("COREOS_RELEASE_VERSION=", releaseF)
	if err != nil {
		return result{err: err}
	}
	rmtRel, err := rmtRel(channel)
	if err != nil {
		return result{err: err}
	}
	if release != rmtRel {
		return result{
			val: fmt.Sprintf("Local release %v different from remote %v", release, rmtRel),
			err: errors.New(fmt.Sprintf("Local release %v different from remote %v", release, rmtRel)),
		}
	}
	return result{val: fmt.Sprintf("Current release %v is latest on %v", release, channel)}
}

func valFromFile(key, path string) (val string, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, key) {
			v := strings.TrimPrefix(line, key)
			return v, nil
		}
	}
	return "", errors.New(fmt.Sprintf("No %s in %s", val, path))
}

func rmtRel(channel string) (release string, err error) {
	uri := fmt.Sprintf(versionUri, channel)
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Got %v requesting %v", resp.StatusCode, uri))
	}
	body, err := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "COREOS_VERSION=") {
			release = strings.TrimPrefix(line, "COREOS_VERSION=")
			return release, nil
		}
	}
	return "", errors.New("No CoreOS version on the page")
}
