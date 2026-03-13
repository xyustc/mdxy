<template>
  <Teleport to="body">
    <div v-if="visible" class="search-overlay">
      <button type="button" class="search-overlay__backdrop" aria-label="关闭搜索弹窗" @click="close"></button>
      <section class="paper-sheet search-modal" role="dialog" aria-modal="true" aria-label="全站搜索">
        <header class="search-modal__header">
          <span class="meta-label">Global Index</span>
          <button type="button" class="search-close" aria-label="关闭搜索弹窗" @click="close">ESC</button>
        </header>

        <div class="search-input-wrap">
          <span class="search-input-wrap__icon">⌕</span>
          <input
            id="search-modal-input"
            ref="inputRef"
            v-model="query"
            type="search"
            name="site-search"
            class="search-input"
            autocomplete="off"
            aria-label="搜索笔记、工具和关键词"
            placeholder="搜索笔记、工具、关键字…"
            @input="onInput"
            @keydown="onInputKeydown"
          />
          <span class="search-shortcut">ENTER</span>
        </div>

        <div class="search-tabs">
          <button
            v-for="tab in tabs"
            :key="tab.value"
            class="search-tab"
            :class="{ active: searchType === tab.value }"
            @click="switchType(tab.value)"
          >
            {{ tab.label }}
          </button>
        </div>

        <div class="search-body" aria-live="polite">
          <template v-if="loading">
            <div class="search-state">正在检索索引…</div>
          </template>

          <template v-else-if="query && searched">
            <div v-if="result && result.total > 0">
              <div v-for="(group, gIndex) in result.groups" :key="group.type" class="result-group">
                <div class="group-label">{{ group.label }}</div>
                <button
                  v-for="(item, i) in group.results"
                  :key="`${group.type}-${i}`"
                  type="button"
                  class="result-item"
                  :class="{ 'result-item--active': selectedKey === `${gIndex}-${i}` }"
                  @mouseenter="selectedKey = `${gIndex}-${i}`"
                  @click="goTo(item)"
                >
                  <span class="result-type-tag">{{ item.type === 'note' ? '笔记' : '工具' }}</span>
                  <div class="result-content">
                    <div class="result-title" v-html="item.title"></div>
                    <div class="result-context" v-html="item.context"></div>
                  </div>
                </button>
              </div>
            </div>
            <div v-else class="search-state">无搜索结果，请更换关键词。</div>
          </template>

          <template v-else>
            <div v-if="history.length" class="search-section">
              <div class="search-section__header">
                <span>搜索历史</span>
                <button class="clear-btn" @click="clearHist">清除</button>
              </div>
              <div class="tag-list">
                <button v-for="h in history" :key="h" type="button" class="tag-item" @click="quickSearch(h)">{{ h }}</button>
              </div>
            </div>
            <div v-if="popular.length" class="search-section">
              <div class="search-section__header"><span>热门搜索</span></div>
              <div class="tag-list">
                <button v-for="p in popular" :key="p.keyword" type="button" class="tag-item" @click="quickSearch(p.keyword)">
                  {{ p.keyword }}
                </button>
              </div>
            </div>
          </template>
        </div>
      </section>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { searchApi } from '@/api/search'
import { useSearchHistory } from '@/composables/useSearchHistory'
import type { UnifiedSearchResult, SearchResultItem, PopularKeyword } from '@/api/types'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: boolean): void }>()
const router = useRouter()
const { getHistory, addHistory, clearHistory } = useSearchHistory()

const visible = ref(false)
const query = ref('')
const searchType = ref('all')
const searched = ref(false)
const loading = ref(false)
const result = ref<UnifiedSearchResult | null>(null)
const popular = ref<PopularKeyword[]>([])
const history = ref<string[]>([])
const inputRef = ref<HTMLInputElement>()
const selectedKey = ref('')

const tabs = [
  { label: '全部', value: 'all' },
  { label: '笔记', value: 'notes' },
  { label: '工具', value: 'tools' }
]

const flatResults = computed(() => {
  if (!result.value) return [] as { key: string; item: SearchResultItem }[]
  return result.value.groups.flatMap((group, gIndex) =>
    group.results.map((item, i) => ({ key: `${gIndex}-${i}`, item }))
  )
})

watch(
  () => props.modelValue,
  async (v) => {
    visible.value = v
    if (v) {
      query.value = ''
      searched.value = false
      loading.value = false
      result.value = null
      selectedKey.value = ''
      history.value = getHistory()
      document.body.style.overflow = 'hidden'
      try {
        const res = await searchApi.getPopular()
        if (res.success) popular.value = res.data || []
      } catch {
        popular.value = []
      }
      nextTick(() => inputRef.value?.focus())
    } else {
      document.body.style.overflow = ''
    }
  }
)

onBeforeUnmount(() => {
  document.body.style.overflow = ''
})

function close() {
  emit('update:modelValue', false)
}

let debounceTimer: ReturnType<typeof setTimeout>
function onInput() {
  searched.value = false
  selectedKey.value = ''
  clearTimeout(debounceTimer)
  if (query.value.trim()) {
    debounceTimer = setTimeout(doSearch, 260)
  }
}

function quickSearch(keyword: string) {
  query.value = keyword
  doSearch()
}

