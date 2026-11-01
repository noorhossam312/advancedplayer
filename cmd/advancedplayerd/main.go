package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

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
	"stop":        "STOP",
	"rvolume":     "RVOLUME",
	"svolume":     "SVOLUME",
	"volume":      "VOLUME",
	"kill":        "KILL",
}

var once sync.Once

func main() {
	signal.Ignore(syscall.SIGHUP)
	listener, err := net.Listen("tcp", "127.0.0.1:3382")
	if err != nil {
		slog.Error("Error while listening on port 3382", "error", err.Error())
		os.Exit(1)
	}
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
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Info("Shutting down safely...")
		listener.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			slog.Error("Failed to accept connection", "error", err)
			continue
		}
		go func(conn net.Conn, m *mpv.Mpv) {
			if handleConnection(conn, m) {
				once.Do(func() {
					listener.Close()
				})
			}
		}(conn, m)
	}
}

func handleConnection(conn net.Conn, m *mpv.Mpv) bool {
	defer conn.Close()
	success := `{"success": true}`
	reader := bufio.NewReader(conn)
	request, err := reader.ReadBytes('\n')
	if err != nil {
		return false
	}

	command := gjson.GetBytes(request, "command")

	_, ok := commands[command.String()]

	if !ok {
		slog.Error("No command", "command", command.String())
		success, _ = sjson.Set(success, "success", false)
		success, _ = sjson.Set(success, "reason", fmt.Sprintf("no such command %q", command.String()))
		conn.Write(append([]byte(success), '\n'))
		return false
	}

	allArgs := gjson.GetBytes(request, "args")

	allArgsMapped := allArgs.Map()
	switch command.String() {
	case "play":
		filePath, ok := allArgsMapped["1"]
		if !ok {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", "expected filepath, found nothing")
			conn.Write(append([]byte(success), '\n'))
			return false
		}

		err = player.Play(m, filePath.String())
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}

	case "search":
		//_, ok := allArgsMapped["1"]
		//if !ok { // if no search arg exists
		//	result, title, err := searcher.QueryYtMusic("", 10)
		//	if err != nil {
		//		success, _ = sjson.Set(success, "success", false)
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
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
		cleanFilename := pathologize.Clean(title) + ".ogg"
		songsDir, err := pathfuncs.SongsDir()
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
		err = downloader.DownloadFile(result, filepath.Join(songsDir, cleanFilename))
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
		err = player.Play(m, filepath.Join(songsDir, cleanFilename))
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
	case "togglepause":
		err := player.TogglePause(m)
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
	case "toggleloop":
		isLooped, err := player.ToggleLoop(m)
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
		success, _ = sjson.Set(success, "looped", isLooped)
	case "seek":
		amountToSeekResult, ok := allArgsMapped["1"]
		if !ok {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", "expected argument after seek, found nothing")
			conn.Write(append([]byte(success), '\n'))
			return false
		}
		amountToSeekString := amountToSeekResult.String()
		err := player.SeekRelative(m, amountToSeekString)
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
	case "timepos":
		timePosition, err := player.GetTimePos(m)
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
		success, _ = sjson.Set(success, "timepos", timePosition)
	case "stop":
		err := player.UnloadFile(m)
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
	case "rvolume":
		currentVolume, err := player.VolumeRelative(m, allArgsMapped["1"].String())
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
		success, _ = sjson.Set(success, "volume", currentVolume)
	case "svolume":
		err := player.VolumeLiteral(m, allArgsMapped["1"].String())
		if err != nil {
			success, _ = sjson.Set(success, "success", false)
			success, _ = sjson.Set(success, "reason", err.Error())
			conn.Write(append([]byte(success), '\n'))
			return false
		}
	case "volume":
		volume := player.GetVolume(m)
		success, _ = sjson.Set(success, "volume", volume)
	case "kill":
		slog.Info("Received \"kill\" command, shutting down gracefully...")
		conn.Write(append([]byte(success), '\n'))
		return true
	}
	conn.Write(append([]byte(success), '\n'))
	return false
}
