import { useState } from 'react'
import { api } from '../api/client'

type Status = 'idle' | 'submitting' | 'done'
type FieldErrors = Partial<Record<'employeeId' | 'name' | 'password', string>>

function validate(v: {
  employeeId: string
  name: string
  password: string
}): FieldErrors {
  const errors: FieldErrors = {}
  if (!/^[0-9]{6}$/.test(v.employeeId)) {
    errors.employeeId = 'Employee ID must be 6 digits'
  }
  if (v.name.length < 1 || v.name.length > 100) {
    errors.name = 'Name must be 1-100 characters'
  }
  if (v.password.length < 1 || v.password.length > 72) {
    errors.password = 'Password must be 1-72 characters'
  }
  return errors
}

export default function Register() {
  const [employeeId, setEmployeeId] = useState('')
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [status, setStatus] = useState<Status>('idle')
  const [formError, setFormError] = useState('')
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({})

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault()
    setFormError('')
    setFieldErrors({})

    const local = validate({ employeeId, name, password })
    if (Object.keys(local).length > 0) {
      setFieldErrors(local)
      return
    }

    setStatus('submitting')
    try {
      const { error } = await api.POST('/api/accounts', {
        body: { employeeId, name, password },
      })
      if (!error) {
        setStatus('done')
        return
      }
      if (error.code === 'EMPLOYEE_ID_TAKEN') {
        setFieldErrors({ employeeId: error.message })
      } else if (error.code === 'VALIDATION_FAILED' && error.fields) {
        const server: FieldErrors = {}
        for (const f of error.fields) {
          server[f.field as keyof FieldErrors] = f.message
        }
        setFieldErrors(server)
      } else {
        setFormError(error.message)
      }
    } catch {
      // fetch itself failed (offline, DNS): openapi-fetch throws here
      setFormError('Network error, please try again')
    }
    setStatus('idle')
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
      {fieldErrors.employeeId && <p role="alert">{fieldErrors.employeeId}</p>}
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
        {fieldErrors.name && <p role="alert">{fieldErrors.name}</p>}
      </label>
      {fieldErrors.password && <p role="alert">{fieldErrors.password}</p>}
      {formError && <p role="alert">{formError}</p>}
      <button type="submit" disabled={status === 'submitting'}>
        {status === 'submitting' ? 'Submitting...' : 'Register'}
      </button>
    </form>
  )
}
