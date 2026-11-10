package utils

import (
	"fmt"

	"github.com/noorhossam312/advancedplayer/internal/timefuncs"
	"github.com/tidwall/gjson"
	// "github.com/google/shlex"
)

func HumanizeStatus(json string) string { // {"volume":50,"timepos":120.5,"looped":true,"paused":false}
	var final string
	results := gjson.GetMany(json, "volume", "timepos", "looped", "paused")
	volume := results[0].Int()
	timePos := results[1].Int()
	timePosHuman := timefuncs.SecondToHHMMSS(timePos)
	looped := results[2].Bool()
	var loopedHuman string
	if looped {
		loopedHuman = "Yes"
	} else {
		loopedHuman = "No"
	}
	paused := results[3].Bool()
	var pausedHuman string
	if paused {
		pausedHuman = "Yes"
	} else {
		pausedHuman = "No"
	}
	final = fmt.Sprintf(
		`Volume: %d
Time: %s
Looped: %s
Paused: %s`,
		volume,
		timePosHuman,
		loopedHuman,
		pausedHuman,
	)
	return final
}

// func ExpandAlias(aliasMap string, alias string) ([]string, bool, error) {
// 	aliasValue := gjson.Get(aliasMap, gjson.Escape(alias))
// 	if !aliasValue.Exists() {
// 		return nil, false, nil
// 	}
// 	result, err := shlex.Split(aliasValue.String())
// 	if err != nil {
// 		return nil, true, fmt.Errorf("bad alias %q: %w", alias, err)
// 	}
// 	return result, true, nil
// }

