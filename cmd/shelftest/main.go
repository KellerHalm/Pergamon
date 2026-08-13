package main

import (
	"fmt"
	"os"
	"path/filepath"

	"mediateka/internal/store"
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

	sh := store.Shelf{
		Name:     "Моя полка",
		Kind:     "read",
		Position: 0,
		TitleIDs: []int64{},
		CreatedAt: "",
	}
	id, err := s.SaveShelf(sh)
	fmt.Println("ID:", id)
	fmt.Println("ERR:", err)

	if err == nil {
		list, lerr := s.ListShelves()
		fmt.Println("LIST count:", len(list), "ERR:", lerr)
		if len(list) > 0 {
			fmt.Printf("Shelf: %+v\n", list[0])
		}
	}
}
