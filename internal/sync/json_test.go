package sync

import (
	"bytes"
	"path/filepath"
	"testing"

	"pergamon/internal/store"
)

// Импорт работает как слияние: уже существующие записи (по имени) обновляются,
// новые — добавляются, и все ссылки (связи, полки, персонажи, заметки)
// перемаппиваются на новые ID.
func TestJSONRoundTripRemapsIDs(t *testing.T) {
	src, err := store.Open(filepath.Join(t.TempDir(), "src.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	idA, err := src.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "original", Value: "Book A"}}, Score: 9.8})
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
	backup := append([]byte(nil), buf.Bytes()...)

	dst, err := store.Open(filepath.Join(t.TempDir(), "dst.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer dst.Close()

	// Заранее существующие записи: «Existing» сдвигает ID импортируемых книг,
	// локальная копия «Book A» должна обновиться, а не задублироваться.
	if _, err := dst.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "original", Value: "Existing"}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := dst.SaveTitle(store.Title{Type: "read", Category: "book", Names: []store.Name{{Kind: "original", Value: "Book A"}}, Score: 3}); err != nil {
		t.Fatal(err)
	}
	res, err := importJSON(dst, bytes.NewReader(backup))
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if res.Added != 4 || res.Updated != 1 {
		t.Fatalf("unexpected counts: %+v", res)
	}

	bookA, bookB := int64(0), int64(0)
	titles, err := dst.ListTitles(store.ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(titles) != 3 {
		t.Fatalf("want 3 titles, got %d", len(titles))
	}
	for _, tt := range titles {
		switch tt.Names[0].Value {
		case "Book A":
			bookA = tt.ID
		case "Book B":
			bookB = tt.ID
		}
	}
	if bookA != 2 || bookB != 3 {
		t.Fatalf("expected IDs 2 and 3, got %d and %d", bookA, bookB)
	}
	gotA, err := dst.GetTitle(bookA)
	if err != nil {
		t.Fatal(err)
	}
	if gotA.Score != 9.8 {
		t.Fatalf("existing title not updated with backup data: score %v", gotA.Score)
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

	// Повторный импорт того же бэкапа не должен ничего дублировать.
	res2, err := importJSON(dst, bytes.NewReader(backup))
	if err != nil {
		t.Fatalf("re-import: %v", err)
	}
	if res2.Added != 0 || res2.Updated != 4 {
		t.Fatalf("unexpected re-import counts: %+v", res2)
	}
	titles2, err := dst.ListTitles(store.ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(titles2) != 3 {
		t.Fatalf("re-import duplicated titles: %d", len(titles2))
	}
	shelves2, err := dst.ListShelves()
	if err != nil {
		t.Fatal(err)
	}
	if len(shelves2) != 1 || len(shelves2[0].TitleIDs) != 2 {
		t.Fatalf("re-import duplicated shelves: %+v", shelves2)
	}
	chars2, err := dst.ListCharacters("")
	if err != nil {
		t.Fatal(err)
	}
	if len(chars2) != 1 {
		t.Fatalf("re-import duplicated characters: %d", len(chars2))
	}
	notes2, err := dst.ListNotes(bookA)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes2) != 1 {
		t.Fatalf("re-import duplicated notes: %d", len(notes2))
	}
	_ = charID
}
