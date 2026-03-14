<template>
  <div ref="detailRef" class="note-detail" :class="detailClass">
    <article class="surface-panel note-sheet">
      <header class="note-sheet__header">
        <div>
          <p class="meta-label">Reading Room / {{ breadcrumb }}</p>
          <h1 class="note-sheet__title">{{ noteTitle }}</h1>
        </div>
        <div class="note-sheet__meta">
          <span>{{ estimatedReadLabel }}</span>
          <span>{{ wordCount }} chars</span>
          <div class="note-sheet__view-modes" role="group" aria-label="阅读宽度">
            <button
              type="button"
              class="note-sheet__view-btn"
              :class="{ 'note-sheet__view-btn--active': readingWidthMode === 'standard' }"
              @click="setReadingWidthMode('standard')"
            >
              标准
            </button>
            <button
              type="button"
              class="note-sheet__view-btn"
              :class="{ 'note-sheet__view-btn--active': readingWidthMode === 'wide' }"
              @click="setReadingWidthMode('wide')"
            >
              宽屏
            </button>
          </div>
        </div>
      </header>

      <div v-if="loading" class="note-state">
        <div class="note-state__spinner"></div>
        <p>正在排版内容…</p>
      </div>

      <div v-else-if="!content" class="note-state">
        <p>这篇笔记暂时不存在。</p>
        <router-link to="/notes" class="note-state__link">返回笔记目录</router-link>
      </div>

      <div v-else ref="articleRef" class="note-prose" v-html="renderedContent"></div>
    </article>

    <aside ref="railRef" class="note-rail" :style="railShellStyle">
      <div ref="railStickyRef" class="note-rail__sticky" :class="railStickyClass" :style="railStickyStyle">
        <section class="surface-panel note-rail__card">
          <div class="note-rail__card-head">
            <span class="meta-label">Progress</span>
            <button type="button" class="note-rail__collapse" @click="toggleRailCollapsed">
              {{ effectiveRailCollapsed ? '展开' : '缩小' }}
            </button>
          </div>
          <div class="note-progress">
            <div v-if="!effectiveRailCollapsed" class="note-progress__track">
              <span class="note-progress__bar" :style="{ width: `${readingProgress}%` }"></span>
            </div>
            <strong>{{ Math.round(readingProgress) }}%</strong>
          </div>
          <router-link v-if="!effectiveRailCollapsed" to="/notes" class="note-rail__link">回到目录索引</router-link>
        </section>

        <section v-if="toc.length && !effectiveRailCollapsed" class="surface-panel note-rail__card">
          <span class="meta-label">On this page</span>
          <nav class="note-toc" aria-label="页面目录">
            <button
              v-for="item in toc"
              :key="item.id"
              type="button"
              class="note-toc__item"
              :class="{ 'note-toc__item--active': activeHeading === item.id }"
              :style="{ '--toc-level': String(item.level) }"
              @click="scrollToHeading(item.id)"
            >
              {{ item.text }}
            </button>
          </nav>
        </section>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import MarkdownIt from 'markdown-it'
import hljs from '@/utils/highlight'
import { noteApi } from '@/api/note'
import 'highlight.js/styles/github.css'

const route = useRoute()
const content = ref('')
const loading = ref(true)
const displayTitle = ref('')
const detailRef = ref<HTMLElement | null>(null)
const articleRef = ref<HTMLElement | null>(null)
const railRef = ref<HTMLElement | null>(null)
const railStickyRef = ref<HTMLElement | null>(null)
const toc = ref<{ id: string; text: string; level: number }[]>([])
const activeHeading = ref('')
const readingProgress = ref(0)
const railMode = ref<'static' | 'fixed' | 'bottom'>('static')
const railWidth = ref(0)
const railLeft = ref(0)
const railHeight = ref(0)
const readingWidthMode = ref<'standard' | 'wide'>('standard')
const railCollapsed = ref(false)
const isCompactViewport = ref(false)
const protectedClipboardText = '内容受保护，请勿复制。'
const NOTE_READING_WIDTH_KEY = 'mdxy.note.reading-width'
const NOTE_RAIL_COLLAPSED_KEY = 'mdxy.note.rail-collapsed'
const NOTES_WIDTH_EVENT = 'mdxy:notes-reading-width'
let proseContextMenuHandler: ((event: Event) => void) | null = null
let proseCopyHandler: ((event: Event) => void) | null = null
let proseCutHandler: ((event: Event) => void) | null = null
let proseDragStartHandler: ((event: Event) => void) | null = null

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight: (str, lang) => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch {
        return ''
      }
    }
    return ''
  }
})

