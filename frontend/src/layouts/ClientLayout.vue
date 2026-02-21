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
            <router-link to="/" class="gradient-text">{{ profileName }}</router-link>
          </div>
          <nav class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/notes">笔记</router-link>
          </nav>
          <div class="actions">
            <n-button text @click="toggleTheme">
              <template #icon>
                <n-icon :component="theme === 'dark' ? SunnyOutline : MoonOutline" />
              </template>
            </n-button>
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
    </div>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { NConfigProvider, NButton, NIcon, darkTheme } from 'naive-ui'
import { MoonOutline, SunnyOutline } from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'
import { profileApi } from '@/api/profile'

const appStore = useAppStore()
const theme = computed(() => appStore.theme)
const naiveTheme = computed(() => (theme.value === 'dark' ? darkTheme : null))
const profileName = ref('Loading...')

const toggleTheme = () => {
  appStore.toggleTheme()
}

onMounted(async () => {
  try {
    const res = await profileApi.get()
    if (res.success && res.data) {
      profileName.value = res.data.name || 'My Site'
    }
  } catch {
    profileName.value = 'My Site'
  }
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
  font-size: 20px;
  font-weight: 700;
}

.logo a {
  font-size: 20px;
  font-weight: 700;
}

.nav {
  display: flex;
  gap: 32px;
}

.nav a {
  color: var(--color-text-secondary);
  transition: var(--transition);
  font-weight: 500;
}

.nav a:hover,
.nav a.router-link-active {
  color: var(--accent-cyan);
}

.actions :deep(.n-button) {
  color: var(--color-text-secondary);
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
}
</style>