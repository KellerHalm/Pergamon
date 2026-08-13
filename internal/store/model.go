package store

import (
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

type Progress struct {
	Volumes int `json:"volumes"`
	Chapters int `json:"chapters"`
	Pages    int `json:"pages"`
	Seasons  int `json:"seasons"`
	Episodes int `json:"episodes"`
	Minutes  int `json:"minutes"`
}

type Title struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	Category   string    `json:"category"`
	Names      []Name    `json:"names"`
	Cover      string    `json:"cover"`
	Synopsis   string    `json:"synopsis"`
	Creators   []Creator `json:"creators"`
	Genres     []string  `json:"genres"`
	Tags       []string  `json:"tags"`
	Score      float64   `json:"score"`
	Status     string    `json:"status"`
	CustomList string    `json:"customList"`
	Progress   Progress  `json:"progress"`
	Notes      string    `json:"notes"`
	SpineColor string `json:"spineColor"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type Shelf struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Kind      string  `json:"kind"`
	Position  int     `json:"position"`
	TitleIDs  []int64 `json:"titleIds"`
	CreatedAt string `json:"createdAt"`
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
		      OR id IN (SELECT title_id FROM title_creators WHERE name LIKE ?)`
		args = append(args, like, like)
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
	t := &Title{Names: []Name{}, Creators: []Creator{}, Genres: []string{}, Tags: []string{}}
	err := s.db.QueryRow(`SELECT id,type,category,cover,synopsis,score,status,custom_list,
			progress_volumes,progress_chapters,progress_pages,progress_seasons,progress_episodes,progress_minutes,
			notes,spine_color,created_at,updated_at FROM titles WHERE id=?`, id).Scan(
			&t.ID, &t.Type, &t.Category, &t.Cover, &t.Synopsis, &t.Score, &t.Status, &t.CustomList,
			&t.Progress.Volumes, &t.Progress.Chapters, &t.Progress.Pages, &t.Progress.Seasons, &t.Progress.Episodes, &t.Progress.Minutes,
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

	return t, nil
}

func (s *Store) SaveTitle(t Title) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if t.ID == 0 {
		res, err := tx.Exec(`INSERT INTO titles(type,category,cover,synopsis,score,status,custom_list,
			progress_volumes,progress_chapters,progress_pages,progress_seasons,progress_episodes,progress_minutes,
			notes,spine_color,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
			t.Type, t.Category, t.Cover, t.Synopsis, t.Score, t.Status, t.CustomList,
			t.Progress.Volumes, t.Progress.Chapters, t.Progress.Pages, t.Progress.Seasons, t.Progress.Episodes, t.Progress.Minutes,
			t.Notes, t.SpineColor)
		if err != nil {
			return 0, err
		}
		t.ID, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else {
		_, err := tx.Exec(`UPDATE titles SET type=?,category=?,cover=?,synopsis=?,score=?,status=?,custom_list=?,
			progress_volumes=?,progress_chapters=?,progress_pages=?,progress_seasons=?,progress_episodes=?,progress_minutes=?,
			notes=?,spine_color=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,
			t.Type, t.Category, t.Cover, t.Synopsis, t.Score, t.Status, t.CustomList,
			t.Progress.Volumes, t.Progress.Chapters, t.Progress.Pages, t.Progress.Seasons, t.Progress.Episodes, t.Progress.Minutes,
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
	for _, c := range t.Creators {
		if strings.TrimSpace(c.Name) == "" {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO title_creators(title_id,role,name) VALUES(?,?,?)`, t.ID, c.Role, c.Name); err != nil {
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

	return t.ID, tx.Commit()
}

func (s *Store) DeleteTitle(id int64) error {
	_, err := s.db.Exec(`DELETE FROM titles WHERE id=?`, id)
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
