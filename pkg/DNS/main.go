package DNS

import (
	"dns-automizer/pkg/PARSER"
	"fmt"
)

type DNS struct{}

func StartDNSService() *DNS {
	return &DNS{}
}

// StartRecordUpdate search for the updated A record lines and updates the new records to the remote DNS registry
func (d *DNS) StartRecordUpdate(updatedAddr string, entries []*PARSER.DomainEntry) error {

	for _, domain := range entries {
		fmt.Printf("--- SEARCHING UPDATED RECORDS FOR %s ---\n", domain.MainDomainUrl)

		records, err := d.getARegistries(domain)
		if err != nil {
			return err
		}

		fmt.Printf("--- UPDATING REMOTE A RECORDS FOR %s---\n", domain.MainDomainUrl)

		d.changeRemoteRegistries(records, updatedAddr, domain)
	}
	return nil
}
