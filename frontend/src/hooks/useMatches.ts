import { useQuery } from '@tanstack/react-query'
import { fetchMatches } from '../api/matches'

export function useMatches(steamId: string | null, numMatches: number) {
  return useQuery({
    // queryKey changes will trigger a refect
    queryKey: ['matches', steamId, numMatches],
    queryFn: () => fetchMatches(steamId!, numMatches),
    enabled: !!steamId,
  })
}
