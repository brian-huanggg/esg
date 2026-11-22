import { useState } from 'react'
import { useNavigate } from 'react-router'
import { api } from '../api/client'

type FieldErrors = Partial<Record<'title' | 'description', string>>

function validate(v: { title: string; description: string }): FieldErrors {
  const errors: FieldErrors = {}
  if (v.title.length < 1 || v.title.length > 100) {
    errors.title = 'Title must be 1-100 characters'
  }
  if (v.description.length > 2000) {
    errors.description = 'Description must be at most 2000 characters'
  }
  return errors
}

export default function CreateProposal() {
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({})

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault()
    setFormError('')
    setFieldErrors({})

    const local = validate({ title, description })
    if (Object.keys(local).length > 0) {
      setFieldErrors(local)
      return
    }

    setSubmitting(true)
    try {
      const { data, error } = await api.POST('/api/proposals', {
        body: { title, description },
      })
      if (data) {
        // Use the body id: Location is an API path, not a frontend route.
        void navigate(`/proposals/${data.id}`)
        return
      }
      if (error.code === 'DUPLICATE_TITLE') {
        setFieldErrors({ title: error.message })
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
    setSubmitting(false)
  }

  return (
    <form onSubmit={handleSubmit}>
      <h1>New proposal</h1>
      <label>
        Title
        <input value={title} onChange={(e) => setTitle(e.target.value)} />
      </label>
      {fieldErrors.title && <p role="alert">{fieldErrors.title}</p>}
      <label>
        Description
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
      </label>
      {fieldErrors.description && <p role="alert">{fieldErrors.description}</p>}
      {formError && <p role="alert">{formError}</p>}
      <button type="submit" disabled={submitting}>
        {submitting ? 'Submitting...' : 'Create'}
      </button>
    </form>
  )
}
