import { AuthProvider, useAuth } from './context/AuthContext'
import LoginPage from './LoginPage/LoginPage'
import Dashboard from './Dashboard/Dashboard'
import styles from './App.module.css'
import { useHeroAssets } from './hooks/useHeroAssets'

function AppShell() {
  const { steamId } = useAuth()
  
  // Grab Hero assets on app load.
  const {data: heroAssets } = useHeroAssets();
  console.log(heroAssets);

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <div className={styles.headerLeft}>
          <div className={styles.logoMark}>DL</div>
          <span className={styles.appName}>Deadlock Stats Tracker</span>
        </div>
        {steamId && (
          <div className={styles.loggedInBadge}>
            <span className={styles.badgeDot} />
            {steamId}
          </div>
        )}
      </header>
      <main>
        {steamId ? <Dashboard /> : <LoginPage />}
      </main>
    </div>
  )
}

function App() {
  return (
    <AuthProvider>
      <AppShell />
    </AuthProvider>
  )
}

export default App
