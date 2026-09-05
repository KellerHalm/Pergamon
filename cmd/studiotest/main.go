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

func titleName(s *store.Store, id int64) string {
	t, err := s.GetTitle(id)
	if err != nil {
		return fmt.Sprintf("ERR:%v", err)
	}
	if len(t.Names) == 0 {
		return "?"
	}
	return t.Names[0].Value
}

func printStudio(s *store.Store, id int64, tag string) {
	st, err := s.GetStudio(id)
	if err != nil {
		fmt.Println(tag, "GET ERR:", err)
		return
	}
	b, _ := json.Marshal(st)
	fmt.Println(tag, string(b))
}

func printPerson(s *store.Store, id int64, tag string) {
	p, err := s.GetPerson(id)
	if err != nil {
		fmt.Println(tag, "GET ERR:", err)
		return
	}
	b, _ := json.Marshal(p)
	fmt.Println(tag, string(b))
}

func printTitleRefs(s *store.Store, id int64, tag string) {
	t, err := s.GetTitle(id)
	if err != nil {
		fmt.Println(tag, "GET ERR:", err)
		return
	}
	b1, _ := json.Marshal(t.People)
	b2, _ := json.Marshal(t.Studios)
	fmt.Printf("%s PEOPLE: %s STUDIOS: %s\n", tag, string(b1), string(b2))
}