const renderedContent = computed(() => {
  return content.value ? md.render(content.value) : ''
})

const rawPath = computed(() => (Array.isArray(route.params.path) ? route.params.path.join('/') : route.params.path || ''))
const noteTitle = computed(() => {
  if (displayTitle.value) {
    return displayTitle.value
  }

  const segments = rawPath.value.split('/').filter(Boolean)
  const last = segments[segments.length - 1] || '未命名笔记'
  return last.replace(/\.md$/i, '')
})
const breadcrumb = computed(() => rawPath.value.replace(/\.md$/i, '') || 'index')
const plainText = computed(() => extractReadableText(content.value))
const wordCount = computed(() => plainText.value.length)
const estimatedReadLabel = computed(() => {
  const minutes = wordCount.value / 350

  if (minutes < 1) {
    return '<1 min read'
  }

  return `${Math.round(minutes)} min read`
})
const effectiveReadingWidthMode = computed<'standard' | 'wide'>(() => {
  return isCompactViewport.value ? 'standard' : readingWidthMode.value
})
const effectiveRailCollapsed = computed(() => {
  return !isCompactViewport.value && railCollapsed.value
})
const detailClass = computed(() => ({
  'note-detail--wide': effectiveReadingWidthMode.value === 'wide',
  'note-detail--rail-collapsed': effectiveRailCollapsed.value
}))
const railStickyClass = computed(() => ({
  'note-rail__sticky--fixed': railMode.value === 'fixed',
  'note-rail__sticky--bottom': railMode.value === 'bottom'
}))
const railStickyStyle = computed(() => {
  if (railMode.value !== 'fixed') {
    return {}
  }

  return {
    left: `${railLeft.value}px`,
    width: `${railWidth.value}px`
  }
})
const railShellStyle = computed(() => {
  if (railHeight.value <= 0) {
    return {}
  }

  return {
    minHeight: `${railHeight.value}px`
  }
})

async function fetchContent(path: string) {
  loading.value = true
  try {
    const response = await noteApi.getContent(path)
    if (response.success && response.data) {
      const parsed = extractDisplayContent(normalizeMarkdownForDisplay(response.data.content))
      displayTitle.value = parsed.title
      content.value = parsed.body
      await nextTick()
      decorateRenderedArticle()
      buildToc()
      bindNoteProtection()
      handleViewportChange()
    }
  } catch (error) {
    console.error('获取笔记内容失败:', error)
    content.value = ''
    displayTitle.value = ''
    toc.value = []
    unbindNoteProtection()
  } finally {
    loading.value = false
  }
}

