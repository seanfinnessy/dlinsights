import { GetMatches, GetMatchInfo } from "../../wailsjs/go/main/App";

export type Match = {
  account_id: number
  player_kills: number
  player_deaths: number
  player_assists: number
  hero_id: number
  match_result: number
  player_team: number
  match_duration_s: number
  start_time: number
  net_worth: number
  game_mode: number
  match_id: number // int64, can exceed js range... just to note
}

export type Player = {
  account_id: number
  team: number
}

export type MatchInfo = {
  players: Player[]
  winning_team: number
}

export type MatchInfoResponse = {
  match_info: MatchInfo
}

export async function fetchMatches(steamId: string, numMatches: number): Promise<Match[]> {
  return GetMatches(steamId, numMatches)
}

export async function fetchMatchInfo(matchId: string): Promise<MatchInfo> {
  const response: MatchInfoResponse = await GetMatchInfo(matchId)
  return response.match_info
}
