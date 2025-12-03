<template>
  <div class="sidebar-content">
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="tree">
      <TreeNode 
        v-for="item in noteTree" 
        :key="item.path" 
        :node="item"
        @note-click="$emit('note-click')"
      />
    </div>
    <div v-if="!loading && noteTree.length === 0" class="empty">
      暂无笔记
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getNoteTree } from '../api/notes'
import TreeNode from './TreeNode.vue'

defineEmits(['note-click'])

const noteTree = ref([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const result = await getNoteTree()
    if (result.success) {
      noteTree.value = result.data
    }
  } catch (e) {
    error.value = '加载目录失败'
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
