package main

import (
	"os"
	"path/filepath"
)

func configBaseDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		base, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	return base, nil
}

func userDataDir(appName string) (string, error) {
	base, err := configBaseDir()
	if err != nil {
		return "", err
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

// migrateLegacyData переносит базу и обложки из папки под старым именем
// приложения (Mediateka), если в новой папке ещё нет своей базы.
func migrateLegacyData(newDir string) {
	newDB := filepath.Join(newDir, "pergamon.db")
	if _, err := os.Stat(newDB); err == nil {
		return
	}
	base, err := configBaseDir()
	if err != nil {
		return
	}
	oldDir := filepath.Join(base, "Mediateka")
	data, err := os.ReadFile(filepath.Join(oldDir, "mediateka.db"))
	if err != nil {
		return
	}
	if err := os.WriteFile(newDB, data, 0644); err != nil {
		return
	}
	copyFiles(filepath.Join(oldDir, "covers"), filepath.Join(newDir, "covers"))
}

func copyFiles(srcDir, dstDir string) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		dst := filepath.Join(dstDir, e.Name())
		if _, err := os.Stat(dst); err == nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, e.Name()))
		if err != nil {
			continue
		}
		_ = os.WriteFile(dst, data, 0644)
	}
}
