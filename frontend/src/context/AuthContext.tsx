import { createContext, useContext } from 'react'
import type { ReactNode } from 'react'

type AuthContextType = {
  steamId: string | null
}

// Context for if user is logged in or not
const AuthContext = createContext<AuthContextType>({ steamId: null })

export function AuthProvider({ children }: { children: ReactNode }) {
  const params = new URLSearchParams(window.location.search)
  const steamId = params.get('steamId')

  return (
    <AuthContext.Provider value={{ steamId }}>
      {children}
    </AuthContext.Provider>
  )
}

// Hook to see if user is logged in or not
export function useAuth() {
  return useContext(AuthContext)
}
