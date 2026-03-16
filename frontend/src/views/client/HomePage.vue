<template>
  <div class="home-page">
    <section class="cover section-shell">
      <div class="app-frame">
        <div class="cover-grid">
          <div class="cover-copy fade-in-up visible">
            <span class="section-kicker">Volume 02 / Personal Publishing System</span>
            <h1 class="cover-title display-title">
              {{ profile?.name || 'MDXY' }}
              <span class="cover-title__support">writes, ships, and curates.</span>
            </h1>
            <p class="cover-lede">{{ mainBio || 'A living archive of engineering notes, tools, and evolving ideas.' }}</p>
            <p v-if="subBio" class="cover-support">{{ subBio }}</p>
            <div v-if="identityTags.length" class="cover-tags">
              <span v-for="tag in identityTags" :key="tag" class="cover-tag">{{ tag }}</span>
            </div>
            <div class="cover-actions">
              <router-link to="/notes" class="cta cta--primary">
                <BookOutline class="cta__icon" />
                在线笔记
              </router-link>
              <a v-if="profile?.github" :href="profile.github" target="_blank" rel="noopener" class="cta cta--secondary">
                <LogoGithub class="cta__icon" />
                GitHub
              </a>
            </div>
          </div>

          <aside class="issue-panel surface-panel fade-in-up visible">
            <div class="issue-panel__header">
              <span class="meta-label">Profile</span>
              <p>简历之外，展示可持续的工程判断与执行力。</p>
            </div>

            <dl class="issue-metrics">
              <div v-for="item in issueMetrics" :key="item.label" class="issue-metric">
                <dt>{{ item.label }}</dt>
                <dd>{{ item.value }}</dd>
              </div>
            </dl>

            <div class="issue-rail">
              <section class="issue-rail__section">
                <div class="issue-rail__head">
                  <span class="meta-label">Featured notes</span>
                  <router-link to="/notes" class="issue-rail__link">查看全部</router-link>
                </div>

                <div class="issue-rail__list">
                  <router-link
                    v-for="note in featuredNotes"
                    :key="note.path"
                    :to="`/notes/${note.path}`"
                    class="issue-rail__item"
                  >
                    <span class="issue-rail__body">
                      <span class="issue-rail__title">{{ formatNoteName(note.name) }}</span>
                      <small class="issue-rail__meta">{{ formatNoteContext(note.path) }}</small>
                    </span>
                    <span class="issue-rail__arrow" aria-hidden="true">→</span>
                  </router-link>
                  <p v-if="!featuredNotes.length" class="issue-empty">更多内容将持续更新于此。</p>
                </div>
              </section>

              <section class="issue-rail__section">
                <div class="issue-rail__head">
                  <span class="meta-label">Field kit</span>
                  <router-link to="/tools" class="issue-rail__link">进入工具箱</router-link>
                </div>

                <div class="issue-rail__list">
                  <template v-for="tool in featuredTools" :key="tool.id">
                    <router-link v-if="isInternalToolUrl(tool.url)" :to="tool.url || '/tools'" class="issue-rail__item">
                      <span class="issue-rail__body">
                        <span class="issue-rail__title">{{ tool.name }}</span>
                        <small class="issue-rail__meta">{{ tool.category || normalizeToolType(tool.type) || 'Tool entry' }}</small>
                      </span>
                      <span class="issue-rail__arrow" aria-hidden="true">→</span>
                    </router-link>
                    <a v-else :href="tool.url || '/tools'" :target="tool.url ? '_blank' : undefined" :rel="tool.url ? 'noopener' : undefined" class="issue-rail__item">
                      <span class="issue-rail__body">
                        <span class="issue-rail__title">{{ tool.name }}</span>
                        <small class="issue-rail__meta">{{ tool.category || normalizeToolType(tool.type) || 'Tool entry' }}</small>
                      </span>
                      <span class="issue-rail__arrow" aria-hidden="true">→</span>
                    </a>
                  </template>
                  <router-link v-if="!featuredTools.length" to="/tools" class="issue-rail__item">
                    <span class="issue-rail__body">
                      <span class="issue-rail__title">浏览工具箱</span>
                      <small class="issue-rail__meta">查看收藏与实验入口</small>
                    </span>
                    <span class="issue-rail__arrow" aria-hidden="true">→</span>
                  </router-link>
                </div>
              </section>
            </div>
          </aside>
        </div>
      </div>
    </section>

    <section v-if="skillsData?.categories?.length" class="section-shell section-divider">
      <div class="app-frame">
        <div class="capability-shell">
          <div class="capability-copy section-intro">
            <span class="section-kicker">Capability Index</span>
            <h2 class="section-title">能力版图</h2>
            <p class="section-description">能力的价值，不在术语堆叠，而在复杂问题中的稳定交付。</p>
            <div class="capability-meta">
              <span><strong>{{ allCapabilityCategories.length }}</strong> 条能力线索</span>
              <span><strong>{{ capabilityKeywordCount }}</strong> 个核心关键词</span>
            </div>
            <button v-if="hiddenCapabilityCount > 0" type="button" class="capability-toggle" @click="showAllCapabilities = !showAllCapabilities">
              {{ showAllCapabilities ? '收起扩展内容' : `查看其余 ${hiddenCapabilityCount} 项` }}
            </button>
          </div>

          <div class="capability-grid">
            <article v-for="cat in displayedCapabilityCategories" :key="cat.name" class="surface-panel capability-card">
              <component :is="getStatIcon(cat.icon)" class="capability-card__icon" />
              <h3>{{ cat.name }}</h3>
              <p>{{ cat.items.join(' / ') }}</p>
            </article>
          </div>
        </div>
      </div>
    </section>

    <section v-if="skillsData?.experience?.length || skillsData?.education?.length" class="section-shell section-divider">
      <div class="app-frame trajectory-grid">
        <div class="section-intro trajectory-grid__heading">
          <span class="section-kicker">Trajectory</span>
          <h2 class="section-title">经历是轨迹，项目是注脚。</h2>
          <p class="section-description">按时间归档项目实践与学术训练，呈现能力演进的关键节点。</p>
        </div>

        <div class="trajectory-column" v-if="skillsData?.experience?.length">
          <span class="meta-label">Experience</span>
          <article v-for="(exp, index) in skillsData.experience" :key="`${exp.company}-${index}`" class="surface-panel trajectory-item">
            <div class="trajectory-item__meta">
              <span class="trajectory-item__kind">Experience</span>
              <span class="trajectory-item__period">{{ exp.period }}</span>
            </div>
            <h3 class="trajectory-item__title">{{ exp.company }}</h3>
            <p class="trajectory-item__subtitle">{{ exp.role }}</p>
            <p class="trajectory-item__detail">{{ exp.detail }}</p>
          </article>
        </div>

        <div class="trajectory-column" v-if="skillsData?.education?.length">
          <span class="meta-label">Education</span>
          <article v-for="(edu, index) in skillsData.education" :key="`${edu.school}-${index}`" class="surface-panel trajectory-item">
            <div class="trajectory-item__meta">
              <span class="trajectory-item__kind">Education</span>
              <span class="trajectory-item__period">{{ edu.period }}</span>
            </div>
            <h3 class="trajectory-item__title">{{ edu.degree }}</h3>
            <p class="trajectory-item__subtitle">{{ edu.school }}</p>
            <p class="trajectory-item__detail">{{ edu.detail }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="section-shell section-divider">
      <div class="app-frame closing-grid">
        <div class="section-intro">
          <span class="section-kicker">Contact</span>
          <h2 class="section-title">保持联系</h2>
          <p class="section-description">
            围绕技术、项目与实践，欢迎继续交流。也欢迎讨论具体笔记背后的实现路径。
          </p>
          <div v-if="skillsData?.hobbies?.length" class="hobby-strip">
            <span v-for="hobby in skillsData.hobbies" :key="hobby">{{ hobby }}</span>
          </div>
        </div>

        <div class="closing-panels">
          <a v-if="profile?.email" :href="`mailto:${profile.email}`" class="surface-panel contact-card">
            <MailOutline class="contact-card__icon" />
            <div>
              <span class="meta-label">Mail</span>
              <strong>{{ profile.email }}</strong>
            </div>
          </a>

          <a v-if="profile?.github" :href="profile.github" target="_blank" rel="noopener" class="surface-panel contact-card">
            <LogoGithub class="contact-card__icon" />
            <div>
              <span class="meta-label">GitHub</span>
              <strong>@{{ profile.github.split('/').pop() }}</strong>
            </div>
          </a>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  LogoGithub,
  MailOutline,
  BookOutline,
  SchoolOutline,
  BriefcaseOutline,
  CodeSlashOutline,
  ServerOutline,
  CloudOutline,
  StorefrontOutline,
  HeartOutline
} from '@vicons/ionicons5'
import { profileApi } from '@/api/profile'
import { noteApi } from '@/api/note'
import { toolApi } from '@/api/tool'
import type { Profile, NoteNode, Tool } from '@/api/types'
import { isResolvableToolRoute } from '@/utils/toolNavigation'

