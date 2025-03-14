package main

import (
	"github-actions-cloudflare-dns/cloudflare"

	log "github.com/sirupsen/logrus"
	"github.com/yzinkovets/utils/env"
)

// How it works:
// Get the zone ID by domain name
// Check if the record exists
// If it exists but target not equal to the new one, update the target, otherwise exit
// If it doesn't exist, create a new CNAME record

// TODO:
// Support for other record types (some of them could have more than one record, e.g. A, AAAA, MX, TXT)

func main() {
	// Set log level
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)

	// Read inputs from environment variables
	cfApiToken := env.Must("INPUT_CLOUDFLARE_API_TOKEN")
	domain := env.Must("INPUT_DOMAIN") // full domain name
	target := env.Must("INPUT_TARGET")
	recordType := "CNAME" // for now only CNAME records are supported
	ttl := 3600           // 1 hour

	cf := cloudflare.NewCFClient(cfApiToken)

	// Get the zone ID by domain name
	zoneId, err := cf.GetZoneForDomain(getTLD(domain))
	if err != nil {
		log.Fatal("can't get zone for domain. Error:", err)
	}

	// Check if the record exists
	record, err := cf.GetDnsRecord(zoneId, domain)
	if err != nil {
		log.Fatal("can't get DNS record. Error:", err)
	}

	// If it exists but target not equal to the new one, update the target, otherwise exit
	if record.Id == "" {
		log.Info("DNS record not found. Creating a new one...")

		record.Name = domain
		record.Type = recordType
		record.Content = target
		record.TTL = ttl
		record.Proxied = true

		if err := cf.CreateDnsRecord(zoneId, record); err != nil {
			log.Fatal("can't create DNS record. Error:", err)
		}
		log.Info("DNS record created")
		return
	}

	if record.Content == target {
		log.Info("DNS record found. Target is the same. Exiting...")
		return
	}

	log.Info("DNS record found. Updating the target...")
	record.Content = target

	if err := cf.UpdateDnsRecord(zoneId, record); err != nil {
		log.Fatal("can't update DNS record. Error:", err)
	}
	log.Info("DNS record updated")
}
