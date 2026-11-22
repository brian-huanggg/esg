import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { api } from '../api/client'
import type { components } from '../api/schema'

type Proposal = components['schemas']['Proposal']

type State =
  | { kind: 'loading' }
  | { kind: 'ok'; proposal: Proposal }
  | { kind: 'notFound' }
  | { kind: 'error'; message: string }

export default function ProposalDetail() {
  const { id } = useParams<{ id: string }>()
  const [state, setState] = useState<State>({ kind: 'loading' })

  useEffect(() => {
    if (!id) return
    let cancelled = false
    api
      .GET('/api/proposals/{id}', { params: { path: { id } } })
      .then(({ data, error, response }) => {
        if (cancelled) return
        if (data) setState({ kind: 'ok', proposal: data })
        // Backend answers 404 for both missing and not-yours.
        else if (response.status === 404) setState({ kind: 'notFound' })
        else setState({ kind: 'error', message: error?.message ?? 'Failed' })
      })
      .catch(() => {
        if (!cancelled) setState({ kind: 'error', message: 'Network error' })
      })
    return () => {
      cancelled = true
    }
  }, [id])

  if (state.kind === 'loading') return <p>Loading...</p>
  if (state.kind === 'notFound') return <p>Not found</p>
  if (state.kind === 'error') return <p role="alert">{state.message}</p>

  const { proposal } = state
  return (
    <article>
      <h1>{proposal.title}</h1>
      <p>{proposal.description}</p>
      <p>Created at {proposal.createdAt}</p>
      <Link to="/">Back</Link>
    </article>
  )
}