interface SkillsData {
  stats?: { label: string; value: string; icon: string; color: string }[]
  categories?: { name: string; icon: string; items: string[] }[]
  education?: { degree: string; school: string; period: string; detail: string }[]
  experience?: { company: string; role: string; period: string; detail: string }[]
  hobbies?: string[]
}

const router = useRouter()
const profile = ref<Profile | null>(null)
const featuredNotes = ref<NoteNode[]>([])
const featuredTools = ref<Tool[]>([])
const noteCount = ref(0)
const toolCount = ref(0)
const showAllCapabilities = ref(false)

const identityTags = computed(() => {
  if (!profile.value?.title) return []
  return profile.value.title
    .split(/[,，]/)
    .map((s) => s.trim())
    .filter(Boolean)
})

const skillsData = computed<SkillsData | null>(() => {
  if (!profile.value?.skills) return null
  try {
    return JSON.parse(profile.value.skills)
  } catch {
    return null
  }
})

const mainBio = computed(() => {
  if (!profile.value?.bio) return ''
  const lines = profile.value.bio.split('\n').filter(Boolean)
  return lines[0] || ''
})

const subBio = computed(() => {
  if (!profile.value?.bio) return ''
  const lines = profile.value.bio.split('\n').filter(Boolean)
  return lines.length > 1 ? lines.slice(1).join(' ') : ''
})

