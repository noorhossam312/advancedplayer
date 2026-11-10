package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"bufio"
	"strconv"

	"github.com/noorhossam312/advancedplayer/internal/network"
	"github.com/noorhossam312/advancedplayer/internal/numfuncs"
	"github.com/noorhossam312/advancedplayer/internal/searcher"
	"github.com/noorhossam312/advancedplayer/internal/timefuncs"
	"github.com/noorhossam312/advancedplayer/internal/utils"
	"github.com/noorhossam312/advancedplayer/internal/pathfuncs"
	"github.com/spf13/cobra"
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
			payload, _ := sjson.Set("{}", "command", "dl")
			payload, _ = sjson.Set(payload, "args", map[string]string{})
			payload, _ = sjson.Set(payload, "args.1", result)
			payload, _ = sjson.Set(payload, "args.2", title)
			response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
			if err != nil {
				return err
			}
			if !gjson.GetBytes(response, "success").Bool() {
				return errors.New(gjson.GetBytes(response, "reason").String())
			}
			pathToFile := gjson.GetBytes(response, "path")
			payload, _ = sjson.Set("{}", "command", "play")
			payload, _ = sjson.Set(payload, "args", map[string]string{})
			payload, _ = sjson.Set(payload, "args.1", pathToFile.String())
			response, err = network.SendAndRecv("127.0.0.1", "3382", payload)
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

// var rvolumeCmd = &cobra.Command{
// 	Use:                "rvolume <volume>",
// 	Short:              "Add/subtract volume relatively",
// 	Args:               cobra.ExactArgs(1),
// 	DisableFlagParsing: true,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		payload, _ := sjson.Set("{}", "command", "rvolume")
// 		payload, _ = sjson.Set(payload, "args", map[string]string{})
// 		payload, _ = sjson.Set(payload, "args.1", args[0])
// 		if !numfuncs.ValidNum(args[0]) {
// 			return fmt.Errorf("rvolume needs a numeric value, got %q", args[0])
// 		}
// 		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
// 		if err != nil {
// 			return err
// 		}
// 		if !gjson.GetBytes(response, "success").Bool() {
// 			return errors.New(gjson.GetBytes(response, "reason").String())
// 		}
// 		fmt.Println("Volume is now " + numfuncs.TruncateFloat(gjson.GetBytes(response, "volume").Float()))
// 		return nil
// 	},
// }

// var svolumeCmd = &cobra.Command{
// 	Use:   "svolume <volume>",
// 	Short: "Set volume",
// 	Args:  cobra.ExactArgs(1),
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		payload, _ := sjson.Set("{}", "command", "svolume")
// 		payload, _ = sjson.Set(payload, "args", map[string]string{})
// 		payload, _ = sjson.Set(payload, "args.1", args[0])
// 		if !numfuncs.ValidNum(args[0]) {
// 			return fmt.Errorf("svolume needs a numeric value, got %q", args[0])
// 		}
// 		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
// 		if err != nil {
// 			return err
// 		}
// 		if !gjson.GetBytes(response, "success").Bool() {
// 			return errors.New(gjson.GetBytes(response, "reason").String())
// 		}
// 		return nil
// 	},
// }

// var volumeCmd = &cobra.Command{
// 	Use:   "volume",
// 	Short: "Get current volume",
// 	Args:  cobra.ExactArgs(0),
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		payload := `{"command":"volume","args":{}}`
// 		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
// 		if err != nil {
// 			return err
// 		}
// 		if !gjson.GetBytes(response, "success").Bool() {
// 			return errors.New(gjson.GetBytes(response, "reason").String())
// 		}
// 		fmt.Println("Volume is " + numfuncs.TruncateFloat(gjson.GetBytes(response, "volume").Float()))
// 		return nil
// 	},
// }

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

