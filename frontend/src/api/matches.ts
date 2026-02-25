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
  const res = await fetch(`http://localhost:3000/get-matches/${steamId}?amount=${numMatches}`);
  if (!res.ok) {
    throw new Error('Failed to fetch matches')
  }
  return res.json()
}