const issueMetrics = computed(() => [
  { label: 'Notes', value: String(noteCount.value).padStart(2, '0') },
  { label: 'Tools', value: String(toolCount.value).padStart(2, '0') },
  { label: 'Tracks', value: String(skillsData.value?.categories?.length || 0).padStart(2, '0') }
])

const allCapabilityCategories = computed(() => skillsData.value?.categories || [])
const capabilityKeywordCount = computed(() =>
  allCapabilityCategories.value.reduce((count, cat) => count + cat.items.length, 0)
)
const displayedCapabilityCategories = computed(() => {
  if (showAllCapabilities.value) return allCapabilityCategories.value
  return allCapabilityCategories.value.slice(0, 4)
})
const hiddenCapabilityCount = computed(() => Math.max(0, allCapabilityCategories.value.length - 4))

const iconMap: Record<string, any> = {
  school: SchoolOutline,
  briefcase: BriefcaseOutline,
  code: CodeSlashOutline,
  heart: HeartOutline,
  server: ServerOutline,
  cloud: CloudOutline,
  database: StorefrontOutline
}

function getStatIcon(name: string) {
  return iconMap[name] || CodeSlashOutline
}

function flattenNotes(nodes: NoteNode[]): NoteNode[] {
  return nodes.flatMap((node) => (node.type === 'file' ? [node] : flattenNotes(node.children || [])))
}

function formatNoteName(name: string) {
  return name.replace(/\.md$/i, '')
}

function formatNotePath(path: string) {
  return path.replace(/\.md$/i, '')
}

function formatNoteContext(path: string) {
  const normalized = formatNotePath(path)
  const parts = normalized.split('/').filter(Boolean)
  return parts.length > 1 ? parts.slice(0, -1).join(' / ') : 'Markdown note'
}

function isInternalToolUrl(url?: string) {
  return isResolvableToolRoute(router, url)
}

function normalizeToolType(type?: string) {
  if (!type) return ''
  return type === 'software' ? 'app' : type
}

onMounted(async () => {
  try {
    const [profileRes, notesRes, toolsRes] = await Promise.allSettled([
      profileApi.get(),
      noteApi.getTree(),
      toolApi.list()
    ])

    if (profileRes.status === 'fulfilled' && profileRes.value.success && profileRes.value.data) {
      profile.value = profileRes.value.data
    }

    if (notesRes.status === 'fulfilled' && notesRes.value.success && notesRes.value.data) {
      const flattened = flattenNotes(notesRes.value.data)
      noteCount.value = flattened.length
      featuredNotes.value = flattened.slice(0, 2)
    }

    if (toolsRes.status === 'fulfilled' && toolsRes.value.success && toolsRes.value.data) {
      const visibleTools = (toolsRes.value.data || []).filter((tool) => tool.is_visible)
      toolCount.value = visibleTools.length
      featuredTools.value = visibleTools.slice(0, 2)
    }
  } catch (error) {
    console.error('初始化首页内容失败:', error)
  }
})
</script>

