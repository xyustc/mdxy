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

    <section class="dashboard-panels dashboard-panels--single">
      <article class="paper-sheet panel-card panel-card--sync">
        <div class="sync-card__header">
          <div>
            <span class="meta-label">External Sync</span>
            <h3>Claude Code 速查表同步</h3>
          </div>
          <button class="sync-btn sync-btn--primary" :disabled="syncing" @click="handleSync">
            {{ syncing ? '同步中...' : '立即同步' }}
          </button>
        </div>

        <div class="sync-overview">
          <div class="sync-overview__item">
            <span>当前发布版本</span>
            <strong>{{ publishedSnapshot?.content?.meta?.version || '尚未发布' }}</strong>
          </div>
          <div class="sync-overview__item">
            <span>源站更新时间</span>
            <strong>{{ publishedSnapshot?.content?.meta?.updated_at || '—' }}</strong>
          </div>
          <div class="sync-overview__item">
            <span>最近发布时间</span>
            <strong>{{ formatDate(publishedSnapshot?.published_at) }}</strong>
          </div>
        </div>

        <p class="sync-card__tip">
          后台会按配置周期自动抓取并生成候选快照。你也可以在这里手动同步，并把某个草稿快照发布到站内页面。
        </p>

        <div class="snapshot-list">
          <article v-for="snapshot in snapshots" :key="snapshot.id" class="snapshot-item">
            <div class="snapshot-item__main">
              <div class="snapshot-item__row">
                <strong>{{ snapshot.source_version || '未识别版本' }}</strong>
                <span class="snapshot-status" :class="`snapshot-status--${snapshot.status}`">{{ snapshot.status }}</span>
              </div>
              <p>源站时间：{{ snapshot.source_updated_at_text || '—' }}</p>
              <p>抓取时间：{{ formatDate(snapshot.created_at) }}</p>
              <p v-if="snapshot.published_at">发布时间：{{ formatDate(snapshot.published_at) }}</p>
              <p v-if="snapshot.sync_error" class="snapshot-item__error">{{ snapshot.sync_error }}</p>
            </div>

            <div class="snapshot-item__actions">
              <button
                v-if="snapshot.status === 'draft'"
                class="sync-btn"
                :disabled="publishingId === snapshot.id"
                @click="handlePublish(snapshot.id)"
              >
                {{ publishingId === snapshot.id ? '发布中...' : '发布' }}
              </button>
              <span v-else class="snapshot-item__tag">已生效</span>
            </div>
          </article>

          <p v-if="!snapshots.length" class="sync-card__empty">还没有同步快照。点击“立即同步”后会在这里出现记录。</p>
        </div>
      </article>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Document as ElIconDocument, View as ElIconView, Collection as ElIconCollection, User as ElIconUser } from '@element-plus/icons-vue'
import { analyticsApi } from '@/api/analytics'
import { cheatSheetApi } from '@/api/cheatSheet'
import type { AnalyticsOverview, CheatSheetPublishedPayload, CheatSheetSnapshotSummary } from '@/api/types'

const overview = ref<AnalyticsOverview>({
  total_pv: 0,
  total_uv: 0,
  today_pv: 0,
  today_uv: 0,
  note_count: 0,
  tool_count: 0
})
const publishedSnapshot = ref<CheatSheetPublishedPayload | null>(null)
const snapshots = ref<CheatSheetSnapshotSummary[]>([])
const syncing = ref(false)
const publishingId = ref<number | null>(null)

const metrics = computed(() => [
  { label: '笔记总数', value: overview.value.note_count, icon: ElIconDocument },
  { label: '总访问量', value: overview.value.total_pv, icon: ElIconView },
  { label: '工具数量', value: overview.value.tool_count, icon: ElIconCollection },
  { label: '访客数量', value: overview.value.total_uv, icon: ElIconUser }
])

async function loadOverview() {
  try {
    const res = await analyticsApi.getOverview()
    if (res.success && res.data) overview.value = res.data
  } catch {
    // keep fallback zeros
  }
}

