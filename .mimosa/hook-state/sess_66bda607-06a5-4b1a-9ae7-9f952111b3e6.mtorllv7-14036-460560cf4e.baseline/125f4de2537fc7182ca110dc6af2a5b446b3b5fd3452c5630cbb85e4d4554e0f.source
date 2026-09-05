package sync

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"pergamon/internal/store"

	"github.com/studio-b12/gowebdav"
)

type Manager struct {
	dbPath string
	store  *store.Store
}

func NewManager(dbPath string, s *store.Store) *Manager {
	return &Manager{dbPath: dbPath, store: s}
}

type WebDAVConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	RemoteDir string `json:"remoteDir"`
}

type SyncResult struct {
	Direction string `json:"direction"`
	Uploaded  int    `json:"uploaded"`
	Message   string `json:"message"`
	Time      string `json:"time"`
}

func (m *Manager) TestConnection(cfg WebDAVConfig) error {
	c := client(cfg)
	return c.Connect()
}

func (m *Manager) SyncNow() (SyncResult, error) {
	res := SyncResult{Time: time.Now().Format(time.RFC3339)}
	cfg, err := m.loadConfig()
	if err != nil {
		return res, err
	}
	if cfg.URL == "" {
		return res, errors.New("webdav not configured")
	}
	c := client(cfg)
	if err := c.Connect(); err != nil {
		return res, fmt.Errorf("connect: %w", err)
	}
	if err := c.Mkdir(cfg.RemoteDir, 0755); err != nil && !isAlreadyExists(err) {
		return res, err
	}

	localMod := m.localModTime()
	remoteMod, rErr := m.remoteModTime(c, cfg)

	upload := true
	if rErr == nil && remoteMod.After(localMod) {
		upload = false
	}

	if upload {
		if err := m.upload(c, cfg); err != nil {
			return res, err
		}
		res.Direction = "up"
		res.Uploaded = 1
		res.Message = "База выгружена в облако"
	} else {
		if err := m.download(c, cfg); err != nil {
			return res, err
		}
		res.Direction = "down"
		res.Message = "База скачана из облака"
	}
	return res, nil
}

func (m *Manager) loadConfig() (WebDAVConfig, error) {
	cfg := WebDAVConfig{RemoteDir: "/Pergamon"}
	cfg.URL, _ = m.store.GetSetting("webdav_url")
	cfg.Username, _ = m.store.GetSetting("webdav_user")
	cfg.Password, _ = m.store.GetSetting("webdav_pass")
	rd, _ := m.store.GetSetting("webdav_remote_dir")
	if rd != "" {
		cfg.RemoteDir = rd
	}
	return cfg, nil
}

func (m *Manager) SaveConfig(cfg WebDAVConfig) error {
	sets := map[string]string{
		"webdav_url":        cfg.URL,
		"webdav_user":       cfg.Username,
		"webdav_pass":       cfg.Password,
		"webdav_remote_dir": cfg.RemoteDir,
	}
	for k, v := range sets {
		if err := m.store.SetSetting(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) GetConfig() (WebDAVConfig, error) {
	return m.loadConfig()
}

func client(cfg WebDAVConfig) *gowebdav.Client {
	c := gowebdav.NewClient(cfg.URL, cfg.Username, cfg.Password)
	c.SetTimeout(30 * time.Second)
	return c
}

func (m *Manager) remotePath(cfg WebDAVConfig) string {
	// Имя файла в облаке оставлено прежним: смена имени оторвала бы новые
	// сборки от базы, уже выгруженной старыми версиями приложения.
	return cfg.RemoteDir + "/mediateka.db"
}

func (m *Manager) localModTime() time.Time {
	st, err := os.Stat(m.dbPath)
	if err != nil {
		return time.Time{}
	}
	return st.ModTime()
}

func (m *Manager) remoteModTime(c *gowebdav.Client, cfg WebDAVConfig) (time.Time, error) {
	info, err := c.Stat(m.remotePath(cfg))
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func (m *Manager) upload(c *gowebdav.Client, cfg WebDAVConfig) error {
	data, err := os.ReadFile(m.dbPath)
	if err != nil {
		return err
	}
	return c.Write(m.remotePath(cfg), data, 0644)
}

func (m *Manager) download(c *gowebdav.Client, cfg WebDAVConfig) error {
	data, err := c.Read(m.remotePath(cfg))
	if err != nil {
		return err
	}
	tmp := m.dbPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, m.dbPath)
}

func isAlreadyExists(err error) bool {
	if err == nil {
		return true
	}
	return err.Error() == "405 Method Not Allowed" || err.Error() == "301 Moved Permanently"
}

func (m *Manager) ExportJSON(w io.Writer) error {
	return exportJSON(m.store, w)
}

func (m *Manager) ImportJSON(r io.Reader) error {
	return importJSON(m.store, r)
}

func (m *Manager) ExportSQLite(dest string) error {
	src, err := os.Open(m.dbPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func (m *Manager) ImportSQLite(src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(m.dbPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return os.Remove(filepath.Base(src))
}
