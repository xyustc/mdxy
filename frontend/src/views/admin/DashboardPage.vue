<template>
  <div class="dashboard-page">
    <section class="metrics-grid">
      <article v-for="metric in metrics" :key="metric.label" class="surface-panel metric-card">
        <component :is="metric.icon" class="metric-card__icon" />
        <div>
          <strong>{{ metric.value }}</strong>
          <span>{{ metric.label }}</span>
        </div>
      </article>
    </section>

    <section class="dashboard-panels">
      <article class="paper-sheet panel-card">
        <span class="meta-label">Today Focus</span>
        <h3>日内访问表现</h3>
        <div class="today-grid">
          <div>
            <span>今日 PV</span>
            <strong>{{ overview.today_pv }}</strong>
          </div>
          <div>
            <span>今日 UV</span>
            <strong>{{ overview.today_uv }}</strong>
          </div>
        </div>
        <p>建议结合「数据统计」页面查看 7 天与 30 天趋势变化，判断内容发布效果。</p>
      </article>

      <article class="paper-sheet panel-card">
        <span class="meta-label">Next Actions</span>
        <h3>编辑台建议流程</h3>
        <ul>
          <li>先更新个人信息中的当前方向与身份标签。</li>
          <li>维护工具箱排序，把常用入口放在前排。</li>
          <li>发布新笔记后观察 PV/UV 波动，优化目录结构。</li>
        </ul>
      </article>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Document as ElIconDocument, View as ElIconView, Collection as ElIconCollection, User as ElIconUser } from '@element-plus/icons-vue'
import { analyticsApi } from '@/api/analytics'
import type { AnalyticsOverview } from '@/api/types'

const overview = ref<AnalyticsOverview>({
  total_pv: 0,
  total_uv: 0,
  today_pv: 0,
  today_uv: 0,
  note_count: 0,
  tool_count: 0
})

const metrics = computed(() => [
  { label: '笔记总数', value: overview.value.note_count, icon: ElIconDocument },
  { label: '总访问量', value: overview.value.total_pv, icon: ElIconView },
  { label: '工具数量', value: overview.value.tool_count, icon: ElIconCollection },
  { label: '访客数量', value: overview.value.total_uv, icon: ElIconUser }
])

onMounted(async () => {
  try {
    const res = await analyticsApi.getOverview()
    if (res.success && res.data) overview.value = res.data
  } catch {
    // keep fallback zeros
  }
})
</script>

<style scoped>
.dashboard-page {
  display: grid;
  gap: var(--space-xl);
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--space-md);
}

.metric-card {
  padding: var(--space-lg);
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.metric-card__icon {
  width: 1.5rem;
  height: 1.5rem;
  color: var(--accent-primary);
}

.metric-card strong {
  display: block;
  font-family: var(--font-display);
  font-size: 2rem;
  letter-spacing: -0.04em;
}

.metric-card span {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.dashboard-panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-lg);
}

.panel-card {
  padding: var(--space-xl);
}

.panel-card h3 {
  margin-top: var(--space-sm);
  font-family: var(--font-display);
  font-size: 1.6rem;
  letter-spacing: -0.03em;
}

.today-grid {
  margin-top: var(--space-lg);
  margin-bottom: var(--space-lg);
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-md);
}

.today-grid div {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-md);
  background: var(--bg-panel);
}

.today-grid span {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.today-grid strong {
  display: block;
  margin-top: 0.25rem;
  font-family: var(--font-display);
  font-size: 1.8rem;
  letter-spacing: -0.04em;
}

.panel-card p,
.panel-card ul {
  color: var(--text-secondary);
  line-height: 1.8;
}

.panel-card ul {
  margin-top: var(--space-md);
  padding-left: 1.2rem;
}

@media (max-width: 1100px) {
  .metrics-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dashboard-panels {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .metrics-grid,
  .today-grid {
    grid-template-columns: 1fr;
  }
}
</style>
