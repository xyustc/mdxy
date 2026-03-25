<template>
  <div class="claude-page section-shell">
    <div class="app-frame claude-page__frame">
      <header class="claude-hero">
        <div class="claude-hero__copy">
          <span class="section-kicker">Field Kit / Internal Replica</span>
          <h1 class="section-title">Claude Code 速查表</h1>
          <p class="section-description">
            Claude Code 键盘快捷键、斜杠命令等高密度功能速查表，持续更新中。
          </p>
        </div>

        <div class="claude-hero__actions">
          <router-link to="/tools" class="hero-action hero-action--secondary">返回工具箱</router-link>
          <a :href="content.meta.source_url" target="_blank" rel="noopener" class="hero-action hero-action--primary">打开原页</a>
          <button type="button" class="hero-action hero-action--ghost" @click="handlePrint">打印 / 导出 PDF</button>
        </div>
      </header>

      <section class="sheet-stage">
        <div class="sheet-noise" aria-hidden="true"></div>

        <article class="cheat-sheet" aria-label="Claude Code 速查表">
          <header class="sheet-header">
            <div>
              <p class="sheet-eyebrow">Quick Reference / Internal Adaptation</p>
              <h2>{{ content.meta.source_title }}</h2>
            </div>

            <div class="sheet-header__right">
              <div class="os-toggle" role="group" aria-label="切换快捷键显示系统">
                <button
                  type="button"
                  class="os-btn"
                  :class="{ 'os-btn--active': osMode === 'mac' }"
                  @click="setOsMode('mac')"
                >
                  Mac
                </button>
                <button
                  type="button"
                  class="os-btn"
                  :class="{ 'os-btn--active': osMode === 'win' }"
                  @click="setOsMode('win')"
                >
                  Windows
                </button>
              </div>

              <div class="sheet-meta">
                <span class="sheet-meta__version">{{ content.meta.version }}</span>
                <span class="sheet-meta__date">最后更新：{{ content.meta.updated_at }}</span>
              </div>
            </div>
          </header>

          <div class="sheet-sync-status">
            <span class="sheet-sync-status__pill" :class="`sheet-sync-status__pill--${dataSource}`">
              {{ dataSource === 'published' ? '自动同步快照' : '本地静态回退' }}
            </span>
            <span v-if="remoteSyncedAt" class="sheet-sync-status__meta">同步时间：{{ formatDate(remoteSyncedAt) }}</span>
            <span v-if="remotePublishedAt" class="sheet-sync-status__meta">发布时间：{{ formatDate(remotePublishedAt) }}</span>
          </div>

          <section v-if="showChangelog" class="sheet-changelog">
            <div class="sheet-changelog__header">
              <span>📋 近期更新</span>
              <button type="button" class="sheet-changelog__dismiss" aria-label="关闭更新提示" @click="dismissChangelog">
                ✕
              </button>
            </div>
            <ul class="sheet-changelog__list">
              <li v-for="item in content.changelog" :key="item.code">
                <code>{{ item.code }}</code>
                <span>{{ item.text }}</span>
              </li>
            </ul>
          </section>

          <main class="sheet-grid">
            <div v-for="(column, columnIndex) in content.columns" :key="columnIndex" class="sheet-column">
              <section
                v-for="section in column"
                :key="section.id"
                class="sheet-section"
                :class="`sheet-section--${section.theme}`"
              >
                <div class="sheet-section__header">{{ section.title }}</div>
                <div class="sheet-section__content">
                  <div v-for="group in section.groups" :key="group.title" class="sheet-group">
                    <h3 class="sheet-group__title">{{ group.title }}</h3>
                    <div v-for="row in group.rows" :key="`${group.title}-${row.key}-${row.desc}`" class="sheet-row">
                      <div class="sheet-key">
                        <template v-for="(segment, segmentIndex) in getKeySegments(row.key, row.key_variant)" :key="segmentIndex">
                          <span v-if="segment.type === 'keycap'" class="keycap">{{ segment.label }}</span>
                          <span v-else class="sheet-key__text">{{ segment.label }}</span>
                        </template>
                      </div>

                      <div v-if="row.desc || row.added_at" class="sheet-desc">
                        <span v-if="row.desc">{{ row.desc }}</span>
                        <span v-if="isRecentBadge(row.added_at)" class="badge-new">NEW</span>
                      </div>
                    </div>
                  </div>
                </div>
              </section>
            </div>
          </main>

          <footer class="sheet-footer">
            <div v-for="row in content.footer" :key="row.label" class="sheet-footer__row">
              <span class="sheet-footer__label">{{ row.label }}：</span>
              <div class="sheet-footer__items">
                <span v-for="item in row.items" :key="`${row.label}-${item.code}`" class="sheet-footer__item">
                  <code>{{ item.code }}</code>
                  <span v-if="item.text">{{ item.text }}</span>
                </span>
              </div>
            </div>
          </footer>
        </article>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { claudeCodeCheatSheetContent } from '@/data/claudeCodeCheatSheet'
