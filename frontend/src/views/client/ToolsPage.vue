<template>
  <div class="tools-page section-shell">
    <div class="app-frame tools-frame">
      <header class="tools-hero">
        <span class="section-kicker">Field Kit</span>
        <h1 class="section-title">工具箱</h1>
        <p class="section-description">收录常用工具、实践资源与轻量演示，按场景筛选并快速直达。</p>
      </header>

      <section class="category-tabs">
        <button class="tab-btn" :class="{ active: activeCategory === '' }" @click="activeCategory = ''">全部</button>
        <button
          v-for="cat in categories"
          :key="cat"
          class="tab-btn"
          :class="{ active: activeCategory === cat }"
          @click="activeCategory = cat"
        >
          {{ cat }}
        </button>
      </section>

      <section class="tools-grid" v-if="filteredTools.length">
        <article v-for="tool in filteredTools" :key="tool.id" class="surface-panel tool-card">
          <div class="tool-card__header">
            <div class="tool-icon">
              <span v-if="tool.icon">{{ tool.icon }}</span>
              <span v-else>{{ typeIcon(tool.type) }}</span>
            </div>
            <div>
              <h3 class="tool-name">{{ tool.name }}</h3>
              <p class="tool-desc">{{ tool.description || '暂无描述' }}</p>
            </div>
          </div>

          <div class="tool-meta">
            <span class="tool-type" :class="`type-${normalizeToolType(tool.type)}`">{{ normalizeToolType(tool.type) }}</span>
            <span class="tool-category">{{ tool.category || '未分类' }}</span>
          </div>

          <router-link v-if="tool.url && isInternalToolUrl(tool.url)" :to="tool.url" class="tool-link">打开工具</router-link>
          <a v-else-if="tool.url" :href="tool.url" target="_blank" rel="noopener" class="tool-link">访问资源</a>
          <button
            v-else-if="tool.type === 'game'"
            class="tool-link tool-play-btn"
            @click="activeGame = activeGame === tool.id ? null : tool.id"
          >
            {{ activeGame === tool.id ? '收起游戏' : '开始游戏' }}
          </button>
        </article>
      </section>

      <Transition name="game-expand">
        <section v-if="activeGame" class="game-area surface-panel">
          <Game2048 v-if="activeGameTool?.name === '2048'" @close="activeGame = null" />
          <GameSnake v-else-if="activeGameTool?.name === '贪吃蛇'" @close="activeGame = null" />
        </section>
      </Transition>

      <section v-if="!filteredTools.length && !loading" class="surface-panel empty-state">
        <p>当前分类暂无工具，试试切换分类看看。</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toolApi } from '@/api/tool'
import type { Tool } from '@/api/types'
import Game2048 from '@/components/Game2048.vue'
import GameSnake from '@/components/GameSnake.vue'
import { isResolvableToolRoute } from '@/utils/toolNavigation'

const router = useRouter()
const tools = ref<Tool[]>([])
const categories = ref<string[]>([])
const activeCategory = ref('')
const loading = ref(true)
const activeGame = ref<number | null>(null)

const filteredTools = computed(() => {
  if (!activeCategory.value) return tools.value
  return tools.value.filter((t) => t.category === activeCategory.value)
})

const activeGameTool = computed(() => tools.value.find((t) => t.id === activeGame.value))

function typeIcon(type: string) {
  const icons: Record<string, string> = { video: '🎬', app: '🧩', game: '🎮', link: '🔗' }
  return icons[normalizeToolType(type)] || '🔧'
}

function normalizeToolType(type: string) {
  return type === 'software' ? 'app' : type
}

function isInternalToolUrl(url: string) {
  return isResolvableToolRoute(router, url)
}

onMounted(async () => {
  try {
    const [toolsRes, catsRes] = await Promise.all([toolApi.list(), toolApi.getCategories()])
    if (toolsRes.success) tools.value = toolsRes.data || []
    if (catsRes.success) categories.value = catsRes.data || []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.tools-frame {
  display: grid;
  gap: var(--space-xl);
}

.tools-hero {
  display: grid;
  gap: var(--space-sm);
}

.tools-hero .section-title {
  max-width: 10ch;
}

.tools-hero .section-description {
  max-width: 38ch;
}

.category-tabs {
  display: flex;
  gap: 0.55rem;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 0.52rem 1rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-secondary);
  background: var(--bg-panel);
}

.tab-btn.active {
  color: var(--text-primary);
  background: var(--accent-soft);
  border-color: var(--border-strong);
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--space-md);
}

.tool-card {
  min-height: 15rem;
  padding: 1.2rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.tool-card__header {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: var(--space-md);
}

.tool-icon {
  font-size: 1.8rem;
}

.tool-name {
  font-family: var(--font-display);
  font-size: 1.36rem;
  font-weight: 600;
  letter-spacing: -0.03em;
}

.tool-desc {
  margin-top: 0.35rem;
  color: var(--text-secondary);
  line-height: 1.65;
}

.tool-meta {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.tool-type,
.tool-category {
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  padding: 0.2rem 0.55rem;
  font-size: 0.76rem;
}

.tool-type {
  color: var(--text-secondary);
}

.type-video {
  color: #9f4f31;
}

.type-app {
  color: #355f54;
}

.type-game {
  color: #8c6822;
}

.type-link {
  color: #4b5d9f;
}

.tool-category {
  color: var(--text-muted);
}

.tool-link {
  margin-top: auto;
  display: inline-flex;
  justify-content: center;
  min-height: 2.6rem;
  align-items: center;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.35);
}

.tool-play-btn {
  width: 100%;
}

.game-area {
  padding: var(--space-xl);
}

.game-expand-enter-active,
.game-expand-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
  overflow: hidden;
}

.game-expand-enter-from,
.game-expand-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}

.empty-state {
  padding: var(--space-2xl);
  text-align: center;
  color: var(--text-muted);
}

@media (max-width: 900px) {
  .tools-grid {
    grid-template-columns: 1fr;
  }

  .game-area {
    padding: var(--space-md);
  }
}
</style>
