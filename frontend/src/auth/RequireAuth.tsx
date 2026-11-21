import { Navigate, Outlet } from 'react-router'
import { useAuth } from './context'

export default function RequireAuth() {
  const { state } = useAuth()
  if (state.kind === 'loading') return <p>Loading...</p>
  if (state.kind === 'anonymous') return <Navigate to="/login" replace />
  return <Outlet />
}
