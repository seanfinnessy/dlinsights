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
	PlayerTeam    int32 `json:"player_team"`
	MatchDuration int32 `json:"match_duration_s"`
	StartTime     int32 `json:"start_time"`
	NetWorth      int32 `json:"net_worth"`
	MatchMode     int32 `json:"game_mode"`
	MatchID       int64 `json:"match_id"`
}

type HeroAssets struct {
	HeroID int32  `json:"id"`
	Name   string `json:"name"`
}
type MatchInfoResponse struct {
	MatchInfo MatchInfo `json:"match_info"`
}

type MatchInfo struct {
	Players     []Player `json:"players"`
	WinningTeam int32    `json:"winning_team"`
}

type Player struct {
	AccountID int32 `json:"account_id"`
	Team      int32 `json:"team"`
}

func ConstructGetMatchHistoryURL(userId string) string {
	return PlayersURL + userId + MatchHistorySuff
}

func GetMatchInfo(matchId string) (MatchInfoResponse, error) {
	matchInfo := MatchInfoResponse{}
	url := MatchesURL + matchId + "/metadata"
	fmt.Println(url)
	res, errGet := http.Get(url)
	if errGet != nil || res.StatusCode != 200 {
		return matchInfo, fmt.Errorf("Issue calling API: %s", url)
	}
	err := json.NewDecoder(res.Body).Decode(&matchInfo)
	if err != nil {
		return matchInfo, err
	}

	for _, player := range matchInfo.MatchInfo.Players {
		fmt.Println(player.AccountID)
		fmt.Println(player.Team)
	}

	return matchInfo, nil
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
