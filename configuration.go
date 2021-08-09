package util

import (
	"encoding/json"
	"os"
)

// Configuration is a function for get info configuration
func Configuration(osFile *os.File, model interface{}) error {
	decoder := json.NewDecoder(osFile)
	err := decoder.Decode(model)
	if err != nil {
		return err
	}

	return err
}
