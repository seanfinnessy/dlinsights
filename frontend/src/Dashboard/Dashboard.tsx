import { useState } from 'react'
import type { Match } from '../api/matches'
import MatchHistory from '../MatchHistory/MatchHistory'
import styles from './Dashboard.module.css'

function Dashboard() {
  const [selectedMatch, setSelectedMatch] = useState<Match | null>(null)

  return (
    <div className={styles.grid}>
      <Section>
        <MatchHistory onSelectMatch={setSelectedMatch} selectedMatch={selectedMatch} />
      </Section>
      <Section title="Match Details">
        <p className={styles.empty}>
          {selectedMatch ? 'Details coming soon.' : 'Select a match to view details.'}
        </p>
      </Section>
      <Section title="Live Match Feed" fullWidth>
        <p className={styles.empty}>No live match detected.</p>
      </Section>
    </div>
  )
}

function Section({ title, fullWidth, children }: { title?: string; fullWidth?: boolean; children: React.ReactNode }) {
  return (
    <section className={`${styles.section}${fullWidth ? ` ${styles.fullWidth}` : ''}`}>
      {title && <h2 className={styles.sectionTitle}>{title}</h2>}
      <div className={styles.sectionBody}>{children}</div>
    </section>
  )
}

export default Dashboard
