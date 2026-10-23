package searcher

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	itgo "github.com/nezbut/innertube-go"
	"github.com/tidwall/gjson"
)

func QueryYtMusic(searchQuery string, querySize uint8) (string, string, error) {
	var interactiveMode bool
	httpClient := http.DefaultClient
	client, err := itgo.NewInnerTube(httpClient, "WEB_REMIX", "2.20230728.00.00", "", "", "", nil, true)

	if err != nil {
		return "", "", err
	}

	ctx := context.Background()

	if searchQuery == "" {
		searchQuery, err = PromptUserSearch()
		if err != nil {
			return "", "", err
		}
		interactiveMode = true
	}

	//goland:noinspection SpellCheckingInspection
	params := "EgWKAQIIAWoSEAQQAxAJEBAQChAFEBUQERAO"

	results, err := client.Search(ctx, &searchQuery, &params, nil)

	if err != nil {
		return "", "", err
	}

	jsonBytes, err := json.Marshal(results)
	if err != nil {
		return "", "", err
	}

	allVids := gjson.GetBytes(jsonBytes, "contents.tabbedSearchResultsRenderer.tabs.0.tabRenderer.content.sectionListRenderer.contents.0.musicShelfRenderer.contents")

	videoIdPath := "musicResponsiveListItemRenderer.playlistItemData.videoId"
	titlePath := "musicResponsiveListItemRenderer.flexColumns.0.musicResponsiveListItemFlexColumnRenderer.text.runs.0.text"

	if !interactiveMode && searchQuery != "" {
		firstId := allVids.Get(fmt.Sprintf("0.%s", videoIdPath))
		title := allVids.Get(fmt.Sprintf("0.%s", titlePath))
		firstIdLink := fmt.Sprintf("https://www.youtube.com/watch?v=%s", firstId.String())
		return firstIdLink, title.String(), nil
	}

	if interactiveMode {
		fullLink, title, err := PromptUserPick(allVids, titlePath, videoIdPath, querySize)
		if err != nil {
			return "", "", err
		}
		return fullLink, title, nil
	}
	return "", "", errors.New("unexpected path, not returning data")
}

func PromptUserSearch() (string, error) {
	fmt.Print("Search query\n> ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", scanner.Err()
	}
	searchQuery := scanner.Text()
	return searchQuery, nil
}

func PromptUserPick(allVids gjson.Result, titlePath string, videoIdPath string, querySize uint8) (string, string, error) {
	allVids.ForEach(func(key, value gjson.Result) bool {
		currentVideoTitle := value.Get(titlePath)
		fmt.Printf("%d) %s\n", key.Int()+1, currentVideoTitle.String())
		if uint8(key.Uint()+1) == querySize {
			return false
		}
		return true
	})

	fmt.Print("\nPlease pick a song.\n> ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", "", scanner.Err()
	}
	selectedSong := scanner.Text()
	selectedSongInt, err := strconv.Atoi(selectedSong)

	if err != nil {
		return "", "", err
	}

	if selectedSongInt < 1 || selectedSongInt > int(querySize) {
		return "", "", fmt.Errorf("selection must be between 1 and %d", querySize)
	}

	videoId := allVids.Get(fmt.Sprintf("%d.%s", selectedSongInt-1, videoIdPath))
	title := allVids.Get(fmt.Sprintf("%d.%s", selectedSongInt-1, titlePath))
	fullLink := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId.String())
	return fullLink, title.String(), nil
}
