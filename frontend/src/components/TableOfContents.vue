<template>
  <div class="toc-container" :class="{ collapsed: isCollapsed }">
    <div class="toc-header" @click="toggleCollapse">
      <span class="toc-title">📑 目录</span>
      <button class="toc-toggle-btn">
        {{ isCollapsed ? '←' : '→' }}
      </button>
    </div>
    
    <!-- 折叠时显示竖排文字 -->
    <div class="toc-collapsed-text" v-if="isCollapsed" @click="toggleCollapse">
      <span>目</span>
      <span>录</span>
      <span>大</span>
      <span>纲</span>
    </div>
    
    <div class="toc-content" v-show="!isCollapsed">
      <div v-if="headings.length === 0" class="toc-empty">
        暂无目录
      </div>
      <ul v-else class="toc-list">
        <li 
          v-for="heading in headings" 
          :key="heading.id"
          :class="['toc-item', `toc-level-${heading.level}`]"
          @click="scrollToHeading(heading.id)"
        >
          <span class="toc-dot">•</span>
          <span class="toc-text">{{ heading.text }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'

const props = defineProps({
  content: {
    type: String,
    default: ''
  }
})

const headings = ref([])
const isCollapsed = ref(false)

// 解析 Markdown 内容中的 h2, h3 标题
const parseHeadings = (content) => {
  if (!content) {
    headings.value = []
    return
  }
  
  const lines = content.split('\n')
  const result = []
  let idCounter = 0
  
  for (const line of lines) {
    // 匹配 ## 和 ### 开头的标题（h2 和 h3）
    const h2Match = line.match(/^##\s+(.+)$/)
    const h3Match = line.match(/^###\s+(.+)$/)
    
    if (h2Match) {
      result.push({
        level: 2,
        text: h2Match[1].trim(),
        id: `heading-${idCounter++}`
      })
    } else if (h3Match) {
      result.push({
        level: 3,
        text: h3Match[1].trim(),
        id: `heading-${idCounter++}`
      })
    }
  }
  
  headings.value = result
}

// 点击跳转到对应标题
const scrollToHeading = (id) => {
  const element = document.getElementById(id)
  if (element) {
    element.scrollIntoView({ 
      behavior: 'smooth',
      block: 'start'
    })
  }
}

// 切换折叠状态
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

// 监听内容变化，重新解析标题
watch(() => props.content, (newContent) => {
  nextTick(() => {
    parseHeadings(newContent)
  })
}, { immediate: true })
</script>
