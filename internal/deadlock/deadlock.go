package deadlock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PlayerMatchHistoryEntry struct {
	PlayerKills int32 `json:"player_kills"` 
}

func GetMatchHistoryURL(userId string) string {
	return "https://api.deadlock-api.com/v1/players/" + userId + "/match-history"
}

func GetMatches(userId string) error {
	var playerHistory[] *PlayerMatchHistoryEntry
	url := GetMatchHistoryURL(userId)
	res, errGet := http.Get(url)
	if errGet != nil || res.StatusCode != 200 {
		return fmt.Errorf("Issue calling API: %s", url)
	}

	// Decode bytes from response
	bodyBytes, errDecode := io.ReadAll(res.Body)
	if errDecode != nil {
		return errDecode
	}

	errUnmarshal := json.Unmarshal(bodyBytes, &playerHistory)
	if errUnmarshal != nil {
		return errUnmarshal
	}
	
	for _, match := range(playerHistory) {
		fmt.Println(match.PlayerKills)
	}
	return nil
}