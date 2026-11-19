import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import App from './App'

function mockFetch(body: string, status: number) {
    const fetchMock = vi.fn().mockResolvedValue(new Response(body, { status }))
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

        render(<App />)

        expect(await screen.findByText('Backend: ok')).toBeTruthy()
        expect(fetchMock).toHaveBeenCalledWith('/api/health')
    })

    it('shows an error when the backend responds with a failure', async () => {
        mockFetch('', 502)

        render(<App />)
        expect(await screen.findByText(/Backend: error .*502/)).toBeTruthy()
    })
})