import { cheatSheetApi } from '@/api/cheatSheet'
import type { CheatSheetContent } from '@/api/types'

type OsMode = 'mac' | 'win'
type KeySegment = { type: 'keycap' | 'text'; label: string }

const OS_STORAGE_KEY = 'mdxy.claude-code-cheatsheet.os'
const CHANGELOG_STORAGE_KEY = 'mdxy.claude-code-cheatsheet.dismissed'

const osMode = ref<OsMode>('mac')
const showChangelog = ref(true)
const content = ref<CheatSheetContent>({
  meta: {
    source_url: claudeCodeCheatSheetContent.meta.sourceUrl,
    source_title: claudeCodeCheatSheetContent.meta.sourceTitle,
    version: claudeCodeCheatSheetContent.meta.version,
    updated_at: claudeCodeCheatSheetContent.meta.updatedAt
  },
  changelog: claudeCodeCheatSheetContent.changelog,
  footer: claudeCodeCheatSheetContent.footer,
  columns: claudeCodeCheatSheetContent.columns.map((column) =>
    column.map((section) => ({
      ...section,
      groups: section.groups.map((group) => ({
        ...group,
        rows: group.rows.map((row) => ({
          key: row.key,
          desc: row.desc,
          added_at: row.addedAt,
          key_variant: row.keyVariant
        }))
      }))
    }))
  )
})
const dataSource = ref<'published' | 'fallback'>('fallback')
const remoteSyncedAt = ref('')
const remotePublishedAt = ref('')

function detectOs(): OsMode {
  const platform = navigator.platform || ''
  const userAgent = navigator.userAgent || ''
  if (/Mac|iPhone|iPod|iPad/.test(platform) || /Mac/.test(userAgent)) {
    return 'mac'
  }
  return 'win'
}

function setOsMode(mode: OsMode) {
  osMode.value = mode
  window.localStorage.setItem(OS_STORAGE_KEY, mode)
}

function dismissChangelog() {
  showChangelog.value = false
  window.localStorage.setItem(CHANGELOG_STORAGE_KEY, '1')
}

function isRecentBadge(addedAt?: string) {
  if (!addedAt) return false
  const now = new Date()
  const created = new Date(addedAt)
  if (Number.isNaN(created.getTime())) return false
  return (now.getTime() - created.getTime()) / 86400000 <= 14
}

function mapKeycapLabel(label: string, variant?: string) {
  if (variant === 'paste-image' && (label === 'Ctrl' || label === '⌘')) {
    return osMode.value === 'mac' ? 'Ctrl' : 'Alt'
  }
  if (label === 'Alt' || label === '⌥') {
    return osMode.value === 'mac' ? '⌥' : 'Alt'
  }
  if (label === 'Shift' || label === '⇧') {
    return osMode.value === 'mac' ? '⇧' : 'Shift'
  }
  if (label === 'Ctrl' || label === '⌘') {
    return 'Ctrl'
  }
  return label
}

