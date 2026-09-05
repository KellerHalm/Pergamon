package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"mediateka/internal/store"
	"mediateka/internal/sync"
)

func check(cond bool, msg string) {
	if !cond {
		panic("FAIL: " + msg)
	}
	fmt.Println("ok:", msg)
}

func main() {
	dir, _ := os.MkdirTemp("", "chartest")
	defer os.RemoveAll(dir)

	s, err := store.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		panic(err)
	}
	defer s.Close()

	t1, err := s.SaveTitle(store.Title{
		Type:  "watch",
		Names: []store.Name{{Kind: "russian", Value: "Аниме Один"}, {Kind: "original", Value: "Anime One"}},
	})
	check(err == nil && t1 > 0, "save title 1")

	t2, err := s.SaveTitle(store.Title{
		Type:  "read",
		Names: []store.Name{{Kind: "russian", Value: "Манга Два"}},
	})
	check(err == nil && t2 > 0, "save title 2")

	c1, err := s.SaveCharacter(store.Character{
		Names:       []store.Name{{Kind: "russian", Value: "Алиса"}, {Kind: "original", Value: "Alice"}},
		Age:         "17",
		Gender:      "female",
		Race:        "Эльф",
		Description: "Главная героиня",
		Fields:      []store.CharacterField{{Name: "Цвет глаз", Value: "зелёный"}},
		TitleIDs:    []int64{t1, t2},
	})
	check(err == nil && c1 > 0, "save character 1")

	c2, err := s.SaveCharacter(store.Character{
		Names: []store.Name{{Kind: "russian", Value: "Борис"}},
		Age:   "300",
	})
	check(err == nil && c2 > 0, "save character 2")

	got, err := s.GetCharacter(c1)
	check(err == nil, "get character")
	check(got.Age == "17" && got.Description == "Главная героиня", "character fields")
	check(got.Gender == "female" && got.Race == "Эльф", "gender and race")
	check(len(got.Fields) == 1 && got.Fields[0].Name == "Цвет глаз" && got.Fields[0].Value == "зелёный", "custom fields")
	check(len(got.Names) == 2 && got.Names[0].Value == "Алиса", "character names")
	check(len(got.TitleIDs) == 2 && got.TitleIDs[0] == t1, "character title ids")
	check(len(got.Titles) == 2 && got.Titles[0].Name == "Аниме Один", "character title refs")

	tg, err := s.GetTitle(t1)
	check(err == nil, "get title 1")
	check(len(tg.Characters) == 1, "title 1 has 1 character")
	check(tg.Characters[0].ID == c1 && tg.Characters[0].Name == "Алиса", "title character ref")

	tg2, err := s.GetTitle(t2)
	check(len(tg2.Characters) == 1 && tg2.Characters[0].ID == c1, "title 2 character ref")

	byName, _ := s.ListCharacters("name")
	check(len(byName) == 2 && byName[0].Names[0].Value == "Алиса", "sort by name")
	byCreated, _ := s.ListCharacters("created")
	check(len(byCreated) == 2 && byCreated[0].ID == c2, "sort by created")

	tg2.Characters = nil
	_, err = s.SaveTitle(*tg2)
	check(err == nil, "save title 2 without characters")
	got, _ = s.GetCharacter(c1)
	check(len(got.TitleIDs) == 1 && got.TitleIDs[0] == t1, "link removed from character side")

	upd, err := s.SaveCharacter(store.Character{
		ID:       c1,
		Names:    []store.Name{{Kind: "russian", Value: "Алиса Updated"}},
		Age:      "18",
		TitleIDs: []int64{t1},
	})
	check(err == nil && upd == c1, "update character")
	got, _ = s.GetCharacter(c1)
	check(got.Age == "18" && len(got.Names) == 1 && got.Names[0].Value == "Алиса Updated", "update fields")
	check(got.Gender == "" && got.Race == "" && len(got.Fields) == 0, "gender/race/fields cleared on update")

	img := filepath.Join(dir, "img1.png")
	_ = os.WriteFile(img, []byte("fake"), 0644)
	_, _ = s.SaveCharacter(store.Character{
		ID:        c1,
		Names:     []store.Name{{Kind: "russian", Value: "Алиса Updated"}},
		MainImage: "img1.png",
		Images:    []string{"img1.png"},
		Gender:    "female",
		Race:      "Эльф",
		Fields:    []store.CharacterField{{Name: "Цвет глаз", Value: "зелёный"}, {Name: "  ", Value: "ignored"}},
	})
	got, _ = s.GetCharacter(c1)
	check(got.MainImage == "img1.png" && len(got.Images) == 1, "images saved")
	check(got.Gender == "female" && got.Race == "Эльф", "gender/race re-saved")
	check(len(got.Fields) == 1 && got.Fields[0].Name == "Цвет глаз", "blank-name field skipped")

	m := sync.NewManager(filepath.Join(dir, "test.db"), s)
	var buf bytes.Buffer
	check(m.ExportJSON(&buf) == nil, "export json")

	s2, err := store.Open(filepath.Join(dir, "import.db"))
	check(err == nil, "open second db")
	defer s2.Close()
	m2 := sync.NewManager(filepath.Join(dir, "import.db"), s2)
	check(m2.ImportJSON(&buf) == nil, "import json")

	chars, _ := s2.ListCharacters("")
	titles, _ := s2.ListTitles(store.ListFilter{})
	check(len(chars) == 2 && len(titles) == 2, "imported counts")
	var impC store.Character
	for _, c := range chars {
		if len(c.Names) > 0 && c.Names[0].Value == "Алиса Updated" {
			impC = c
		}
	}
	check(impC.ID > 0 && len(impC.TitleIDs) == 1, "imported character keeps link")
	impT, _ := s2.GetTitle(impC.TitleIDs[0])
	check(len(impT.Characters) == 1 && impT.Characters[0].ID == impC.ID, "imported title references imported character")

	check(s.DeleteCharacter(c1) == nil, "delete character")
	tg, _ = s.GetTitle(t1)
	check(len(tg.Characters) == 0, "title characters cleared after delete")

	fmt.Println("ALL PASSED")
}
