package config

import (
	"path/filepath"
	"strconv"

	"github.com/btcsuite/btcutil"
)

var (
	DefaultName    = "playground.config.json"
	DefaultCompose = "cmd/docker/plebnet-playground-docker/docker-compose.yaml"

	DefaultDatadir = btcutil.AppDataDir("playground", false)
	DefaultPath    = filepath.Join(DefaultDatadir, DefaultName)

	InitialState = map[string]string{
		"network": "regtest",
		"ready":   strconv.FormatBool(false),
		"running": strconv.FormatBool(false),
	}
)