function extractDisplayContent(markdown: string) {
  const normalized = markdown.replace(/\r\n/g, '\n')
  const match = normalized.match(/^\s*#\s+(.+?)\s*\n+/)

  if (!match) {
    return {
      title: '',
      body: markdown
    }
  }

  return {
    title: match[1].trim(),
    body: normalized.slice(match[0].length).trimStart()
  }
}

function normalizeMarkdownForDisplay(markdown: string) {
  const invisibleChars = /[\u00a0\u200b\u200c\u200d\ufeff\u2060\u3000]/g

  return markdown
    .replace(/\r\n/g, '\n')
    .split('\n')
    .map((line) => {
      const cleaned = line.replace(invisibleChars, '')
      return cleaned.trim() ? cleaned : ''
    })
    .join('\n')
}

function extractReadableText(markdown: string) {
  if (!markdown) {
    return ''
  }

  const tokens = md.parse(markdown, {})
  let text = ''

  for (const token of tokens) {
    if (token.type === 'fence' || token.type === 'code_block') {
      text += token.content || ''
      continue
    }

    if (token.type !== 'inline') {
      continue
    }

    for (const child of token.children || []) {
      if (child.type === 'text' || child.type === 'code_inline' || child.type === 'image') {
        text += child.content || ''
      }
    }
  }

  const matches = text.match(/[\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaffA-Za-z]/g)
  return matches ? matches.join('') : ''
}

function decorateRenderedArticle() {
  if (!articleRef.value) {
    return
  }

  const root = articleRef.value

  // Remove empty spacer paragraphs that may slip in from pasted markdown/html.
  Array.from(root.querySelectorAll('p')).forEach((paragraph) => {
    const element = paragraph as HTMLParagraphElement
    const normalizedText = (element.textContent || '').replace(/[\s\u00a0\u200b\u200c\u200d\ufeff\u2060\u3000]/g, '')
    const onlySoftBreaks = Array.from(element.childNodes).every((node) => {
      if (node.nodeType === Node.TEXT_NODE) {
        return !(node.textContent || '').replace(/[\s\u00a0\u200b\u200c\u200d\ufeff\u2060\u3000]/g, '')
      }

      return node instanceof HTMLBRElement
    })

    if (!normalizedText && onlySoftBreaks) {
      element.remove()
    }
  })

  const blocks = Array.from(root.children) as HTMLElement[]

  blocks.forEach((block, index) => {
    block.classList.add('note-block')

    if (index === 0) {
      block.classList.add('note-block--first')
    }

    if (/^H[1-4]$/.test(block.tagName)) {
      block.classList.add('note-heading', `note-heading--level-${block.tagName.slice(1)}`)
    }
  })
}

function bindNoteProtection() {
  unbindNoteProtection()
  if (!articleRef.value) {
    return
  }

  const prose = articleRef.value

  proseContextMenuHandler = (event) => {
    event.preventDefault()
  }

  proseCopyHandler = (event) => {
    event.preventDefault()
    const copyEvent = event as ClipboardEvent
    copyEvent.clipboardData?.setData('text/plain', protectedClipboardText)
  }

  proseCutHandler = (event) => {
    event.preventDefault()
  }

  proseDragStartHandler = (event) => {
    const dragEvent = event as DragEvent
    const target = dragEvent.target as Element | null
    if (target?.closest('img, pre, code, table')) {
      event.preventDefault()
    }
  }

  prose.addEventListener('contextmenu', proseContextMenuHandler)
  prose.addEventListener('copy', proseCopyHandler)
  prose.addEventListener('cut', proseCutHandler)
  prose.addEventListener('dragstart', proseDragStartHandler)
}

function unbindNoteProtection() {
  if (!articleRef.value) {
    return
  }

  if (proseContextMenuHandler) {
    articleRef.value.removeEventListener('contextmenu', proseContextMenuHandler)
    proseContextMenuHandler = null
  }

  if (proseCopyHandler) {
    articleRef.value.removeEventListener('copy', proseCopyHandler)
    proseCopyHandler = null
  }

  if (proseCutHandler) {
    articleRef.value.removeEventListener('cut', proseCutHandler)
    proseCutHandler = null
  }

  if (proseDragStartHandler) {
    articleRef.value.removeEventListener('dragstart', proseDragStartHandler)
    proseDragStartHandler = null
  }
}

function buildToc() {
  if (!articleRef.value) {
    toc.value = []
    return
  }

  const headings = Array.from(articleRef.value.querySelectorAll('h1, h2, h3')) as HTMLElement[]
  toc.value = headings.map((heading, index) => {
    const id = `section-${index + 1}`
    heading.id = id
    return {
      id,
      text: heading.innerText.trim(),
      level: Number(heading.tagName.slice(1))
    }
  })
  activeHeading.value = toc.value[0]?.id || ''
}

function updateReadingState() {
  if (!articleRef.value || !content.value) {
    readingProgress.value = 0
    return
  }

  const rect = articleRef.value.getBoundingClientRect()
  const top = window.scrollY + rect.top
  const height = articleRef.value.offsetHeight
  const viewport = window.innerHeight
  const current = window.scrollY + viewport * 0.25 - top
  const progress = height <= 0 ? 0 : (current / Math.max(height - viewport * 0.6, 1)) * 100
  readingProgress.value = Math.max(0, Math.min(100, progress))

  const headingPositions = toc.value.map((item) => {
    const target = document.getElementById(item.id)
    return {
      id: item.id,
      top: target ? target.getBoundingClientRect().top : Infinity
    }
  })

  const currentHeading = headingPositions
    .filter((item) => item.top <= 140)
    .sort((a, b) => b.top - a.top)[0]

  if (currentHeading) {
    activeHeading.value = currentHeading.id
  } else if (toc.value[0]) {
    activeHeading.value = toc.value[0].id
  }
}

function getRailOffsetTop() {
  const rootStyles = getComputedStyle(document.documentElement)
  const headerHeight = Number.parseFloat(rootStyles.getPropertyValue('--header-height')) || 84
  const rootFontSize = Number.parseFloat(rootStyles.fontSize) || 16
  return headerHeight + rootFontSize * 1.2
}

function syncRailMetrics() {
  if (!railRef.value || !railStickyRef.value) {
    railWidth.value = 0
    railLeft.value = 0
    railHeight.value = 0
    return
  }

  const rect = railRef.value.getBoundingClientRect()
  railWidth.value = rect.width
  railLeft.value = rect.left
  railHeight.value = railStickyRef.value.offsetHeight
}

function updateRailPosition() {
  syncRailMetrics()

  if (!detailRef.value || !railStickyRef.value || window.innerWidth <= 980) {
    railMode.value = 'static'
    return
  }

  const offsetTop = getRailOffsetTop()
  const detailRect = detailRef.value.getBoundingClientRect()
  const detailTop = window.scrollY + detailRect.top
  const detailBottom = detailTop + detailRef.value.offsetHeight
  const viewportTop = window.scrollY + offsetTop
  const stickyHeight = railStickyRef.value.offsetHeight

  if (viewportTop <= detailTop) {
    railMode.value = 'static'
    return
  }

  if (viewportTop + stickyHeight >= detailBottom) {
    railMode.value = 'bottom'
    return
  }

  railMode.value = 'fixed'
}

function handleViewportChange() {
  isCompactViewport.value = window.innerWidth <= 980
  updateReadingState()
  updateRailPosition()
}

function scrollToHeading(id: string) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function setReadingWidthMode(mode: 'standard' | 'wide') {
  if (isCompactViewport.value) {
    return
  }

  if (readingWidthMode.value === mode) {
    return
  }

  readingWidthMode.value = mode
  window.localStorage.setItem(NOTE_READING_WIDTH_KEY, mode)
  window.dispatchEvent(new CustomEvent(NOTES_WIDTH_EVENT, { detail: mode }))
  nextTick(handleViewportChange)
}

function toggleRailCollapsed() {
  if (isCompactViewport.value) {
    return
  }

  railCollapsed.value = !railCollapsed.value
  window.localStorage.setItem(NOTE_RAIL_COLLAPSED_KEY, railCollapsed.value ? '1' : '0')
  nextTick(handleViewportChange)
}

function restoreLayoutPreference() {
  const savedWidth = window.localStorage.getItem(NOTE_READING_WIDTH_KEY)
  if (savedWidth === 'standard' || savedWidth === 'wide') {
    readingWidthMode.value = savedWidth
  }

  railCollapsed.value = window.localStorage.getItem(NOTE_RAIL_COLLAPSED_KEY) === '1'
}

onMounted(() => {
  restoreLayoutPreference()
  window.dispatchEvent(new CustomEvent(NOTES_WIDTH_EVENT, { detail: readingWidthMode.value }))
  window.addEventListener('scroll', handleViewportChange, { passive: true })
  window.addEventListener('resize', handleViewportChange)
  nextTick(handleViewportChange)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleViewportChange)
  window.removeEventListener('resize', handleViewportChange)
  unbindNoteProtection()
})

