import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import css from 'highlight.js/lib/languages/css'
import go from 'highlight.js/lib/languages/go'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import python from 'highlight.js/lib/languages/python'
import sql from 'highlight.js/lib/languages/sql'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import yaml from 'highlight.js/lib/languages/yaml'

const languages: Array<[string, (hljs: any) => any]> = [
  ['bash', bash],
  ['sh', bash],
  ['shell', bash],
  ['css', css],
  ['go', go],
  ['javascript', javascript],
  ['js', javascript],
  ['json', json],
  ['python', python],
  ['py', python],
  ['sql', sql],
  ['typescript', typescript],
  ['ts', typescript],
  ['html', xml],
  ['xml', xml],
  ['yaml', yaml],
  ['yml', yaml]
]

for (const [name, language] of languages) {
  hljs.registerLanguage(name, language)
}

export default hljs

/**
 * markdown-it highlight callback — wraps output in <pre class="hljs"><code>
 */
export function markdownHighlight(str: string, lang: string): string {
  if (lang && hljs.getLanguage(lang)) {
    try {
      return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang }).value}</code></pre>`
    } catch { /* fall through */ }
  }
  // escapeHtml via a one-time div trick to avoid importing MarkdownIt just for escaping
  return `<pre class="hljs"><code>${str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')}</code></pre>`
}
