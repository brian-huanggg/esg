import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { api } from '../api/client'
import { AuthProvider } from './AuthProvider'
import RequireAuth from './RequireAuth'

const ACCOUNT =
  '{"id":"7c9e6679-7425-40de-944b-e07fc1f90ae7","employeeId":"123456","name":"Test","createdAt":"2026-01-01T00:00:00Z"}'

function mockFetch(meStatus: number) {
  vi.stubGlobal(
    'fetch',
    vi.fn((input: RequestInfo | URL) => {
      const url = input instanceof Request ? input.url : String(input)
      const isMe = url.endsWith('/api/accounts/me')
      const status = isMe ? meStatus : 401
      const body = isMe && meStatus === 200 ? ACCOUNT : '{}'
      return Promise.resolve(
        new Response(body, {
          status,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
    }),
  )
}

function renderProtected() {
  return render(
    <MemoryRouter initialEntries={['/']}>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<p>Login page</p>} />
          <Route element={<RequireAuth />}>
            <Route
              path="/"
              element={
                <>
                  <p>Secret page</p>
                  <button onClick={() => void api.GET('/api/health')}>
                    Call API
                  </button>
                </>
              }
            />
          </Route>
        </Routes>
      </AuthProvider>
    </MemoryRouter>,
  )
}

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('RequireAuth', () => {
  it('redirects an anonymous visitor to the login page', async () => {
    mockFetch(401)

    renderProtected()

    expect(await screen.findByText('Login page')).toBeTruthy()
    expect(screen.queryByText('Secret page')).toBeNull()
  })

  it('shows the protected page to a logged-in user', async () => {
    mockFetch(200)

    renderProtected()

    expect(await screen.findByText('Secret page')).toBeTruthy()
  })

  it('redirects to login when a protected call returns 401', async () => {
    mockFetch(200)
    renderProtected()
    fireEvent.click(await screen.findByRole('button', { name: 'Call API' }))

    expect(await screen.findByText('Login page')).toBeTruthy()
  })
})