watch(
  () => route.params.path,
  (newPath) => {
    if (newPath) {
      fetchContent(Array.isArray(newPath) ? newPath.join('/') : newPath)
    }
  },
  { immediate: true }
)

watch([readingWidthMode, railCollapsed], () => {
  nextTick(handleViewportChange)
})
</script>

<style scoped>
.note-detail {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 260px;
  gap: var(--space-xl);
  min-width: 0;
}

.note-detail--wide {
  gap: var(--space-2xl);
}

.note-detail--rail-collapsed {
  grid-template-columns: minmax(0, 1fr) 104px;
}

.note-sheet {
  padding: clamp(1.15rem, 2.6vw, 2.1rem);
  min-width: 0;
}

.note-sheet__header {
  display: flex;
  justify-content: space-between;
  gap: var(--space-lg);
  padding-bottom: 1.1rem;
  margin-bottom: 1.2rem;
  border-bottom: 1px solid var(--border-primary);
}

.note-sheet__title {
  margin-top: 0.45rem;
  font-family: var(--font-display);
  font-size: clamp(1.8rem, 3.3vw, 2.8rem);
  line-height: 1.08;
  letter-spacing: -0.04em;
}

.note-sheet__meta {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  align-items: flex-end;
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.note-sheet__view-modes {
  margin-top: 0.45rem;
  display: inline-flex;
  align-items: center;
  padding: 0.2rem;
  gap: 0.2rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  background: var(--bg-panel-muted);
  text-transform: none;
  letter-spacing: normal;
}

.note-sheet__view-btn {
  border: none;
  background: transparent;
  cursor: pointer;
  min-width: 3.7rem;
  padding: 0.26rem 0.62rem;
  border-radius: var(--radius-pill);
  color: var(--text-secondary);
  font-family: var(--font-body);
  font-size: 0.78rem;
  line-height: 1.2;
}

.note-sheet__view-btn:hover {
  color: var(--text-primary);
}

.note-sheet__view-btn--active {
  background: var(--bg-panel-strong);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
}

.note-state {
  min-height: 320px;
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  text-align: center;
}

.note-state__spinner {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  border: 2px solid var(--border-primary);
  border-top-color: var(--accent-primary);
  animation: spin 1s linear infinite;
}

.note-state__link {
  font-weight: 600;
}

.note-prose {
  max-width: min(100%, 82ch);
  color: var(--text-primary);
  font-size: 1rem;
  line-height: 1.78;
  user-select: none;
  -webkit-user-select: none;
  -webkit-touch-callout: none;
}

.note-detail--wide .note-prose {
  max-width: min(100%, 112ch);
  font-size: 1.03rem;
  line-height: 1.82;
}

.note-prose :deep(.note-heading) {
  scroll-margin-top: 140px;
  font-family: var(--font-display);
  line-height: 1.18;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.note-prose :deep(.note-block--first) {
  margin-top: 0 !important;
}

.note-prose :deep(.note-heading--level-1) {
  font-size: clamp(1.75rem, 3vw, 2.45rem);
  margin-top: 0;
  margin-bottom: 0.95rem;
}

.note-prose :deep(.note-heading--level-2) {
  margin-top: 1.8rem;
  margin-bottom: 0.55rem;
  font-size: clamp(1.4rem, 2.4vw, 1.9rem);
  padding-left: 0.75rem;
  border-left: 2px solid var(--accent-primary);
}

.note-prose :deep(.note-heading--level-3) {
  margin-top: 0.95rem;
  margin-bottom: 0.45rem;
  font-size: 1.1rem;
}

.note-prose :deep(.note-heading--level-4) {
  margin-top: 0.8rem;
  margin-bottom: 0.35rem;
  font-size: 1rem;
}

.note-prose :deep(p),
.note-prose :deep(ul),
.note-prose :deep(ol),
.note-prose :deep(blockquote),
.note-prose :deep(table) {
  margin-bottom: 0.95rem;
}

.note-prose :deep(code) {
  padding: 0.08rem 0.38rem;
  border-radius: 0.35rem;
  background: var(--bg-panel-muted);
  border: 1px solid var(--border-soft);
  font-family: var(--font-mono);
  font-size: 0.88em;
}

.note-prose :deep(pre) {
  margin: 0.95rem 0 1.2rem;
  padding: 0.95rem 1.05rem;
  border-radius: var(--radius-lg);
  background: var(--bg-contrast);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08);
  overflow-x: auto;
}

.note-prose :deep(pre code) {
  background: none;
  padding: 0;
  border: none;
  color: #f8f4ec;
  font-size: 0.9rem;
  line-height: 1.65;
}

.note-prose :deep(blockquote) {
  padding: 0.95rem 1.1rem;
  margin: 0.95rem 0 1.1rem;
  border-left: 2px solid var(--accent-secondary);
  background: var(--accent-secondary-soft);
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
  color: var(--text-secondary);
  line-height: 1.72;
}

.note-prose :deep(.note-heading--level-2 + .note-heading--level-3) {
  margin-top: 0.5rem;
}

.note-prose :deep(.note-heading--level-3 + pre),
.note-prose :deep(.note-heading--level-3 + table),
.note-prose :deep(.note-heading--level-3 + p),
.note-prose :deep(.note-heading--level-2 + pre),
.note-prose :deep(.note-heading--level-2 + table),
.note-prose :deep(.note-heading--level-2 + p) {
  margin-top: 0.55rem;
}

.note-prose :deep(a) {
  color: var(--accent-primary);
  text-decoration: underline;
  text-decoration-color: rgba(41, 70, 58, 0.24);
  text-underline-offset: 0.18rem;
}

.note-prose :deep(table) {
  width: 100%;
  border-collapse: collapse;
  overflow: hidden;
  border-radius: var(--radius-lg);
  border-style: hidden;
  box-shadow: 0 0 0 1px var(--border-primary);
}

.note-prose :deep(th),
.note-prose :deep(td) {
  border: 1px solid var(--border-primary);
  padding: 0.72rem 0.85rem;
  text-align: left;
  font-size: 0.94rem;
}

.note-prose :deep(th) {
  background: var(--bg-panel-muted);
  font-weight: 700;
}

.note-prose :deep(img) {
  max-width: 100%;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft);
}

