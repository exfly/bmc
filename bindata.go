package bmc

import (
	"embed"
)

var (
	//go:embed bin/runc
	//go:embed config.json
	BinDataFs embed.FS
)
