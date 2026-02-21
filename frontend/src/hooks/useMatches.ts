import { useQuery } from '@tanstack/react-query'
import { fetchMatches } from '../api/matches'

export function useMatches(steamId: string | null, numMatches: number) {
  return useQuery({
    queryKey: ['matches', steamId],
    queryFn: () => fetchMatches(steamId!, numMatches),
    enabled: !!steamId,
  })
}
