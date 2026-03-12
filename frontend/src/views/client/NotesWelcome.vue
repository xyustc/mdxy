<template>
  <div class="notes-welcome">
    <div class="welcome-content surface-panel">
      <BookOutline class="welcome-icon" />
      <h2>开始阅读</h2>
      <p>从左侧索引选择一篇笔记，进入阅读与复盘。</p>
    </div>
  </div>
</template>

<script setup lang="ts">
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
  padding: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.welcome-icon {
  width: 48px;
  height: 48px;
  color: var(--accent-primary);
}

.welcome-content h2 {
  font-family: var(--font-display);
  font-size: clamp(1.8rem, 3.4vw, 2.2rem);
  font-weight: 600;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.welcome-content p {
  max-width: 30ch;
  color: var(--text-secondary);
  font-size: 15px;
  line-height: 1.8;
}
</style>
