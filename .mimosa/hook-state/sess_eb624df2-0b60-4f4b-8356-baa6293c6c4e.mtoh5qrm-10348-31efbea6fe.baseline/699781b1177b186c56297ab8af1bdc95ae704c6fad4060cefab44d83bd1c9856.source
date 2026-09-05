package cover

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFromBytes(coverDir string, data []byte, ext string) (string, error) {
	if err := os.MkdirAll(coverDir, 0755); err != nil {
		return "", err
	}
	if ext == "" || !isImageExt(ext) {
		ext = ".png"
	}
	sum := sha256.Sum256(data)
	name := hex.EncodeToString(sum[:]) + ext
	full := filepath.Join(coverDir, name)
	if _, err := os.Stat(full); os.IsNotExist(err) {
		if err := os.WriteFile(full, data, 0644); err != nil {
			return "", err
		}
	}
	return name, nil
}

func SaveFromURL(coverDir, rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", errors.New("empty url")
	}
	if strings.HasPrefix(rawURL, "data:") {
		return "", errors.New("data urls handled on frontend")
	}
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mediateka/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("bad status: " + resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 15<<20))
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(rawURL)
	if ext == "" {
		switch resp.Header.Get("Content-Type") {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			ext = ".png"
		}
	}
	return SaveFromBytes(coverDir, data, ext)
}

func DominantColorHex(coverDir, coverFile string) string {
	if coverFile == "" {
		return defaultColor
	}
	full := filepath.Join(coverDir, coverFile)
	f, err := os.Open(full)
	if err != nil {
		return defaultColor
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return defaultColor
	}
	bounds := img.Bounds()
	var r, g, b, count uint64
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 4 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 4 {
			cr, cg, cb, _ := img.At(x, y).RGBA()
			r += uint64(cr >> 8)
			g += uint64(cg >> 8)
			b += uint64(cb >> 8)
			count++
		}
	}
	if count == 0 {
		return defaultColor
	}
	r /= count
	g /= count
	b /= count
	r = clamp(r)
	g = clamp(g)
	b = clamp(b)
	return toHex(r, g, b)
}

func ReadBytes(coverDir, coverFile string) ([]byte, error) {
	if coverFile == "" {
		return nil, nil
	}
	return os.ReadFile(filepath.Join(coverDir, coverFile))
}

func toHex(r, g, b uint64) string {
	const hexd = "0123456789abcdef"
	buf := make([]byte, 7)
	buf[0] = '#'
	buf[1] = hexd[(r>>4)&0xf]
	buf[2] = hexd[r&0xf]
	buf[3] = hexd[(g>>4)&0xf]
	buf[4] = hexd[g&0xf]
	buf[5] = hexd[(b>>4)&0xf]
	buf[6] = hexd[b&0xf]
	return string(buf)
}

func clamp(v uint64) uint64 {
	if v > 255 {
		return 255
	}
	if v < 30 {
		return 30
	}
	return v
}

func isImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	}
	return false
}

const defaultColor = "#5b6470"
