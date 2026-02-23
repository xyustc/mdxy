<template>
  <Teleport to="body">
    <div v-if="visible" class="search-overlay" @click.self="close">
      <div class="search-modal glass-card">
        <div class="search-input-wrap">
          <span class="search-icon">&#128269;</span>
          <input
            ref="inputRef"
            v-model="query"
            class="search-input"
            placeholder="搜索笔记、工具..."
            @input="onInput"
            @keydown.enter="doSearch"
            @keydown.escape="close"
          />
          <span class="search-shortcut">ESC</span>
        </div>

        <div class="search-tabs">
          <button
            v-for="tab in tabs"
            :key="tab.value"
            class="search-tab"
            :class="{ active: searchType === tab.value }"
            @click="searchType = tab.value; doSearch()"
          >{{ tab.label }}</button>
        </div>

        <div class="search-body">
          <!-- 搜索结果 -->
          <template v-if="query && searched">
            <div v-if="result && result.total > 0">
              <div v-for="group in result.groups" :key="group.type" class="result-group">
                <div class="group-label">{{ group.label }}</div>
                <div
                  v-for="(item, i) in group.results"
                  :key="i"
                  class="result-item"
                  @click="goTo(item)"
                >
                  <span class="result-type-tag">{{ item.type === 'note' ? '笔记' : '工具' }}</span>
                  <div class="result-content">
                    <div class="result-title" v-html="item.title"></div>
                    <div class="result-context" v-html="item.context"></div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="empty-hint">无搜索结果</div>
          </template>

          <!-- 无输入时：热门搜索 + 搜索历史 -->
          <template v-if="!query">
            <div v-if="history.length" class="section">
              <div class="section-header">
                <span>搜索历史</span>
                <button class="clear-btn" @click="clearHist">清除</button>
              </div>
              <div class="tag-list">
                <span v-for="h in history" :key="h" class="tag-item" @click="query = h; doSearch()">{{ h }}</span>
              </div>
            </div>
            <div v-if="popular.length" class="section">
              <div class="section-header"><span>热门搜索</span></div>
              <div class="tag-list">
                <span v-for="p in popular" :key="p.keyword" class="tag-item" @click="query = p.keyword; doSearch()">{{ p.keyword }}</span>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
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
const result = ref<UnifiedSearchResult | null>(null)
const popular = ref<PopularKeyword[]>([])
const history = ref<string[]>([])
const inputRef = ref<HTMLInputElement>()

const tabs = [
  { label: '全部', value: 'all' },
  { label: '笔记', value: 'notes' },
  { label: '工具', value: 'tools' }
]

watch(() => props.modelValue, async (v) => {
  visible.value = v
  if (v) {
    query.value = ''
    searched.value = false
    result.value = null
    history.value = getHistory()
    try {
      const res = await searchApi.getPopular()
      if (res.success) popular.value = res.data || []
    } catch {}
    nextTick(() => inputRef.value?.focus())
  }
})

function close() { emit('update:modelValue', false) }

let debounceTimer: ReturnType<typeof setTimeout>
function onInput() {
  searched.value = false
  clearTimeout(debounceTimer)
  if (query.value.trim()) {
    debounceTimer = setTimeout(doSearch, 300)
  }
}

async function doSearch() {
  const q = query.value.trim()
  if (!q) return
  addHistory(q)
  history.value = getHistory()
  try {
    const res = await searchApi.search(q, searchType.value)
    if (res.success) result.value = res.data || null
  } catch {}
  searched.value = true
}

function clearHist() { clearHistory(); history.value = [] }

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
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5); z-index: 9999;
  display: flex; justify-content: center; padding-top: 12vh;
}
.search-modal {
  width: 600px; max-height: 70vh; display: flex; flex-direction: column;
  border-radius: 16px; overflow: hidden;
}
.search-input-wrap {
  display: flex; align-items: center; padding: 16px 20px;
  border-bottom: 1px solid var(--color-border); gap: 12px;
}
.search-icon { font-size: 20px; }
.search-input {
  flex: 1; border: none; outline: none; font-size: 16px;
  background: transparent; color: var(--color-text-primary);
}
.search-shortcut {
  font-size: 12px; padding: 2px 6px; border-radius: 4px;
  background: var(--color-bg-tertiary); color: var(--color-text-tertiary);
}
.search-tabs {
  display: flex; gap: 8px; padding: 12px 20px;
  border-bottom: 1px solid var(--color-border);
}
.search-tab {
  padding: 4px 12px; border: 1px solid var(--color-border); border-radius: 6px;
  background: transparent; cursor: pointer; font-size: 13px;
  color: var(--color-text-secondary); transition: var(--transition);
}
.search-tab.active { color: var(--accent-cyan); border-color: var(--accent-cyan); }
.search-body { flex: 1; overflow-y: auto; padding: 16px 20px; }
.result-group { margin-bottom: 16px; }
.group-label { font-size: 12px; color: var(--color-text-tertiary); margin-bottom: 8px; font-weight: 600; }
.result-item {
  display: flex; align-items: flex-start; gap: 10px; padding: 10px 12px;
  border-radius: 8px; cursor: pointer; transition: var(--transition);
}
.result-item:hover { background: var(--color-bg-secondary); }
.result-type-tag {
  font-size: 11px; padding: 2px 6px; border-radius: 4px; flex-shrink: 0;
  background: var(--color-bg-tertiary); color: var(--color-text-tertiary);
}
.result-title { font-size: 14px; font-weight: 500; color: var(--color-text-primary); }
.result-context { font-size: 13px; color: var(--color-text-secondary); margin-top: 4px; line-height: 1.5; }
.result-context :deep(mark) { background: rgba(0,212,170,0.3); color: inherit; border-radius: 2px; }
.result-title :deep(mark) { background: rgba(0,212,170,0.3); color: inherit; border-radius: 2px; }
.empty-hint { text-align: center; padding: 40px; color: var(--color-text-tertiary); }
.section { margin-bottom: 16px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.section-header span { font-size: 13px; color: var(--color-text-tertiary); font-weight: 600; }
.clear-btn {
  font-size: 12px; color: var(--color-text-tertiary); background: none;
  border: none; cursor: pointer;
}
.clear-btn:hover { color: var(--accent-coral); }
.tag-list { display: flex; flex-wrap: wrap; gap: 8px; }
.tag-item {
  padding: 4px 12px; border-radius: 6px; font-size: 13px; cursor: pointer;
  background: var(--color-bg-tertiary); color: var(--color-text-secondary);
  transition: var(--transition);
}
.tag-item:hover { color: var(--accent-cyan); }
@media (max-width: 768px) {
  .search-modal { width: 95vw; }
}
</style>
