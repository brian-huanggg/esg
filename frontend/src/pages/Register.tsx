import { useState } from 'react'
import { api } from '../api/client'

type Status = 'idle' | 'submitting' | 'done'

export default function Register() {
  const [employeeId, setEmployeeId] = useState('')
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [status, setStatus] = useState<Status>('idle')
  const [error, setError] = useState('')

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault()
    setStatus('submitting')
    setError('')

    const { error } = await api.POST('/api/accounts', {
      body: { employeeId, name, password },
    })

    if (error) {
      // Field-level errors come in step 2.10
      setError(error.message)
      setStatus('idle')
      return
    }
    setStatus('done')
  }

  if (status === 'done') {
    return <p>Registration succeeded.</p>
  }

  return (
    <form onSubmit={handleSubmit}>
      <h1>Register</h1>
      <label>
        Employee ID
        <input
          value={employeeId}
          onChange={(e) => setEmployeeId(e.target.value)}
        />
      </label>
      <label>
        Name
        <input value={name} onChange={(e) => setName(e.target.value)} />
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
      <button type="submit" disabled={status === 'submitting'}>
        {status === 'submitting' ? 'Submitting...' : 'Register'}
      </button>
    </form>
  )
}