function getKeySegments(source: string, variant?: string): KeySegment[] {
  const segments: KeySegment[] = []
  const matcher = /\[(.+?)\]/g
  let lastIndex = 0
  let match: RegExpExecArray | null

  while ((match = matcher.exec(source))) {
    if (match.index > lastIndex) {
      const text = source.slice(lastIndex, match.index).trim()
      if (text) {
        segments.push({ type: 'text', label: text })
      }
    }

    segments.push({
      type: 'keycap',
      label: mapKeycapLabel(match[1], variant)
    })
    lastIndex = match.index + match[0].length
  }

  if (lastIndex < source.length) {
    const text = source.slice(lastIndex).trim()
    if (text) {
      segments.push({ type: 'text', label: text })
    }
  }

  return segments.length ? segments : [{ type: 'text', label: source }]
}

function handlePrint() {
  window.print()
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function loadPublishedContent() {
  try {
    const res = await cheatSheetApi.getPublished('claude-code')
    if (res.success && res.data) {
      content.value = res.data.content
      dataSource.value = 'published'
      remoteSyncedAt.value = res.data.synced_at || ''
      remotePublishedAt.value = res.data.published_at || ''
    }
  } catch {
    dataSource.value = 'fallback'
  }
}

onMounted(() => {
  const stored = window.localStorage.getItem(OS_STORAGE_KEY) as OsMode | null
  osMode.value = stored === 'mac' || stored === 'win' ? stored : detectOs()
  showChangelog.value = window.localStorage.getItem(CHANGELOG_STORAGE_KEY) !== '1'
  loadPublishedContent()
})
</script>

<style scoped>
.claude-page__frame {
  width: min(100%, 1720px);
  max-width: none;
  margin: 0 auto;
  padding-inline: clamp(18px, 2.6vw, 42px);
  display: grid;
  gap: var(--space-2xl);
}

.claude-hero {
  display: flex;
  justify-content: space-between;
  gap: var(--space-xl);
  align-items: flex-end;
}

.claude-hero__copy {
  display: grid;
  gap: var(--space-sm);
}

.claude-hero__copy .section-title {
  max-width: 12ch;
}

.claude-hero__copy .section-description {
  max-width: 70ch;
}

.claude-hero__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  justify-content: flex-end;
}

.hero-action {
  min-height: 3rem;
  padding: 0 1rem;
  border-radius: var(--radius-pill);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  background: var(--bg-panel);
  transition: var(--transition);
}

.hero-action:hover {
  color: var(--text-primary);
  border-color: var(--border-strong);
}

.hero-action--primary {
  background: var(--bg-contrast);
  border-color: var(--bg-contrast);
  color: var(--text-inverse);
}

.hero-action--primary:hover {
  color: var(--text-inverse);
  opacity: 0.92;
}

.hero-action--ghost {
  background: rgba(255, 255, 255, 0.28);
}

.sheet-stage {
  position: relative;
  isolation: isolate;
  width: 100%;
}

.sheet-noise {
  position: absolute;
  inset: 2rem 1.2rem auto;
  height: 88%;
  border-radius: 2rem;
  background:
    radial-gradient(circle at 15% 20%, rgba(7, 89, 133, 0.14), transparent 26%),
    radial-gradient(circle at 85% 15%, rgba(220, 38, 38, 0.14), transparent 22%),
    radial-gradient(circle at 50% 88%, rgba(5, 150, 105, 0.14), transparent 24%);
  filter: blur(48px);
  opacity: 0.72;
  pointer-events: none;
  z-index: 0;
}

.cheat-sheet {
  position: relative;
  z-index: 1;
  color-scheme: light;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.88), rgba(252, 248, 240, 0.96)),
    linear-gradient(135deg, rgba(255, 255, 255, 0.72), rgba(255, 255, 255, 0.92));
  border: 1px solid rgba(31, 41, 55, 0.12);
  border-radius: 2rem;
  padding: clamp(0.95rem, 1.5vw, 1.45rem);
  box-shadow: 0 28px 90px rgba(18, 24, 28, 0.14);
  overflow: hidden;
}

.cheat-sheet::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(rgba(31, 41, 55, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 41, 55, 0.025) 1px, transparent 1px);
  background-size: 100% 30px, 30px 100%;
  opacity: 0.4;
}

.sheet-header,
.sheet-changelog,
.sheet-grid,
.sheet-footer {
  position: relative;
  z-index: 1;
}

