<template>
  <div class="search-view">
    <h2>搜索结果: "{{ searchQuery }}"</h2>
    
    <div v-if="loading" class="loading-state">
      <span>搜索中...</span>
    </div>
    
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
    </div>
    
    <div v-else-if="results.length === 0" class="empty-state">
      <p>没有找到相关笔记</p>
    </div>
    
    <div v-else class="search-results">
      <p class="result-count">找到 {{ results.length }} 条结果</p>
      
      <div 
        v-for="item in results" 
        :key="item.path" 
        class="result-item"
        @click="goToNote(item.path)"
      >
        <div class="result-title">📄 {{ item.name }}</div>
        <div class="result-path">{{ item.path }}</div>
        <div class="result-context">{{ item.context }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { searchNotes } from '../api/notes'

const route = useRoute()
const router = useRouter()

const searchQuery = ref('')
const results = ref([])
const loading = ref(false)
const error = ref('')

const doSearch = async () => {
  const query = route.query.q
  if (!query) return
  
  searchQuery.value = query
  loading.value = true
  error.value = ''
  
  try {
    const result = await searchNotes(query)
    if (result.success) {
      results.value = result.data
    }
  } catch (e) {
    error.value = '搜索失败'
    console.error(e)
  } finally {
    loading.value = false
  }
}

const goToNote = (path) => {
  router.push({ name: 'note', params: { path } })
}

watch(() => route.query.q, () => {
  doSearch()
})

onMounted(() => {
  doSearch()
})
</script>
