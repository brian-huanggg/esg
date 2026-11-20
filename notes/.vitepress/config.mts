import { defineConfig } from 'vitepress'
import { readdirSync } from 'node:fs'

// Derive the sidebar from filenames so a new note shows up without touching this config.
// Resolve relative to this file so it works regardless of the cwd.
const notesDir = new URL('..', import.meta.url)

const items = readdirSync(notesDir)
    .filter((f) => /^\d+\.\d+.*\.md$/.test(f))
    // numeric: true makes 1.10 sort after 1.9 instead of after 1.1.
    .sort((a, b) => a.localeCompare(b, 'en', { numeric: true }))
    .map((f) => {
        const name = f.replace(/\.md$/, '')
        return { text: name, link: `/${name}` }
    })

export default defineConfig ({
    base: '/esg/',
    title: 'Learning Notes',
    lang: 'zh-TW',
    themeConfig: {
        sidebar: items,
        search: { provider: 'local' },
    },
})