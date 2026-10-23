package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/noorhossam312/advancedplayer/internal/downloader"
	"github.com/noorhossam312/advancedplayer/internal/pathfuncs"
	"github.com/noorhossam312/advancedplayer/internal/searcher"
	"github.com/noorhossam312/advancedplayer/internal/timefuncs"
	"github.com/spf13/cobra"
	"github.com/spf13/pathologize"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var rootCmd = &cobra.Command{
	Use:   "adpctl",
	Short: "Control advancedplayerd",
}

var playCmd = &cobra.Command{
	Use:   "play <file>",
	Short: "Play a local file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := sjson.Set("{}", "command", "play")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		client, err := net.Dial("tcp", "127.0.0.1:3382")
		if err != nil {
			return fmt.Errorf("could not connect, is advancedplayerd on? %w", err)
		}
		defer client.Close()
		client.Write([]byte(payload))
		buf := make([]byte, 4096)
		client.Read(buf)
		if !gjson.Get(string(buf), "success").Bool() {
			return errors.New(gjson.Get(string(buf), "reason").String())
		}
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search YouTube and play the result",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			result, title, err := searcher.QueryYtMusic("", 10)
			songsDir, err := pathfuncs.SongsDir()
			if err != nil {
				return err
			}
			cleanFilename := pathologize.Clean(title) + ".ogg"
			err = downloader.DownloadFile(result, filepath.Join(songsDir, cleanFilename))
			if err != nil {
				return err
			}
			payload, _ := sjson.Set("{}", "command", "play")
			payload, _ = sjson.Set(payload, "args", map[string]string{})
			payload, _ = sjson.Set(payload, "args.1", filepath.Join(songsDir, cleanFilename))
			client, err := net.Dial("tcp", "127.0.0.1:3382")
			if err != nil {
				return fmt.Errorf("could not connect, is advancedplayerd on? %w", err)
			}
			defer client.Close()
			client.Write([]byte(payload))
			buf := make([]byte, 4096)
			client.Read(buf)
			if !gjson.Get(string(buf), "success").Bool() {
				return errors.New(gjson.Get(string(buf), "reason").String())
			}
			return nil
		}
		payload, _ := sjson.Set("{}", "command", "search")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		client, err := net.Dial("tcp", "127.0.0.1:3382")
		if err != nil {
			return fmt.Errorf("could not connect, is advancedplayerd on? %w", err)
		}
		defer client.Close()
		client.Write([]byte(payload))
		buf := make([]byte, 4096)
		client.Read(buf)
		if !gjson.Get(string(buf), "success").Bool() {
			return errors.New(gjson.Get(string(buf), "reason").String())
		}
		return nil
	},
}

var togglepauseCmd = &cobra.Command{
	Use:   "togglepause",
	Short: "Toggle pause state",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command":"togglepause","args":{}}`
		client, err := net.Dial("tcp", "127.0.0.1:3382")
		if err != nil {
			return fmt.Errorf("could not connect, is advancedplayerd on? %w", err)
		}
		defer client.Close()
		client.Write([]byte(payload))
		buf := make([]byte, 4096)
		client.Read(buf)
		if !gjson.Get(string(buf), "success").Bool() {
			return errors.New(gjson.Get(string(buf), "reason").String())
		}
		return nil
	},
}

var toggleloopCmd = &cobra.Command{
	Use:   "toggleloop",
	Short: "Toggle loop state",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command":"toggleloop","args":{}}`
		client, err := net.Dial("tcp", "127.0.0.1:3382")
		if err != nil {
			return fmt.Errorf("could not connect, is advancedplayerd on? %w", err)
		}
		defer client.Close()
		client.Write([]byte(payload))
		buf := make([]byte, 4096)
		client.Read(buf)
		if !gjson.Get(string(buf), "success").Bool() {
			return errors.New(gjson.Get(string(buf), "reason").String())
		}
		var loopedHuman string
		if gjson.Get(string(buf), "looped").Bool() {
			loopedHuman = "on"
		} else {
			loopedHuman = "off"
		}
		fmt.Println("Looping is " + loopedHuman)
		return nil
	},
}

var seekCmd = &cobra.Command{
	Use:                "seek <seconds>",
	Short:              "Seek relative to the current time position",
	Args:               cobra.ExactArgs(1),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := sjson.Set("{}", "command", "seek")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		if _, err := strconv.ParseFloat(args[0], 64); err != nil {
			return fmt.Errorf("seek needs a number of seconds, got %q", args[0])
		}
		client, err := net.Dial("tcp", "127.0.0.1:3382")
		if err != nil {
			return fmt.Errorf("could not connect, is advancedplayerd on? %w", err)
		}
		defer client.Close()
		client.Write([]byte(payload))
		buf := make([]byte, 4096)
		client.Read(buf)
		if !gjson.Get(string(buf), "success").Bool() {
			return errors.New(gjson.Get(string(buf), "reason").String())
		}
		return nil
	},
}

var timeposCmd = &cobra.Command{
	Use:   "timepos",
	Short: "Get current time position",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command":"timepos","args":{}}`
		client, err := net.Dial("tcp", "127.0.0.1:3382")
		if err != nil {
			return fmt.Errorf("could not connect, is advancedplayerd on? %w", err)
		}
		defer client.Close()
		client.Write([]byte(payload))
		buf := make([]byte, 4096)
		client.Read(buf)
		if !gjson.Get(string(buf), "success").Bool() {
			return errors.New(gjson.Get(string(buf), "reason").String())
		}
		int64Time := gjson.Get(string(buf), "timepos").Int()
		fmt.Println(timefuncs.SecondToHHMMSS(int64Time))
		return nil
	},
}

func main() {
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(playCmd, searchCmd, togglepauseCmd, toggleloopCmd, seekCmd, timeposCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
