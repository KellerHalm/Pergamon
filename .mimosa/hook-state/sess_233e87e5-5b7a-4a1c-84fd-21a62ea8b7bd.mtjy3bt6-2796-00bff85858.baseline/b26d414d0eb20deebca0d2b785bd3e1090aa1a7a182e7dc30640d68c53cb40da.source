package store

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS titles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL DEFAULT 'read',
			category TEXT NOT NULL DEFAULT 'book',
			cover TEXT DEFAULT '',
			synopsis TEXT DEFAULT '',
			score REAL NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'planned',
			custom_list TEXT DEFAULT '',
			progress_volumes INTEGER NOT NULL DEFAULT 0,
			progress_chapters INTEGER NOT NULL DEFAULT 0,
			progress_pages INTEGER NOT NULL DEFAULT 0,
			progress_seasons INTEGER NOT NULL DEFAULT 0,
			progress_episodes INTEGER NOT NULL DEFAULT 0,
			progress_minutes INTEGER NOT NULL DEFAULT 0,
			notes TEXT DEFAULT '',
			spine_color TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS title_names (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			kind TEXT NOT NULL,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS title_creators (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			role TEXT NOT NULL,
			name TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS title_genres (
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			genre TEXT NOT NULL,
			PRIMARY KEY (title_id, genre)
		)`,
		`CREATE TABLE IF NOT EXISTS title_tags (
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			tag TEXT NOT NULL,
			PRIMARY KEY (title_id, tag)
		)`,
		`CREATE TABLE IF NOT EXISTS shelves (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			kind TEXT NOT NULL DEFAULT 'read',
			position INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS shelf_items (
			shelf_id INTEGER NOT NULL REFERENCES shelves(id) ON DELETE CASCADE,
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			position INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (shelf_id, title_id)
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE INDEX IF NOT EXISTS idx_titles_status ON titles(status)`,
		`CREATE INDEX IF NOT EXISTS idx_titles_type ON titles(type)`,
		`CREATE INDEX IF NOT EXISTS idx_titles_category ON titles(category)`,
		`CREATE INDEX IF NOT EXISTS idx_titles_updated ON titles(updated_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_title_names_value ON title_names(value)`,
		`CREATE INDEX IF NOT EXISTS idx_title_creators_name ON title_creators(name)`,
		`CREATE INDEX IF NOT EXISTS idx_title_tags_tag ON title_tags(tag)`,
	}
	for _, st := range stmts {
		if _, err := s.db.Exec(st); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) GetSetting(key string) (string, error) {
	var v string
	err := s.db.QueryRow(`SELECT value FROM settings WHERE key=?`, key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return v, err
}

func (s *Store) SetSetting(key, value string) error {
	_, err := s.db.Exec(`INSERT INTO settings(key,value) VALUES(?,?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}
