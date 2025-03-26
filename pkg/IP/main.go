package IP

import (
	"fmt"
	"strings"
)

type IP struct{}

func StartIPService() *IP {
	return &IP{}
}

func (i *IP) StartIPComprobation() (string, string, bool, error) {
	fmt.Println("--- CHEKING LAST IP RECORD ---")

	var match bool = true

	currIP, err := i.readIPAddr()
	if err != nil {
		return "", "", match, err
	}

	serverIP, err := i.getPublicIPAddr()
	if err != nil {
		return "", "", match, err
	}

	serverIP = []byte(strings.ReplaceAll(string(serverIP), "\n", ""))

	// check if old ip is not yet registered
	if len(currIP) == 0 {
		fmt.Println("---NOT PREVIOUSE SAVED STATE DETECTED. RECORDING...---")
		i.saveAddress(serverIP)
	}

	if string(currIP) != string(serverIP) {
		fmt.Println("---IP CHANGE DETECTED SAVING NEW STATE---")
		// saving new IP
		i.saveAddress(serverIP)
		match = false
	} else {
		fmt.Println("---IP STATE REMAINS UNTOUCHED---")
	}

	return string(serverIP), string(currIP), match, nil

}
