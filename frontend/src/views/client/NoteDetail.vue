<template>
  <div class="note-detail glass-card">
    <n-spin :show="loading">
      <div v-if="content" class="markdown-body" v-html="renderedContent"></div>
      <n-empty v-else-if="!loading" description="笔记不存在" />
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
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
  html: false,
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
      fetchContent(Array.isArray(newPath) ? newPath.join('/') : newPath)
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.note-detail {
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

.markdown-body :deep(h2) {
  padding-left: 12px;
  border-left: 4px solid var(--accent-cyan);
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
}

.markdown-body :deep(code) {
  background: var(--color-bg-secondary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 0.9em;
  border: 1px solid var(--color-border);
}

.markdown-body :deep(pre) {
  background: var(--color-bg-tertiary);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin-bottom: 16px;
}

.markdown-body :deep(pre code) {
  background: none;
  padding: 0;
  border: none;
  color: var(--color-text-primary);
}

.markdown-body :deep(blockquote) {
  margin: 16px 0;
  padding: 12px 20px;
  border-left: 4px solid var(--accent-cyan);
  background: var(--glass-bg);
  color: var(--color-text-secondary);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
}

.markdown-body :deep(a) {
  color: var(--accent-cyan);
}

.markdown-body :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
}

.markdown-body :deep(th),
.markdown-body :deep(td) {
  border: 1px solid var(--color-border);
  padding: 10px 16px;
  text-align: left;
}

.markdown-body :deep(th) {
  background: var(--glass-bg);
  font-weight: 600;
}

.markdown-body :deep(img) {
  max-width: 100%;
  border-radius: 8px;
}

.markdown-body :deep(hr) {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: 24px 0;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin-bottom: 16px;
  padding-left: 2em;
}

.markdown-body :deep(li) {
  margin-bottom: 4px;
}

@media (max-width: 768px) {
  .note-detail {
    padding: 20px 16px;
    min-height: 300px;
  }

  .markdown-body :deep(pre) {
    padding: 12px;
    font-size: 0.85em;
  }

  .markdown-body :deep(table) {
    display: block;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  .markdown-body :deep(img) {
    max-width: 100%;
    height: auto;
  }
}
</style>

<style>
/* Dark mode syntax highlighting - must be unscoped to match [data-theme] on <html> */
[data-theme='dark'] .markdown-body .hljs-keyword { color: #c678dd; }
[data-theme='dark'] .markdown-body .hljs-string { color: #98c379; }
[data-theme='dark'] .markdown-body .hljs-number { color: #d19a66; }
[data-theme='dark'] .markdown-body .hljs-comment { color: #7f848e; font-style: italic; }
[data-theme='dark'] .markdown-body .hljs-function { color: #61afef; }
[data-theme='dark'] .markdown-body .hljs-title { color: #61afef; }
[data-theme='dark'] .markdown-body .hljs-class { color: #e5c07b; }
[data-theme='dark'] .markdown-body .hljs-variable { color: #e06c75; }
[data-theme='dark'] .markdown-body .hljs-attr { color: #d19a66; }
[data-theme='dark'] .markdown-body .hljs-built_in { color: #e5c07b; }
[data-theme='dark'] .markdown-body .hljs-type { color: #e5c07b; }
[data-theme='dark'] .markdown-body .hljs-params { color: #abb2bf; }
[data-theme='dark'] .markdown-body .hljs-literal { color: #56b6c2; }
[data-theme='dark'] .markdown-body .hljs-symbol { color: #56b6c2; }
[data-theme='dark'] .markdown-body .hljs-meta { color: #61afef; }
[data-theme='dark'] .markdown-body pre code { color: #abb2bf; }
</style>