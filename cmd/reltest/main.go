package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"pergamon/internal/store"
	"pergamon/internal/sync"
)

func relJSON(rels []store.TitleRelation) string {
	b, _ := json.Marshal(rels)
	return string(b)
}

func printTitle(s *store.Store, id int64, tag string) {
	t, err := s.GetTitle(id)
	if err != nil {
		fmt.Println(tag, "GET ERR:", err)
		return
	}
	fmt.Printf("%s OWN: %s\n", tag, relJSON(t.Relations))
	fmt.Printf("%s REVERSE: %s\n", tag, relJSON(t.ReverseRelations))
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

	idA, err := s.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "russian", Value: "Книга"}}})
	fmt.Println("SAVE A:", idA, err)
	idB, err := s.SaveTitle(store.Title{Type: "watch", Category: "movie", Names: []store.Name{{Kind: "russian", Value: "Фильм"}}})
	fmt.Println("SAVE B:", idB, err)
	idC, err := s.SaveTitle(store.Title{Type: "watch", Category: "anime", Names: []store.Name{{Kind: "english", Value: "Title C"}}})
	fmt.Println("SAVE C:", idC, err)

	_, err = s.SaveTitle(store.Title{
		ID: idB, Type: "watch", Category: "movie",
		Names: []store.Name{{Kind: "russian", Value: "Фильм"}},
		Relations: []store.TitleRelation{
			{RelatedID: idA, Label: "Источник", ReverseLabel: "Экранизация"},
			{RelatedID: idC, Label: "Сиквел"},
			{RelatedID: idB, Label: "сам на себя"},
			{RelatedID: 0, Label: "пусто"},
		},
	})
	fmt.Println("SAVE B WITH RELS:", err)

	printTitle(s, idB, "B")
	printTitle(s, idA, "A")
	printTitle(s, idC, "C")

	if err := s.UpdateIncomingRelations(idA, []store.TitleRelation{{RelatedID: idB, Label: "Снято по книге", ReverseLabel: "Первичный источник"}}); err != nil {
		fmt.Println("UPD INCOMING ERR:", err)
		return
	}
	printTitle(s, idA, "A AFTER INCOMING EDIT")
	printTitle(s, idB, "B AFTER INCOMING EDIT")

	if err := s.UpdateIncomingRelations(idA, nil); err != nil {
		fmt.Println("DEL INCOMING ERR:", err)
		return
	}
	printTitle(s, idB, "B AFTER INCOMING DELETED")

	if err := s.DeleteTitle(idC); err != nil {
		fmt.Println("DEL C ERR:", err)
		return
	}
	printTitle(s, idB, "B AFTER C DELETED")

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
	if _, err := m2.ImportJSON(bytes.NewReader(buf.Bytes())); err != nil {
		fmt.Println("IMPORT ERR:", err)
		return
	}

	imported, err := s2.ListTitles(store.ListFilter{})
	if err != nil {
		fmt.Println("LIST DST ERR:", err)
		return
	}
	for _, it := range imported {
		fmt.Printf("IMPORTED id=%d name=%q own=%s reverse=%s\n", it.ID, it.Names[0].Value, relJSON(it.Relations), relJSON(it.ReverseRelations))
	}
}
