import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router'
import CreateProposal from './CreateProposal'

function mockJson(body: unknown, status: number) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify(body), {
      status,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function fillAndSubmit(title: string) {
  render(
    <MemoryRouter initialEntries={['/proposals/new']}>
      <Routes>
        <Route path="/proposals/new" element={<CreateProposal />} />
        <Route path="/proposals/:id" element={<p>detail page</p>} />
      </Routes>
    </MemoryRouter>,
  )
  fireEvent.change(screen.getByLabelText('Title'), { target: { value: title } })
  fireEvent.change(screen.getByLabelText('Description'), {
    target: { value: 'Solar roof' },
  })
  fireEvent.click(screen.getByRole('button', { name: 'Create' }))
}

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('CreateProposal', () => {
  it('blocks an empty title without calling the API', async () => {
    const fetchMock = mockJson({}, 201)

    fillAndSubmit('')

    expect(
      await screen.findByText('Title must be 1-100 characters'),
    ).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('navigates to the detail page after 201', async () => {
    mockJson(
      {
        id: '11111111-1111-1111-1111-111111111111',
        title: 'Solar',
        description: 'Solar roof',
        createdAt: '2026-01-01T00:00:00Z',
        updatedAt: '2026-01-01T00:00:00Z',
      },
      201,
    )

    fillAndSubmit('Solar')

    expect(await screen.findByText('detail page')).toBeTruthy()
  })

  it('shows a 409 as a title field error', async () => {
    mockJson(
      {
        code: 'DUPLICATE_TITLE',
        message: 'Title already used',
        requestId: 'r1',
      },
      409,
    )

    fillAndSubmit('Solar')

    expect(await screen.findByText('Title already used')).toBeTruthy()
  })
})
