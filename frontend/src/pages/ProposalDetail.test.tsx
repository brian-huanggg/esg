import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router'
import ProposalDetail from './ProposalDetail'

function mockJson(body: unknown, status: number) {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue(
      new Response(JSON.stringify(body), {
        status,
        headers: { 'Content-Type': 'application/json' },
      }),
    ),
  )
}

function renderAt(id: string) {
  render(
    <MemoryRouter initialEntries={[`/proposals/${id}`]}>
      <Routes>
        <Route path="/proposals/:id" element={<ProposalDetail />} />
      </Routes>
    </MemoryRouter>,
  )
}

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('ProposalDetail', () => {
  it('shows the proposal on 200', async () => {
    mockJson(
      {
        id: '11111111-1111-1111-1111-111111111111',
        title: 'Solar',
        description: 'Solar roof',
        createdAt: '2026-01-01T00:00:00Z',
        updatedAt: '2026-01-01T00:00:00Z',
      },
      200,
    )

    renderAt('11111111-1111-1111-1111-111111111111')

    expect(await screen.findByText('Solar')).toBeTruthy()
    expect(screen.getByText('Solar roof')).toBeTruthy()
  })

  it('shows Not found on 404', async () => {
    mockJson({ code: 'NOT_FOUND', message: 'Not found', requestId: 'r1' }, 404)

    renderAt('22222222-2222-2222-2222-222222222222')

    expect(await screen.findByText('Not found')).toBeTruthy()
  })
})
