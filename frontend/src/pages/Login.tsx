import { useState } from 'react'
import { useNavigate } from 'react-router'
import { api } from '../api/client'
import { useAuth } from '../auth/context'

export default function Login() {
  const [employeeId, setEmployeeId] = useState('')
  const [password, setPassword] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login } = useAuth()

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault()
    setError('')
    setSubmitting(true)
    try {
      const { data, error: apiError } = await api.POST('/api/sessions', {
        body: { employeeId, password },
      })
      if (data) {
        login(data)
        navigate('/')
        return
      }
      // 401 and 422 both carry a user-facing message from the backend
      setError(apiError.message)
    } catch {
      setError('Network error, please try again')
    }
    setSubmitting(false)
  }

  return (
    <form onSubmit={handleSubmit}>
      <h1>Login</h1>
      <label>
        Employee ID
        <input
          value={employeeId}
          onChange={(e) => setEmployeeId(e.target.value)}
        />
      </label>
      <label>
        Password
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </label>
      {error && <p role="alert">{error}</p>}
      <button type="submit" disabled={submitting}>
        {submitting ? 'Logging in...' : 'Login'}
      </button>
    </form>
  )
}
