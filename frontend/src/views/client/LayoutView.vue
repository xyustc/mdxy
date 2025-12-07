<template>
  <div class="app-container">
    <header class="app-header">
      <div class="header-left">
        <!-- 移动端菜单按钮 -->
        <button class="mobile-menu-btn" @click="sidebarCollapsed = !sidebarCollapsed">
          ☰
        </button>
        <div class="logo">
          <span class="logo-icon">📝</span>
          <span class="logo-text">XingYu的笔记</span>
        </div>
      </div>
      <div class="search-box">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="搜索笔记..." 
          @keyup.enter="handleSearch"
        />
        <button @click="handleSearch" class="search-btn">🔍</button>
      </div>
    </header>
    
    <!-- 移动端遮罩 -->
    <div 
      class="sidebar-overlay" 
      v-if="!sidebarCollapsed" 
      @click="sidebarCollapsed = true"
    ></div>
    
    <main class="app-main">
      <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
        <div class="sidebar-header">
          <span>📁 目录</span>
          <button @click="sidebarCollapsed = !sidebarCollapsed" class="toggle-btn">
            ✕
          </button>
        </div>
        <Sidebar v-if="!sidebarCollapsed" @note-click="handleNoteClick" />
      </aside>
      
      <section class="content">
        <router-view />
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../../components/Sidebar.vue'

const router = useRouter()
const searchQuery = ref('')
const sidebarCollapsed = ref(false)
const isMobile = ref(false)

// 检测屏幕尺寸
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
  // 移动端默认收起侧边栏
  if (isMobile.value) {
    sidebarCollapsed.value = true
  }
}

// 点击笔记时，移动端自动收起侧边栏
const handleNoteClick = () => {
  if (isMobile.value) {
    sidebarCollapsed.value = true
  }
}

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ name: 'search', query: { q: searchQuery.value.trim() } })
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
/* 这里可以添加特定的样式 */
</style>