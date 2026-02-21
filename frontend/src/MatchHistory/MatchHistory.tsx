import { useAuth } from '../context/AuthContext'
import { useMatches } from '../hooks/useMatches'
import styles from './MatchHistory.module.css'

function MatchHistory() {
  const { steamId } = useAuth()

  // useQuery allows us to return data of response, if its loading, and if an error occurred
  const { data: matches, isLoading, isError } = useMatches(steamId, 10)

  if (isLoading) return <p className={styles.status}>Loading matches...</p>
  if (isError) return <p className={styles.status}>Failed to load matches.</p>
  if (!matches?.length) return <p className={styles.status}>No matches found.</p>

  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>#</th>
            <th>Hero ID</th>
            <th>KDA</th>
          </tr>
        </thead>
        <tbody>
          {matches.map((match, i) => (
            <tr key={i}>
              <td className={styles.index}>{i + 1}</td>
              <td>{match.hero_id}</td>
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
