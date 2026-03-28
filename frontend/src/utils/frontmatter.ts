export interface FrontmatterMeta {
  featured?: boolean
  [key: string]: unknown
}

/**
 * Parse YAML frontmatter delimited by `---`.
 * Only handles simple `key: value` lines — no nested structures.
 */
export function parseFrontmatter(content: string): { meta: FrontmatterMeta; body: string } {
  if (!content.startsWith('---')) {
    return { meta: {}, body: content }
  }

  const end = content.indexOf('\n---', 3)
  if (end === -1) {
    return { meta: {}, body: content }
  }

  const rawMeta = content.slice(4, end) // skip opening "---\n"
  const body = content.slice(end + 4).replace(/^\n/, '') // skip closing "---\n"

  const meta: FrontmatterMeta = {}
  for (const line of rawMeta.split('\n')) {
    const colon = line.indexOf(':')
    if (colon === -1) continue
    const key = line.slice(0, colon).trim()
    const val = line.slice(colon + 1).trim()
    if (key === 'featured') {
      meta.featured = val === 'true'
    } else {
      meta[key] = val
    }
  }

  return { meta, body }
}

/**
 * Serialize meta + body back to a markdown string with frontmatter.
 * Omits keys with undefined values. Omits the frontmatter block entirely if meta is empty.
 */
export function stringifyFrontmatter(meta: FrontmatterMeta, body: string): string {
  const entries = Object.entries(meta).filter(([, v]) => v !== undefined)
  if (entries.length === 0) return body

  const lines = entries.map(([k, v]) => `${k}: ${v}`)
  return `---\n${lines.join('\n')}\n---\n${body}`
}
