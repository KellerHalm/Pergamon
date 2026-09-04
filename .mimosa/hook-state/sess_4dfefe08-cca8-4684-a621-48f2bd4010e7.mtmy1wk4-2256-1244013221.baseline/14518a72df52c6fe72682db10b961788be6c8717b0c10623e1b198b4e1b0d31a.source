package main

import (
	"os"
	"path/filepath"
)

func userDataDir(appName string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		base, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	dir := filepath.Join(base, appName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(dir, "covers"), 0755); err != nil {
		return "", err
	}
	return dir, nil
}
