package sync

import (
	"bytes"
	"path/filepath"
	"testing"

	"pergamon/internal/store"
)

// Импорт в непустую базу сдвигает ID: записи из бэкапа получают новые ID,
// и все ссылки (связи, полки, персонажи, заметки) должны перемаппиться.
func TestJSONRoundTripRemapsIDs(t *testing.T) {
	src, err := store.Open(filepath.Join(t.TempDir(), "src.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	idA, err := src.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "original", Value: "Book A"}}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = src.SaveTitle(store.Title{
		Type:      "read",
		Category:  "book",
		Names:     []store.Name{{Kind: "original", Value: "Book B"}},
		Relations: []store.TitleRelation{{RelatedID: idA, Label: "сиквел"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	charID, err := src.SaveCharacter(store.Character{Names: []store.Name{{Kind: "original", Value: "Hero"}}, TitleIDs: []int64{idA}})
	if err != nil {
		t.Fatal(err)
	}
	shelfID, err := src.SaveShelf(store.Shelf{Name: "Прочитано", TitleIDs: []int64{idA}})
	if err != nil {
		t.Fatal(err)
	}
	if err := src.SetShelfItems(shelfID, []int64{idA, idA + 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := src.SaveNote(store.Note{TitleID: idA, Heading: "Заметка", Content: "текст"}); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := exportJSON(src, &buf); err != nil {
		t.Fatal(err)
	}

	dst, err := store.Open(filepath.Join(t.TempDir(), "dst.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer dst.Close()

	// Заранее существующая запись сдвигает ID импортированных книг с 1,2 на 2,3.
	if _, err := dst.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "original", Value: "Existing"}}}); err != nil {
		t.Fatal(err)
	}
	if err := importJSON(dst, &buf); err != nil {
		t.Fatalf("import: %v", err)
	}

	titles, err := dst.ListTitles(store.ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(titles) != 3 {
		t.Fatalf("want 3 titles, got %d", len(titles))
	}
	bookA, bookB := int64(0), int64(0)
	for _, tt := range titles {
		switch tt.Names[0].Value {
		case "Book A":
			bookA = tt.ID
		case "Book B":
			bookB = tt.ID
		}
	}
	if bookA != 2 || bookB != 3 {
		t.Fatalf("expected imported IDs 2 and 3, got %d and %d", bookA, bookB)
	}
	gotB, err := dst.GetTitle(bookB)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotB.Relations) != 1 || gotB.Relations[0].RelatedID != bookA {
		t.Fatalf("relation not remapped: %+v", gotB.Relations)
	}

	shelves, err := dst.ListShelves()
	if err != nil {
		t.Fatal(err)
	}
	if len(shelves) != 1 || len(shelves[0].TitleIDs) != 2 || shelves[0].TitleIDs[0] != bookA || shelves[0].TitleIDs[1] != bookB {
		t.Fatalf("shelf items not remapped: %+v", shelves)
	}

	chars, err := dst.ListCharacters("")
	if err != nil {
		t.Fatal(err)
	}
	if len(chars) != 1 || len(chars[0].TitleIDs) != 1 || chars[0].TitleIDs[0] != bookA {
		t.Fatalf("character titles not remapped: %+v", chars)
	}

	notes, err := dst.ListNotes(bookA)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 || notes[0].Heading != "Заметка" {
		t.Fatalf("note not remapped: %+v", notes)
	}
	_ = charID
}
