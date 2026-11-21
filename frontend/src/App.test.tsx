import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import App from './App'
import { MemoryRouter } from 'react-router'
import { AuthProvider } from './auth/AuthProvider'

const ACCOUNT =
  '{"id":"7c9e6679-7425-40de-944b-e07fc1f90ae7","employeeId":"123456","name":"Test","createdAt":"2026-01-01T00:00:00Z"}'

function renderApp() {
  return render(
    <MemoryRouter initialEntries={['/']}>
      <AuthProvider>
        <App />
      </AuthProvider>
    </MemoryRouter>,
  )
}

function mockFetch(body: string, status: number) {
  const fetchMock = vi.fn((input: RequestInfo | URL) => {
    const url = input instanceof Request ? input.url : String(input)
    if (url.endsWith('/api/accounts/me')) {
      return Promise.resolve(new Response(ACCOUNT, { status: 200 }))
    }
    return Promise.resolve(new Response(body, { status }))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('App', () => {
  it('show backend status when the health check succeeds', async () => {
    const fetchMock = mockFetch('{"status": "ok"}', 200)

    renderApp()

    expect(await screen.findByText('Backend: ok')).toBeTruthy()
    expect(fetchMock).toHaveBeenCalledWith('/api/health')
  })

  it('shows an error when the backend responds with a failure', async () => {
    mockFetch('', 502)

    renderApp()
    expect(await screen.findByText(/Backend: error .*502/)).toBeTruthy()
  })
})
