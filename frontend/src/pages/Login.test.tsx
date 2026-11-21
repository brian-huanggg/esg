import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { AuthProvider } from '../auth/AuthProvider'
import Login from './Login'

const ACCOUNT =
  '{"id":"7c9e6679-7425-40de-944b-e07fc1f90ae7","employeeId":"123456","name":"Test","createdAt":"2026-01-01T00:00:00Z"}'

// /me is the on-load "am I logged in" check and always answers 401 here;
// only the POST /api/sessions answer changes between tests.
function mockFetch(sessionBody: string, sessionStatus: number) {
  vi.stubGlobal(
    'fetch',
    vi.fn((input: RequestInfo | URL) => {
      const url = input instanceof Request ? input.url : String(input)
      const isMe = url.endsWith('/api/accounts/me')
      return Promise.resolve(
        new Response(isMe ? '{}' : sessionBody, {
          status: isMe ? 401 : sessionStatus,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
    }),
  )
}

function fillAndSubmit() {
  render(
    <MemoryRouter initialEntries={['/login']}>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<p>Home page</p>} />
        </Routes>
      </AuthProvider>
    </MemoryRouter>,
  )
  fireEvent.change(screen.getByLabelText('Employee ID'), {
    target: { value: '123456' },
  })
  fireEvent.change(screen.getByLabelText('Password'), {
    target: { value: 'secret' },
  })
  fireEvent.click(screen.getByRole('button', { name: 'Login' }))
}

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('Login', () => {
  it('shows the backend message when the password is wrong', async () => {
    mockFetch(
      '{"code":"INVALID_CREDENTIALS","message":"Wrong employee ID or password","requestId":"r1"}',
      401,
    )

    fillAndSubmit()

    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toBe('Wrong employee ID or password')
    expect(screen.queryByText('Home page')).toBeNull()
  })

  it('goes to the home page after a successful login', async () => {
    mockFetch(ACCOUNT, 200)

    fillAndSubmit()

    expect(await screen.findByText('Home page')).toBeTruthy()
  })
})