var dlCmd = &cobra.Command{
	Use: "dl <link> [filename]",
	Short: "Download a file without playback",
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := sjson.Set("{}", "command", "dl")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", args[0])
		if len(args) == 2 {
			payload, _ = sjson.Set(payload, "args.2", args[1])
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

var dlpCmd = &cobra.Command{
	Use: "dlp <link>",
	Short: "Download a file and play it",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, _ := sjson.Set("{}", "command", "dlp")
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

var statusCmd = &cobra.Command{
	Use: "status",
	Short: "Print the daemon's audio player status. A file must be loaded",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command": "status"}`
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		humanizedStatus := utils.HumanizeStatus(string(response))
		fmt.Println(humanizedStatus)
		return nil
	},
}

var statusjCmd = &cobra.Command{
	Use: "statusj",
	Short: "Print the daemon's audio player status in JSON, ending with a newline character. A file must be loaded",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := `{"command": "status"}`
		response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
		if err != nil {
			return err
		}
		if !gjson.GetBytes(response, "success").Bool() {
			return errors.New(gjson.GetBytes(response, "reason").String())
		}
		response, _ = sjson.DeleteBytes(response, "success")
		fmt.Print(string(response))
		return nil
	},
}

var bvolumeCmd = &cobra.Command{
	Use: "volume",
	Short: "Volume controls, using with no args prints the volume, using it with +5 increments the volume by 5, using it with -5 decrements the volume by 5, using it with 5 sets the volume to 5",
	Args: cobra.MaximumNArgs(1),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := "{}"
		var mode string
		if len(args) == 0 {
			mode = "pvol"
		} else {
			if args[0][:1] != "+" && args[0][:1] != "-" {
				mode = "svol"		
			}
		}
		switch mode {
		case "pvol":
			payload, _ = sjson.Set(payload, "command", "volume")
			response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
			if err != nil {
				return err
			}
			if !gjson.GetBytes(response, "success").Bool() {
				return errors.New(gjson.GetBytes(response, "reason").String())
			}
			fmt.Println(numfuncs.TruncateFloat(gjson.GetBytes(response, "volume").Float()))
			return nil
		case "svol":
			payload, _ := sjson.Set("{}", "command", "svolume")
			payload, _ = sjson.Set(payload, "args", map[string]string{})
			payload, _ = sjson.Set(payload, "args.1", args[0])
			if !numfuncs.ValidNum(args[0]) {
				return fmt.Errorf("need a numeric value, got %q", args[0])
			}
			response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
			if err != nil {
				return err
			}
			if !gjson.GetBytes(response, "success").Bool() {
				return errors.New(gjson.GetBytes(response, "reason").String())
			}
		default:
			payload, _ := sjson.Set("{}", "command", "rvolume")
			payload, _ = sjson.Set(payload, "args", map[string]string{})
			payload, _ = sjson.Set(payload, "args.1", args[0])
			if !numfuncs.ValidNum(args[0]) {
				return fmt.Errorf("need a numeric value, got %q", args[0])
			}
			response, err := network.SendAndRecv("127.0.0.1", "3382", payload)
			if err != nil {
				return err
			}
			if !gjson.GetBytes(response, "success").Bool() {
				return errors.New(gjson.GetBytes(response, "reason").String())
			}
			fmt.Println("Volume is now " + numfuncs.TruncateFloat(gjson.GetBytes(response, "volume").Float()))
		}
		return nil
	},
}

var libraryCmd = &cobra.Command{
	Use: "library",
	Short: "Browse the music directory for music to play",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		songsDir, err := pathfuncs.SongsDir()
		if err != nil {
			return err
		}
		files, err := os.ReadDir(songsDir)
		if len(files) == 0 {
			return errors.New("no music in music directory")
		}
		for index, file := range files {
			fmt.Printf("(%d) %s\n", index+1, file.Name())
		}
		fmt.Print("\nPlease pick a file.\n> ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return scanner.Err()
		}
		input := scanner.Text()
		inputInt, err := strconv.Atoi(input)
		if err != nil {
			return err
		}
		if inputInt > len(files) || inputInt < 1 {
			return fmt.Errorf("out of range %q", input)
		}
		selectedFile := files[inputInt-1]
		payload, _ := sjson.Set("{}", "command", "play")
		payload, _ = sjson.Set(payload, "args", map[string]string{})
		payload, _ = sjson.Set(payload, "args.1", filepath.Join(songsDir, selectedFile.Name()))
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
	rootCmd.AddCommand(playCmd, searchCmd, togglepauseCmd, toggleloopCmd, seekCmd, timeposCmd, stopCmd, bvolumeCmd, killCmd, dlCmd, dlpCmd, statusCmd, statusjCmd, libraryCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// func commandExists(root *cobra.Command, name string) bool {
// 	if name == "help" || name == "completion" {
// 		return false
// 	}
// 	cmd, _, err := root.Find([]string{name})
// 	return err == nil && cmd != root
// }