<style scoped>
.home-page {
  position: relative;
  z-index: 1;
}

.cover {
  padding-top: calc(var(--space-3xl) + 1rem);
}

.cover-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(340px, 0.9fr);
  gap: var(--space-2xl);
  align-items: start;
}

.cover-copy {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
}

.cover-title {
  font-size: clamp(3.8rem, 11vw, 7rem);
  max-width: 9ch;
  font-weight: 600;
}

.cover-title__support {
  display: block;
  margin-top: var(--space-md);
  font-family: var(--font-body);
  font-size: clamp(1.25rem, 2vw, 1.6rem);
  font-weight: 500;
  letter-spacing: 0;
  line-height: 1.35;
  color: var(--text-secondary);
}

.cover-lede {
  max-width: 42rem;
  font-size: 1.18rem;
  line-height: 1.85;
  color: var(--text-secondary);
  text-wrap: pretty;
}

.cover-support {
  max-width: 40rem;
  color: var(--text-muted);
  line-height: 1.85;
  text-wrap: pretty;
}

.cover-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-sm);
}

.cover-tag {
  padding: 0.6rem 1rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  background: rgba(255, 255, 255, 0.36);
  color: var(--text-secondary);
  font-size: 0.92rem;
}

.cover-actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-md);
  padding-top: var(--space-sm);
}

.cta {
  display: inline-flex;
  align-items: center;
  gap: 0.75rem;
  min-height: 3.25rem;
  padding: 0 1.15rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  font-weight: 600;
  letter-spacing: -0.01em;
}

.cta__icon {
  width: 1rem;
  height: 1rem;
}

.cta--primary {
  background: var(--bg-contrast);
  color: var(--text-inverse);
}

.cta--primary:hover {
  color: var(--text-inverse);
  opacity: 0.92;
}

.cta--secondary {
  background: var(--bg-panel);
  color: var(--text-primary);
}

.issue-panel {
  margin-top: -0.95rem;
  padding: clamp(1.3rem, 2vw, 1.8rem);
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.issue-panel__header {
  display: grid;
  gap: 0.35rem;
}

.issue-panel__header p {
  color: var(--text-secondary);
  line-height: 1.58;
}

.issue-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.72rem;
}

.issue-metric {
  padding-top: 0.62rem;
  border-top: 1px solid var(--border-primary);
}

.issue-metric dt {
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.issue-metric dd {
  margin-top: 0.28rem;
  font-family: var(--font-display);
  font-size: clamp(1.58rem, 2.7vw, 1.9rem);
  letter-spacing: -0.04em;
}

.issue-rail {
  display: grid;
  gap: 0.85rem;
}

.issue-rail__section {
  display: grid;
  gap: 0.75rem;
  padding-top: 0.78rem;
  border-top: 1px solid var(--border-primary);
}

.issue-rail__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-sm);
}

.issue-rail__link {
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.issue-rail__list {
  display: grid;
  gap: 0.58rem;
}

.issue-rail__item {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  gap: 0.8rem;
  min-height: 3.85rem;
  padding: 0.66rem 0.9rem;
  border-radius: calc(var(--radius-lg) - 0.2rem);
  border: 1px solid rgba(62, 53, 39, 0.12);
  background: rgba(255, 255, 255, 0.2);
  color: var(--text-primary);
  transition:
    transform var(--duration-fast) var(--ease-standard),
    border-color var(--duration-fast) var(--ease-standard),
    background-color var(--duration-fast) var(--ease-standard);
}

.issue-rail__item:hover {
  transform: translateX(2px);
  border-color: rgba(62, 53, 39, 0.18);
  background: rgba(255, 255, 255, 0.34);
}

.issue-rail__body {
  min-width: 0;
  display: grid;
  gap: 0.12rem;
}

.issue-rail__title {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 1rem;
  font-weight: 600;
  line-height: 1.25;
}

.issue-rail__meta,
.issue-empty {
  color: var(--text-muted);
}

.issue-rail__meta {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.78rem;
  letter-spacing: 0.04em;
}

.issue-rail__arrow {
  color: var(--text-muted);
  font-size: 1rem;
  line-height: 1;
}

.issue-empty {
  padding: 0.45rem 0;
}

.section-intro {
  display: grid;
  gap: var(--space-md);
}

.section-intro .section-title {
  max-width: 12ch;
}

.section-intro .section-description {
  max-width: 38ch;
}

.capability-shell {
  display: grid;
  grid-template-columns: minmax(0, 0.95fr) minmax(0, 1.25fr);
  gap: var(--space-2xl);
  align-items: start;
}

.capability-copy {
  display: grid;
  gap: var(--space-md);
}

.capability-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-md);
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.capability-meta strong {
  margin-right: 0.25rem;
  font-family: var(--font-display);
  font-size: 1.2rem;
  letter-spacing: -0.03em;
  color: var(--text-primary);
}

