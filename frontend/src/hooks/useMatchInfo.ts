import { useQuery } from '@tanstack/react-query'
import { fetchMatchInfo } from '../api/matches'

export function useMatchInfo(matchId: string) {
  return useQuery({
    // queryKey changes will trigger a refect
    queryKey: ['matchInfo', matchId],
    queryFn: () => fetchMatchInfo(matchId),
    enabled: !!matchId,
  })
}
