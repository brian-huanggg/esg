import { afterEach, describe, expect, it, vi } from 'vitest'
import { api, setUnauthorizedHandler } from './client'

function stubFetch(status: number) {
  // A fresh Response per call: a body can only be read once
  vi.stubGlobal(
    'fetch',
    vi.fn(() =>
      Promise.resolve(
        new Response('{}', {
          status,
          headers: { 'Content-Type': 'application/json' },
        }),
      ),
    ),
  )
}

afterEach(() => {
  setUnauthorizedHandler(null)
  vi.unstubAllGlobals()
})

describe('401 handler', () => {
  it('fires on a 401 from a protected path', async () => {
    const handler = vi.fn()
    setUnauthorizedHandler(handler)
    stubFetch(401)

    await api.GET('/api/health')

    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('ignores the expected 401 from POST /api/sessions', async () => {
    const handler = vi.fn()
    setUnauthorizedHandler(handler)
    stubFetch(401)

    await api.POST('/api/sessions', {
      body: { employeeId: '123456', password: 'wrong' },
    })

    expect(handler).not.toHaveBeenCalled()
  })

  it('ignores the expected 401 from GET /api/accounts/me', async () => {
    const handler = vi.fn()
    setUnauthorizedHandler(handler)
    stubFetch(401)

    await api.GET('/api/accounts/me')

    expect(handler).not.toHaveBeenCalled()
  })

  it('does not fire on a non-401 response', async () => {
    const handler = vi.fn()
    setUnauthorizedHandler(handler)
    stubFetch(200)

    await api.GET('/api/health')

    expect(handler).not.toHaveBeenCalled()
  })
})