.capability-toggle {
  justify-self: start;
  margin-top: var(--space-sm);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  padding: 0.58rem 1rem;
  color: var(--text-secondary);
  background: var(--bg-panel);
}

.capability-toggle:hover {
  color: var(--text-primary);
  border-color: var(--border-strong);
}

.capability-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--space-md);
  align-content: start;
}

.capability-card {
  min-height: 12.5rem;
  padding: 1.35rem;
}

.capability-card__icon {
  width: 1.7rem;
  height: 1.7rem;
  color: var(--accent-primary);
}

.capability-card h3 {
  margin-top: var(--space-lg);
  font-family: var(--font-display);
  font-size: 1.5rem;
  letter-spacing: -0.03em;
}

.capability-card p {
  margin-top: var(--space-sm);
  color: var(--text-secondary);
  line-height: 1.8;
}

.trajectory-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: var(--space-xl);
  align-items: start;
}

.trajectory-grid__heading {
  margin-bottom: 0;
}

.trajectory-column {
  display: grid;
  gap: var(--space-md);
  align-content: start;
}

.trajectory-item {
  min-height: 14.5rem;
  padding: 1.35rem 1.5rem;
  display: grid;
  gap: var(--space-sm);
}

.trajectory-item__meta {
  display: flex;
  justify-content: space-between;
  gap: var(--space-sm);
  align-items: center;
}

.trajectory-item__kind {
  display: inline-flex;
  align-items: center;
  padding: 0.2rem 0.56rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 0.64rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.trajectory-item__period {
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.trajectory-item__title {
  margin-top: 0.15rem;
  font-family: var(--font-display);
  font-size: clamp(1.62rem, 1.95vw, 1.9rem);
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.16;
}

.trajectory-item__subtitle {
  font-size: 1.06rem;
  font-weight: 600;
  color: var(--text-primary);
}

.trajectory-item__detail {
  color: var(--text-secondary);
  line-height: 1.8;
}

.closing-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(320px, 0.85fr);
  gap: var(--space-2xl);
  align-items: start;
}

.hobby-strip {
  margin-top: var(--space-xl);
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-sm);
}

.hobby-strip span {
  padding: 0.55rem 0.9rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-secondary);
}

.closing-panels {
  display: grid;
  gap: var(--space-md);
}

.contact-card {
  display: flex;
  align-items: center;
  gap: var(--space-lg);
  padding: var(--space-lg);
  color: var(--text-primary);
}

.contact-card__icon {
  width: 1.5rem;
  height: 1.5rem;
  color: var(--accent-primary);
}

.contact-card strong {
  display: block;
  margin-top: 0.45rem;
  font-family: var(--font-display);
  font-size: 1.4rem;
  letter-spacing: -0.03em;
}

@media (max-width: 1100px) {
  .cover-grid,
  .capability-shell,
  .closing-grid,
  .trajectory-grid {
    grid-template-columns: 1fr;
  }

  .capability-grid {
    grid-template-columns: 1fr;
  }

  .issue-panel {
    margin-top: 0;
  }
}

@media (max-width: 768px) {
  .issue-metrics {
    grid-template-columns: 1fr;
  }

  .cover-actions,
  .closing-panels {
    grid-template-columns: 1fr;
  }

  .cta {
    width: 100%;
    justify-content: center;
  }

  .trajectory-item {
    min-height: 0;
  }

  .trajectory-item__meta {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
