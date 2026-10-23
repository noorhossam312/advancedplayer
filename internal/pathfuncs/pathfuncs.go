package pathfuncs

import (
	"os"
	"path/filepath"
)

func SongsDir() (string, error) {
	if dir := os.Getenv("ADP_MUSIC_DIR"); dir != "" {
		return dir, os.MkdirAll(dir, 0o755)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Music", "advancedplayer")
	return dir, os.MkdirAll(dir, 0o755)
}
