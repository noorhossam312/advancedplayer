package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"path/filepath"

	"github.com/gen2brain/go-mpv"
	"github.com/noorhossam312/advancedplayer/internal/downloader"
	"github.com/noorhossam312/advancedplayer/internal/pathfuncs"
	"github.com/noorhossam312/advancedplayer/internal/player"
	"github.com/noorhossam312/advancedplayer/internal/searcher"
	"github.com/spf13/pathologize"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var commands = map[string]string{
	"play":        "PLAY",
	"search":      "SEARCH",
	"togglepause": "TOGGLEPAUSE",
	"toggleloop":  "TOGGLELOOP",
	"seek":        "SEEK",
	"timepos":     "TIMEPOS",
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:3382")
	dieOnErr(err)
	defer listener.Close()

	m := mpv.New()
	defer m.TerminateDestroy()
	err = m.SetOptionString("idle", "yes")
	if err != nil {
		log.Fatal(err)
	}
	err = m.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	for {
		success := `{"success": true}`

		data := make([]byte, 0)
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			continue
		}

		buffer := make([]byte, 4096)
		n, _ := conn.Read(buffer)

		data = append(data, buffer[:n]...)

		stringData := string(data)

		command := gjson.Get(stringData, "command")

		_, ok := commands[command.String()]

		if !ok {
			slog.Error("No command", "command", command.String())
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", fmt.Sprintf("no such command %q", command.String()))
			conn.Write([]byte(success))
			conn.Close()
			continue
		}

		allArgs := gjson.Get(stringData, "args")

		allArgsMapped := allArgs.Map()
		switch command.String() {
		case "play":
			filePath, ok := allArgsMapped["1"]
			if !ok {
				success, _ = sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", "expected filepath, found nothing")
				conn.Write([]byte(success))
				conn.Close()
				continue
			}

			err = player.Play(filePath.String(), m)
			if err != nil {
				success, _ = sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", err.Error())
				conn.Write([]byte(success))
				conn.Close()
				continue
			}

			if err != nil {
				success, _ := sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", err.Error())
				conn.Write([]byte(success))
				conn.Close()
				continue
			}

		case "search":
			//_, ok := allArgsMapped["1"]
			//if !ok { // if no search arg exists
			//	result, title, err := searcher.QueryYtMusic("", 10)
			//	if err != nil {
			//		success, _ := sjson.Set(success, "success", false)
			//		success, _ = sjson.Set(success, "reason", err.Error())
			//		conn.Write([]byte(success))
			//		conn.Close()
			//		continue
			//	}
			//	go func() {
			//		cleanFilename := pathologize.Clean(title) + ".ogg"
			//		songsDir, err := pathfuncs.SongsDir()
			//		if err != nil {
			//			slog.Error("Failed to get music library directory", "error", err)
			//		}
			//		err = downloader.DownloadFile(result, filepath.Join(songsDir, cleanFilename))
			//		if err != nil {
			//			slog.Error("Failed to download file", "error", err)
			//		}
			//		err = player.Play(filepath.Join(songsDir, cleanFilename), m)
			//		if err != nil {
			//			slog.Error("Failed to play file", "error", err)
			//		}
			//	}()
			//} else { // if it DOES exist
			result, title, err := searcher.QueryYtMusic(allArgsMapped["1"].String(), 10)
			if err != nil {
				success, _ := sjson.Set(success, "success", false)
				conn.Write([]byte(success))
				conn.Close()
				continue
			}
			go func() {
				cleanFilename := pathologize.Clean(title) + ".ogg"
				songsDir, err := pathfuncs.SongsDir()
				if err != nil {
					slog.Error("Failed to get music library directory", "error", err)
				}
				err = downloader.DownloadFile(result, filepath.Join(songsDir, cleanFilename))
				if err != nil {
					slog.Error("Failed to download file", "error", err)
				}
				err = player.Play(filepath.Join(songsDir, cleanFilename), m)
				if err != nil {
					slog.Error("Failed to play file", "error", err)
				}
			}()
		case "togglepause":
			err := player.TogglePause(m)
			if err != nil {
				success, _ = sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", err.Error())
				conn.Write([]byte(success))
				conn.Close()
				continue
			}
		case "toggleloop":
			isLooped, err := player.ToggleLoop(m)
			if err != nil {
				success, _ = sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", err.Error())
				conn.Write([]byte(success))
				conn.Close()
				continue
			}
			success, _ = sjson.Set(success, "looped", isLooped)
		case "seek":
			amountToSeekResult, ok := allArgsMapped["1"]
			if !ok {
				success, _ = sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", "expected argument after seek, found nothing")
				conn.Write([]byte(success))
				conn.Close()
				continue
			}
			amountToSeekString := amountToSeekResult.String()
			err := player.SeekRelative(m, amountToSeekString)
			if err != nil {
				success, _ = sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", err.Error())
				conn.Write([]byte(success))
				conn.Close()
				continue
			}
		case "timepos":
			timePosition, err := player.GetTimePos(m)
			if err != nil {
				success, _ = sjson.Set(success, "success", false)
				success, _ = sjson.Set(success, "reason", err.Error())
				conn.Write([]byte(success))
				conn.Close()
				continue
			}
			success, _ = sjson.Set(success, "timepos", timePosition)
		}

		conn.Write([]byte(success))
		conn.Close()
	}
}

func dieOnErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}
