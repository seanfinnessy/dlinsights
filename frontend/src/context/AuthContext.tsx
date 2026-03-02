import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import { EventsOn } from '../../wailsjs/runtime/runtime'

type AuthContextType = {
  steamId: string | null
  authError: string | null
}

const AuthContext = createContext<AuthContextType>({ steamId: null, authError: null })

export function AuthProvider({ children }: { children: ReactNode }) {
  const [steamId, setSteamId] = useState<string | null>(null)
  const [authError, setAuthError] = useState<string | null>(null)

  // handle the events emitted from Go
  useEffect(() => {
    EventsOn('steam:authed', (id: string) => setSteamId(id))
    EventsOn('steam:error', () => setAuthError('Login failed. Please try again.'))
  }, [])

  return (
    <AuthContext.Provider value={{ steamId, authError }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