.note-prose :deep(hr) {
  margin: 1.5rem 0;
  border: none;
  border-top: 1px solid var(--border-primary);
}

.note-prose :deep(ul),
.note-prose :deep(ol) {
  padding-left: 1.25rem;
}

.note-prose :deep(li) {
  margin-bottom: 0.25rem;
}

.note-rail {
  position: relative;
  align-self: stretch;
  height: 100%;
  min-width: 0;
}

.note-rail__sticky {
  display: grid;
  gap: var(--space-md);
  align-content: start;
}

.note-rail__sticky--fixed {
  position: fixed;
  top: calc(var(--header-height) + 1.2rem);
}

.note-rail__sticky--bottom {
  position: absolute;
  inset: auto 0 0 0;
}

.note-rail__card {
  padding: var(--space-lg);
}

.note-rail__card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-sm);
}

.note-rail__collapse {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 0.3rem 0.65rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  font-family: var(--font-mono);
  font-size: 0.68rem;
  letter-spacing: 0.08em;
  color: var(--text-secondary);
  background: var(--bg-panel-muted);
}

.note-rail__collapse:hover {
  border-color: var(--border-strong);
  color: var(--text-primary);
}

.note-rail__card + .note-rail__card {
  margin-top: 0;
}

.note-progress {
  margin-top: var(--space-md);
  display: grid;
  gap: var(--space-sm);
}

