import type { Match } from '../api/matches'
import { useMatchInfo } from '../hooks/useMatchInfo'
import { useAuth } from '../context/AuthContext'
import styles from './MatchDetails.module.css'

interface MatchDetailsProps {
  match: Match
}

function MatchDetails({ match }: MatchDetailsProps) {
  const { accountId } = useAuth()
  const { data: matchInfo, isLoading, isError } = useMatchInfo(String(match.match_id))

  if (isLoading) return <p className={styles.status}>Loading...</p>
  if (isError) return <p className={styles.status}>Failed to load details.</p>
  if (!matchInfo) return null

  const teamA = matchInfo.players.filter(p => p.team === 0)
  const teamB = matchInfo.players.filter(p => p.team === 1)

  const renderPlayer = (p: { account_id: number }) => {
    const isYou = p.account_id === accountId
    return (
      <div key={p.account_id} className={`${styles.player}${isYou ? ` ${styles.you}` : ''}`}>
        <span className={styles.accountId}>{p.account_id}</span>
        {isYou && <span className={styles.youBadge}>you</span>}
      </div>
    )
  }

  const winningTeam = matchInfo.winning_team

  return (
    <div className={styles.teams}>
      <div className={styles.team}>
        <div className={`${styles.teamHeader} ${winningTeam === 0 ? styles.win : styles.loss}`}>Team 1</div>
        {teamA.map(renderPlayer)}
      </div>
      <div className={styles.teamDivider} />
      <div className={styles.team}>
        <div className={`${styles.teamHeader} ${winningTeam === 1 ? styles.win : styles.loss}`}>Team 2</div>
        {teamB.map(renderPlayer)}
      </div>
    </div>
  )
}

export default MatchDetails
