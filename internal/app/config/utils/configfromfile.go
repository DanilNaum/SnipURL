package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfigFromFile[T any](filename string, config *T) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := json.Unmarshal(file, config); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil
}
