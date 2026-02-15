<template>
  <n-config-provider :theme="naiveTheme">
    <n-layout class="client-layout">
      <n-layout-header class="header" bordered>
        <div class="header-content">
          <div class="logo">
            <router-link to="/">{{ profileName }}</router-link>
          </div>
          <nav class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/notes">八股笔记</router-link>
          </nav>
          <div class="actions">
            <n-button text @click="toggleTheme">
              <template #icon>
                <n-icon :component="theme === 'dark' ? SunnyOutline : MoonOutline" />
              </template>
            </n-button>
          </div>
        </div>
      </n-layout-header>

      <n-layout-content class="content">
        <router-view />
      </n-layout-content>

      <n-layout-footer class="footer" bordered>
        <div class="footer-content">
          <p>&copy; 2024 {{ profileName }}. All rights reserved.</p>
        </div>
      </n-layout-footer>
    </n-layout>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NLayout, NLayoutHeader, NLayoutContent, NLayoutFooter, NButton, NIcon, darkTheme } from 'naive-ui'
import { MoonOutline, SunnyOutline } from '@vicons/ionicons5'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const theme = computed(() => appStore.theme)
const naiveTheme = computed(() => (theme.value === 'dark' ? darkTheme : null))

const profileName = 'Your Name'

const toggleTheme = () => {
  appStore.toggleTheme()
}
</script>

<style scoped>
.client-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--color-bg-primary);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  font-size: 20px;
  font-weight: 600;
}

.nav {
  display: flex;
  gap: 32px;
}

.nav a {
  color: var(--color-text-secondary);
  transition: var(--transition);
}

.nav a:hover,
.nav a.router-link-active {
  color: var(--color-primary);
}

.content {
  flex: 1;
  padding: 40px 20px;
}

.footer {
  background: var(--color-bg-secondary);
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
  color: var(--color-text-secondary);
  font-size: 14px;
}
</style>
