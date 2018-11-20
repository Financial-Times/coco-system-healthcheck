package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckCertificateWrongStatus(t *testing.T) {
	expectedStatus := 500
	ts := httptest.NewTLSServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(expectedStatus)
		res.Write([]byte("dummy"))
	}))
	defer ts.Close()

	sc := apiServerCheckerImpl{ts.URL}
	msg, err := sc.CheckCertificate()
	assert.EqualError(t, err, fmt.Sprintf("could not get the API server certificate. Reason: the api server returned status %d", expectedStatus),
		"An error should be thrown when the API server responds with other status than 200")
	assert.Equal(t, "", msg)
}

func TestCheckCertificateServerDown(t *testing.T) {
	notExistingServer := "https://dummyserver44ffdadf.not.existing.ddd"
	sc := apiServerCheckerImpl{notExistingServer}

	msg, err := sc.CheckCertificate()
	assert.NotNil(t, err, "Error should be thrown on not existing server")
	assert.Equal(t, "", msg)
}

func TestCheckCertificateValidCert(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("dummy"))
	}))
	defer ts.Close()

	sc := apiServerCheckerImpl{ts.URL}
	msg, err := sc.CheckCertificate()
	assert.Nil(t, err, "No error should be thrown on valid certificate")
	assert.Contains(t, msg, "Expiry date:", "The expiry date should be present in the output.")
}
