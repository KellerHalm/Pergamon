package sync

import (
	"encoding/json"
	"io"

	"mediateka/internal/store"
)

type backupFile struct {
	Version  int           `json:"version"`
	Exported string        `json:"exported"`
	Titles   []store.Title `json:"titles"`
	Shelves  []store.Shelf `json:"shelves"`
}

func exportJSON(s *store.Store, w io.Writer) error {
	titles, err := s.ListTitles(store.ListFilter{Sort: "updated"})
	if err != nil {
		return err
	}
	shelves, err := s.ListShelves()
	if err != nil {
		return err
	}
	b := backupFile{
		Version: 1,
		Titles:  titles,
		Shelves: shelves,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

func importJSON(s *store.Store, r io.Reader) error {
	var b backupFile
	dec := json.NewDecoder(r)
	if err := dec.Decode(&b); err != nil {
		return err
	}
	for _, t := range b.Titles {
		t.ID = 0
		if _, err := s.SaveTitle(t); err != nil {
			return err
		}
	}
	for _, sh := range b.Shelves {
		id, err := s.SaveShelf(sh)
		if err != nil {
			return err
		}
		if err := s.SetShelfItems(id, sh.TitleIDs); err != nil {
			return err
		}
	}
	return nil
}
