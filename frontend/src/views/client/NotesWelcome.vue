<template>
  <div class="notes-welcome">
    <div class="welcome-content glass-card">
      <n-icon :component="BookOutline" :size="48" class="welcome-icon" />
      <h2>欢迎来到笔记</h2>
      <p>从左侧选择一篇笔记开始阅读</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NIcon } from 'naive-ui'
import { BookOutline } from '@vicons/ionicons5'
</script>

<style scoped>
.notes-welcome {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.welcome-content {
  text-align: center;
  padding: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.welcome-icon {
  color: var(--accent-cyan);
}

.welcome-content h2 {
  font-size: 22px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.welcome-content p {
  color: var(--color-text-secondary);
  font-size: 15px;
}
</style>
