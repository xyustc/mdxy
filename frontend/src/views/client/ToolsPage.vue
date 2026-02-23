<template>
  <div class="tools-page">
    <section class="hero">
      <h1 class="section-title gradient-text">工具箱</h1>
      <p class="section-desc">收藏的实用工具和资源</p>
    </section>

    <div class="category-tabs">
      <button
        class="tab-btn glass-card"
        :class="{ active: activeCategory === '' }"
        @click="activeCategory = ''"
      >全部</button>
      <button
        v-for="cat in categories"
        :key="cat"
        class="tab-btn glass-card"
        :class="{ active: activeCategory === cat }"
        @click="activeCategory = cat"
      >{{ cat }}</button>
    </div>

    <div class="tools-grid">
      <div v-for="tool in filteredTools" :key="tool.id" class="glass-card tool-card">
        <div class="tool-icon">
          <span v-if="tool.icon">{{ tool.icon }}</span>
          <span v-else>{{ typeIcon(tool.type) }}</span>
        </div>
        <div class="tool-info">
          <h3 class="tool-name">{{ tool.name }}</h3>
          <p class="tool-desc">{{ tool.description }}</p>
          <div class="tool-meta">
            <span class="tool-type" :class="`type-${tool.type}`">{{ tool.type }}</span>
            <span class="tool-category">{{ tool.category }}</span>
          </div>
        </div>
        <a v-if="tool.url" :href="tool.url" target="_blank" rel="noopener" class="tool-link">
          访问 →
        </a>
        <button
          v-else-if="tool.type === 'game'"
          class="tool-link tool-play-btn"
          @click="activeGame = activeGame === tool.id ? null : tool.id"
        >
          {{ activeGame === tool.id ? '收起' : '开始游戏' }} →
        </button>
      </div>
    </div>

    <Transition name="game-expand">
      <div v-if="activeGame" class="game-area">
        <Game2048 v-if="activeGameTool?.name === '2048'" @close="activeGame = null" />
        <GameSnake v-else-if="activeGameTool?.name === '贪吃蛇'" @close="activeGame = null" />
      </div>
    </Transition>

    <div v-if="filteredTools.length === 0 && !loading" class="empty-state">
      <p>暂无工具</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { toolApi } from '@/api/tool'
import type { Tool } from '@/api/types'
import Game2048 from '@/components/Game2048.vue'
import GameSnake from '@/components/GameSnake.vue'

const tools = ref<Tool[]>([])
const categories = ref<string[]>([])
const activeCategory = ref('')
const loading = ref(true)
const activeGame = ref<number | null>(null)

const filteredTools = computed(() => {
  if (!activeCategory.value) return tools.value
  return tools.value.filter(t => t.category === activeCategory.value)
})

const activeGameTool = computed(() =>
  tools.value.find(t => t.id === activeGame.value)
)

function typeIcon(type: string) {
  const icons: Record<string, string> = { video: '🎬', software: '💻', game: '🎮', link: '🔗' }
  return icons[type] || '🔧'
}

onMounted(async () => {
  try {
    const [toolsRes, catsRes] = await Promise.all([
      toolApi.list(),
      toolApi.getCategories()
    ])
    if (toolsRes.success) tools.value = toolsRes.data || []
    if (catsRes.success) categories.value = catsRes.data || []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.tools-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 60px 32px;
}

.hero {
  text-align: center;
  margin-bottom: 40px;
}

.section-title {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 12px;
}

.section-desc {
  color: var(--color-text-secondary);
  font-size: 16px;
}

.category-tabs {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 36px;
}

.tab-btn {
  padding: 8px 20px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-secondary);
  transition: var(--transition);
}

.tab-btn.active {
  color: var(--accent-cyan);
  border-color: var(--accent-cyan);
}


.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.tool-card {
  display: flex;
  flex-direction: column;
  padding: 24px;
  transition: transform 0.2s;
}

.tool-card:hover {
  transform: translateY(-4px);
}

.tool-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.tool-info {
  flex: 1;
}

.tool-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 8px;
}

.tool-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin-bottom: 12px;
  flex: 1;
}

.tool-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.tool-type {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  background: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
}

.type-video { color: var(--accent-coral); }
.type-software { color: var(--accent-cyan); }
.type-game { color: var(--accent-orange); }
.type-link { color: var(--accent-teal); }

.tool-category {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.tool-link {
  display: inline-block;
  color: var(--accent-cyan);
  font-weight: 500;
  font-size: 14px;
  transition: var(--transition);
}

.tool-link:hover {
  color: var(--accent-teal);
}

.tool-play-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  font-size: 14px;
  font-weight: 500;
}

.game-area {
  margin-top: 32px;
  padding: 32px 0;
}

.game-expand-enter-active,
.game-expand-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.game-expand-enter-from,
.game-expand-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}

.empty-state {
  text-align: center;
  padding: 60px;
  color: var(--color-text-tertiary);
}

@media (max-width: 768px) {
  .tools-page { padding: 40px 16px; }
  .section-title { font-size: 28px; }
  .tools-grid { grid-template-columns: 1fr; }
}
</style>
