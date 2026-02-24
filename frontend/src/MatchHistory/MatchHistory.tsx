import { useMemo } from 'react'
import { useAuth } from '../context/AuthContext'
import { useMatches } from '../hooks/useMatches'
import { useHeroAssets } from '../hooks/useHeroAssets'
import styles from './MatchHistory.module.css'

function MatchHistory() {
  const { steamId } = useAuth()

  const { data: matches, isLoading, isError } = useMatches(steamId, 20)
  const { data: heroAssets } = useHeroAssets()

  // We want to use useMemo here to cache the heroMap, that way it doesnt
  // get remade on each re-render.. only when heroAssets updates.
  const heroMap = useMemo(() => {
    if (!heroAssets) return {}
    return Object.fromEntries(heroAssets.map(h => [h.id, h.name]))
  }, [heroAssets])

  if (isLoading) return <p className={styles.status}>Loading matches...</p>
  if (isError) return <p className={styles.status}>Failed to load matches.</p>
  if (!matches?.length) return <p className={styles.status}>No matches found.</p>

  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>#</th>
            <th>Hero</th>
            <th>KDA</th>
          </tr>
        </thead>
        <tbody>
          {matches.map((match, i) => (
            <tr key={i} className={match.match_result === 1 ? styles.win : styles.loss}>
              <td className={styles.index}>{i + 1}</td>
              <td>{heroMap[match.hero_id] ?? match.hero_id}</td>
              <td className={styles.kda}>
                <span className={styles.kills}>{match.player_kills}</span>
                <span className={styles.separator}>/</span>
                <span className={styles.deaths}>{match.player_deaths}</span>
                <span className={styles.separator}>/</span>
                <span className={styles.assists}>{match.player_assists}</span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default MatchHistory
