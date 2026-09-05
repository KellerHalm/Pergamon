package store

import (
	"database/sql"
	"strings"

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
			release_status TEXT NOT NULL DEFAULT '',
			custom_list TEXT DEFAULT '',
			progress_volumes INTEGER NOT NULL DEFAULT 0,
			progress_chapters INTEGER NOT NULL DEFAULT 0,
			progress_pages INTEGER NOT NULL DEFAULT 0,
			progress_seasons INTEGER NOT NULL DEFAULT 0,
			progress_episodes INTEGER NOT NULL DEFAULT 0,
			progress_minutes INTEGER NOT NULL DEFAULT 0,
			progress_total_chapters INTEGER NOT NULL DEFAULT 0,
			progress_total_episodes INTEGER NOT NULL DEFAULT 0,
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
		`CREATE TABLE IF NOT EXISTS title_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			file TEXT NOT NULL,
			position INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS title_relations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			related_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			label TEXT NOT NULL DEFAULT '',
			reverse_label TEXT NOT NULL DEFAULT ''
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
		`CREATE TABLE IF NOT EXISTS title_notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			heading TEXT DEFAULT '',
			content TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS characters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			main_image TEXT DEFAULT '',
			age TEXT DEFAULT '',
			gender TEXT NOT NULL DEFAULT '',
			race TEXT DEFAULT '',
			description TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS character_fields (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			value TEXT DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS character_names (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
			kind TEXT NOT NULL,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS character_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
			file TEXT NOT NULL,
			position INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS title_characters (
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
			position INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (title_id, character_id)
		)`,
		`CREATE TABLE IF NOT EXISTS studios (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			main_image TEXT DEFAULT '',
			founded TEXT DEFAULT '',
			description TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS studio_names (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			studio_id INTEGER NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
			kind TEXT NOT NULL,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS studio_founders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			studio_id INTEGER NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			position INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS studio_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			studio_id INTEGER NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
			file TEXT NOT NULL,
			position INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS title_studios (
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			studio_id INTEGER NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
			position INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (title_id, studio_id)
		)`,
		`CREATE TABLE IF NOT EXISTS people (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			main_image TEXT DEFAULT '',
			age TEXT DEFAULT '',
			birth_date TEXT DEFAULT '',
			death_date TEXT DEFAULT '',
			gender TEXT NOT NULL DEFAULT '',
			role TEXT NOT NULL DEFAULT 'author',
			description TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS person_names (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
			kind TEXT NOT NULL,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS person_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
			file TEXT NOT NULL,
			position INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS title_people (
			title_id INTEGER NOT NULL REFERENCES titles(id) ON DELETE CASCADE,
			person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
			position INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (title_id, person_id)
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
		`CREATE INDEX IF NOT EXISTS idx_title_images_title ON title_images(title_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_notes_title ON title_notes(title_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_relations_title ON title_relations(title_id)`,
		`CREATE INDEX IF NOT EXISTS idx_character_names_value ON character_names(value)`,
		`CREATE INDEX IF NOT EXISTS idx_character_images_character ON character_images(character_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_characters_title ON title_characters(title_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_characters_character ON title_characters(character_id)`,
		`CREATE INDEX IF NOT EXISTS idx_studio_names_value ON studio_names(value)`,
		`CREATE INDEX IF NOT EXISTS idx_studio_images_studio ON studio_images(studio_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_studios_title ON title_studios(title_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_studios_studio ON title_studios(studio_id)`,
		`CREATE INDEX IF NOT EXISTS idx_person_names_value ON person_names(value)`,
		`CREATE INDEX IF NOT EXISTS idx_person_images_person ON person_images(person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_people_title ON title_people(title_id)`,
		`CREATE INDEX IF NOT EXISTS idx_title_people_person ON title_people(person_id)`,
	}
	for _, st := range stmts {
		if _, err := s.db.Exec(st); err != nil {
			return err
		}
	}
	if err := s.addColumnIfMissing("titles", "release_status", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.addColumnIfMissing("titles", "progress_total_chapters", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.addColumnIfMissing("titles", "progress_total_episodes", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.addColumnIfMissing("title_relations", "reverse_label", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.addColumnIfMissing("characters", "gender", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.addColumnIfMissing("characters", "race", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if _, err := s.db.Exec(`INSERT INTO title_notes(title_id, heading, content)
		SELECT id, 'Заметки', notes FROM titles WHERE TRIM(notes) != ''`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`UPDATE titles SET notes='' WHERE TRIM(notes) != ''`); err != nil {
		return err
	}
	if err := s.migrateCreatorsToEntities(); err != nil {
		return err
	}
	return nil
}

func (s *Store) migrateCreatorsToEntities() error {
	rows, err := s.db.Query(`SELECT title_id, role, name FROM title_creators ORDER BY id`)
	if err != nil {
		return err
	}
	type creatorRow struct {
		titleID int64
		role    string
		name    string
	}
	var rowsData []creatorRow
	for rows.Next() {
		var r creatorRow
		if err := rows.Scan(&r.titleID, &r.role, &r.name); err != nil {
			rows.Close()
			return err
		}
		rowsData = append(rowsData, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	if len(rowsData) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, r := range rowsData {
		name := strings.TrimSpace(r.name)
		if name == "" {
			continue
		}
		if r.role == "studio" {
			id, err := findOrCreateStudioTx(tx, name)
			if err != nil {
				return err
			}
			if _, err := tx.Exec(`INSERT OR IGNORE INTO title_studios(title_id,studio_id,position) VALUES(?,?,0)`, r.titleID, id); err != nil {
				return err
			}
		} else {
			id, err := findOrCreatePersonTx(tx, name, r.role)
			if err != nil {
				return err
			}
			if _, err := tx.Exec(`INSERT OR IGNORE INTO title_people(title_id,person_id,position) VALUES(?,?,0)`, r.titleID, id); err != nil {
				return err
			}
		}
	}
	if _, err := tx.Exec(`DELETE FROM title_creators`); err != nil {
		return err
	}
	return tx.Commit()
}

func findOrCreateStudioTx(tx *sql.Tx, name string) (int64, error) {
	rows, err := tx.Query(`SELECT studio_id, value FROM studio_names`)
	if err != nil {
		return 0, err
	}
	var found int64
	for rows.Next() {
		var id int64
		var v string
		if err := rows.Scan(&id, &v); err != nil {
			rows.Close()
			return 0, err
		}
		if strings.EqualFold(strings.TrimSpace(v), name) {
			found = id
			break
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if found != 0 {
		return found, nil
	}
	res, err := tx.Exec(`INSERT INTO studios(main_image,founded,description) VALUES('','','')`)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if _, err := tx.Exec(`INSERT INTO studio_names(studio_id,kind,value) VALUES(?,?,?)`, id, "original", name); err != nil {
		return 0, err
	}
	return id, nil
}

func findOrCreatePersonTx(tx *sql.Tx, name, role string) (int64, error) {
	rows, err := tx.Query(`SELECT p.id, p.role, pn.value FROM people p JOIN person_names pn ON pn.person_id=p.id`)
	if err != nil {
		return 0, err
	}
	var found int64
	for rows.Next() {
		var id int64
		var pRole, v string
		if err := rows.Scan(&id, &pRole, &v); err != nil {
			rows.Close()
			return 0, err
		}
		if pRole == role && strings.EqualFold(strings.TrimSpace(v), name) {
			found = id
			break
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if found != 0 {
		return found, nil
	}
	res, err := tx.Exec(`INSERT INTO people(role) VALUES(?)`, role)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if _, err := tx.Exec(`INSERT INTO person_names(person_id,kind,value) VALUES(?,?,?)`, id, "original", name); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) addColumnIfMissing(table, column, decl string) error {
	rows, err := s.db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.RawBytes
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return err
		}
		if name == column {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	rows.Close()
	_, err = s.db.Exec(`ALTER TABLE ` + table + ` ADD COLUMN ` + column + ` ` + decl)
	return err
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
