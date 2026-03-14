<template>
  <div class="notes-page section-shell" :class="{ 'notes-page--wide': isWideLayout }">
    <div class="app-frame notes-frame">
      <header class="notes-heading">
        <span class="section-kicker">Reading Index</span>
        <h1 class="section-title">笔记目录</h1>
        <p class="section-description">从目录进入，从结构理解。这里收纳我的技术笔记与持续更新的知识沉淀。</p>
      </header>

      <div class="notes-layout">
        <aside class="surface-panel notes-sidebar">
          <div class="notes-search">
            <SearchOutline class="notes-search__icon" />
            <input
              v-model="searchKeyword"
              type="search"
              name="note-search"
              autocomplete="off"
              aria-label="搜索并直达笔记"
              class="notes-search__input"
              placeholder="搜索并直达笔记…"
              @keydown.enter="handleSearch"
            />
            <button type="button" class="notes-search__button" @click="handleSearch">检索</button>
          </div>
          <p v-if="searchHint" class="notes-search__hint">{{ searchHint }}</p>

          <div class="notes-sidebar__meta">
            <span><strong>{{ fileCount }}</strong> 篇笔记</span>
            <span><strong>{{ folderCount }}</strong> 个目录</span>
          </div>

          <div class="tree-container">
            <n-tree :data="filteredTreeData" :node-props="nodeProps" block-line />
          </div>
        </aside>

        <main class="notes-content">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NTree } from 'naive-ui'
import type { TreeOption } from 'naive-ui'
import { SearchOutline } from '@vicons/ionicons5'
import { noteApi } from '@/api/note'
import type { NoteNode } from '@/api/types'

interface TreeNode extends TreeOption {
  label: string
  key: string
  isLeaf: boolean
  rawName: string
  children?: TreeNode[]
}

const router = useRouter()
const route = useRoute()
const NOTE_READING_WIDTH_KEY = 'mdxy.note.reading-width'
const NOTES_WIDTH_EVENT = 'mdxy:notes-reading-width'
const searchKeyword = ref('')
const searchHint = ref('')
const treeData = ref<TreeNode[]>([])
const isCompactViewport = ref(false)
const readingWidthMode = ref<'standard' | 'wide'>('standard')
const isWideLayout = computed(() => !isCompactViewport.value && readingWidthMode.value === 'wide')
const currentNotePath = computed(() =>
  Array.isArray(route.params.path) ? route.params.path.join('/') : (route.params.path as string | undefined) || ''
)

function applyReadingWidthMode(mode: string | null | undefined) {
  readingWidthMode.value = mode === 'wide' ? 'wide' : 'standard'
}

function handleReadingWidthModeChange(event: Event) {
  const customEvent = event as CustomEvent<'standard' | 'wide'>
  applyReadingWidthMode(customEvent.detail)
}

function syncViewportMode() {
  isCompactViewport.value = window.innerWidth <= 980
}

function convertToTreeData(nodes: NoteNode[]): TreeNode[] {
  return nodes.map((node) => ({
    label: node.name,
    rawName: node.name,
    key: node.path,
    isLeaf: node.type === 'file',
    children: node.children ? convertToTreeData(node.children) : undefined
  }))
}

const filteredTreeData = computed<TreeOption[]>(() => {
  const keyword = searchKeyword.value.trim().toLowerCase()
  if (!keyword) return treeData.value

  const filterNodes = (nodes: TreeNode[]): TreeNode[] => {
    return nodes
      .map((node) => {
        const matchedSelf = node.rawName.toLowerCase().includes(keyword)
        const children = node.children ? filterNodes(node.children) : undefined
        if (matchedSelf || (children && children.length)) {
          return { ...node, children }
        }
        return null
      })
      .filter(Boolean) as TreeNode[]
  }

  return filterNodes(treeData.value)
})

const fileCount = computed(() => {
  const count = (nodes: TreeNode[]): number => {
    return nodes.reduce((total, node) => total + (node.isLeaf ? 1 : 0) + (node.children ? count(node.children) : 0), 0)
  }
  return count(treeData.value)
})

const folderCount = computed(() => {
  const count = (nodes: TreeNode[]): number => {
    return nodes.reduce((total, node) => total + (node.isLeaf ? 0 : 1) + (node.children ? count(node.children) : 0), 0)
  }
  return count(treeData.value)
})

const nodeProps = ({ option }: { option: TreeOption }) => {
  const node = option as TreeNode
  const classes = [node.isLeaf ? 'note-node note-node--leaf' : 'note-node note-node--folder']

  if (node.isLeaf && currentNotePath.value === node.key) {
    classes.push('note-node--active')
  }

  return {
    class: classes.join(' '),
    onClick() {
      if (node.isLeaf) {
        router.push(`/notes/${node.key}`)
      }
    }
  }
}

async function fetchTree() {
  try {
    const response = await noteApi.getTree()
    if (response.success && response.data) {
      treeData.value = convertToTreeData(response.data)
    }
  } catch (error) {
    console.error('获取笔记目录失败:', error)
  }
}

async function handleSearch() {
  const q = searchKeyword.value.trim()
  if (!q) {
    searchHint.value = ''
    return
  }

  try {
    const response = await noteApi.search(q)
    if (response.success && response.data && response.data.total > 0) {
      const first = response.data.results[0]
      searchHint.value = `已定位到：${first.name}`
      router.push(`/notes/${first.path}`)
    } else {
      searchHint.value = '未找到匹配笔记，请尝试更换关键词。'
    }
  } catch (error) {
    console.error('笔记搜索失败:', error)
    searchHint.value = '搜索服务暂时不可用。'
  }
}

