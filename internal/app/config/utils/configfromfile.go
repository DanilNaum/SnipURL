package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfigFromFile reads a JSON configuration file specified by filename and unmarshals its contents into the provided config struct pointer.
// Returns an error if the file cannot be read or if the JSON is invalid.
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