.sheet-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: center;
  padding-bottom: 0.9rem;
  margin-bottom: 0.95rem;
  border-bottom: 1px solid rgba(107, 114, 128, 0.2);
}

.sheet-eyebrow {
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #6b7280;
}

.sheet-header h2 {
  margin-top: 0.45rem;
  font-family: var(--font-display);
  font-size: clamp(1.75rem, 3vw, 2.5rem);
  color: #111827;
  letter-spacing: -0.04em;
}

.sheet-header__right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.os-toggle {
  display: inline-flex;
  padding: 0.25rem;
  border-radius: 999px;
  border: 1px solid rgba(55, 65, 81, 0.14);
  background: rgba(249, 250, 251, 0.96);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.78);
}

.os-btn {
  min-width: 5.2rem;
  padding: 0.52rem 0.9rem;
  border-radius: 999px;
  font-family: var(--font-mono);
  font-size: 0.76rem;
  font-weight: 700;
  color: #6b7280;
  transition: 180ms ease;
}

.os-btn--active {
  background: #111827;
  color: #f9fafb;
}

.sheet-meta {
  text-align: right;
  display: grid;
  gap: 0.22rem;
}

.sheet-meta__version {
  color: #374151;
  font-size: 0.84rem;
  font-weight: 700;
}

.sheet-meta__date {
  color: #9ca3af;
  font-size: 0.76rem;
}

.sheet-changelog {
  margin-bottom: 1rem;
  padding: 0.85rem 1rem;
  border: 1px solid rgba(245, 158, 11, 0.45);
  border-radius: 1rem;
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
}

.sheet-sync-status {
  margin-bottom: 0.9rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem 0.9rem;
  align-items: center;
}

.sheet-sync-status__pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 2rem;
  padding: 0 0.75rem;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  border: 1px solid rgba(17, 24, 39, 0.1);
}

.sheet-sync-status__pill--published {
  background: rgba(5, 150, 105, 0.12);
  color: #047857;
}

.sheet-sync-status__pill--fallback {
  background: rgba(217, 119, 6, 0.12);
  color: #b45309;
}

.sheet-sync-status__meta {
  color: #6b7280;
  font-size: 0.78rem;
}

.sheet-changelog__header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: center;
  color: #92400e;
  font-weight: 700;
}

.sheet-changelog__dismiss {
  font-size: 0.9rem;
  color: #92400e;
  opacity: 0.72;
}

.sheet-changelog__dismiss:hover {
  opacity: 1;
}

.sheet-changelog__list {
  margin-top: 0.75rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem 1rem;
  list-style: none;
  color: #78350f;
}

.sheet-changelog__list li {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.78rem;
}

.sheet-changelog__list code,
.sheet-footer code {
  font-family: var(--font-mono);
  font-size: 0.74rem;
  background: rgba(255, 255, 255, 0.42);
  padding: 0.18rem 0.42rem;
  border-radius: 0.45rem;
}

.sheet-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.85rem;
}

.sheet-column {
  display: grid;
  gap: 1rem;
  align-content: start;
}

.sheet-section {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  border-radius: 1rem;
  overflow: hidden;
  border: 1px solid rgba(17, 24, 39, 0.08);
  box-shadow: 0 12px 26px rgba(17, 24, 39, 0.06);
}