.note-progress__track {
  height: 0.5rem;
  border-radius: var(--radius-pill);
  background: var(--bg-panel-muted);
  overflow: hidden;
}

.note-progress__bar {
  display: block;
  height: 100%;
  background: var(--line-accent);
}

.note-progress strong {
  font-family: var(--font-display);
  font-size: 2rem;
  letter-spacing: -0.04em;
}

.note-rail__link {
  display: inline-flex;
  margin-top: var(--space-md);
  color: var(--accent-primary);
  font-weight: 600;
}

.note-detail--rail-collapsed .note-rail__card {
  padding: 0.75rem;
}

.note-detail--rail-collapsed .note-rail__card-head {
  justify-content: center;
}

.note-detail--rail-collapsed .note-rail__card-head .meta-label {
  display: none;
}

.note-detail--rail-collapsed .note-progress {
  margin-top: 0.15rem;
  justify-items: center;
}

.note-detail--rail-collapsed .note-progress strong {
  font-size: 1.35rem;
}

.note-detail--rail-collapsed .note-rail__collapse {
  width: 100%;
  justify-content: center;
}

.note-toc {
  margin-top: var(--space-md);
  display: grid;
  gap: 0.45rem;
  max-height: calc(100vh - var(--header-height) - 18rem);
  overflow-y: auto;
  padding-right: 0.2rem;
}

