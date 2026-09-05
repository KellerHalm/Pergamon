package store

import (
	"database/sql"
	"sort"
	"strings"
)

type Name struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Creator struct {
	Role string `json:"role"`
	Name string `json:"name"`
}

type TitleRelation struct {
	RelatedID    int64  `json:"relatedId"`
	Label        string `json:"label"`
	ReverseLabel string `json:"reverseLabel"`
	Name         string `json:"name"`
	Cover        string `json:"cover"`
	Status       string `json:"status"`
}

type CharacterRef struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	MainImage string `json:"mainImage"`
}

type TitleRef struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Cover  string `json:"cover"`
	Status string `json:"status"`
}

type CharacterField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Character struct {
	ID          int64            `json:"id"`
	Names       []Name           `json:"names"`
	MainImage   string           `json:"mainImage"`
	Age         string           `json:"age"`
	Gender      string           `json:"gender"`
	Race        string           `json:"race"`
	Description string           `json:"description"`
	Fields      []CharacterField `json:"fields"`
	Images      []string         `json:"images"`
	Titles      []TitleRef       `json:"titles"`
	TitleIDs    []int64          `json:"titleIds"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

type StudioRef struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PersonRef struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
	Name string `json:"name"`
}

type Studio struct {
	ID          int64      `json:"id"`
	Names       []Name     `json:"names"`
	MainImage   string     `json:"mainImage"`
	Founded     string     `json:"founded"`
	Founders    []string   `json:"founders"`
	Description string     `json:"description"`
	Images      []string   `json:"images"`
	Titles      []TitleRef `json:"titles"`
	TitleIDs    []int64    `json:"titleIds"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

type Person struct {
	ID          int64      `json:"id"`
	Names       []Name     `json:"names"`
	MainImage   string     `json:"mainImage"`
	Age         string     `json:"age"`
	BirthDate   string     `json:"birthDate"`
	DeathDate   string     `json:"deathDate"`
	Gender      string     `json:"gender"`
	Role        string     `json:"role"`
	Description string     `json:"description"`
	Images      []string   `json:"images"`
	Titles      []TitleRef `json:"titles"`
	TitleIDs    []int64    `json:"titleIds"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

type Progress struct {
	Volumes       int `json:"volumes"`
	Chapters      int `json:"chapters"`
	Pages         int `json:"pages"`
	Seasons       int `json:"seasons"`
	Episodes      int `json:"episodes"`
	Minutes       int `json:"minutes"`
	TotalChapters int `json:"totalChapters"`
	TotalEpisodes int `json:"totalEpisodes"`
}

type Title struct {
	ID               int64           `json:"id"`
	Type             string          `json:"type"`
	Category         string          `json:"category"`
	Names            []Name          `json:"names"`
	Cover            string          `json:"cover"`
	Images           []string        `json:"images"`
	Synopsis         string          `json:"synopsis"`
	Creators         []Creator       `json:"creators"`
	Genres           []string        `json:"genres"`
	Tags             []string        `json:"tags"`
	Relations        []TitleRelation `json:"relations"`
	ReverseRelations []TitleRelation `json:"reverseRelations"`
	Characters       []CharacterRef  `json:"characters"`
	Studios          []StudioRef     `json:"studios"`
	People           []PersonRef     `json:"people"`
	Score            float64         `json:"score"`
	Status           string          `json:"status"`
	ReleaseStatus    string          `json:"releaseStatus"`
	CustomList       string          `json:"customList"`
	Progress         Progress        `json:"progress"`
	Notes            string          `json:"notes"`
	SpineColor       string          `json:"spineColor"`
	CreatedAt        string          `json:"createdAt"`
	UpdatedAt        string          `json:"updatedAt"`
}

type Note struct {
	ID        int64  `json:"id"`
	TitleID   int64  `json:"titleId"`
	Heading   string `json:"heading"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Shelf struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Kind      string  `json:"kind"`
	Position  int     `json:"position"`
	TitleIDs  []int64 `json:"titleIds"`
	CreatedAt string  `json:"createdAt"`
}

type ListFilter struct {
	Sort     string
	Type     string
	Category string
	Status   string
	Tags     []string
	Search   string
}

func (s *Store) ListTitles(f ListFilter) ([]Title, error) {
	q := `SELECT id FROM titles WHERE 1=1`
	args := []interface{}{}
	if f.Type != "" {
		q += ` AND type=?`
		args = append(args, f.Type)
	}
	if f.Category != "" {
		q += ` AND category=?`
		args = append(args, f.Category)
	}
	if f.Status != "" {
		q += ` AND status=?`
		args = append(args, f.Status)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q += ` AND id IN (SELECT title_id FROM title_names WHERE value LIKE ?)
		      OR id IN (SELECT title_id FROM title_creators WHERE name LIKE ?)
		      OR id IN (SELECT tp.title_id FROM title_people tp JOIN person_names pn ON pn.person_id=tp.person_id WHERE pn.value LIKE ?)
		      OR id IN (SELECT ts.title_id FROM title_studios ts JOIN studio_names sn ON sn.studio_id=ts.studio_id WHERE sn.value LIKE ?)`
		args = append(args, like, like, like, like)
	}
	if len(f.Tags) > 0 {
		ph := make([]string, len(f.Tags))
		for i, t := range f.Tags {
			ph[i] = "?"
			args = append(args, t)
		}
		q += ` AND id IN (SELECT title_id FROM title_tags WHERE tag IN (` + strings.Join(ph, ",") + `))`
	}

	switch f.Sort {
	case "title":
		q += ` ORDER BY (SELECT value FROM title_names WHERE title_id=titles.id ORDER BY kind LIMIT 1)`
	case "score":
		q += ` ORDER BY score DESC, updated_at DESC`
	case "updated":
		q += ` ORDER BY updated_at DESC`
	default:
		q += ` ORDER BY created_at DESC, id DESC`
	}

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()

	out := make([]Title, 0, len(ids))
	for _, id := range ids {
		t, err := s.GetTitle(id)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return out, nil
}

func (s *Store) GetTitle(id int64) (*Title, error) {
	t := &Title{Names: []Name{}, Creators: []Creator{}, Genres: []string{}, Tags: []string{}, Images: []string{},
		Relations: []TitleRelation{}, ReverseRelations: []TitleRelation{}, Characters: []CharacterRef{},
		Studios: []StudioRef{}, People: []PersonRef{}}
	err := s.db.QueryRow(`SELECT id,type,category,cover,synopsis,score,status,release_status,custom_list,
			progress_volumes,progress_chapters,progress_pages,progress_seasons,progress_episodes,progress_minutes,
			progress_total_chapters,progress_total_episodes,
			notes,spine_color,created_at,updated_at FROM titles WHERE id=?`, id).Scan(
		&t.ID, &t.Type, &t.Category, &t.Cover, &t.Synopsis, &t.Score, &t.Status, &t.ReleaseStatus, &t.CustomList,
		&t.Progress.Volumes, &t.Progress.Chapters, &t.Progress.Pages, &t.Progress.Seasons, &t.Progress.Episodes, &t.Progress.Minutes,
		&t.Progress.TotalChapters, &t.Progress.TotalEpisodes,
		&t.Notes, &t.SpineColor, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	nRows, err := s.db.Query(`SELECT kind,value FROM title_names WHERE title_id=? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	for nRows.Next() {
		var n Name
		if err := nRows.Scan(&n.Kind, &n.Value); err != nil {
			nRows.Close()
			return nil, err
		}
		t.Names = append(t.Names, n)
	}
	nRows.Close()

	cRows, err := s.db.Query(`SELECT role,name FROM title_creators WHERE title_id=? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	for cRows.Next() {
		var c Creator
		if err := cRows.Scan(&c.Role, &c.Name); err != nil {
			cRows.Close()
			return nil, err
		}
		t.Creators = append(t.Creators, c)
	}
	cRows.Close()

	gRows, err := s.db.Query(`SELECT genre FROM title_genres WHERE title_id=?`, id)
	if err != nil {
		return nil, err
	}
	for gRows.Next() {
		var g string
		if err := gRows.Scan(&g); err != nil {
			gRows.Close()
			return nil, err
		}
		t.Genres = append(t.Genres, g)
	}
	gRows.Close()

	tRows, err := s.db.Query(`SELECT tag FROM title_tags WHERE title_id=?`, id)
	if err != nil {
		return nil, err
	}
	for tRows.Next() {
		var tg string
		if err := tRows.Scan(&tg); err != nil {
			tRows.Close()
			return nil, err
		}
		t.Tags = append(t.Tags, tg)
	}
	tRows.Close()

	iRows, err := s.db.Query(`SELECT file FROM title_images WHERE title_id=? ORDER BY position, id`, id)
	if err != nil {
		return nil, err
	}
	for iRows.Next() {
		var f string
		if err := iRows.Scan(&f); err != nil {
			iRows.Close()
			return nil, err
		}
		t.Images = append(t.Images, f)
	}
	iRows.Close()

	rRows, err := s.db.Query(`SELECT r.related_id, r.label, r.reverse_label,
			COALESCE((SELECT status FROM titles WHERE id=r.related_id),''),
			COALESCE((SELECT value FROM title_names WHERE title_id=r.related_id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),''),
			COALESCE((SELECT cover FROM titles WHERE id=r.related_id),'')
		FROM title_relations r WHERE r.title_id=? ORDER BY r.id`, id)
	if err != nil {
		return nil, err
	}
	for rRows.Next() {
		var r TitleRelation
		if err := rRows.Scan(&r.RelatedID, &r.Label, &r.ReverseLabel, &r.Status, &r.Name, &r.Cover); err != nil {
			rRows.Close()
			return nil, err
		}
		t.Relations = append(t.Relations, r)
	}
	rRows.Close()

	vRows, err := s.db.Query(`SELECT r.title_id, COALESCE(NULLIF(r.reverse_label,''), r.label), r.label,
			COALESCE((SELECT status FROM titles WHERE id=r.title_id),''),
			COALESCE((SELECT value FROM title_names WHERE title_id=r.title_id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),''),
			COALESCE((SELECT cover FROM titles WHERE id=r.title_id),'')
		FROM title_relations r WHERE r.related_id=? ORDER BY r.id`, id)
	if err != nil {
		return nil, err
	}
	for vRows.Next() {
		var r TitleRelation
		if err := vRows.Scan(&r.RelatedID, &r.Label, &r.ReverseLabel, &r.Status, &r.Name, &r.Cover); err != nil {
			vRows.Close()
			return nil, err
		}
		t.ReverseRelations = append(t.ReverseRelations, r)
	}
	vRows.Close()

	chRows, err := s.db.Query(`SELECT c.id,
			COALESCE((SELECT value FROM character_names WHERE character_id=c.id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),''),
			COALESCE(c.main_image,'')
		FROM title_characters tc JOIN characters c ON c.id=tc.character_id
		WHERE tc.title_id=? ORDER BY tc.position, tc.rowid`, id)
	if err != nil {
		return nil, err
	}
	for chRows.Next() {
		var cr CharacterRef
		if err := chRows.Scan(&cr.ID, &cr.Name, &cr.MainImage); err != nil {
			chRows.Close()
			return nil, err
		}
		t.Characters = append(t.Characters, cr)
	}
	chRows.Close()

	pRows, err := s.db.Query(`SELECT p.id, p.role,
			COALESCE((SELECT value FROM person_names WHERE person_id=p.id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),'')
		FROM title_people tp JOIN people p ON p.id=tp.person_id
		WHERE tp.title_id=? ORDER BY tp.position, tp.rowid`, id)
	if err != nil {
		return nil, err
	}
	for pRows.Next() {
		var pr PersonRef
		if err := pRows.Scan(&pr.ID, &pr.Role, &pr.Name); err != nil {
			pRows.Close()
			return nil, err
		}
		t.People = append(t.People, pr)
	}
	pRows.Close()

	stRows, err := s.db.Query(`SELECT st.id,
			COALESCE((SELECT value FROM studio_names WHERE studio_id=st.id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),'')
		FROM title_studios ts JOIN studios st ON st.id=ts.studio_id
		WHERE ts.title_id=? ORDER BY ts.position, ts.rowid`, id)
	if err != nil {
		return nil, err
	}
	for stRows.Next() {
		var sr StudioRef
		if err := stRows.Scan(&sr.ID, &sr.Name); err != nil {
			stRows.Close()
			return nil, err
		}
		t.Studios = append(t.Studios, sr)
	}
	stRows.Close()

	return t, nil
}

func (s *Store) SaveTitle(t Title) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if t.ID == 0 {
		res, err := tx.Exec(`INSERT INTO titles(type,category,cover,synopsis,score,status,release_status,custom_list,
			progress_volumes,progress_chapters,progress_pages,progress_seasons,progress_episodes,progress_minutes,
			progress_total_chapters,progress_total_episodes,
			notes,spine_color,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
			t.Type, t.Category, t.Cover, t.Synopsis, t.Score, t.Status, t.ReleaseStatus, t.CustomList,
			t.Progress.Volumes, t.Progress.Chapters, t.Progress.Pages, t.Progress.Seasons, t.Progress.Episodes, t.Progress.Minutes,
			t.Progress.TotalChapters, t.Progress.TotalEpisodes,
			t.Notes, t.SpineColor)
		if err != nil {
			return 0, err
		}
		t.ID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		_, err := tx.Exec(`UPDATE titles SET type=?,category=?,cover=?,synopsis=?,score=?,status=?,release_status=?,custom_list=?,
			progress_volumes=?,progress_chapters=?,progress_pages=?,progress_seasons=?,progress_episodes=?,progress_minutes=?,
			progress_total_chapters=?,progress_total_episodes=?,
			notes=?,spine_color=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,
			t.Type, t.Category, t.Cover, t.Synopsis, t.Score, t.Status, t.ReleaseStatus, t.CustomList,
			t.Progress.Volumes, t.Progress.Chapters, t.Progress.Pages, t.Progress.Seasons, t.Progress.Episodes, t.Progress.Minutes,
			t.Progress.TotalChapters, t.Progress.TotalEpisodes,
			t.Notes, t.SpineColor, t.ID)
		if err != nil {
			return 0, err
		}
	}

	if _, err := tx.Exec(`DELETE FROM title_names WHERE title_id=?`, t.ID); err != nil {
		return 0, err
	}
	for _, n := range t.Names {
		if strings.TrimSpace(n.Value) == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO title_names(title_id,kind,value) VALUES(?,?,?)`, t.ID, n.Kind, n.Value); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM title_creators WHERE title_id=?`, t.ID); err != nil {
		return 0, err
	}
	if _, err := tx.Exec(`DELETE FROM title_studios WHERE title_id=?`, t.ID); err != nil {
		return 0, err
	}
	if _, err := tx.Exec(`DELETE FROM title_people WHERE title_id=?`, t.ID); err != nil {
		return 0, err
	}
	for _, c := range t.Creators {
		name := strings.TrimSpace(c.Name)
		if name == "" {
			continue
		}
		if c.Role == "studio" {
			sid, err := findOrCreateStudioTx(tx, name)
			if err != nil {
				return 0, err
			}
			if _, err := tx.Exec(`INSERT OR IGNORE INTO title_studios(title_id,studio_id,position) VALUES(?,?,0)`, t.ID, sid); err != nil {
				return 0, err
			}
		} else {
			pid, err := findOrCreatePersonTx(tx, name, c.Role)
			if err != nil {
				return 0, err
			}
			if _, err := tx.Exec(`INSERT OR IGNORE INTO title_people(title_id,person_id,position) VALUES(?,?,0)`, t.ID, pid); err != nil {
				return 0, err
			}
		}
	}
	for i, sr := range t.Studios {
		if sr.ID != 0 {
			if _, err := tx.Exec(`INSERT OR IGNORE INTO title_studios(title_id,studio_id,position)
				SELECT ?,?,? WHERE EXISTS(SELECT 1 FROM studios WHERE id=?)`, t.ID, sr.ID, i, sr.ID); err != nil {
				return 0, err
			}
			continue
		}
		name := strings.TrimSpace(sr.Name)
		if name == "" {
			continue
		}
		sid, err := findOrCreateStudioTx(tx, name)
		if err != nil {
			return 0, err
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_studios(title_id,studio_id,position) VALUES(?,?,?)`, t.ID, sid, i); err != nil {
			return 0, err
		}
	}
	for i, pr := range t.People {
		if pr.ID != 0 {
			if _, err := tx.Exec(`INSERT OR IGNORE INTO title_people(title_id,person_id,position)
				SELECT ?,?,? WHERE EXISTS(SELECT 1 FROM people WHERE id=?)`, t.ID, pr.ID, i, pr.ID); err != nil {
				return 0, err
			}
			continue
		}
		name := strings.TrimSpace(pr.Name)
		if name == "" {
			continue
		}
		pid, err := findOrCreatePersonTx(tx, name, pr.Role)
		if err != nil {
			return 0, err
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_people(title_id,person_id,position) VALUES(?,?,?)`, t.ID, pid, i); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM title_genres WHERE title_id=?`, t.ID); err != nil {
		return 0, err
	}
	for _, g := range t.Genres {
		g = strings.TrimSpace(g)
		if g == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_genres(title_id,genre) VALUES(?,?)`, t.ID, g); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM title_tags WHERE title_id=?`, t.ID); err != nil {
		return 0, err
	}
	for _, tg := range t.Tags {
		tg = strings.TrimSpace(tg)
		if tg == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_tags(title_id,tag) VALUES(?,?)`, t.ID, tg); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM title_images WHERE title_id=?`, t.ID); err != nil {
		return 0, err
	}
	for i, f := range t.Images {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO title_images(title_id,file,position) VALUES(?,?,?)`, t.ID, f, i); err != nil {
			return 0, err
		}
	}
	if err := saveRelations(tx, t.ID, t.Relations); err != nil {
		return 0, err
	}
	if err := saveTitleCharacters(tx, t.ID, t.Characters); err != nil {
		return 0, err
	}

	return t.ID, tx.Commit()
}