async function loadCheatSheetState() {
  const [publishedRes, snapshotsRes] = await Promise.allSettled([
    cheatSheetApi.getPublished('claude-code'),
    cheatSheetApi.adminListSnapshots('claude-code', 8)
  ])

  if (publishedRes.status === 'fulfilled' && publishedRes.value.success && publishedRes.value.data) {
    publishedSnapshot.value = publishedRes.value.data
  } else {
    publishedSnapshot.value = null
  }

  if (snapshotsRes.status === 'fulfilled' && snapshotsRes.value.success) {
    snapshots.value = snapshotsRes.value.data || []
  } else {
    snapshots.value = []
  }
}

function formatDate(value?: string) {
  if (!value) return '—'
  return new Date(value).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function handleSync() {
  syncing.value = true
  try {
    const res = await cheatSheetApi.adminSync('claude-code')
    if (res.success) {
      if (res.data?.changed) {
        ElMessage.success(res.data.auto_published ? '同步成功，并已自动发布首个版本' : '同步成功，已生成新的候选快照')
      } else {
        ElMessage.success('源站暂无新变化')
      }
      await loadCheatSheetState()
    }
  } catch {
    ElMessage.error('同步失败，请查看失败快照或后端日志')
    await loadCheatSheetState()
  } finally {
    syncing.value = false
  }
}

async function handlePublish(id: number) {
  publishingId.value = id
  try {
    const res = await cheatSheetApi.adminPublish('claude-code', id)
    if (res.success) {
      ElMessage.success('已发布选中的快照')
      await loadCheatSheetState()
    }
  } catch {
    ElMessage.error('发布失败')
  } finally {
    publishingId.value = null
  }
}

onMounted(async () => {
  await Promise.all([loadOverview(), loadCheatSheetState()])
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

.dashboard-panels--single {
  grid-template-columns: 1fr;
}

.panel-card {
  padding: var(--space-xl);
}

.panel-card--sync {
  display: grid;
  gap: var(--space-lg);
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

.sync-card__header {
  display: flex;
  justify-content: space-between;
  gap: var(--space-lg);
  align-items: flex-start;
}

.sync-card__header h3 {
  margin-top: var(--space-sm);
}

.sync-btn {
  min-height: 2.75rem;
  padding: 0 1rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  background: var(--bg-panel);
}

.sync-btn--primary {
  background: var(--bg-contrast);
  color: var(--text-inverse);
  border-color: var(--bg-contrast);
}

.sync-overview {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--space-md);
}

.sync-overview__item {
  padding: var(--space-md);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-primary);
  background: var(--bg-panel);
}

.sync-overview__item span {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.sync-overview__item strong {
  display: block;
  margin-top: 0.35rem;
  font-family: var(--font-display);
  font-size: 1.2rem;
  letter-spacing: -0.03em;
}

.sync-card__tip,
.sync-card__empty,
.snapshot-item p {
  color: var(--text-secondary);
}

.snapshot-list {
  display: grid;
  gap: var(--space-md);
}

.snapshot-item {
  display: flex;
  justify-content: space-between;
  gap: var(--space-lg);
  align-items: flex-start;
  padding: var(--space-md);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-primary);
  background: rgba(255, 255, 255, 0.28);
}

.snapshot-item__main {
  display: grid;
  gap: 0.3rem;
}

.snapshot-item__row {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.snapshot-item__row strong {
  font-size: 1rem;
}

.snapshot-status,
.snapshot-item__tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 1.75rem;
  padding: 0 0.6rem;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.snapshot-status--published,
.snapshot-item__tag {
  background: rgba(5, 150, 105, 0.12);
  color: #047857;
}

.snapshot-status--draft {
  background: rgba(217, 119, 6, 0.12);
  color: #b45309;
}

.snapshot-status--failed {
  background: rgba(220, 38, 38, 0.12);
  color: #b91c1c;
}

.snapshot-item__error {
  color: #b91c1c;
}

.snapshot-item__actions {
  flex-shrink: 0;
}

@media (max-width: 1100px) {
  .metrics-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dashboard-panels {
    grid-template-columns: 1fr;
  }

  .sync-overview {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .metrics-grid,
  .today-grid {
    grid-template-columns: 1fr;
  }

  .snapshot-item,
  .sync-card__header {
    flex-direction: column;
  }
}
</style>
