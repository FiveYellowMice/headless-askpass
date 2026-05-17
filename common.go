package main

import (
	"os"
	"path/filepath"
)

func socketPath() string {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		dir = "/tmp"
	}
	return filepath.Join(dir, "headless-askpass.socket")
}
