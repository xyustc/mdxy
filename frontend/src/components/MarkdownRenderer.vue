<template>
  <div class="markdown-body" v-html="renderedContent"></div>
</template>

<script setup>
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
})

// 用于生成标题 id 的计数器
let headingIdCounter = 0

// 创建 markdown-it 实例，配置代码高亮
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang }).value}</code></pre>`
      } catch (__) {}
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  }
})

// 自定义标题渲染，为 h2 和 h3 添加 id
const defaultHeadingRenderer = md.renderer.rules.heading_open || function(tokens, idx, options, env, self) {
  return self.renderToken(tokens, idx, options)
}

md.renderer.rules.heading_open = function(tokens, idx, options, env, self) {
  const token = tokens[idx]
  const level = token.tag // h1, h2, h3 等
  
  // 只为 h2 和 h3 添加 id
  if (level === 'h2' || level === 'h3') {
    token.attrSet('id', `heading-${headingIdCounter++}`)
  }
  
  return defaultHeadingRenderer(tokens, idx, options, env, self)
}

const renderedContent = computed(() => {
  // 每次渲染前重置计数器
  headingIdCounter = 0
  return md.render(props.content)
})
</script>
