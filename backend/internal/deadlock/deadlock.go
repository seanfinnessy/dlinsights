package deadlock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PlayerMatchHistoryEntry struct {
	AccountID int32 `json:"account_id"`
	PlayerKills int32 `json:"player_kills"` 
	PlayerDeaths int32 `json:"player_deaths"` 
	PlayerAssists int32 `json:"player_assists"` 
	HeroID int32 `json:"hero_id"`
}

func GetMatchHistoryURL(userId string) string {
	return "https://api.deadlock-api.com/v1/players/" + userId + "/match-history?only_stored_history=true"
}

func GetMatches(userId string, numMatches int) ([]PlayerMatchHistoryEntry, error) {
	var playerHistory[] PlayerMatchHistoryEntry
	url := GetMatchHistoryURL(userId)
	res, errGet := http.Get(url)
	if errGet != nil || res.StatusCode != 200 {
		return []PlayerMatchHistoryEntry{}, fmt.Errorf("Issue calling API: %s", url)
	}

	// read bytes
	bodyBytes, errDecode := io.ReadAll(res.Body)
	if errDecode != nil {
		return []PlayerMatchHistoryEntry{}, errDecode
	}

	// unmarshal bytes into player history
	errUnmarshal := json.Unmarshal(bodyBytes, &playerHistory)
	if errUnmarshal != nil {
		return []PlayerMatchHistoryEntry{}, errUnmarshal
	}
	
	return playerHistory[:numMatches], nil
}