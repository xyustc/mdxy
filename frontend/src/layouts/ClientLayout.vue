<template>
  <n-config-provider :theme="naiveTheme">
    <div class="client-layout">
      <!-- Animated background -->
      <div class="animated-bg">
        <div class="bg-orb bg-orb-1"></div>
        <div class="bg-orb bg-orb-2"></div>
        <div class="bg-orb bg-orb-3"></div>
      </div>

      <header class="header glass-card">
        <div class="header-content">
          <div class="logo">
            <router-link to="/" class="gradient-text">{{ profileName || '\u00A0' }}</router-link>
          </div>
          <nav class="nav" ref="navEl">
            <router-link to="/" exact-active-class="nav-active">首页</router-link>
            <router-link to="/notes" active-class="nav-active">笔记</router-link>
            <router-link to="/tools" active-class="nav-active">工具箱</router-link>
            <span class="nav-indicator" :style="indicatorStyle"></span>
          </nav>
          <div class="actions">
            <button class="search-btn" @click="showSearch = true" title="搜索 (Ctrl+K)">
              <n-icon :size="18" :component="SearchOutline" />
            </button>
          </div>
        </div>
      </header>

      <main class="content">
        <router-view />
      </main>

      <footer class="footer glass-card">
        <div class="footer-content">
          <p>&copy; {{ new Date().getFullYear() }} {{ profileName }}</p>
        </div>
      </footer>

      <SearchModal v-model="showSearch" />

      <button
        class="theme-fab glass-card"
        :class="{ 'theme-fab-spin': isThemeAnimating }"
        @click="handleToggleTheme"
        :title="theme === 'dark' ? '切换亮色' : '切换暗色'"
      >
        <n-icon :size="20" :component="theme === 'dark' ? SunnyOutline : MoonOutline" />
      </button>
    </div>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, nextTick, watch, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { NConfigProvider, NButton, NIcon, darkTheme } from 'naive-ui'
import { MoonOutline, SunnyOutline, SearchOutline } from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'
import { profileApi } from '@/api/profile'
import SearchModal from '@/components/SearchModal.vue'

const appStore = useAppStore()
const route = useRoute()
const theme = computed(() => appStore.theme)
const naiveTheme = computed(() => (theme.value === 'dark' ? darkTheme : null))
const profileName = ref('')
const showSearch = ref(false)
const navEl = ref<HTMLElement>()
const indicatorStyle = reactive({ left: '0px', width: '0px', opacity: '0' })
const isThemeAnimating = ref(false)

function updateIndicator() {
  if (!navEl.value) return
  const activeLink = navEl.value.querySelector('.nav-active') as HTMLElement
  if (activeLink) {
    indicatorStyle.left = activeLink.offsetLeft + 'px'
    indicatorStyle.width = activeLink.offsetWidth + 'px'
    indicatorStyle.opacity = '1'
  } else {
    indicatorStyle.opacity = '0'
  }
}

watch(() => route.path, () => nextTick(updateIndicator))

const toggleTheme = () => {
  appStore.toggleTheme()
}

const handleToggleTheme = () => {
  isThemeAnimating.value = true
  // 先播放动画，在动画中间切换主题（视觉更流畅）
  setTimeout(() => {
    toggleTheme()
  }, 150)
  setTimeout(() => {
    isThemeAnimating.value = false
  }, 500)
}

const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    showSearch.value = true
  }
}

onMounted(async () => {
  document.addEventListener('keydown', handleKeydown)
  window.addEventListener('resize', updateIndicator)
  nextTick(updateIndicator)
  try {
    const res = await profileApi.get()
    if (res.success && res.data) {
      profileName.value = res.data.name || 'My Site'
    }
  } catch {
    profileName.value = 'My Site'
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('resize', updateIndicator)
})
</script>

<style scoped>
.client-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--page-gradient);
  position: relative;
  overflow-x: hidden;
}

/* Animated background orbs */
.animated-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
}

