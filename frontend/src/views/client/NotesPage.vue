<template>
  <div class="notes-page">
    <div class="container">
      <h1 class="page-title">八股笔记</h1>

      <div class="notes-layout">
        <aside class="sidebar">
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索笔记..."
            clearable
            @keyup.enter="handleSearch"
          />

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
  min-height: 100%;
}

.page-title {
  margin-bottom: 32px;
}

.notes-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 32px;
  min-height: 600px;
}

.sidebar {
  position: sticky;
  top: 84px;
  height: fit-content;
}

.tree-container {
  margin-top: 16px;
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}

.content {
  min-height: 400px;
}

@media (max-width: 768px) {
  .notes-layout {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: static;
  }
}
</style>
