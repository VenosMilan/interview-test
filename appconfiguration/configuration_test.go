package appconfiguration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfiguration(t *testing.T) {
	configWithDefaultValue := NewAppConfiguration()

	assert.Equal(t, "8080", configWithDefaultValue.ServerPort)
	assert.Equal(t, false, configWithDefaultValue.LogDebug)
	assert.Equal(t, "./records.bin", configWithDefaultValue.BinaryFilePath)
}

func TestCustomConfiguration(t *testing.T) {
	err := os.Setenv("LOG_DEBUG", "true")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("PORT", "9090")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("BINARY_FILE_PATH", "/opt/records.bin")
	if err != nil {
		t.Fatal(err)
	}

	config := NewAppConfiguration()

	assert.Equal(t, "9090", config.ServerPort)
	assert.Equal(t, true, config.LogDebug)
	assert.Equal(t, "/opt/records.bin", config.BinaryFilePath)
}
