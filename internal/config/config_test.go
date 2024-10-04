package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
)

func TestConfig_LoadInvalidYAML(t *testing.T) {
	cfg := New()
	tmpFile, err := os.CreateTemp(os.TempDir(), "")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("usernamd:sdkfo\npassword: foo")
	assert.NoError(t, err)

	cfg.GlobalFlags.ConfigFile = tmpFile.Name()
	err = cfg.Load()
	assert.EqualError(t, err, fmt.Sprintf("unable to parse config from file '%s': While parsing config: yaml: line 2: mapping values are not allowed in this context", tmpFile.Name()))
}

func TestConfig_Load(t *testing.T) {
	cfg := New()
	tmpFile, err := os.CreateTemp(os.TempDir(), "")
	assert.NoError(t, err)
	_, err = tmpFile.WriteString("username: sdkfo\npassword: foo")
	assert.NoError(t, err)

	cfg.GlobalFlags.ConfigFile = tmpFile.Name()
	err = cfg.Load()
	assert.NoError(t, err)
	assert.NotEmpty(t, cfg.GetString("username"))
	assert.NotEmpty(t, cfg.GetString("password"))
}

func TestConfig_Keyring(t *testing.T) {
	keyring.MockInit()
	err := keyring.Set("UPCLOUD", "username", "sdkfo")
	assert.NoError(t, err)
	err = keyring.Set("UPCLOUD", "password", "foo")
	assert.NoError(t, err)

	cfg := New()
	cfg.GlobalFlags.ConfigType = "keyring"
	err = cfg.Load()
	assert.NoError(t, err)
	assert.Equal(t, "sdkfo", cfg.GetString("username"))
	assert.Equal(t, "foo", cfg.GetString("password"))
}
