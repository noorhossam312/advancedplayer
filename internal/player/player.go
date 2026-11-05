package player

import (
	"errors"
	"fmt"

	"github.com/gen2brain/go-mpv"
	"github.com/tidwall/sjson"
)

func Play(m *mpv.Mpv, path string) error {
	if err := m.Command([]string{"loadfile", path}); err != nil {
		return err
	}
	return nil
}

func ToggleLoop(m *mpv.Mpv) (any, error) {
	err := m.Command([]string{"cycle-values", "loop-file", "inf", "no"})
	loopedString := m.GetPropertyString("loop-file")
	isLooped := loopedString == "inf"
	return isLooped, err
}

func SeekRelative(m *mpv.Mpv, seekAmount string) error {
	return m.Command([]string{"seek", seekAmount, "relative"})
}

func TogglePause(m *mpv.Mpv) error {
	return m.Command([]string{"cycle", "pause"})
}

func GetTimePos(m *mpv.Mpv) (int64, error) {
	if m.GetPropertyString("path") == "" {
		return 0, errors.New("no loaded file")
	}
	unassertedTimePosFloat, err := m.GetProperty("time-pos", mpv.FormatInt64)
	if err != nil {
		return 0, err
	}
	timePosInt, ok := unassertedTimePosFloat.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected type %T for time-pos", unassertedTimePosFloat)
	}

	return timePosInt, nil
}

func UnloadFile(m *mpv.Mpv) error {
	return m.Command([]string{"stop"})
}

func VolumeRelative(m *mpv.Mpv, volumeAmount string) (string, error) {
	err := m.Command([]string{"add", "volume", volumeAmount})
	if err != nil {
		return "", err
	}
	currentVolume := m.GetPropertyString("volume")
	return currentVolume, nil
}

func VolumeLiteral(m *mpv.Mpv, volume string) error {
	return m.SetPropertyString("volume", volume)
}

func GetVolume(m *mpv.Mpv) string {
	return m.GetPropertyString("volume")
}

func Status(m *mpv.Mpv) (string, error) {
	volume := GetVolume(m)
	timePos, err := GetTimePos(m)
	if err != nil {
		return "", err
	}
	looped := m.GetPropertyString("loop-file") == "yes"
	paused := isPaused(m)
	final, _ := sjson.Set("{}", "volume", volume)
	final, _ = sjson.Set(final, "timepos", timePos)
	final, _ = sjson.Set(final, "looped", looped)
	final, _ = sjson.Set(final, "paused", paused)
	return final, nil
}

// Unexported functions
func isPaused(m *mpv.Mpv) bool {
	return m.GetPropertyString("idle-active") == "yes"
}
