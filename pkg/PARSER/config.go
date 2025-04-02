package PARSER

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const DNS_LOOKUP_URL = "json-api/cpanel?cpanel_jsonapi_module=ZoneEdit&cpanel_jsonapi_func=fetchzone&domain"
const DNS_EDITOR_URL = "json-api/cpanel?cpanel_jsonapi_module=ZoneEdit&cpanel_jsonapi_func=edit_zone_record&domain"
const ENTRIES_PATH = "./domains"

type DomainEntry struct {
	DomainName    string     `json:"domain_name"`
	MainDomainUrl string     `json:"main_domain_url"`
	TargetDomains []string   `json:"target_domains"`
	Auth          DomainAuth `json:"auth"`
	CpanelPort    int16      `json:"cpanel_port"`
}

type DomainAuth struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func readEntries() ([]*DomainEntry, error) {

	fmt.Println("--- READING DOMAINS ---")

	var res []*DomainEntry

	// read file entries
	entries, err := os.ReadDir(ENTRIES_PATH)
	if err != nil {
		return nil, err
	}

	// loop through and check for *.json files
	for _, f := range entries {
		if !f.IsDir() {
			// validate file extension
			if validateFileExtension(f.Name()) {

				var doc DomainEntry

				file, err := os.Open(fmt.Sprintf("%s/%s", ENTRIES_PATH, f.Name()))
				if err != nil {
					return nil, err
				}

				err = json.NewDecoder(file).Decode(&doc)
				if err != nil {
					return nil, err
				}

				fmt.Printf("--- WATCHING %s DOMAIN ---\n", doc.DomainName)
				res = append(res, &doc)

				file.Close()
			}
		}
	}

	fmt.Printf("--- %d DOMAINS RECORDED ---\n", len(res))
	return res, nil
}

// validateFileExtension validates that the config
// file is of type json
func validateFileExtension(name string) bool {
	entries := strings.Split(name, ".")

	lenght := len(entries)
	return entries[lenght-1] == "json"
}