.sheet-section__header {
  padding: 0.8rem 1rem;
  color: #fff;
  font-size: 0.84rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.sheet-section__content {
  padding: 0.82rem 0.88rem;
  display: grid;
  gap: 0.75rem;
  flex: 1;
}

.sheet-section--keyboard {
  background: #eff6ff;
}

.sheet-section--keyboard .sheet-section__header {
  background: #2563eb;
}

.sheet-section--mcp {
  background: #ecfeff;
}

.sheet-section--mcp .sheet-section__header {
  background: #0891b2;
}

.sheet-section--slash {
  background: #f5f3ff;
}

.sheet-section--slash .sheet-section__header {
  background: #7c3aed;
}

.sheet-section--memory {
  background: #fffbeb;
}

.sheet-section--memory .sheet-section__header {
  background: #d97706;
}

.sheet-section--workflows {
  background: #fef2f2;
}

.sheet-section--workflows .sheet-section__header {
  background: #dc2626;
}

.sheet-section--config {
  background: #fff7ed;
}

.sheet-section--config .sheet-section__header {
  background: #c2410c;
}

.sheet-section--skills {
  background: #ecfdf5;
}

.sheet-section--skills .sheet-section__header {
  background: #059669;
}

.sheet-section--cli {
  background: #f9fafb;
}

.sheet-section--cli .sheet-section__header {
  background: #4b5563;
}

.sheet-group {
  display: grid;
  gap: 0.35rem;
}

.sheet-group__title {
  font-size: 0.78rem;
  font-style: italic;
  font-weight: 700;
  color: #374151;
  padding-bottom: 0.25rem;
  border-bottom: 1px dotted rgba(55, 65, 81, 0.18);
}

.sheet-row {
  display: flex;
  gap: 0.62rem;
  align-items: flex-start;
  min-width: 0;
}

.sheet-key {
  min-width: 0;
  flex-shrink: 0;
  font-family: var(--font-mono);
  font-size: 0.76rem;
  font-weight: 700;
  color: #111827;
  line-height: 1.22;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.18rem;
}

.sheet-key__text {
  display: inline-flex;
  align-items: center;
}

.sheet-desc {
  flex: 1;
  min-width: 0;
  font-size: 0.76rem;
  color: #4b5563;
  line-height: 1.32;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.3rem;
}

.keycap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.5rem;
  padding: 0.14rem 0.42rem;
  border-radius: 0.42rem;
  border: 1px solid #d1d5db;
  background: linear-gradient(180deg, #fafafa 0%, #e5e7eb 100%);
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.08), inset 0 -1px 0 rgba(0, 0, 0, 0.05);
}

.badge-new {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.12rem 0.4rem;
  border-radius: 999px;
  background: #ef4444;
  color: #fff;
  font-size: 0.6rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.sheet-footer {
  margin-top: 1rem;
  padding: 1rem 1.1rem;
  border-radius: 1rem;
  background: #1f2937;
  color: #f9fafb;
  display: grid;
  gap: 0.65rem;
}

.sheet-footer__row {
  display: grid;
  grid-template-columns: max-content minmax(0, 1fr);
  gap: 0.7rem;
  align-items: flex-start;
}

.sheet-footer__label {
  color: #9ca3af;
  font-weight: 700;
  flex-shrink: 0;
}

.sheet-footer__items {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem 0.65rem;
  min-width: 0;
}

.sheet-footer__item {
  display: inline-flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.35rem;
  font-size: 0.78rem;
  min-width: 0;
}

.sheet-footer__item code {
  overflow-wrap: anywhere;
}

@media (max-width: 1480px) {
  .sheet-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 960px) {
  .claude-hero,
  .sheet-header,
  .sheet-footer__row {
    flex-direction: column;
    align-items: flex-start;
  }

  .claude-hero__actions {
    justify-content: flex-start;
  }

  .sheet-header__right {
    width: 100%;
    justify-content: space-between;
  }
}

@media (max-width: 760px) {
  .sheet-grid {
    grid-template-columns: 1fr;
  }

  .cheat-sheet {
    padding: 1rem;
    border-radius: 1.35rem;
  }

  .claude-page__frame {
    width: 100%;
    padding-inline: 18px;
  }

  .sheet-header__right,
  .sheet-footer__items {
    width: 100%;
  }

  .os-toggle {
    width: 100%;
  }

  .os-btn {
    flex: 1;
  }
}

@media print {
  :global(.site-header),
  :global(.site-footer),
  :global(.layout-ornament),
  .claude-hero,
  .sheet-noise {
    display: none !important;
  }

  .claude-page {
    padding: 0 !important;
  }

  .claude-page__frame {
    gap: 0 !important;
  }

  .cheat-sheet {
    box-shadow: none;
    border-radius: 0;
    border: none;
    padding: 0;
  }

  .sheet-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .sheet-section {
    break-inside: avoid;
    box-shadow: none;
  }
}
</style>
