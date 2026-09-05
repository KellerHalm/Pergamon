package sync

import (
	"encoding/json"
	"io"

	"mediateka/internal/store"
)

type backupFile struct {
	Version    int                `json:"version"`
	Exported   string             `json:"exported"`
	Titles     []store.Title      `json:"titles"`
	Shelves    []store.Shelf      `json:"shelves"`
	Notes      []store.Note       `json:"notes,omitempty"`
	Characters []store.Character  `json:"characters,omitempty"`
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
	notes, err := s.ListAllNotes()
	if err != nil {
		return err
	}
	characters, err := s.ListCharacters("")
	if err != nil {
		return err
	}
	b := backupFile{
		Version:    2,
		Titles:     titles,
		Shelves:    shelves,
		Notes:      notes,
		Characters: characters,
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
	type pendingRelations struct {
		newID int64
		rels  []store.TitleRelation
	}
	idMap := make(map[int64]int64)
	var pending []pendingRelations
	for _, t := range b.Titles {
		oldID := t.ID
		rels := t.Relations
		t.ID = 0
		t.Relations = nil
		t.Characters = nil
		newID, err := s.SaveTitle(t)
		if err != nil {
			return err
		}
		if oldID != 0 {
			idMap[oldID] = newID
		}
		if len(rels) > 0 {
			pending = append(pending, pendingRelations{newID: newID, rels: rels})
		}
	}
	for _, p := range pending {
		remapped := make([]store.TitleRelation, 0, len(p.rels))
		for _, rel := range p.rels {
			if id, ok := idMap[rel.RelatedID]; ok {
				remapped = append(remapped, store.TitleRelation{RelatedID: id, Label: rel.Label, ReverseLabel: rel.ReverseLabel})
			}
		}
		if len(remapped) > 0 {
			if err := s.ReplaceTitleRelations(p.newID, remapped); err != nil {
				return err
			}
		}
	}
	type pendingCharacter struct {
		newID    int64
		titleIDs []int64
	}
	var pendingChars []pendingCharacter
	for _, c := range b.Characters {
		titleIDs := c.TitleIDs
		c.ID = 0
		c.TitleIDs = nil
		newID, err := s.SaveCharacter(c)
		if err != nil {
			return err
		}
		pendingChars = append(pendingChars, pendingCharacter{newID: newID, titleIDs: titleIDs})
	}
	for _, pc := range pendingChars {
		remapped := make([]int64, 0, len(pc.titleIDs))
		for _, tid := range pc.titleIDs {
			if id, ok := idMap[tid]; ok {
				remapped = append(remapped, id)
			}
		}
		if len(remapped) > 0 {
			if err := s.SetCharacterTitles(pc.newID, remapped); err != nil {
				return err
			}
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
	for _, n := range b.Notes {
		n.ID = 0
		if newID, ok := idMap[n.TitleID]; ok {
			n.TitleID = newID
		}
		if _, err := s.SaveNote(n); err != nil {
			return err
		}
	}
	return nil
}
