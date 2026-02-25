import { useMemo, useState } from 'react'
import type { Match } from '../api/matches'
import { useAuth } from '../context/AuthContext'
import { useMatches } from '../hooks/useMatches'
import { useHeroAssets } from '../hooks/useHeroAssets'
import { formatDuration, formatStartTime, formatSouls, formatGameMode } from '../services/conversions'
import styles from './MatchHistory.module.css'

interface MatchHistoryProps {
  onSelectMatch: (match: Match) => void
  selectedMatch: Match | null
}

function MatchHistory({ onSelectMatch, selectedMatch }: MatchHistoryProps) {
  const { steamId } = useAuth()

  const [matchAmount, setMatchAmount] = useState(10)
  const { data: matches, isLoading, isError } = useMatches(steamId, matchAmount)
  const { data: heroAssets } = useHeroAssets()

  // We want to use useMemo here to cache the heroMap, that way it doesnt
  // get remade on each re-render.. only when heroAssets updates.
  const heroMap = useMemo(() => {
    if (!heroAssets) return {}
    return Object.fromEntries(heroAssets.map(h => [h.id, h.name]))
  }, [heroAssets])

  const handleSelectMatchAmount = (value: string) => {
    setMatchAmount(Number(value));
  }
  
  const renderBody = () => {
    if (matches) {
      return matches?.map((match, i) => {
        const { date, time } = formatStartTime(match.start_time)
        return <tr
          key={i}
          className={[
            match.match_result === 1 ? styles.win : styles.loss,
            match === selectedMatch ? styles.selected : ''
          ].join(' ')}
          onClick={() => onSelectMatch(match)}
        >
          <td className={styles.date}>
            <div>{date}</div>
            <div className={styles.dateTime}>{time}</div>
          </td>
          <td>{heroMap[match.hero_id] ?? match.hero_id}</td>
          <td>{(() => { const m = formatGameMode(match.game_mode); return m ? <span className={`${styles.modeBadge} ${styles[m.style]}`}>{m.label}</span> : null })()}</td>
          <td className={styles.kda}>
            <span className={styles.kills}>{match.player_kills}</span>
            <span className={styles.separator}>/</span>
            <span className={styles.deaths}>{match.player_deaths}</span>
            <span className={styles.separator}>/</span>
            <span className={styles.assists}>{match.player_assists}</span>
          </td>
          <td className={styles.kills}>{formatSouls(match.net_worth)}</td>
          <td className={styles.duration}>{formatDuration(match.match_duration_s)}</td>
        </tr>
     })
    }
    return null;
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2 className={styles.title}>Match History</h2>
        <select
          className={styles.select}
          value={matchAmount}
          onChange={e => handleSelectMatchAmount(e.target.value)}
        >
          <option value={10}>10</option>
          <option value={25}>25</option>
          <option value={50}>50</option>
        </select>
      </div>
      {isLoading && <p className={styles.status}>Loading matches...</p>}
      {isError && <p className={styles.status}>Failed to load matches.</p>}
      {!isLoading && !isError && !matches?.length && <p className={styles.status}>No matches found.</p>}
      {matches && matches.length > 0 && (
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Date</th>
              <th>Hero</th>
              <th>Mode</th>
              <th>KDA</th>
              <th>Net Souls</th>
              <th>Duration</th>
            </tr>
          </thead>
          <tbody>
            {renderBody()}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default MatchHistory
