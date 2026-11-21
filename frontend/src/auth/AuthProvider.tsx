import { useEffect, useState } from 'react'
import { api, setUnauthorizedHandler } from '../api/client'
import { AuthContext, type AuthState } from './context'

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>({ kind: 'loading' })

  useEffect(() => {
    let cancelled = false
    api
      .GET('/api/accounts/me')
      .then(({ data }) => {
        // 401 is the normal "not logged in" answer, so data is just undefined
        if (cancelled) return
        setState(
          data
            ? { kind: 'authenticated', account: data }
            : { kind: 'anonymous' },
        )
      })
      .catch(() => {
        if (!cancelled) setState({ kind: 'anonymous' })
      })
    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    // Session expired mid-use: flipping to anonymous makes RequireAuth redirect.
    setUnauthorizedHandler(() => setState({ kind: 'anonymous' }))
    return () => setUnauthorizedHandler(null)
  }, [])

  async function logout() {
    await api.DELETE('/api/sessions/current')
    setState({ kind: 'anonymous' })
  }

  return (
    <AuthContext.Provider
      value={{
        state,
        login: (account) => setState({ kind: 'authenticated', account }),
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
