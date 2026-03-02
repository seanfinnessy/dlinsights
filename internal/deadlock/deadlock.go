package deadlock

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PlayerMatchHistoryEntry struct {
	AccountID     int32 `json:"account_id"`
	PlayerKills   int32 `json:"player_kills"`
	PlayerDeaths  int32 `json:"player_deaths"`
	PlayerAssists int32 `json:"player_assists"`
	HeroID        int32 `json:"hero_id"`
	MatchResult   int32 `json:"match_result"`
	MatchDuration int32 `json:"match_duration_s"`
	StartTime     int32 `json:"start_time"`
	NetWorth      int32 `json:"net_worth"`
	MatchMode int32 	`json:"game_mode"`
}

type HeroAssets struct {
	HeroID int32  `json:"id"`
	Name   string `json:"name"`
}

func ConstructGetMatchHistoryURL(userId string) string {
	return PlayersURL + userId + MatchHistorySuff
}

func GetMatches(userId string, numMatches int) ([]PlayerMatchHistoryEntry, error) {
	playerHistory := []PlayerMatchHistoryEntry{}
	url := ConstructGetMatchHistoryURL(userId)
	res, errGet := http.Get(url)
	if errGet != nil || res.StatusCode != 200 {
		return []PlayerMatchHistoryEntry{}, fmt.Errorf("Issue calling API: %s", url)
	}

	err := json.NewDecoder(res.Body).Decode(&playerHistory)
	if err != nil {
		return playerHistory, err
	}

	return playerHistory[:numMatches], nil
}

func GetHeroAssets() ([]HeroAssets, error) {
	// Either grab from DB, or if version changes, we would update then return those values
	url := AssetsURL + HeroAssetsSuff

	heroAssets := []HeroAssets{}
	res, ok := http.Get(url)
	if ok != nil || res.StatusCode != 200 {
		return heroAssets, fmt.Errorf("Issue calling API: %s", url)
	}

	err := json.NewDecoder(res.Body).Decode(&heroAssets)
	if err != nil {
		return heroAssets, err
	}

	// TODO:
	// Write values to database
	// Implement some logic eventually that grabs values from DB..
	// If version changes, we would need to recall this API, update DB.

	return heroAssets, nil
}
