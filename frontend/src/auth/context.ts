import { createContext, useContext } from 'react'
import type { components } from '../api/schema'

type Account = components['schemas']['Account']

export type AuthState =
  | { kind: 'loading' }
  | { kind: 'anonymous' }
  | { kind: 'authenticated'; account: Account }

export type AuthContextValue = {
  state: AuthState
  login: (account: Account) => void
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthContextValue | null>(null)

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  // Fail loudly: a missing Provider is a wiring bug, not a runtime state.
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider')
  return ctx
}
