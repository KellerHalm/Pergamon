package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"pergamon/internal/store"
)

func main() {
	tmp, _ := os.MkdirTemp("", "shelf-test")
	defer os.RemoveAll(tmp)

	s, err := store.Open(filepath.Join(tmp, "test.db"))
	if err != nil {
		fmt.Println("OPEN ERR:", err)
		return
	}
	defer s.Close()

	empty, err := s.ListShelves()
	if err != nil {
		fmt.Println("EMPTY LIST ERR:", err)
		return
	}
	b, _ := json.Marshal(empty)
	fmt.Println("EMPTY LIST JSON:", string(b))
	fmt.Println("EMPTY LIST == nil:", empty == nil)

	id, err := s.SaveShelf(store.Shelf{Name: "Полка 1", Kind: "read", TitleIDs: []int64{}})
	fmt.Println("SAVE id:", id, "err:", err)

	list, err := s.ListShelves()
	if err != nil {
		fmt.Println("LIST ERR:", err)
		return
	}
	b2, _ := json.Marshal(list)
	fmt.Println("AFTER SAVE JSON:", string(b2))

	titles, err := s.ListTitles(store.ListFilter{})
	if err != nil {
		fmt.Println("TITLES ERR:", err)
		return
	}
	b3, _ := json.Marshal(titles)
	fmt.Println("TITLES EMPTY JSON:", string(b3))
}
