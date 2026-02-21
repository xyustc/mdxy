<template>
  <div class="notes-page">
    <div class="notes-inner">
      <h1 class="page-title gradient-text">笔记</h1>

      <div class="notes-layout">
        <aside class="sidebar glass-card">
          <div class="search-wrap">
            <n-input
              v-model:value="searchKeyword"
              placeholder="搜索笔记..."
              clearable
              @keyup.enter="handleSearch"
            />
          </div>

          <div class="tree-container">
            <n-tree
              :data="treeData"
              :node-props="nodeProps"
              block-line
            />
          </div>
        </aside>

        <main class="content">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NInput, NTree } from 'naive-ui'
import { noteApi } from '@/api/note'
import type { NoteNode } from '@/api/types'

const router = useRouter()
const searchKeyword = ref('')
const treeData = ref<any[]>([])

const fetchTree = async () => {
  try {
    const response = await noteApi.getTree()
    if (response.success && response.data) {
      treeData.value = convertToTreeData(response.data)
    }
  } catch (error) {
    console.error('获取笔记目录失败:', error)
  }
}

const convertToTreeData = (nodes: NoteNode[]): any[] => {
  return nodes.map(node => ({
    label: node.name,
    key: node.path,
    isLeaf: node.type === 'file',
    children: node.children ? convertToTreeData(node.children) : undefined
  }))
}

const nodeProps = ({ option }: any) => ({
  onClick() {
    if (option.isLeaf) {
      router.push(`/notes/${option.key}`)
    }
  }
})

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push(`/notes/search?q=${encodeURIComponent(searchKeyword.value)}`)
  }
}

onMounted(() => {
  fetchTree()
})
</script>

<style scoped>
.notes-page {
  min-height: calc(100vh - 72px);
  padding: 40px 24px;
}

.notes-inner {
  max-width: 1400px;
  margin: 0 auto;
}

.page-title {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 32px;
}

.notes-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 24px;
  min-height: 600px;
}

.sidebar {
  position: sticky;
  top: 92px;
  height: fit-content;
  padding: 16px;
}

.search-wrap {
  margin-bottom: 12px;
}

.tree-container {
  max-height: calc(100vh - 260px);
  overflow-y: auto;
}

.content {
  min-height: 400px;
}

@media (max-width: 768px) {
  .notes-page {
    padding: 24px 16px;
  }

  .page-title {
    font-size: 28px;
    margin-bottom: 20px;
  }

  .notes-layout {
    grid-template-columns: 1fr;
    gap: 16px;
    min-height: auto;
  }

  .sidebar {
    position: static;
    max-height: 40vh;
    overflow-y: auto;
  }

  .tree-container {
    max-height: 30vh;
  }
}
</style>