onMounted(() => {
  syncViewportMode()
  applyReadingWidthMode(window.localStorage.getItem(NOTE_READING_WIDTH_KEY))
  window.addEventListener(NOTES_WIDTH_EVENT, handleReadingWidthModeChange as EventListener)
  window.addEventListener('resize', syncViewportMode)
  fetchTree()
})

onBeforeUnmount(() => {
  window.removeEventListener(NOTES_WIDTH_EVENT, handleReadingWidthModeChange as EventListener)
  window.removeEventListener('resize', syncViewportMode)
})
</script>

<style scoped>
.notes-page {
  overflow-x: clip;
}

.notes-frame {
  display: grid;
  gap: var(--space-2xl);
}

.notes-page--wide .notes-frame {
  width: min(100%, 1580px);
}

.notes-heading {
  display: grid;
  gap: var(--space-md);
}

.notes-heading .section-title {
  max-width: 10ch;
}

.notes-heading .section-description {
  max-width: 38ch;
}

.notes-layout {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: var(--space-lg);
  align-items: start;
}

.notes-page--wide .notes-layout {
  grid-template-columns: 320px minmax(0, 1fr);
  gap: var(--space-2xl);
}

.notes-sidebar {
  position: static;
  transform: translateX(-0.75rem);
  padding: 1.1rem;
  display: grid;
  gap: var(--space-md);
}

.notes-page--wide .notes-sidebar {
  transform: translateX(-1.1rem);
}

.notes-search {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: var(--space-sm);
  padding: 0.35rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  background: var(--bg-panel);
}

.notes-search:focus-within {
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 2px rgba(41, 70, 58, 0.14);
}

.notes-search__icon {
  width: 1rem;
  height: 1rem;
  color: var(--text-muted);
  margin-left: 0.45rem;
}

.notes-search__input {
  min-width: 0;
  border: none;
  background: transparent;
  color: var(--text-primary);
  height: 2.3rem;
}

.notes-search__button {
  height: 2.1rem;
  padding: 0 0.85rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  background: rgba(255, 255, 255, 0.36);
  color: var(--text-secondary);
  font-size: 0.85rem;
}

.notes-search__button:hover {
  color: var(--text-primary);
  border-color: var(--border-strong);
}

.notes-search__hint {
  color: var(--text-muted);
  font-size: 0.84rem;
}

.notes-sidebar__meta {
  display: flex;
  gap: var(--space-md);
  color: var(--text-secondary);
  font-size: 0.85rem;
}

.notes-sidebar__meta strong {
  font-family: var(--font-display);
  font-size: 1.2rem;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.tree-container {
  min-height: 280px;
  overflow-y: auto;
  padding-top: var(--space-sm);
  padding-right: 0.25rem;
  border-top: 1px solid var(--border-soft);
}

.tree-container :deep(.n-tree-node-content) {
  border-radius: var(--radius-md);
  min-height: 2rem;
  padding-inline: 0.35rem;
  transition: background-color var(--duration-fast) var(--ease-standard), box-shadow var(--duration-fast) var(--ease-standard), color var(--duration-fast) var(--ease-standard);
}

.tree-container :deep(.note-node--leaf .n-tree-node-content:hover),
.tree-container :deep(.note-node--active .n-tree-node-content) {
  background: var(--accent-soft) !important;
  box-shadow: inset 0 0 0 1px rgba(41, 70, 58, 0.12);
}

.tree-container :deep(.n-tree-node-content__text) {
  color: var(--text-secondary);
}

.tree-container :deep(.note-node--leaf .n-tree-node-content:hover .n-tree-node-content__text),
.tree-container :deep(.note-node--active .n-tree-node-content__text) {
  color: var(--text-primary);
}

.tree-container :deep(.note-node--leaf .n-tree-node-content__text) {
  color: var(--text-primary);
}

.notes-content {
  min-width: 0;
}

@media (max-width: 980px) {
  .notes-page--wide .notes-frame {
    width: 100%;
  }

  .notes-page--wide .notes-layout {
    gap: var(--space-lg);
  }

  .notes-layout {
    grid-template-columns: 1fr;
  }

  .notes-page--wide .notes-sidebar,
  .notes-sidebar {
    position: static;
    transform: none;
  }
}
</style>

<style>
[data-theme='dark'] .tree-container .note-node--leaf .n-tree-node-content:hover,
[data-theme='dark'] .tree-container .note-node--active .n-tree-node-content {
  background: linear-gradient(90deg, rgba(53, 95, 79, 0.34), rgba(35, 59, 77, 0.2)) !important;
  box-shadow: inset 0 0 0 1px rgba(143, 182, 163, 0.18);
}

[data-theme='dark'] .tree-container .note-node--leaf .n-tree-node-content:hover .n-tree-node-content__text,
[data-theme='dark'] .tree-container .note-node--active .n-tree-node-content__text {
  color: #f1eadf;
}

[data-theme='dark'] .tree-container .note-node--folder .n-tree-node-content:hover {
  background: rgba(255, 255, 255, 0.04) !important;
}
</style>
