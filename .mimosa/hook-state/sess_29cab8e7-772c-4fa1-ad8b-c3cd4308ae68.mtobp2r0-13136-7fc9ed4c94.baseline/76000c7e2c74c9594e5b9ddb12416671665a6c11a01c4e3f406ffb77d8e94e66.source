package store

func (s *Store) ListShelves() ([]Shelf, error) {
	rows, err := s.db.Query(`SELECT id,name,kind,position,created_at FROM shelves ORDER BY position,id`)
	if err != nil {
		return nil, err
	}
	out := make([]Shelf, 0)
	for rows.Next() {
			var sh Shelf
			if err := rows.Scan(&sh.ID, &sh.Name, &sh.Kind, &sh.Position, &sh.CreatedAt); err != nil {
				rows.Close()
				return nil, err
			}
			sh.TitleIDs = []int64{}
		out = append(out, sh)
	}
	rows.Close()

	for i := range out {
		ids, err := s.shelfItems(out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].TitleIDs = ids
	}
	return out, nil
}

func (s *Store) shelfItems(shelfID int64) ([]int64, error) {
	rows, err := s.db.Query(`SELECT title_id FROM shelf_items WHERE shelf_id=? ORDER BY position`, shelfID)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if ids == nil {
		ids = []int64{}
	}
	return ids, nil
}

func (s *Store) SaveShelf(sh Shelf) (int64, error) {
	if sh.ID == 0 {
		res, err := s.db.Exec(`INSERT INTO shelves(name,kind,position) VALUES(?,?,?)`, sh.Name, sh.Kind, sh.Position)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	_, err := s.db.Exec(`UPDATE shelves SET name=?,kind=?,position=? WHERE id=?`, sh.Name, sh.Kind, sh.Position, sh.ID)
	if err != nil {
		return 0, err
	}
	return sh.ID, nil
}

func (s *Store) DeleteShelf(id int64) error {
	_, err := s.db.Exec(`DELETE FROM shelves WHERE id=?`, id)
	return err
}

func (s *Store) SetShelfItems(shelfID int64, titleIDs []int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM shelf_items WHERE shelf_id=?`, shelfID); err != nil {
		return err
	}
	for i, id := range titleIDs {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO shelf_items(shelf_id,title_id,position) VALUES(?,?,?)`, shelfID, id, i); err != nil {
			return err
		}
	}
	return tx.Commit()
}
