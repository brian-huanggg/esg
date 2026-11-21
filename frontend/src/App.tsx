import { Route, Routes } from 'react-router' // new
import Register from './pages/Register'
import { useEffect, useState } from 'react'
import { getHealth } from './api/client'
import Login from './pages/Login'

type HealthState =
  | { kind: 'loading' }
  | { kind: 'ok'; status: string }
  | { kind: 'error'; message: string }

export default function App() {
  const [health, setHealth] = useState<HealthState>({ kind: 'loading' })

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
      <Route
        path="/"
        element={
          <>
            <h1>ESG Idea Contest</h1>
            {health.kind === 'loading' && <p>Backend: loading...</p>}
            {health.kind === 'ok' && <p>Backend: {health.status}</p>}
            {health.kind === 'error' && <p>Backend: error {health.message}</p>}
          </>
        }
      />
    </Routes>
  )
}
