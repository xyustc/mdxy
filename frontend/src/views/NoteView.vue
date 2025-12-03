<template>
  <div class="note-view-container">
    <div class="note-view">
      <div v-if="loading" class="loading-state">
        <span>加载中...</span>
      </div>
      
      <div v-else-if="error" class="error-state">
        <h2>😢 加载失败</h2>
        <p>{{ error }}</p>
        <button @click="loadNote">重试</button>
      </div>
      
      <div v-else class="note-content">
        <div class="note-header">
          <h1 class="note-title">{{ noteTitle }}</h1>
        </div>
        <MarkdownRenderer :content="noteContent" />
      </div>
    </div>
    
    <!-- 右侧大纲 -->
    <TableOfContents v-if="!loading && !error" :content="noteContent" />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getNoteContent } from '../api/notes'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import TableOfContents from '../components/TableOfContents.vue'

const route = useRoute()

const noteContent = ref('')
const loading = ref(true)
const error = ref('')

const notePath = computed(() => route.params.path || '')
const noteTitle = computed(() => {
  const path = notePath.value
  if (!path) return ''
  const parts = path.split('/')
  const filename = parts[parts.length - 1]
  return filename.replace(/\.md$/i, '')
})

const loadNote = async () => {
  if (!notePath.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    const result = await getNoteContent(notePath.value)
    if (result.success) {
      noteContent.value = result.data.content
    } else {
      error.value = '笔记不存在'
    }
  } catch (e) {
    error.value = e.response?.data?.detail || '加载笔记失败'
    console.error(e)
  } finally {
    loading.value = false
  }
}

watch(() => route.params.path, () => {
  loadNote()
})

onMounted(() => {
  loadNote()
})
</script>
