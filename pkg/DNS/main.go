package DNS

import (
	"fmt"
)

type DNS struct{}

func StartDNSService() *DNS {
	return &DNS{}
}

// StartRecordUpdate search for the updated A record lines and updates the new records to the remote DNS registry
func (d *DNS) StartRecordUpdate(updatedAddr string) error {

	fmt.Println("--- SEARCHING FOR UPDATED RECORDS ---")

	records, err := d.getARegistries()
	if err != nil {
		return err
	}

	fmt.Println("--- UPDATING REMOTE A RECORDS ---")

	return d.changeRemoteRegistries(records, updatedAddr)
}