func main() {
	tmp, _ := os.MkdirTemp("", "studio-test")
	defer os.RemoveAll(tmp)

	s, err := store.Open(filepath.Join(tmp, "test.db"))
	if err != nil {
		fmt.Println("OPEN ERR:", err)
		return
	}
	defer s.Close()

	idA, err := s.SaveTitle(store.Title{Type: "read", Category: "manga", Names: []store.Name{{Kind: "russian", Value: "Манга А"}}})
	fmt.Println("SAVE A:", idA, err)
	idB, err := s.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "russian", Value: "Книга Б"}}})
	fmt.Println("SAVE B:", idB, err)

	_, err = s.SaveTitle(store.Title{
		ID:     idA,
		Type:   "read", Category: "manga",
		Names:  []store.Name{{Kind: "russian", Value: "Манга А"}},
		Creators: []store.Creator{
			{Role: "author", Name: "Араки"},
			{Role: "studio", Name: "WIT"},
			{Role: "director", Name: "Миядзаки"},
		},
	})
	fmt.Println("SAVE A LEGACY CREATORS:", err)
	printTitleRefs(s, idA, "A AFTER LEGACY MIGRATION VIA SAVE")
	people, _ := s.ListPeople("")
	studios, _ := s.ListStudios("")
	fmt.Println("PEOPLE COUNT:", len(people), "STUDIOS COUNT:", len(studios))

	_, err = s.SaveTitle(store.Title{
		ID:     idB,
		Type:   "read", Category: "book",
		Names:  []store.Name{{Kind: "russian", Value: "Книга Б"}},
		Studios: []store.StudioRef{{Name: "wit"}},
		People: []store.PersonRef{
			{Name: "араки", Role: "author"},
			{Name: "Новый Художник", Role: "artist"},
		},
	})
	fmt.Println("SAVE B WITH NAME REFS (dedupe check):", err)
	printTitleRefs(s, idB, "B")
	people, _ = s.ListPeople("")
	studios, _ = s.ListStudios("")
	fmt.Println("PEOPLE COUNT:", len(people), "STUDIOS COUNT:", len(studios))

	sid := studios[0].ID
	pid := people[0].ID
	printStudio(s, sid, "STUDIO0")
	printPerson(s, pid, "PERSON0")

	savedStudioID, err := s.SaveStudio(store.Studio{
		Names:       []store.Name{{Kind: "russian", Value: "Студия Тест"}},
		Founded:     "1986-12-24",
		Founders:    []string{"Основатель Один", "Основатель Два"},
		Description: "Описание студии",
		TitleIDs:    []int64{idA, idB},
	})
	fmt.Println("SAVE STUDIO:", savedStudioID, err)
	printStudio(s, savedStudioID, "STUDIO SAVED")

	savedPersonID, err := s.SavePerson(store.Person{
		Names:       []store.Name{{Kind: "original", Value: "Test Person"}, {Kind: "russian", Value: "Тестовый Деятель"}},
		Age:         "54 года",
		BirthDate:   "1970-05-05",
		DeathDate:   "2020-01-01",
		Gender:      "male",
		Role:        "director",
		Description: "Описание деятеля",
		TitleIDs:    []int64{idB},
	})
	fmt.Println("SAVE PERSON:", savedPersonID, err)
	printPerson(s, savedPersonID, "PERSON SAVED")

	titles, err := s.ListTitles(store.ListFilter{Search: "Араки"})
	fmt.Println("SEARCH 'Араки' COUNT:", len(titles), err)
	titles, _ = s.ListTitles(store.ListFilter{Search: "WIT"})
	fmt.Println("SEARCH 'WIT' COUNT:", len(titles), err)

	idC, err := s.SaveTitle(store.Title{
		Type: "read", Category: "manga",
		Names:      []store.Name{{Kind: "russian", Value: "Манга В"}},
		Characters: []store.CharacterRef{{Name: "Джонатан Джостар"}},
	})
	fmt.Println("SAVE C WITH NEW CHARACTER:", idC, err)
	printCChars := func() {
		t, _ := s.GetTitle(idC)
		b, _ := json.Marshal(t.Characters)
		fmt.Println("C CHARACTERS:", string(b))
	}
	printCChars()
	characters, _ := s.ListCharacters("")
	fmt.Println("CHARACTERS COUNT:", len(characters))

	_, err = s.SaveTitle(store.Title{
		ID:         idC,
		Type:       "read", Category: "manga",
		Names:      []store.Name{{Kind: "russian", Value: "Манга В"}},
		Characters: []store.CharacterRef{{Name: "джонатан джостар"}, {Name: "Брендовый Герой"}},
	})
	fmt.Println("SAVE C DEDUPE + SECOND CHAR:", err)
	printCChars()
	characters, _ = s.ListCharacters("")
	fmt.Println("CHARACTERS COUNT:", len(characters))

	autoChar := characters[0]
	autoCharID := autoChar.ID
	autoChar.Names = []store.Name{{Kind: "original", Value: "Jonathan Joestar"}, {Kind: "russian", Value: "Джонатан Джостар"}}
	autoChar.Age = "20 лет"
	autoChar.Gender = "male"
	autoChar.Description = "Джентльмен."
	_, err = s.SaveCharacter(autoChar)
	fmt.Println("FULL EDIT AUTO CHARACTER:", err)
	if c, err := s.GetCharacter(autoCharID); err != nil {
		fmt.Println("GET AUTO CHAR ERR:", err)
	} else {
		b, _ := json.Marshal(c)
		fmt.Printf("AUTO CHAR FULL: %s\n", string(b))
	}

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
	impStudios, _ := s2.ListStudios("")
	impPeople, _ := s2.ListPeople("")
	impTitles, _ := s2.ListTitles(store.ListFilter{})
	for _, st := range impStudios {
		fmt.Printf("IMPORTED STUDIO id=%d names=%d titles=%d founders=%d\n", st.ID, len(st.Names), len(st.Titles), len(st.Founders))
	}
	for _, p := range impPeople {
		fmt.Printf("IMPORTED PERSON id=%d names=%d role=%s titles=%d\n", p.ID, len(p.Names), p.Role, len(p.Titles))
	}
	for _, t := range impTitles {
		fmt.Printf("IMPORTED TITLE id=%d name=%s people=%d studios=%d\n", t.ID, titleName(s2, t.ID), len(t.People), len(t.Studios))
	}

	if err := s.DeleteStudio(savedStudioID); err != nil {
		fmt.Println("DEL STUDIO ERR:", err)
	}
	if _, err := s.GetStudio(savedStudioID); err != nil {
		fmt.Println("STUDIO DELETED OK")
	}
	printTitleRefs(s, idA, "A AFTER STUDIO DELETE")
	if err := s.DeletePerson(savedPersonID); err != nil {
		fmt.Println("DEL PERSON ERR:", err)
	}
	printTitleRefs(s, idB, "B AFTER PERSON DELETE")
}
