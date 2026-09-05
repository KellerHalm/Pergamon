package main

import (
	"fmt"
	"os"

	"mediateka/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: legacydb <dbpath>")
		os.Exit(1)
	}
	s, err := store.Open(os.Args[1])
	if err != nil {
		fmt.Println("OPEN ERR:", err)
		os.Exit(1)
	}
	defer s.Close()
	id1, err := s.SaveTitle(store.Title{
		Type: "read", Category: "manga",
		Names:    []store.Name{{Kind: "russian", Value: "Манга А"}},
		Creators: []store.Creator{{Role: "author", Name: "Араки"}, {Role: "studio", Name: "WIT"}},
	})
	fmt.Println("SAVE1:", id1, err)
	id2, err := s.SaveTitle(store.Title{
		Type: "watch", Category: "anime",
		Names:    []store.Name{{Kind: "english", Value: "Anime B"}},
		Creators: []store.Creator{{Role: "director", Name: "Миядзаки"}, {Role: "author", Name: "араки"}, {Role: "studio", Name: "wit"}},
	})
	fmt.Println("SAVE2:", id2, err)
}
