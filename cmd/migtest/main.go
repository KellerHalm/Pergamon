package main

import (
	"encoding/json"
	"fmt"
	"os"

	"mediateka/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: migtest <dbpath>")
		os.Exit(1)
	}
	s, err := store.Open(os.Args[1])
	if err != nil {
		fmt.Println("OPEN ERR:", err)
		os.Exit(1)
	}
	defer s.Close()
	titles, _ := s.ListTitles(store.ListFilter{})
	for _, t := range titles {
		b1, _ := json.Marshal(t.People)
		b2, _ := json.Marshal(t.Studios)
		fmt.Printf("TITLE id=%d PEOPLE: %s STUDIOS: %s\n", t.ID, string(b1), string(b2))
	}
	people, _ := s.ListPeople("")
	studios, _ := s.ListStudios("")
	fmt.Println("PEOPLE:", len(people), "STUDIOS:", len(studios))
	for _, p := range people {
		b, _ := json.Marshal(p.Names)
		fmt.Printf("  PERSON id=%d role=%s names=%s\n", p.ID, p.Role, string(b))
	}
	for _, st := range studios {
		b, _ := json.Marshal(st.Names)
		fmt.Printf("  STUDIO id=%d names=%s\n", st.ID, string(b))
	}
}
