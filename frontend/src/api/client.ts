import createClient from 'openapi-fetch'
import type { components, paths } from './schema'

// Request() rejects relative URLs outside a browser (e.g. in vitest),
// so anchor to the current origin. Same-origin in prod, so no CORS.
// fetch is looked up per call because createClient would otherwise
// capture globalThis.fetch once, and tests could not stub it.
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

api.use({
  onRequest({ request }) {
    if (SAFE_METHODS.has(request.method)) return
    const token = readCookie('csrf_token')
    if (token) request.headers.set('X-CSRF-Token', token)
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
