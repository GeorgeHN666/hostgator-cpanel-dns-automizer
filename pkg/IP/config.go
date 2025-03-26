package IP

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func (i *IP) getPublicIPAddr() ([]byte, error) {

	res, err := http.Get(os.Getenv("PUBLIC_IP_CHECKER"))
	if err != nil {
		return []byte(""), err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte(""), err
	}

	return data, nil
}

func (i *IP) saveAddress(addr []byte) error {

	// Ensure the directory exists
	if err := os.MkdirAll(os.Getenv("CONFIG_IP_REGISTRY_PATH"), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file with more restrictive permissions (0600 for sensitive data)
	if err := os.WriteFile(os.Getenv("CONFIG_IP_REGISTRY_FILE"), addr, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (i *IP) readIPAddr() ([]byte, error) {

	file, err := os.ReadFile(os.Getenv("CONFIG_IP_REGISTRY_FILE"))
	if err != nil {
		return []byte(""), nil
	}

	return file, err
}
