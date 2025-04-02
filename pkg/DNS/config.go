package DNS

import (
	"dns-automizer/pkg/PARSER"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
)

type dns_regs struct {
	CpanelResult struct {
		Func  string `json:"func"`
		Event struct {
			Result int `json:"return"`
		} `json:"event"`
		ApiVersion int         `json:"apiversion"`
		Data       []*dns_data `json:"data"`
		Preevent   struct {
			Result int `json:"result"`
		} `json:"preevent"`
		Module    string `json:"module"`
		Postevent struct {
			Result int `json:"result"`
		} `json:"postevent"`
	} `json:"cpanelresult"`
}

type dns_data struct {
	StatusMsg string     `json:"statusmsg"`
	SerialNum string     `json:"serialnum"`
	Status    int        `json:"status"`
	Record    []*dns_rec `json:"record"`
}

type dns_rec struct {
	Line   int    `json:"line"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Adress string `json:"address"`
}

// getARegistries check the registries for the target A records
func (d *DNS) getARegistries(domain *PARSER.DomainEntry) ([]*dns_rec, error) {

	fmt.Println("--- Getting A Registries ---")

	var res []*dns_rec

	var entry dns_regs

	client := http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s:%d/%s=%s", domain.MainDomainUrl, domain.CpanelPort, PARSER.DNS_LOOKUP_URL, domain.MainDomainUrl), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("cpanel %s:%s", domain.Auth.Username, domain.Auth.Token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
	}

	err = json.Unmarshal(data, &entry)
	if err != nil {
		return nil, err
	}

	for _, record := range entry.CpanelResult.Data[0].Record {

		// validate agains target domains
		if record.Type == "A" {
			if slices.Contains(domain.TargetDomains, record.Name) {
				fmt.Printf("--- getting %s entry from remote DNS ---\n", record.Name)
				res = append(res, record)
			}
		}

	}

	fmt.Println("--- Finished getting A Records ---")

	return res, nil
}

// changeRemoteRegistries changes the address on the remote A records
func (d *DNS) changeRemoteRegistries(records []*dns_rec, updatedAddr string, domain *PARSER.DomainEntry) error {

	fmt.Println("--- UPDATING REMOTE A RECORDS ---")

	for _, record := range records {
		fmt.Printf("--- UPDATING A RECORD FOR %s ---\n", record.Name)

		client := http.Client{}

		url := fmt.Sprintf("https://%s:%d/%s=%s&line=%d&type=A&address=%s", domain.MainDomainUrl, domain.CpanelPort, PARSER.DNS_EDITOR_URL, domain.MainDomainUrl, record.Line, updatedAddr)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", fmt.Sprintf("cpanel %s:%s", domain.Auth.Username, domain.Auth.Token))

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to perform update, error:::%d", resp.StatusCode)
		}
		resp.Body.Close()
	}

	fmt.Println("--- RECORDS SUCCESSFULY UPDATED ---")

	return nil
}