.note-toc__item {
  padding: 0.4rem 0.4rem 0.4rem calc(0.55rem + (var(--toc-level) - 1) * 0.75rem);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  text-align: left;
  transition: var(--transition);
}

.note-toc__item:hover,
.note-toc__item--active {
  background: var(--accent-soft);
  color: var(--text-primary);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .note-sheet__header {
    flex-direction: column;
  }

  .note-sheet__meta {
    align-items: flex-start;
  }
}

@media (max-width: 980px) {
  .note-sheet__view-modes {
    display: none;
  }

  .note-rail__collapse {
    display: none;
  }

  .note-detail {
    grid-template-columns: 1fr;
  }

  .note-sheet {
    overflow-x: hidden;
  }

  .note-prose {
    max-width: 100%;
  }

  .note-rail {
    align-self: auto;
    height: auto;
    min-width: auto;
  }

  .note-rail__sticky {
    position: static;
    display: grid;
    grid-template-columns: 1fr 1fr;
  }

  .note-rail__sticky--fixed,
  .note-rail__sticky--bottom {
    position: static;
    inset: auto;
  }

  .note-toc {
    max-height: none;
    overflow-y: visible;
    padding-right: 0;
  }
}

@media (max-width: 640px) {
  .note-rail__sticky {
    grid-template-columns: 1fr;
  }
}
</style>

<style>
[data-theme='dark'] .note-prose pre {
  background: linear-gradient(180deg, #161b22, #10151b);
  box-shadow:
    inset 0 0 0 1px rgba(143, 182, 163, 0.16),
    0 16px 36px rgba(0, 0, 0, 0.24);
}

[data-theme='dark'] .note-prose pre code {
  color: #d8e0ea;
}

[data-theme='dark'] .note-prose code:not(pre code) {
  background: rgba(143, 182, 163, 0.14);
  border-color: rgba(143, 182, 163, 0.12);
  color: #e8f0f8;
}

[data-theme='dark'] .note-prose .hljs-keyword { color: #c792ea; }
[data-theme='dark'] .note-prose .hljs-string { color: #ecc48d; }
[data-theme='dark'] .note-prose .hljs-number { color: #f78c6c; }
[data-theme='dark'] .note-prose .hljs-comment { color: #7e8793; font-style: italic; }
[data-theme='dark'] .note-prose .hljs-function { color: #82aaff; }
[data-theme='dark'] .note-prose .hljs-title { color: #82aaff; }
[data-theme='dark'] .note-prose .hljs-class { color: #c3e88d; }
[data-theme='dark'] .note-prose .hljs-variable { color: #f07178; }
[data-theme='dark'] .note-prose .hljs-attr { color: #ffcb6b; }
[data-theme='dark'] .note-prose .hljs-built_in { color: #89ddff; }
[data-theme='dark'] .note-prose .hljs-type { color: #c3e88d; }
[data-theme='dark'] .note-prose .hljs-params { color: #d6deeb; }
[data-theme='dark'] .note-prose .hljs-literal { color: #89ddff; }
[data-theme='dark'] .note-prose .hljs-symbol { color: #89ddff; }
</style>