func saveTitleCharacters(tx *sql.Tx, titleID int64, chars []CharacterRef) error {
	if _, err := tx.Exec(`DELETE FROM title_characters WHERE title_id=?`, titleID); err != nil {
		return err
	}
	for i, c := range chars {
		if c.ID == 0 {
			name := strings.TrimSpace(c.Name)
			if name == "" {
				continue
			}
			var err error
			c.ID, err = findOrCreateCharacterTx(tx, name)
			if err != nil {
				return err
			}
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_characters(title_id,character_id,position)
			SELECT ?,?,? WHERE EXISTS(SELECT 1 FROM characters WHERE id=?)`, titleID, c.ID, i, c.ID); err != nil {
			return err
		}
	}
	return nil
}

func replaceCharacterTitles(tx *sql.Tx, characterID int64, titleIDs []int64) error {
	if _, err := tx.Exec(`DELETE FROM title_characters WHERE character_id=?`, characterID); err != nil {
		return err
	}
	for i, id := range titleIDs {
		if id == 0 {
			continue
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_characters(title_id,character_id,position)
			SELECT ?,?,? WHERE EXISTS(SELECT 1 FROM titles WHERE id=?)`, id, characterID, i, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListCharacters(sort string) ([]Character, error) {
	q := `SELECT id FROM characters`
	switch sort {
	case "name":
		q += ` ORDER BY (SELECT value FROM character_names WHERE character_id=characters.id ORDER BY kind LIMIT 1)`
	case "updated":
		q += ` ORDER BY updated_at DESC`
	default:
		q += ` ORDER BY created_at DESC, id DESC`
	}
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()

	out := make([]Character, 0, len(ids))
	for _, id := range ids {
		c, err := s.GetCharacter(id)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, nil
}

func (s *Store) GetCharacter(id int64) (*Character, error) {
	c := &Character{Names: []Name{}, Fields: []CharacterField{}, Images: []string{}, Titles: []TitleRef{}, TitleIDs: []int64{}}
	err := s.db.QueryRow(`SELECT id,COALESCE(main_image,''),COALESCE(age,''),COALESCE(gender,''),COALESCE(race,''),COALESCE(description,''),
			created_at,updated_at FROM characters WHERE id=?`, id).Scan(
		&c.ID, &c.MainImage, &c.Age, &c.Gender, &c.Race, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	nRows, err := s.db.Query(`SELECT kind,value FROM character_names WHERE character_id=? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	for nRows.Next() {
		var n Name
		if err := nRows.Scan(&n.Kind, &n.Value); err != nil {
			nRows.Close()
			return nil, err
		}
		c.Names = append(c.Names, n)
	}
	nRows.Close()

	fRows, err := s.db.Query(`SELECT name,value FROM character_fields WHERE character_id=? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	for fRows.Next() {
		var f CharacterField
		if err := fRows.Scan(&f.Name, &f.Value); err != nil {
			fRows.Close()
			return nil, err
		}
		c.Fields = append(c.Fields, f)
	}
	fRows.Close()

	iRows, err := s.db.Query(`SELECT file FROM character_images WHERE character_id=? ORDER BY position, id`, id)
	if err != nil {
		return nil, err
	}
	for iRows.Next() {
		var f string
		if err := iRows.Scan(&f); err != nil {
			iRows.Close()
			return nil, err
		}
		c.Images = append(c.Images, f)
	}
	iRows.Close()

	tRows, err := s.db.Query(`SELECT tc.title_id,
			COALESCE((SELECT value FROM title_names WHERE title_id=tc.title_id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),''),
			COALESCE((SELECT cover FROM titles WHERE id=tc.title_id),''),
			COALESCE((SELECT status FROM titles WHERE id=tc.title_id),'')
		FROM title_characters tc WHERE tc.character_id=? ORDER BY tc.position, tc.rowid`, id)
	if err != nil {
		return nil, err
	}
	for tRows.Next() {
		var tr TitleRef
		if err := tRows.Scan(&tr.ID, &tr.Name, &tr.Cover, &tr.Status); err != nil {
			tRows.Close()
			return nil, err
		}
		c.Titles = append(c.Titles, tr)
		c.TitleIDs = append(c.TitleIDs, tr.ID)
	}
	tRows.Close()

	return c, nil
}

func (s *Store) SaveCharacter(c Character) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if c.ID == 0 {
		res, err := tx.Exec(`INSERT INTO characters(main_image,age,gender,race,description,created_at,updated_at)
			VALUES(?,?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`, c.MainImage, c.Age, c.Gender, c.Race, c.Description)
		if err != nil {
			return 0, err
		}
		c.ID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		_, err := tx.Exec(`UPDATE characters SET main_image=?,age=?,gender=?,race=?,description=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,
			c.MainImage, c.Age, c.Gender, c.Race, c.Description, c.ID)
		if err != nil {
			return 0, err
		}
	}

	if _, err := tx.Exec(`DELETE FROM character_names WHERE character_id=?`, c.ID); err != nil {
		return 0, err
	}
	for _, n := range c.Names {
		if strings.TrimSpace(n.Value) == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO character_names(character_id,kind,value) VALUES(?,?,?)`, c.ID, n.Kind, n.Value); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM character_fields WHERE character_id=?`, c.ID); err != nil {
		return 0, err
	}
	for _, f := range c.Fields {
		if strings.TrimSpace(f.Name) == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO character_fields(character_id,name,value) VALUES(?,?,?)`, c.ID, strings.TrimSpace(f.Name), f.Value); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM character_images WHERE character_id=?`, c.ID); err != nil {
		return 0, err
	}
	for i, f := range c.Images {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO character_images(character_id,file,position) VALUES(?,?,?)`, c.ID, f, i); err != nil {
			return 0, err
		}
	}
	if c.TitleIDs != nil {
		if err := replaceCharacterTitles(tx, c.ID, c.TitleIDs); err != nil {
			return 0, err
		}
	}

	return c.ID, tx.Commit()
}

func (s *Store) SetCharacterTitles(characterID int64, titleIDs []int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := replaceCharacterTitles(tx, characterID, titleIDs); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) DeleteCharacter(id int64) error {
	_, err := s.db.Exec(`DELETE FROM characters WHERE id=?`, id)
	return err
}

func saveRelations(tx *sql.Tx, titleID int64, rels []TitleRelation) error {
	if _, err := tx.Exec(`DELETE FROM title_relations WHERE title_id=?`, titleID); err != nil {
		return err
	}
	for _, r := range rels {
		if r.RelatedID == 0 || r.RelatedID == titleID {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO title_relations(title_id,related_id,label,reverse_label) VALUES(?,?,?,?)`,
			titleID, r.RelatedID, strings.TrimSpace(r.Label), strings.TrimSpace(r.ReverseLabel)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ReplaceTitleRelations(titleID int64, rels []TitleRelation) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := saveRelations(tx, titleID, rels); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) UpdateIncomingRelations(titleID int64, rels []TitleRelation) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM title_relations WHERE related_id=?`, titleID); err != nil {
		return err
	}
	for _, r := range rels {
		if r.RelatedID == 0 || r.RelatedID == titleID {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO title_relations(title_id,related_id,label,reverse_label) VALUES(?,?,?,?)`,
			r.RelatedID, titleID, strings.TrimSpace(r.ReverseLabel), strings.TrimSpace(r.Label)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) DeleteTitle(id int64) error {
	_, err := s.db.Exec(`DELETE FROM titles WHERE id=?`, id)
	return err
}

func (s *Store) ListNotes(titleID int64) ([]Note, error) {
	rows, err := s.db.Query(`SELECT id,title_id,heading,content,created_at,updated_at
		FROM title_notes WHERE title_id=? ORDER BY created_at DESC, id DESC`, titleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.TitleID, &n.Heading, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (s *Store) ListAllNotes() ([]Note, error) {
	rows, err := s.db.Query(`SELECT id,title_id,heading,content,created_at,updated_at
		FROM title_notes ORDER BY title_id, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.TitleID, &n.Heading, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (s *Store) SaveNote(n Note) (int64, error) {
	if n.ID == 0 {
		res, err := s.db.Exec(`INSERT INTO title_notes(title_id,heading,content,created_at,updated_at)
			VALUES(?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`, n.TitleID, n.Heading, n.Content)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	_, err := s.db.Exec(`UPDATE title_notes SET heading=?,content=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		n.Heading, n.Content, n.ID)
	return n.ID, err
}

func (s *Store) DeleteNote(id int64) error {
	_, err := s.db.Exec(`DELETE FROM title_notes WHERE id=?`, id)
	return err
}

func (s *Store) AdjustProgress(id int64, field string, delta int) error {
	var col string
	switch field {
	case "volumes":
		col = "progress_volumes"
	case "chapters":
		col = "progress_chapters"
	case "pages":
		col = "progress_pages"
	case "seasons":
		col = "progress_seasons"
	case "episodes":
		col = "progress_episodes"
	case "minutes":
		col = "progress_minutes"
	case "totalChapters":
		col = "progress_total_chapters"
	case "totalEpisodes":
		col = "progress_total_episodes"
	default:
		return nil
	}
	_, err := s.db.Exec(`UPDATE titles SET `+col+`=MAX(0,`+col+`+?), updated_at=CURRENT_TIMESTAMP WHERE id=?`, delta, id)
	return err
}

func (s *Store) AllTags() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT tag FROM title_tags ORDER BY tag`)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, t)
	}
	rows.Close()
	sort.Strings(out)
	return out, nil
}

func (s *Store) AllGenres() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT genre FROM title_genres ORDER BY genre`)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for rows.Next() {
		var g string
		if err := rows.Scan(&g); err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, g)
	}
	rows.Close()
	sort.Strings(out)
	return out, nil
}

func replaceStudioTitles(tx *sql.Tx, studioID int64, titleIDs []int64) error {
	if _, err := tx.Exec(`DELETE FROM title_studios WHERE studio_id=?`, studioID); err != nil {
		return err
	}
	for i, id := range titleIDs {
		if id == 0 {
			continue
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_studios(title_id,studio_id,position)
			SELECT ?,?,? WHERE EXISTS(SELECT 1 FROM titles WHERE id=?)`, id, studioID, i, id); err != nil {
			return err
		}
	}
	return nil
}

func replacePersonTitles(tx *sql.Tx, personID int64, titleIDs []int64) error {
	if _, err := tx.Exec(`DELETE FROM title_people WHERE person_id=?`, personID); err != nil {
		return err
	}
	for i, id := range titleIDs {
		if id == 0 {
			continue
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO title_people(title_id,person_id,position)
			SELECT ?,?,? WHERE EXISTS(SELECT 1 FROM titles WHERE id=?)`, id, personID, i, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListStudios(sort string) ([]Studio, error) {
	q := `SELECT id FROM studios`
	switch sort {
	case "name":
		q += ` ORDER BY (SELECT value FROM studio_names WHERE studio_id=studios.id ORDER BY kind LIMIT 1)`
	case "updated":
		q += ` ORDER BY updated_at DESC`
	default:
		q += ` ORDER BY created_at DESC, id DESC`
	}
	return s.scanStudios(q)
}

func (s *Store) scanStudios(q string, args ...interface{}) ([]Studio, error) {
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	out := make([]Studio, 0, len(ids))
	for _, id := range ids {
		st, err := s.GetStudio(id)
		if err != nil {
			return nil, err
		}
		out = append(out, *st)
	}
	return out, nil
}

func (s *Store) GetStudio(id int64) (*Studio, error) {
	st := &Studio{Names: []Name{}, Founders: []string{}, Images: []string{}, Titles: []TitleRef{}, TitleIDs: []int64{}}
	err := s.db.QueryRow(`SELECT id,COALESCE(main_image,''),COALESCE(founded,''),COALESCE(description,''),created_at,updated_at
			FROM studios WHERE id=?`, id).Scan(
		&st.ID, &st.MainImage, &st.Founded, &st.Description, &st.CreatedAt, &st.UpdatedAt)
	if err != nil {
		return nil, err
	}

	nRows, err := s.db.Query(`SELECT kind,value FROM studio_names WHERE studio_id=? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	for nRows.Next() {
		var n Name
		if err := nRows.Scan(&n.Kind, &n.Value); err != nil {
			nRows.Close()
			return nil, err
		}
		st.Names = append(st.Names, n)
	}
	nRows.Close()

	fRows, err := s.db.Query(`SELECT name FROM studio_founders WHERE studio_id=? ORDER BY position, id`, id)
	if err != nil {
		return nil, err
	}
	for fRows.Next() {
		var f string
		if err := fRows.Scan(&f); err != nil {
			fRows.Close()
			return nil, err
		}
		st.Founders = append(st.Founders, f)
	}
	fRows.Close()

	iRows, err := s.db.Query(`SELECT file FROM studio_images WHERE studio_id=? ORDER BY position, id`, id)
	if err != nil {
		return nil, err
	}
	for iRows.Next() {
		var f string
		if err := iRows.Scan(&f); err != nil {
			iRows.Close()
			return nil, err
		}
		st.Images = append(st.Images, f)
	}
	iRows.Close()

	tRows, err := s.db.Query(`SELECT ts.title_id,
			COALESCE((SELECT value FROM title_names WHERE title_id=ts.title_id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),''),
			COALESCE((SELECT cover FROM titles WHERE id=ts.title_id),''),
			COALESCE((SELECT status FROM titles WHERE id=ts.title_id),'')
		FROM title_studios ts WHERE ts.studio_id=? ORDER BY ts.position, ts.rowid`, id)
	if err != nil {
		return nil, err
	}
	for tRows.Next() {
		var tr TitleRef
		if err := tRows.Scan(&tr.ID, &tr.Name, &tr.Cover, &tr.Status); err != nil {
			tRows.Close()
			return nil, err
		}
		st.Titles = append(st.Titles, tr)
		st.TitleIDs = append(st.TitleIDs, tr.ID)
	}
	tRows.Close()

	return st, nil
}

func (s *Store) SaveStudio(st Studio) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if st.ID == 0 {
		res, err := tx.Exec(`INSERT INTO studios(main_image,founded,description,created_at,updated_at)
			VALUES(?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`, st.MainImage, st.Founded, st.Description)
		if err != nil {
			return 0, err
		}
		st.ID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		_, err := tx.Exec(`UPDATE studios SET main_image=?,founded=?,description=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,
			st.MainImage, st.Founded, st.Description, st.ID)
		if err != nil {
			return 0, err
		}
	}

	if _, err := tx.Exec(`DELETE FROM studio_names WHERE studio_id=?`, st.ID); err != nil {
		return 0, err
	}
	for _, n := range st.Names {
		if strings.TrimSpace(n.Value) == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO studio_names(studio_id,kind,value) VALUES(?,?,?)`, st.ID, n.Kind, n.Value); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM studio_founders WHERE studio_id=?`, st.ID); err != nil {
		return 0, err
	}
	for i, f := range st.Founders {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO studio_founders(studio_id,name,position) VALUES(?,?,?)`, st.ID, f, i); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM studio_images WHERE studio_id=?`, st.ID); err != nil {
		return 0, err
	}
	for i, f := range st.Images {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO studio_images(studio_id,file,position) VALUES(?,?,?)`, st.ID, f, i); err != nil {
			return 0, err
		}
	}
	if st.TitleIDs != nil {
		if err := replaceStudioTitles(tx, st.ID, st.TitleIDs); err != nil {
			return 0, err
		}
	}

	return st.ID, tx.Commit()
}

func (s *Store) SetStudioTitles(studioID int64, titleIDs []int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := replaceStudioTitles(tx, studioID, titleIDs); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) DeleteStudio(id int64) error {
	_, err := s.db.Exec(`DELETE FROM studios WHERE id=?`, id)
	return err
}

func (s *Store) ListPeople(sort string) ([]Person, error) {
	q := `SELECT id FROM people`
	switch sort {
	case "name":
		q += ` ORDER BY (SELECT value FROM person_names WHERE person_id=people.id ORDER BY kind LIMIT 1)`
	case "updated":
		q += ` ORDER BY updated_at DESC`
	default:
		q += ` ORDER BY created_at DESC, id DESC`
	}
	return s.scanPeople(q)
}

func (s *Store) scanPeople(q string, args ...interface{}) ([]Person, error) {
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	out := make([]Person, 0, len(ids))
	for _, id := range ids {
		p, err := s.GetPerson(id)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, nil
}

func (s *Store) GetPerson(id int64) (*Person, error) {
	p := &Person{Names: []Name{}, Images: []string{}, Titles: []TitleRef{}, TitleIDs: []int64{}}
	err := s.db.QueryRow(`SELECT id,COALESCE(main_image,''),COALESCE(age,''),COALESCE(birth_date,''),COALESCE(death_date,''),
			COALESCE(gender,''),COALESCE(role,'author'),COALESCE(description,''),created_at,updated_at
			FROM people WHERE id=?`, id).Scan(
		&p.ID, &p.MainImage, &p.Age, &p.BirthDate, &p.DeathDate, &p.Gender, &p.Role, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	nRows, err := s.db.Query(`SELECT kind,value FROM person_names WHERE person_id=? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	for nRows.Next() {
		var n Name
		if err := nRows.Scan(&n.Kind, &n.Value); err != nil {
			nRows.Close()
			return nil, err
		}
		p.Names = append(p.Names, n)
	}
	nRows.Close()

	iRows, err := s.db.Query(`SELECT file FROM person_images WHERE person_id=? ORDER BY position, id`, id)
	if err != nil {
		return nil, err
	}
	for iRows.Next() {
		var f string
		if err := iRows.Scan(&f); err != nil {
			iRows.Close()
			return nil, err
		}
		p.Images = append(p.Images, f)
	}
	iRows.Close()

	tRows, err := s.db.Query(`SELECT tp.title_id,
			COALESCE((SELECT value FROM title_names WHERE title_id=tp.title_id
				ORDER BY CASE kind WHEN 'russian' THEN 0 WHEN 'english' THEN 1 ELSE 2 END, id LIMIT 1),''),
			COALESCE((SELECT cover FROM titles WHERE id=tp.title_id),''),
			COALESCE((SELECT status FROM titles WHERE id=tp.title_id),'')
		FROM title_people tp WHERE tp.person_id=? ORDER BY tp.position, tp.rowid`, id)
	if err != nil {
		return nil, err
	}
	for tRows.Next() {
		var tr TitleRef
		if err := tRows.Scan(&tr.ID, &tr.Name, &tr.Cover, &tr.Status); err != nil {
			tRows.Close()
			return nil, err
		}
		p.Titles = append(p.Titles, tr)
		p.TitleIDs = append(p.TitleIDs, tr.ID)
	}
	tRows.Close()

	return p, nil
}

func (s *Store) SavePerson(p Person) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if p.Role == "" {
		p.Role = "author"
	}
	if p.ID == 0 {
		res, err := tx.Exec(`INSERT INTO people(main_image,age,birth_date,death_date,gender,role,description,created_at,updated_at)
			VALUES(?,?,?,?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
			p.MainImage, p.Age, p.BirthDate, p.DeathDate, p.Gender, p.Role, p.Description)
		if err != nil {
			return 0, err
		}
		p.ID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		_, err := tx.Exec(`UPDATE people SET main_image=?,age=?,birth_date=?,death_date=?,gender=?,role=?,description=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,
			p.MainImage, p.Age, p.BirthDate, p.DeathDate, p.Gender, p.Role, p.Description, p.ID)
		if err != nil {
			return 0, err
		}
	}

	if _, err := tx.Exec(`DELETE FROM person_names WHERE person_id=?`, p.ID); err != nil {
		return 0, err
	}
	for _, n := range p.Names {
		if strings.TrimSpace(n.Value) == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO person_names(person_id,kind,value) VALUES(?,?,?)`, p.ID, n.Kind, n.Value); err != nil {
			return 0, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM person_images WHERE person_id=?`, p.ID); err != nil {
		return 0, err
	}
	for i, f := range p.Images {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO person_images(person_id,file,position) VALUES(?,?,?)`, p.ID, f, i); err != nil {
			return 0, err
		}
	}
	if p.TitleIDs != nil {
		if err := replacePersonTitles(tx, p.ID, p.TitleIDs); err != nil {
			return 0, err
		}
	}

	return p.ID, tx.Commit()
}

func (s *Store) SetPersonTitles(personID int64, titleIDs []int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := replacePersonTitles(tx, personID, titleIDs); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) DeletePerson(id int64) error {
	_, err := s.db.Exec(`DELETE FROM people WHERE id=?`, id)
	return err
}
