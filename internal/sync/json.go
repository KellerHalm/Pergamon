package sync

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"pergamon/internal/store"
)

type backupFile struct {
	Version    int               `json:"version"`
	Exported   string            `json:"exported"`
	Titles     []store.Title     `json:"titles"`
	Shelves    []store.Shelf     `json:"shelves"`
	Notes      []store.Note      `json:"notes,omitempty"`
	Characters []store.Character `json:"characters,omitempty"`
	Studios    []store.Studio    `json:"studios,omitempty"`
	People     []store.Person    `json:"people,omitempty"`
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
	studios, err := s.ListStudios("")
	if err != nil {
		return err
	}
	people, err := s.ListPeople("")
	if err != nil {
		return err
	}
	b := backupFile{
		Version:    2,
		Titles:     titles,
		Shelves:    shelves,
		Notes:      notes,
		Characters: characters,
		Studios:    studios,
		People:     people,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

type ImportResult struct {
	Added   int `json:"added"`
	Updated int `json:"updated"`
}

// Импорт работает как слияние: записи, уже существующие в базе (определяются
// по имени), обновляются данными из бэкапа вместо дублирования.
func importJSON(s *store.Store, r io.Reader) (ImportResult, error) {
	var res ImportResult
	var b backupFile
	dec := json.NewDecoder(r)
	if err := dec.Decode(&b); err != nil {
		return res, err
	}

	titleIdx, err := buildTitleIndex(s)
	if err != nil {
		return res, err
	}
	charIdx := map[string]int64{}
	chars, err := s.ListCharacters("")
	if err != nil {
		return res, err
	}
	for _, c := range chars {
		indexNames(charIdx, c.ID, c.Names)
	}
	studioIdx := map[string]int64{}
	studios, err := s.ListStudios("")
	if err != nil {
		return res, err
	}
	for _, st := range studios {
		indexNames(studioIdx, st.ID, st.Names)
	}
	peopleIdx := map[string]int64{}
	people, err := s.ListPeople("")
	if err != nil {
		return res, err
	}
	for _, p := range people {
		indexNames(peopleIdx, p.ID, p.Names)
	}
	shelfIdx := map[string]int64{}
	existingShelves, err := s.ListShelves()
	if err != nil {
		return res, err
	}
	for _, sh := range existingShelves {
		if k := normName(sh.Name); k != "" {
			shelfIdx[k] = sh.ID
		}
	}
	existingNotes, err := s.ListAllNotes()
	if err != nil {
		return res, err
	}
	noteSet := map[string]bool{}
	for _, n := range existingNotes {
		noteSet[noteKey(n.TitleID, n.Heading, n.Content)] = true
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
		t.ID = matchTitle(titleIdx, t)
		if t.ID != 0 {
			res.Updated++
		} else {
			res.Added++
		}
		t.Relations = nil
		t.Characters = nil
		newID, err := s.SaveTitle(t)
		if err != nil {
			return res, err
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
				return res, err
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
		c.ID = matchByName(charIdx, c.Names)
		if c.ID != 0 {
			res.Updated++
		} else {
			res.Added++
		}
		c.TitleIDs = nil
		newID, err := s.SaveCharacter(c)
		if err != nil {
			return res, err
		}
		pendingChars = append(pendingChars, pendingCharacter{newID: newID, titleIDs: titleIDs})
	}
	for _, pc := range pendingChars {
		remapped := remapIDs(pc.titleIDs, idMap)
		if len(remapped) > 0 {
			if err := s.SetCharacterTitles(pc.newID, remapped); err != nil {
				return res, err
			}
		}
	}
	type pendingEntity struct {
		newID    int64
		titleIDs []int64
	}
	var pendingStudios []pendingEntity
	for _, st := range b.Studios {
		titleIDs := st.TitleIDs
		st.ID = matchByName(studioIdx, st.Names)
		if st.ID != 0 {
			res.Updated++
		} else {
			res.Added++
		}
		st.TitleIDs = nil
		newID, err := s.SaveStudio(st)
		if err != nil {
			return res, err
		}
		pendingStudios = append(pendingStudios, pendingEntity{newID: newID, titleIDs: titleIDs})
	}
	for _, ps := range pendingStudios {
		remapped := remapIDs(ps.titleIDs, idMap)
		if len(remapped) > 0 {
			if err := s.SetStudioTitles(ps.newID, remapped); err != nil {
				return res, err
			}
		}
	}
	var pendingPeople []pendingEntity
	for _, p := range b.People {
		titleIDs := p.TitleIDs
		p.ID = matchByName(peopleIdx, p.Names)
		if p.ID != 0 {
			res.Updated++
		} else {
			res.Added++
		}
		p.TitleIDs = nil
		newID, err := s.SavePerson(p)
		if err != nil {
			return res, err
		}
		pendingPeople = append(pendingPeople, pendingEntity{newID: newID, titleIDs: titleIDs})
	}
	for _, pp := range pendingPeople {
		remapped := remapIDs(pp.titleIDs, idMap)
		if len(remapped) > 0 {
			if err := s.SetPersonTitles(pp.newID, remapped); err != nil {
				return res, err
			}
		}
	}
	for _, sh := range b.Shelves {
		sh.ID = shelfIdx[normName(sh.Name)]
		if sh.ID != 0 {
			res.Updated++
		} else {
			res.Added++
		}
		sh.TitleIDs = remapIDs(sh.TitleIDs, idMap)
		id, err := s.SaveShelf(sh)
		if err != nil {
			return res, err
		}
		if err := s.SetShelfItems(id, sh.TitleIDs); err != nil {
			return res, err
		}
	}
	for _, n := range b.Notes {
		newTitleID, ok := idMap[n.TitleID]
		if !ok {
			continue
		}
		if noteSet[noteKey(newTitleID, n.Heading, n.Content)] {
			continue
		}
		n.ID = 0
		n.TitleID = newTitleID
		res.Added++
		if _, err := s.SaveNote(n); err != nil {
			return res, err
		}
	}
	return res, nil
}

func normName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func indexNames(idx map[string]int64, id int64, names []store.Name) {
	for _, n := range names {
		if k := normName(n.Value); k != "" {
			idx[k] = id
		}
	}
}

func matchByName(idx map[string]int64, names []store.Name) int64 {
	for _, n := range names {
		if id, ok := idx[normName(n.Value)]; ok {
			return id
		}
	}
	return 0
}

// Ключ книги — любое из её имён в паре с категорией: одна и та же книга
// в разных категориях (например роман и его экранизация) не склеиваются.
func buildTitleIndex(s *store.Store) (map[string]int64, error) {
	titles, err := s.ListTitles(store.ListFilter{})
	if err != nil {
		return nil, err
	}
	idx := make(map[string]int64)
	for _, t := range titles {
		for _, n := range t.Names {
			if k := normName(n.Value); k != "" {
				idx[t.Category+"\x00"+k] = t.ID
			}
		}
	}
	return idx, nil
}

func matchTitle(idx map[string]int64, t store.Title) int64 {
	for _, n := range t.Names {
		if k := normName(n.Value); k != "" {
			if id, ok := idx[t.Category+"\x00"+k]; ok {
				return id
			}
		}
	}
	return 0
}

func noteKey(titleID int64, heading, content string) string {
	return strconv.FormatInt(titleID, 10) + "\x00" + heading + "\x00" + content
}

func remapIDs(ids []int64, idMap map[int64]int64) []int64 {
	remapped := make([]int64, 0, len(ids))
	for _, id := range ids {
		if newID, ok := idMap[id]; ok {
			remapped = append(remapped, newID)
		}
	}
	return remapped
}