function switchType(type: string) {
  searchType.value = type
  if (query.value.trim()) {
    doSearch()
  }
}

function moveSelection(direction: 1 | -1) {
  if (!flatResults.value.length) return
  const index = flatResults.value.findIndex((r) => r.key === selectedKey.value)
  const next = index < 0 ? 0 : (index + direction + flatResults.value.length) % flatResults.value.length
  selectedKey.value = flatResults.value[next].key
}

function onInputKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    e.preventDefault()
    close()
    return
  }

  if (e.key === 'ArrowDown') {
    e.preventDefault()
    moveSelection(1)
    return
  }

  if (e.key === 'ArrowUp') {
    e.preventDefault()
    moveSelection(-1)
    return
  }

  if (e.key === 'Enter') {
    e.preventDefault()
    const selected = flatResults.value.find((r) => r.key === selectedKey.value)
    if (selected) {
      goTo(selected.item)
      return
    }
    doSearch()
  }
}

async function doSearch() {
  const q = query.value.trim()
  if (!q) return
  addHistory(q)
  history.value = getHistory()

  loading.value = true
  try {
    const res = await searchApi.search(q, searchType.value)
    if (res.success) {
      result.value = res.data || null
      selectedKey.value = flatResults.value[0]?.key || ''
    }
  } catch {
    result.value = null
  } finally {
    searched.value = true
    loading.value = false
  }
}

function clearHist() {
  clearHistory()
  history.value = []
}

function goTo(item: SearchResultItem) {
  close()
  if (item.type === 'tool' && item.url) {
    window.open(item.url, '_blank')
  } else {
    router.push(item.path)
  }
}
</script>

<style scoped>
.search-overlay {
  position: fixed;
  inset: 0;
  background: rgba(9, 8, 7, 0.6);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: min(11vh, 92px) 18px 18px;
  overscroll-behavior: contain;
}

.search-overlay__backdrop {
  position: absolute;
  inset: 0;
}

.search-modal {
  position: relative;
  z-index: 1;
  width: min(880px, 100%);
  max-height: min(78vh, 820px);
  display: flex;
  flex-direction: column;
  border-radius: var(--radius-xl);
  overflow: hidden;
}

.search-modal__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-primary);
}

.search-close {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  padding: 0.28rem 0.7rem;
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  color: var(--text-muted);
}

.search-input-wrap {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 0.8rem;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-primary);
}

.search-input-wrap:focus-within {
  background: var(--accent-soft);
  box-shadow: inset 0 0 0 1px rgba(41, 70, 58, 0.16);
}

.search-input-wrap__icon {
  color: var(--text-muted);
  font-size: 1.1rem;
}

.search-input {
  border: none;
  background: transparent;
  color: var(--text-primary);
  font-size: 1.05rem;
}

.search-shortcut {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  letter-spacing: 0.13em;
  color: var(--text-muted);
}

.search-tabs {
  display: flex;
  gap: 0.45rem;
  padding: 0.75rem 1.25rem;
  border-bottom: 1px solid var(--border-primary);
}

.search-tab {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  padding: 0.4rem 0.95rem;
  color: var(--text-secondary);
  font-size: 0.86rem;
}

.search-tab.active {
  color: var(--text-primary);
  background: var(--accent-soft);
  border-color: var(--accent-primary);
}

.search-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.1rem 1.25rem 1.25rem;
}

.search-state {
  min-height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

.result-group {
  margin-bottom: 1rem;
}

.group-label {
  margin-bottom: 0.5rem;
  font-family: var(--font-mono);
  font-size: 0.7rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.result-item {
  width: 100%;
  display: flex;
  align-items: flex-start;
  gap: 0.8rem;
  padding: 0.8rem 0.9rem;
  border-radius: var(--radius-lg);
  border: 1px solid transparent;
  text-align: left;
}

.result-item:hover,
.result-item--active {
  background: var(--accent-soft);
  border-color: var(--border-primary);
}

.result-type-tag {
  font-size: 0.68rem;
  font-family: var(--font-mono);
  letter-spacing: 0.11em;
  text-transform: uppercase;
  color: var(--text-muted);
  padding: 0.28rem 0.5rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
}

.result-content {
  min-width: 0;
}

.result-title {
  color: var(--text-primary);
  font-size: 0.96rem;
  font-weight: 600;
}

.result-context {
  margin-top: 0.3rem;
  color: var(--text-secondary);
  font-size: 0.88rem;
  line-height: 1.6;
}

.result-context :deep(mark),
.result-title :deep(mark) {
  background: var(--accent-secondary-soft);
  border-radius: 0.2rem;
}

.search-section {
  margin-bottom: 1rem;
}

.search-section__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.55rem;
}

.search-section__header span {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.clear-btn {
  color: var(--text-muted);
  font-size: 0.78rem;
}

.clear-btn:hover {
  color: var(--text-primary);
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
}

.tag-item {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  padding: 0.35rem 0.8rem;
  color: var(--text-secondary);
  font-size: 0.84rem;
  background: var(--bg-panel);
}

.tag-item:hover {
  color: var(--text-primary);
  background: var(--accent-soft);
}

@media (max-width: 768px) {
  .search-tabs {
    overflow-x: auto;
  }

  .result-item {
    flex-direction: column;
  }
}
</style>
