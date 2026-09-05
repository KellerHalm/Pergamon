package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateLegacyData(t *testing.T) {
	base := t.TempDir()
	t.Setenv("AppData", base)

	oldDir := filepath.Join(base, "Mediateka")
	if err := os.MkdirAll(filepath.Join(oldDir, "covers"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oldDir, "mediateka.db"), []byte("db-data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oldDir, "covers", "a.png"), []byte("cover-a"), 0644); err != nil {
		t.Fatal(err)
	}

	newDir := filepath.Join(base, "Pergamon")
	if err := os.MkdirAll(newDir, 0755); err != nil {
		t.Fatal(err)
	}

	migrateLegacyData(newDir)

	got, err := os.ReadFile(filepath.Join(newDir, "pergamon.db"))
	if err != nil || string(got) != "db-data" {
		t.Fatalf("pergamon.db not migrated: %v %q", err, got)
	}
	gotCover, err := os.ReadFile(filepath.Join(newDir, "covers", "a.png"))
	if err != nil || string(gotCover) != "cover-a" {
		t.Fatalf("cover not migrated: %v %q", err, gotCover)
	}
}

func TestMigrateLegacyDataKeepsExistingDB(t *testing.T) {
	base := t.TempDir()
	t.Setenv("AppData", base)

	oldDir := filepath.Join(base, "Mediateka")
	if err := os.MkdirAll(oldDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oldDir, "mediateka.db"), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	newDir := filepath.Join(base, "Pergamon")
	if err := os.MkdirAll(newDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(newDir, "pergamon.db"), []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}

	migrateLegacyData(newDir)

	got, _ := os.ReadFile(filepath.Join(newDir, "pergamon.db"))
	if string(got) != "new" {
		t.Fatalf("existing db was overwritten: %q", got)
	}
}

func TestMigrateLegacyDataNoOldData(t *testing.T) {
	base := t.TempDir()
	t.Setenv("AppData", base)

	newDir := filepath.Join(base, "Pergamon")
	if err := os.MkdirAll(newDir, 0755); err != nil {
		t.Fatal(err)
	}

	migrateLegacyData(newDir)

	if _, err := os.Stat(filepath.Join(newDir, "pergamon.db")); !os.IsNotExist(err) {
		t.Fatalf("db created without legacy data: %v", err)
	}
}
