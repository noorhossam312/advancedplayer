package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/noorhossam312/advancedplayer/internal/downloader"
	"github.com/noorhossam312/advancedplayer/internal/network"
	"github.com/noorhossam312/advancedplayer/internal/numfuncs"
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
		var err error
		args[0], err = filepath.Abs(args[0])
		if err != nil {
			return err
		}
		payload, _ := sjson.Set("{}", "command", "play")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
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
			if err != nil {
				return err
			}
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
			response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
			if err != nil {
				return err
			}
			if !gjson.GetBytes(response, "success").Bool() {
				return errors.New(gjson.GetBytes(response, "reason").String())
			}
			return nil
		}
		payload, _ := sjson.Set("{}", "command", "search")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
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
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
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
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		var loopedHuman string
		if gjson.GetBytes(response, "looped").Bool() {
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
		if !numfuncs.ValidNum(args[0]) {
			return fmt.Errorf("seek needs a number of seconds, got %q", args[0])
		}
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
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
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		int64Time := gjson.GetBytes(response, "timepos").Int()
		fmt.Println(timefuncs.SecondToHHMMSS(int64Time))
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Unload current file and stop playback",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command":"stop","args":{}}`
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		return nil
	},
}

var rvolumeCmd = &cobra.Command{
	Use:                "rvolume <volume>",
	Short:              "Add/subtract volume relatively",
	Args:               cobra.ExactArgs(1),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := sjson.Set("{}", "command", "rvolume")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		if !numfuncs.ValidNum(args[0]) {
			return fmt.Errorf("rvolume needs a numeric value, got %q", args[0])
		}
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		fmt.Println("Volume is now " + numfuncs.TruncateFloat(gjson.GetBytes(response, "volume").Float()))
		return nil
	},
}

var svolumeCmd = &cobra.Command{
	Use:   "svolume <volume>",
	Short: "Set volume",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := sjson.Set("{}", "command", "svolume")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		if !numfuncs.ValidNum(args[0]) {
			return fmt.Errorf("svolume needs a numeric value, got %q", args[0])
		}
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		return nil
	},
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Get current volume",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command":"volume","args":{}}`
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		fmt.Println("Volume is " + numfuncs.TruncateFloat(gjson.GetBytes(response, "volume").Float()))
		return nil
	},
}

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill daemon",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command":"kill","args":{}}`
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		return nil
	},
}

func main() {
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(playCmd, searchCmd, togglepauseCmd, toggleloopCmd, seekCmd, timeposCmd, stopCmd, rvolumeCmd, svolumeCmd, volumeCmd, killCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
