import createClient from 'openapi-fetch'
import type { components, paths } from './schema'

export const api = createClient<paths>({
  baseUrl: window.location.origin,
  fetch: (req) => globalThis.fetch(req),
})

const SAFE_METHODS = new Set(['GET', 'HEAD', 'OPTIONS'])

function readCookie(name: string): string | undefined {
  const prefix = `${name}=`
  const pair = document.cookie.split('; ').find((c) => c.startsWith(prefix))
  return pair?.slice(prefix.length)
}

const EXPECTED_401_PATHS = new Set(['/api/sessions', '/api/accounts/me'])

let onUnauthorized: (() => void) | null = null

export function setUnauthorizedHandler(handler: (() => void) | null) {
  onUnauthorized = handler
}

api.use({
  onRequest({ request }) {
    if (SAFE_METHODS.has(request.method)) return
    const token = readCookie('csrf_token')
    if (token) request.headers.set('X-CSRF-Token', token)
  },
  onResponse({ request, response }) {
    if (response.status !== 401) return
    if (EXPECTED_401_PATHS.has(new URL(request.url).pathname)) return
    onUnauthorized?.()
  },
})

export type Health = components['schemas']['Health']

export async function getHealth(): Promise<Health> {
  const res = await fetch('/api/health')
  if (!res.ok) {
    throw new Error(`GET /api/health failed: ${res.status}`)
  }
  return (await res.json()) as Health
}
