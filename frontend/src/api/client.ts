import type { components } from './schema'

export type Health = components['schemas']['Health']

export async function getHealth(): Promise<Health> {
    const res = await fetch('/api/health')
    if (!res.ok) {
        throw new Error(`GET /api/health failed: ${res.status}`)
    }
    return (await res.json()) as Health
}