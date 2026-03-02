import { GetMatches } from "../../wailsjs/go/main/App";

export type Match = {
  account_id: number
  player_kills: number
  player_deaths: number
  player_assists: number
  hero_id: number
  match_result: number
  match_duration_s: number
  start_time: number
  net_worth: number
  game_mode: number
}

export async function fetchMatches(steamId: string, numMatches: number): Promise<Match[]> {
  return GetMatches(steamId, numMatches)
}
