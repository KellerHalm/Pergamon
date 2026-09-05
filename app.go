package main

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mediateka/internal/cover"
	"mediateka/internal/store"
	"mediateka/internal/sync"
)

type App struct {
	ctx      context.Context
	store    *store.Store
	syncer   *sync.Manager
	dataDir  string
	coverDir string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir, err := userDataDir("Mediateka")
	if err != nil {
		dataDir = "."
	}
	a.dataDir = dataDir
	a.coverDir = filepath.Join(dataDir, "covers")

	s, err := store.Open(filepath.Join(dataDir, "mediateka.db"))
	if err != nil {
		panic(err)
	}
	a.store = s
	a.syncer = sync.NewManager(filepath.Join(dataDir, "mediateka.db"), s)

	go func() {
		auto, _ := s.GetSetting("sync_on_startup")
		if auto == "1" {
			_, _ = a.SyncNow()
		}
	}()
}

func (a *App) shutdown(ctx context.Context) {
	auto, _ := a.store.GetSetting("sync_on_shutdown")
	if auto == "1" {
		_, _ = a.SyncNow()
	}
	a.store.Close()
}

type ListQuery struct {
	Sort     string   `json:"sort"`
	Type     string   `json:"type"`
	Category string   `json:"category"`
	Status   string   `json:"status"`
	Tags     []string `json:"tags"`
	Search   string   `json:"search"`
}

func (a *App) ListTitles(q ListQuery) ([]store.Title, error) {
	return a.store.ListTitles(store.ListFilter{
		Sort:     q.Sort,
		Type:     q.Type,
		Category: q.Category,
		Status:   q.Status,
		Tags:     q.Tags,
		Search:   q.Search,
	})
}

func (a *App) GetTitle(id int64) (*store.Title, error) {
	return a.store.GetTitle(id)
}

func (a *App) SaveTitle(t store.Title) (int64, error) {
	if t.Cover != "" && t.SpineColor == "" {
		t.SpineColor = cover.DominantColorHex(a.coverDir, t.Cover)
	}
	return a.store.SaveTitle(t)
}

func (a *App) UpdateIncomingRelations(titleID int64, rels []store.TitleRelation) error {
	return a.store.UpdateIncomingRelations(titleID, rels)
}

func (a *App) DeleteTitle(id int64) error {
	return a.store.DeleteTitle(id)
}

func (a *App) AdjustProgress(id int64, field string, delta int) (*store.Title, error) {
	if err := a.store.AdjustProgress(id, field, delta); err != nil {
		return nil, err
	}
	return a.store.GetTitle(id)
}

func (a *App) ListNotes(titleID int64) ([]store.Note, error) {
	return a.store.ListNotes(titleID)
}

func (a *App) SaveNote(n store.Note) (int64, error) {
	return a.store.SaveNote(n)
}

func (a *App) DeleteNote(id int64) error {
	return a.store.DeleteNote(id)
}

func (a *App) AllTags() ([]string, error) {
	return a.store.AllTags()
}

func (a *App) AllGenres() ([]string, error) {
	return a.store.AllGenres()
}

func (a *App) ListShelves() ([]store.Shelf, error) {
	return a.store.ListShelves()
}

func (a *App) SaveShelf(s store.Shelf) (int64, error) {
	return a.store.SaveShelf(s)
}

func (a *App) DeleteShelf(id int64) error {
	return a.store.DeleteShelf(id)
}

func (a *App) SetShelfItems(shelfID int64, titleIDs []int64) error {
	return a.store.SetShelfItems(shelfID, titleIDs)
}

func (a *App) GetSetting(key string) string {
	v, _ := a.store.GetSetting(key)
	return v
}

func (a *App) SetSetting(key, value string) error {
	return a.store.SetSetting(key, value)
}

type CoverUpload struct {
	DataBase64 string `json:"dataBase64"`
	Ext        string `json:"ext"`
}

func (a *App) UploadCoverDataURL(dataURL string) (string, error) {
	if dataURL == "" {
		return "", nil
	}
	comma := strings.Index(dataURL, ",")
	if comma < 0 {
		return "", nil
	}
	head := dataURL[:comma]
	body := dataURL[comma+1:]
	ext := ".png"
	if strings.Contains(head, "image/jpeg") {
		ext = ".jpg"
	} else if strings.Contains(head, "image/gif") {
		ext = ".gif"
	} else if strings.Contains(head, "image/webp") {
		ext = ".webp"
	}
	data, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return "", err
	}
	return cover.SaveFromBytes(a.coverDir, data, ext)
}

func (a *App) UploadCoverFromURL(url string) (string, error) {
	return cover.SaveFromURL(a.coverDir, url)
}

func (a *App) CoverDataURL(filename string) (string, error) {
	data, err := cover.ReadBytes(a.coverDir, filename)
	if err != nil || data == nil {
		return "", nil
	}
	mime := "image/png"
	switch {
	case strings.HasSuffix(strings.ToLower(filename), ".jpg"), strings.HasSuffix(strings.ToLower(filename), ".jpeg"):
		mime = "image/jpeg"
	case strings.HasSuffix(strings.ToLower(filename), ".gif"):
		mime = "image/gif"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

func (a *App) CoverColor(filename string) string {
	return cover.DominantColorHex(a.coverDir, filename)
}

type SyncResult = sync.SyncResult

func (a *App) SyncNow() (SyncResult, error) {
	return a.syncer.SyncNow()
}

func (a *App) TestWebDAV(cfg sync.WebDAVConfig) error {
	return a.syncer.TestConnection(cfg)
}

func (a *App) SaveWebDAVConfig(cfg sync.WebDAVConfig) error {
	return a.syncer.SaveConfig(cfg)
}

func (a *App) GetWebDAVConfig() (sync.WebDAVConfig, error) {
	return a.syncer.GetConfig()
}

func (a *App) ExportJSON() (string, error) {
	path := filepath.Join(a.dataDir, "mediateka-export-"+time.Now().Format("20060102-150405")+".json")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err := a.syncer.ExportJSON(f); err != nil {
		return "", err
	}
	return path, nil
}

func (a *App) ImportJSONPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return a.syncer.ImportJSON(f)
}

func (a *App) ExportSQLitePath(dest string) error {
	return a.syncer.ExportSQLite(dest)
}

func (a *App) ImportSQLitePath(src string) error {
	return a.syncer.ImportSQLite(src)
}

var _ = time.Now
