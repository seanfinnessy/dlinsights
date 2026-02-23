package deadlock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PlayerMatchHistoryEntry struct {
	AccountID     int32 `json:"account_id"`
	PlayerKills   int32 `json:"player_kills"`
	PlayerDeaths  int32 `json:"player_deaths"`
	PlayerAssists int32 `json:"player_assists"`
	HeroID        int32 `json:"hero_id"`
}

type HeroAssets struct {
	HeroID int32  `json:"id"`
	Name   string `json:"name"`
}

func ConstructGetMatchHistoryURL(userId string) string {
	return PlayersURL + userId + MatchHistorySuff
}

func GetMatches(userId string, numMatches int) ([]PlayerMatchHistoryEntry, error) {
	var playerHistory []PlayerMatchHistoryEntry
	url := ConstructGetMatchHistoryURL(userId)
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

func GetHeroAssets() error {
	// Either grab from DB, or if version changes, we would update then return those values
	url := AssetsURL + HeroAssetsSuff
	res, ok := http.Get(url)
	if ok != nil || res.StatusCode != 200 {
		return fmt.Errorf("Issue calling API: %s", url)
	}

	bodyBytes, errDecode := io.ReadAll(res.Body)
	if errDecode != nil {
		return errDecode
	}

	var heroAssets HeroAssets
	errUnmarshal := json.Unmarshal(bodyBytes, &heroAssets)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	// Write values to database
	// Implement some logic eventually that grabs values from DB..
	// If version changes, we would need to recall this API, update DB.
}
