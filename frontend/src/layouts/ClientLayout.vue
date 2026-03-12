<template>
  <div class="client-layout">
    <div class="layout-ornament layout-ornament-left"></div>
    <div class="layout-ornament layout-ornament-right"></div>

    <header class="site-header">
      <div class="app-frame site-header__inner">
        <router-link to="/" class="brand-mark" aria-label="返回首页">
          <span class="brand-mark__index">VOL 02</span>
          <div class="brand-mark__copy">
            <strong>{{ profileName || 'MDXY' }}</strong>
            <small>{{ activeSection.caption }}</small>
          </div>
        </router-link>

        <nav class="site-nav" ref="navEl" aria-label="主导航">
          <router-link
            v-for="(item, index) in navItems"
            :key="item.to"
            :to="item.to"
            class="site-nav__link"
            :class="{ 'site-nav__link--active': isActiveSection(item.to) }"
          >
            <span class="site-nav__index">0{{ index + 1 }}</span>
            <span>{{ item.label }}</span>
          </router-link>
          <span class="site-nav__indicator" :style="indicatorStyle"></span>
        </nav>

        <div class="site-actions">
          <button class="site-action" type="button" @click="showSearch = true" title="搜索 (Ctrl+K)">
            <SearchOutline class="site-action__icon" />
            <span>Index</span>
            <small>Ctrl+K</small>
          </button>
          <button
            class="site-action site-action--theme"
            type="button"
            @click="handleToggleTheme"
            :title="theme === 'dark' ? '切换亮色' : '切换暗色'"
          >
            <component :is="themeIcon" class="site-action__icon" />
            <span>{{ theme === 'dark' ? 'Daylight' : 'After Dark' }}</span>
          </button>
        </div>
      </div>
    </header>

    <main class="site-main">
      <router-view />
    </main>

    <footer class="site-footer">
      <div class="app-frame site-footer__grid">
        <div>
          <span class="section-kicker">Colophon</span>
          <p class="site-footer__title">{{ profileName || 'MDXY' }}</p>
          <p class="site-footer__text">一个持续更新的个人出版系统，记录技术、工具与长期兴趣。</p>
        </div>
        <div>
          <span class="section-kicker">Sections</span>
          <div class="site-footer__links">
            <router-link v-for="item in navItems" :key="item.to" :to="item.to">{{ item.label }}</router-link>
          </div>
        </div>
        <div>
          <span class="section-kicker">Now Reading</span>
          <p class="site-footer__text">{{ activeSection.caption }}</p>
          <p class="site-footer__copyright">&copy; {{ new Date().getFullYear() }} {{ profileName || 'MDXY' }}</p>
        </div>
      </div>
    </footer>

    <SearchModal v-model="showSearch" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, nextTick, watch, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { MoonOutline, SunnyOutline, SearchOutline } from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'
import { profileApi } from '@/api/profile'
import SearchModal from '@/components/SearchModal.vue'

const appStore = useAppStore()
const route = useRoute()
const theme = computed(() => appStore.theme)
const profileName = ref('')
const showSearch = ref(false)
const navEl = ref<HTMLElement>()
const indicatorStyle = reactive({ left: '0px', width: '0px', opacity: '0' })
const navItems = [
  { to: '/', label: '首页', caption: 'Cover Story' },
  { to: '/notes', label: '笔记', caption: 'Reading Room' },
  { to: '/tools', label: '工具箱', caption: 'Field Kit' }
]

const activeSection = computed(() => {
  return navItems.find((item) => (item.to === '/' ? route.path === '/' : route.path.startsWith(item.to))) || navItems[0]
})

const themeIcon = computed(() => (theme.value === 'dark' ? SunnyOutline : MoonOutline))

function isActiveSection(path: string) {
  return activeSection.value.to === path
}

function updateIndicator() {
  if (!navEl.value) return
  const activeLink = navEl.value.querySelector('.site-nav__link--active') as HTMLElement | null
  if (activeLink) {
    indicatorStyle.left = activeLink.offsetLeft + 'px'
    indicatorStyle.width = activeLink.offsetWidth + 'px'
    indicatorStyle.opacity = '1'
  } else {
    indicatorStyle.opacity = '0'
  }
}

watch(
  () => route.path,
  () => nextTick(updateIndicator)
)

function handleToggleTheme() {
  appStore.toggleTheme()
}

function handleKeydown(e: KeyboardEvent) {
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
  position: relative;
  overflow-x: hidden;
}

.layout-ornament {
  position: fixed;
  width: 32rem;
  height: 32rem;
  border-radius: 50%;
  filter: blur(64px);
  opacity: 0.28;
  pointer-events: none;
  z-index: 0;
}

.layout-ornament-left {
  top: -12rem;
  left: -9rem;
  background: rgba(139, 94, 60, 0.18);
}

