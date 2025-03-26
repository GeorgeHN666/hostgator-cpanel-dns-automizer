package DNS

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
func (d *DNS) getARegistries() ([]*dns_rec, error) {

	fmt.Println("--- Getting A Registries ---")

	var res []*dns_rec

	var entry dns_regs

	client := http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("REMOTE_REGISTRY_LIST"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("cpanel %s", os.Getenv("REMOTE_REGISTRY_AUTH")))

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
		fmt.Printf("REG::: %s\n", record.Name)
		if record.Name == os.Getenv("REMOTE_REGISTRY_TARGET") || record.Name == fmt.Sprintf("www.%s", os.Getenv("REMOTE_REGISTRY_TARGET")) {
			if record.Type == "A" {
				res = append(res, record)
			}
		}
	}

	fmt.Println("--- Finished getting A Records ---")

	return res, nil
}

// changeRemoteRegistries changes the address on the remote A records
func (d *DNS) changeRemoteRegistries(records []*dns_rec, updatedAddr string) error {

	fmt.Println("--- UPDATING REMOTE A RECORDS ---")

	for _, record := range records {
		fmt.Printf("--- UPDATING A RECORD FOR %s ---\n", record.Name)

		client := http.Client{}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s&line=%d&type=A&address=%s", os.Getenv("REMOTE_RECORDS_REGISTRY_PATH"), os.Getenv("REMOTE_REGISTRY_DOMAIN"), record.Line, updatedAddr), nil)
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", fmt.Sprintf("cpanel %s", os.Getenv("REMOTE_REGISTRY_AUTH")))

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
