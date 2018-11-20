package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
)

type apiServerChecker interface {
	Checks() []fthealth.Check
	CheckCertificate() (string, error)
}

type apiServerCheckerImpl struct {
	url string
}

func (checker *apiServerCheckerImpl) Checks() []fthealth.Check {
	check := fthealth.Check{
		BusinessImpact:   "Kubernetes API server will become unavailable",
		Name:             "Kubernetes API server certificate check",
		PanicGuide:       "https://github.com/Financial-Times/content-k8s-provisioner#rotating-the-tls-assets-for-a-cluster",
		Severity:         2,
		TechnicalSummary: "Rotate the TLS assets on the Kubernetes cluster before the API server certificate expires",
		Checker:          checker.CheckCertificate,
	}

	return []fthealth.Check{check}
}

func (checker *apiServerCheckerImpl) CheckCertificate() (string, error) {
	expInOneMonth, expDate, err := checker.certificateExpiresInOneMonth()
	if err != nil {
		return "", fmt.Errorf("could not get the API server certificate. Reason: %s", err.Error())
	}

	expMsg := fmt.Sprintf("Expiry date: %s", expDate.Format(time.RFC822))
	if expInOneMonth {
		return expMsg, fmt.Errorf("the API server certificate expires in less than one month")
	}
	return expMsg, nil
}

func (checker *apiServerCheckerImpl) certificateExpiresInOneMonth() (bool, *time.Time, error) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConnsPerHost: 100,
		},
	}

	req, err := http.NewRequest("HEAD", checker.url, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, nil, err
	}

	if resp.StatusCode == 200 || resp.StatusCode == 401 {
		expiryDate := resp.TLS.PeerCertificates[0].NotAfter

		oneMonthFromNow := time.Now().AddDate(0, 1, 0)
		return expiryDate.Before(oneMonthFromNow), &expiryDate, nil
	}

	return false, nil, fmt.Errorf("the api server returned status %d", resp.StatusCode)
}
