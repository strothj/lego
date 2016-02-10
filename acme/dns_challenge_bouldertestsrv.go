package acme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DNSProviderBoulderTestSRV is an implementation of the DNSProvider interface
// that manages TXT records on a Boulder DNS test server.
type DNSProviderBoulderTestSRV struct {
	boulderDNSBaseURL string
}

// NewDNSProviderBoulderTestSRV returns a new DNSProviderBoulderTestSRV
// instance.
func NewDNSProviderBoulderTestSRV(dnsURL string) (*DNSProviderBoulderTestSRV, error) {
	return &DNSProviderBoulderTestSRV{
		boulderDNSBaseURL: dnsURL,
	}, nil
}

// Present creates a TXT record to fulfil the dns-01 challenge
func (b *DNSProviderBoulderTestSRV) Present(domain, token, keyAuth string) error {
	// txtRecordRequest represents the request body to DO's API to make a TXT record
	type txtRecordRequest struct {
		Host  string `json:"host"`
		Value string `json:"value"`
	}

	fqdn, value, _ := DNS01Record(domain, keyAuth)

	reqURL := fmt.Sprintf("%s/set-txt", b.boulderDNSBaseURL)
	reqData := txtRecordRequest{Host: fqdn, Value: value}
	body, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// CleanUp removes the TXT record matching the specified parameters
func (b *DNSProviderBoulderTestSRV) CleanUp(domain, token, keyAuth string) error {
	return nil
}
