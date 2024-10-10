package config

import (
	"os"
)

var (
	DEBUG_MODE bool = false
)

func Init() {
	if os.Getenv("DEBUG_MODE") == "true" {
		DEBUG_MODE = true
	}
}
