package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CFClient struct {
	baseApiUrl string
	apiToken   string
}

func NewCFClient(apiToken string) *CFClient {
	return &CFClient{
		baseApiUrl: "https://api.cloudflare.com/client/v4",
		apiToken:   apiToken,
	}
}

type CloudFlareResponse struct {
	Success bool `json:"success"`
}

type Zone struct {
	ID string `json:"id"`
}

type ZoneResponse struct {
	Result []Zone `json:"result"`
	CloudFlareResponse
}

func (c *CFClient) GetZoneForDomain(domain string) (string, error) {
	client := &http.Client{}
	reqURL := fmt.Sprintf("%s/zones?name=%s", c.baseApiUrl, domain)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var zoneResponse ZoneResponse
	if err := json.Unmarshal(body, &zoneResponse); err != nil {
		return "", err
	}

	if !zoneResponse.Success {
		return "", fmt.Errorf("failed to fetch zone")
	}

	return zoneResponse.Result[0].ID, nil
}

type DnsRecord struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
	Proxied bool   `json:"proxied"`
}

type DnsRecordsResponse struct {
	Result []DnsRecord `json:"result"`
	CloudFlareResponse
}

func (c *CFClient) GetDnsRecord(zoneID string, name string) (DnsRecord, error) {
	client := &http.Client{}
	reqURL := fmt.Sprintf("%s/zones/%s/dns_records?name=%s", c.baseApiUrl, zoneID, name)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Error("failed to create request")
		return DnsRecord{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to send request")
		return DnsRecord{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body")
		return DnsRecord{}, err
	}

	var dnsRecordsResponse DnsRecordsResponse
	if err := json.Unmarshal(body, &dnsRecordsResponse); err != nil {
		log.Error("failed to unmarshal response")
		return DnsRecord{}, err
	}

	if !dnsRecordsResponse.Success {
		log.Error("response not successful")
		return DnsRecord{}, fmt.Errorf("failed to fetch DNS record")
	}

	if len(dnsRecordsResponse.Result) == 0 {
		log.Debug("record not found")

		// Record not found
		return DnsRecord{}, nil
	}

	return dnsRecordsResponse.Result[0], nil
}

func (c *CFClient) CreateDnsRecord(zoneID string, record DnsRecord) error {
	client := &http.Client{}
	reqURL := fmt.Sprintf("%s/zones/%s/dns_records", c.baseApiUrl, zoneID)

	recordJson, err := json.Marshal(record)
	if err != nil {
		log.Error("failed to marshal record")
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(recordJson))
	if err != nil {
		log.Error("failed to create request")
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to send request")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body")
		return err
	}

	var cloudFlareResponse CloudFlareResponse
	if err := json.Unmarshal(body, &cloudFlareResponse); err != nil {
		log.Error("failed to unmarshal response")
		return err
	}

	if !cloudFlareResponse.Success {
		log.Error("response not successful")
		return fmt.Errorf("failed to create DNS record")
	}

	return nil
}

func (c *CFClient) UpdateDnsRecord(zoneID string, record DnsRecord) error {

	client := &http.Client{}
	reqURL := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.baseApiUrl, zoneID, record.Id)

	updateReq := DnsRecord{
		Type:    record.Type,
		Name:    record.Name,
		Content: record.Content,
		TTL:     record.TTL,
		Proxied: record.Proxied,
	}

	jsonReq, err := json.Marshal(updateReq)
	if err != nil {
		log.Error("failed to marshal record")
		return err
	}

	req, err := http.NewRequest("PUT", reqURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Error("failed to create request")
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to send request")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("response not successful")
		return fmt.Errorf("failed to update DNS record: %s", body)
	}

	return nil
}
