import { Route, Routes } from 'react-router'
import Register from './pages/Register'
import Login from './pages/Login'
import RequireAuth from './auth/RequireAuth'
import { useAuth } from './auth/context'
import { useEffect, useState } from 'react'
import { getHealth } from './api/client'

type HealthState =
  | { kind: 'loading' }
  | { kind: 'ok'; status: string }
  | { kind: 'error'; message: string }

export default function App() {
  const [health, setHealth] = useState<HealthState>({ kind: 'loading' })
  const { logout } = useAuth()

  useEffect(() => {
    let cancelled = false
    getHealth()
      .then((h) => {
        if (!cancelled) setHealth({ kind: 'ok', status: h.status })
      })
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : String(err)
        if (!cancelled) setHealth({ kind: 'error', message })
      })
    return () => {
      cancelled = true
    }
  }, [])

  return (
    <Routes>
      <Route path="/register" element={<Register />} />
      <Route path="/login" element={<Login />} />
      <Route element={<RequireAuth />}>
        <Route
          path="/"
          element={
            <>
              <h1>ESG Idea Contest v2</h1>
              <button onClick={() => void logout()}>Logout</button>
              {health.kind === 'loading' && <p>Backend: loading...</p>}
              {health.kind === 'ok' && <p>Backend: {health.status}</p>}
              {health.kind === 'error' && (
                <p>Backend: error {health.message}</p>
              )}
            </>
          }
        />
      </Route>
    </Routes>
  )
}
