<template>
  <div class="note-detail">
    <n-spin :show="loading">
      <div v-if="content" class="markdown-body" v-html="renderedContent"></div>
      <n-empty v-else description="笔记不存在" />
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { NSpin, NEmpty } from 'naive-ui'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import { noteApi } from '@/api/note'
import 'highlight.js/styles/github.css'

const route = useRoute()
const content = ref('')
const loading = ref(true)

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: (str, lang) => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch {}
    }
    return ''
  }
})

const renderedContent = computed(() => {
  return content.value ? md.render(content.value) : ''
})

const fetchContent = async (path: string) => {
  loading.value = true
  try {
    const response = await noteApi.getContent(path)
    if (response.success && response.data) {
      content.value = response.data.content
    }
  } catch (error) {
    console.error('获取笔记内容失败:', error)
    content.value = ''
  } finally {
    loading.value = false
  }
}

watch(
  () => route.params.path,
  (newPath) => {
    if (newPath) {
      fetchContent(newPath as string)
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.note-detail {
  background: var(--color-bg-primary);
  border-radius: var(--radius-md);
  padding: 32px;
  min-height: 400px;
}

.markdown-body {
  line-height: 1.8;
  color: var(--color-text-primary);
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
}

.markdown-body :deep(code) {
  background: var(--color-bg-secondary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.markdown-body :deep(pre) {
  background: var(--color-bg-secondary);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin-bottom: 16px;
}

.markdown-body :deep(pre code) {
  background: none;
  padding: 0;
}
</style>
