package config_test

import (
	"testing"

	"github.com/dezh-tech/immortal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadfromFile(t *testing.T) {
	cfg, err := config.LoadConfig("./config.yml")
	require.NoError(t, err, "error must be nil.")

	assert.Equal(t, uint16(7777), cfg.ServerConf.Port)
	assert.Equal(t, "127.0.0.1", cfg.ServerConf.Bind)
}
