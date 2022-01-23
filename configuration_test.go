package goutil

import (
	"os"
	"testing"
)

type ConfigModel struct {
	Port        int    `json:"port" env:"PORT"`
	ServiceName string `json:"service_name" env:"SERVICE_NAME"`
	Model       struct {
		ModelName string `json:"model_name" env:"MODEL_NAME"`
	} `json:"model"`
}

func init() {
	os.Setenv("PORT", "8089")
	os.Setenv("SERVICE_NAME", "your service")
	os.Setenv("MODEL_NAME", "your model")
}

// go test -v -run=TestConfiguration
func TestConfiguration(t *testing.T) {
	var cfg ConfigModel
	osFile, err := OpenFile(".", "configuration.json")
	if err != nil {
		t.Error(err)
		return
	}

	if err := Configuration(osFile, &cfg); err != nil {
		t.Error(err)
		return
	}

	t.Logf("Result: %+v", cfg)
}