.bg-orb-1 {
  width: 400px;
  height: 400px;
  background: var(--accent-cyan);
  top: -100px;
  right: -100px;
  animation: float1 20s ease-in-out infinite;
}

.bg-orb-2 {
  width: 350px;
  height: 350px;
  background: var(--accent-coral);
  bottom: -80px;
  left: -80px;
  animation: float2 25s ease-in-out infinite;
}

.bg-orb-3 {
  width: 300px;
  height: 300px;
  background: var(--accent-teal);
  top: 50%;
  left: 50%;
  animation: float3 22s ease-in-out infinite;
}

@keyframes float1 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-80px, 80px); }
}

@keyframes float2 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(60px, -60px); }
}

@keyframes float3 {
  0%, 100% { transform: translate(-50%, -50%); }
  50% { transform: translate(-50%, -60%) translateX(40px); }
}

/* Header */
.header {
  position: sticky;
  top: 0;
  z-index: 100;
  border-radius: 0;
  border-top: none;
  border-left: none;
  border-right: none;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 32px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
/* STYLE_CHUNK_2 */.logo {
  font-size: 24px;
  font-weight: 700;
  min-width: 80px;
}

.logo a {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 2px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background-size: 200% 200%;
  background-position: 0% 50%;
  transition: background-position 0.6s ease;
}

.logo a:hover {
  background-position: 100% 50%;
}

.nav {
  display: flex;
  gap: 32px;
  position: relative;
}

.nav a {
  color: var(--color-text-secondary);
  font-weight: 500;
  padding: 4px 0;
  transition: color 0.25s ease, transform 0.25s ease;
}

.nav a:hover {
  color: var(--accent-cyan);
  transform: translateY(-1px);
}

.nav a.nav-active {
  color: var(--accent-cyan);
}

.nav-indicator {
  position: absolute;
  bottom: -6px;
  height: 2px;
  border-radius: 1px;
  background: linear-gradient(90deg, var(--accent-cyan), var(--accent-teal));
  transition: left 0.3s cubic-bezier(0.4, 0, 0.2, 1), width 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.2s ease;
  pointer-events: none;
}

.actions :deep(.n-button) {
  color: var(--color-text-secondary);
  font-size: 18px;
}

/* Search button */
.search-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: color 0.25s ease, background-color 0.25s ease;
  -webkit-tap-highlight-color: transparent;
}

.search-btn:hover {
  color: var(--accent-cyan);
  background: var(--glass-bg);
}

.search-btn:active {
  transform: scale(0.92);
}

/* Theme floating button */
.theme-fab {
  position: fixed;
  bottom: 32px;
  right: 32px;
  z-index: 999;
  width: 44px;
  height: 44px;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
  transition: transform 0.25s ease, color 0.25s ease, box-shadow 0.25s ease;
  -webkit-tap-highlight-color: transparent;
}

.theme-fab:hover {
  transform: scale(1.1);
  color: var(--accent-cyan);
}

.theme-fab:active {
  transform: scale(0.95);
}

.theme-fab-spin {
  animation: theme-spin 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes theme-spin {
  0% { transform: rotate(0deg) scale(1); }
  50% { transform: rotate(180deg) scale(0.8); }
  100% { transform: rotate(360deg) scale(1); }
}

/* Content */
.content {
  flex: 1;
  position: relative;
  z-index: 1;
}

/* Footer */
.footer {
  position: relative;
  z-index: 1;
  border-radius: 0;
  border-bottom: none;
  border-left: none;
  border-right: none;
  
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 10px;
  text-align: center;
  color: var(--color-text-tertiary);
  font-size: 14px;
}

.footer-motto {
  margin-bottom: 6px;
  font-size: 13px;
}

@media (max-width: 768px) {
  .nav {
    gap: 20px;
  }

  .header-content {
    padding: 0 16px;
  }

  .theme-fab {
    bottom: 20px;
    right: 20px;
  }
}
</style>