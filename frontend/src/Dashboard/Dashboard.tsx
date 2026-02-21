import MatchHistory from '../MatchHistory/MatchHistory'
import styles from './Dashboard.module.css'

function Dashboard() {
  return (
    <div className={styles.grid}>
      <Section title="Match History">
        <MatchHistory />
      </Section>
      <Section title="Live Match">
        <p className={styles.empty}>No live match detected.</p>
      </Section>
    </div>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className={styles.section}>
      <h2 className={styles.sectionTitle}>{title}</h2>
      <div className={styles.sectionBody}>{children}</div>
    </section>
  )
}

export default Dashboard