.layout-ornament-right {
  right: -11rem;
  bottom: -14rem;
  background: rgba(41, 70, 58, 0.2);
}

.site-header {
  position: sticky;
  top: 0;
  z-index: 20;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  background: linear-gradient(180deg, rgba(246, 241, 232, 0.92), rgba(246, 241, 232, 0.66));
  border-bottom: 1px solid var(--border-primary);
}

[data-theme='dark'] .site-header {
  background: linear-gradient(180deg, rgba(19, 22, 24, 0.92), rgba(19, 22, 24, 0.66));
}

.site-header__inner {
  min-height: var(--header-height);
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: var(--space-xl);
  position: relative;
  z-index: 1;
}

.brand-mark {
  display: inline-grid;
  grid-template-columns: auto 1fr;
  gap: var(--space-md);
  align-items: center;
  color: var(--text-primary);
}

.brand-mark__index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 4.5rem;
  padding: 0.55rem 0.8rem;
  border-radius: var(--radius-pill);
  background: var(--accent-soft);
  color: var(--accent-primary);
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.brand-mark__copy {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.brand-mark__copy strong {
  font-family: var(--font-display);
  font-size: 1.35rem;
  font-weight: 600;
  letter-spacing: -0.04em;
}

.brand-mark__copy small {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.site-nav {
  justify-self: center;
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.35rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  background: rgba(255, 255, 255, 0.26);
}

[data-theme='dark'] .site-nav {
  background: rgba(255, 255, 255, 0.03);
}

.site-nav__link {
  position: relative;
  z-index: 1;
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  padding: 0.7rem 1rem;
  border-radius: var(--radius-pill);
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.site-nav__link:hover,
.site-nav__link--active {
  color: var(--text-primary);
}

.site-nav__index {
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.16em;
  color: var(--text-muted);
}

.site-nav__indicator {
  position: absolute;
  inset: auto auto 0.35rem 0;
  height: calc(100% - 0.7rem);
  background: var(--bg-panel-strong);
  border-radius: var(--radius-pill);
  box-shadow: var(--shadow-sm);
  transition: left var(--duration-base) var(--ease-standard), width var(--duration-base) var(--ease-standard), opacity var(--duration-base) var(--ease-standard);
}

.site-actions {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.site-action {
  display: inline-flex;
  align-items: center;
  gap: 0.65rem;
  padding: 0.8rem 1rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  background: rgba(255, 255, 255, 0.28);
  color: var(--text-secondary);
  transition: transform var(--duration-fast) var(--ease-standard), border-color var(--duration-fast) var(--ease-standard), color var(--duration-fast) var(--ease-standard);
}

[data-theme='dark'] .site-action {
  background: rgba(255, 255, 255, 0.03);
}

.site-action:hover {
  transform: translateY(-1px);
  border-color: var(--border-strong);
  color: var(--text-primary);
}

.site-action small {
  padding-left: 0.6rem;
  border-left: 1px solid var(--border-primary);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 0.72rem;
}

.site-action--theme small {
  display: none;
}

.site-action__icon {
  width: 1rem;
  height: 1rem;
}

.site-main {
  position: relative;
  z-index: 1;
}

.site-footer {
  position: relative;
  z-index: 1;
  margin-top: auto;
  padding: var(--space-3xl) 0 var(--space-2xl);
  border-top: 1px solid var(--border-primary);
}

.site-footer__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--space-2xl);
}

.site-footer__title {
  margin-top: var(--space-md);
  font-family: var(--font-display);
  font-size: 1.4rem;
  letter-spacing: -0.03em;
}

.site-footer__text,
.site-footer__copyright {
  margin-top: var(--space-sm);
  color: var(--text-secondary);
  line-height: 1.8;
}

.site-footer__links {
  margin-top: var(--space-md);
  display: flex;
  flex-wrap: wrap;
  gap: 0.85rem;
}

@media (max-width: 1100px) {
  .site-header__inner {
    grid-template-columns: 1fr;
    padding-top: var(--space-lg);
    padding-bottom: var(--space-lg);
  }

  .site-nav {
    justify-self: start;
    overflow-x: auto;
    max-width: 100%;
  }

  .site-actions {
    justify-self: start;
  }

  .site-footer__grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .brand-mark {
    grid-template-columns: 1fr;
    gap: var(--space-sm);
  }

  .site-nav {
    width: 100%;
    justify-content: flex-start;
  }

  .site-nav__link {
    white-space: nowrap;
  }

  .site-actions {
    width: 100%;
    flex-direction: column;
    align-items: stretch;
  }

  .site-action {
    justify-content: space-between;
  }

  .site-action--theme {
    justify-content: center;
  }
}
</style>
