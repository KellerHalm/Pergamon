package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"mediateka/internal/store"
	"mediateka/internal/sync"
)

func relJSON(t *store.Title) string {
	b, _ := json.Marshal(t.Relations)
	return string(b)
}

func main() {
	tmp, _ := os.MkdirTemp("", "rel-test")
	defer os.RemoveAll(tmp)

	s, err := store.Open(filepath.Join(tmp, "test.db"))
	if err != nil {
		fmt.Println("OPEN ERR:", err)
		return
	}
	defer s.Close()

	idA, err := s.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "russian", Value: "Тайтл А"}, {Kind: "original", Value: "Title A"}}})
	fmt.Println("SAVE A:", idA, err)
	idB, err := s.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "russian", Value: "Тайтл Б"}}})
	fmt.Println("SAVE B:", idB, err)
	idC, err := s.SaveTitle(store.Title{Type: "watch", Category: "anime", Names: []store.Name{{Kind: "english", Value: "Title C"}}})
	fmt.Println("SAVE C:", idC, err)

	idB2, err := s.SaveTitle(store.Title{
		ID: idB, Type: "read", Category: "book",
		Names: []store.Name{{Kind: "russian", Value: "Тайтл Б"}},
		Relations: []store.TitleRelation{
			{RelatedID: idA, Label: "Приквел"},
			{RelatedID: idC, Label: "Сиквел"},
			{RelatedID: idB, Label: "сам на себя"},
			{RelatedID: 0, Label: "пусто"},
		},
	})
	fmt.Println("SAVE B WITH RELS:", idB2, err)

	t, err := s.GetTitle(idB)
	if err != nil {
		fmt.Println("GET B ERR:", err)
		return
	}
	fmt.Println("B RELATIONS:", relJSON(t))

	if err := s.DeleteTitle(idA); err != nil {
		fmt.Println("DEL A ERR:", err)
		return
	}
	t, _ = s.GetTitle(idB)
	fmt.Println("B RELATIONS AFTER A DELETED:", relJSON(t))

	if err := s.ReplaceTitleRelations(idB, []store.TitleRelation{{RelatedID: idC, Label: "Продолжение"}}); err != nil {
		fmt.Println("REPLACE ERR:", err)
		return
	}
	t, _ = s.GetTitle(idB)
	fmt.Println("B RELATIONS AFTER REPLACE:", relJSON(t))

	m1 := sync.NewManager(filepath.Join(tmp, "src.db"), s)
	var buf bytes.Buffer
	if err := m1.ExportJSON(&buf); err != nil {
		fmt.Println("EXPORT ERR:", err)
		return
	}

	s2, err := store.Open(filepath.Join(tmp, "dst.db"))
	if err != nil {
		fmt.Println("OPEN DST ERR:", err)
		return
	}
	defer s2.Close()
	m2 := sync.NewManager(filepath.Join(tmp, "dst.db"), s2)
	if err := m2.ImportJSON(bytes.NewReader(buf.Bytes())); err != nil {
		fmt.Println("IMPORT ERR:", err)
		return
	}

	imported, err := s2.ListTitles(store.ListFilter{})
	if err != nil {
		fmt.Println("LIST DST ERR:", err)
		return
	}
	for _, it := range imported {
		fmt.Printf("IMPORTED id=%d name=%q rels=%s\n", it.ID, it.Names[0].Value, relJSON(&it))
	}
}
