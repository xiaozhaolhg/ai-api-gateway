package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultStreamingTokenInterval(t *testing.T) {
	yaml := `
server:
  port: "8080"
`
	tmpFile := "/tmp/test_config.yaml"
	os.WriteFile(tmpFile, []byte(yaml), 0644)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), *cfg.StreamingTokenInterval)
}

func TestCustomStreamingTokenIntervalViaYAML(t *testing.T) {
	yaml := `
server:
  port: "8080"
streaming_token_interval: 500
`
	tmpFile := "/tmp/test_config_custom.yaml"
	os.WriteFile(tmpFile, []byte(yaml), 0644)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(500), *cfg.StreamingTokenInterval)
}

func TestStreamingTokenIntervalViaEnvVar(t *testing.T) {
	os.Setenv("STREAMING_TOKEN_INTERVAL", "2000")
	defer os.Unsetenv("STREAMING_TOKEN_INTERVAL")

	yaml := `
server:
  port: "8080"
streaming_token_interval: 500
`
	tmpFile := "/tmp/test_config_env.yaml"
	os.WriteFile(tmpFile, []byte(yaml), 0644)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(2000), *cfg.StreamingTokenInterval)
}

func TestStreamingTokenIntervalZeroDisables(t *testing.T) {
	yaml := `
server:
  port: "8080"
streaming_token_interval: 0
`
	tmpFile := "/tmp/test_config_zero.yaml"
	os.WriteFile(tmpFile, []byte(yaml), 0644)
	defer os.Remove(tmpFile)

	cfg, err := Load(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), *cfg.StreamingTokenInterval)
}
