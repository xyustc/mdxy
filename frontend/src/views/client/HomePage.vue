<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-content fade-in-up visible">
        <div class="hero-tags">
          <span v-for="tag in identityTags" :key="tag" class="hero-tag glass-card">{{ tag }}</span>
        </div>
        <h1 class="hero-title">
          Hi，我是 <a v-if="profile?.github" :href="profile.github" target="_blank" rel="noopener" class="gradient-text name-link">{{ profile?.name || '...' }}</a>
          <span v-else class="gradient-text">{{ profile?.name || '...' }}</span>
        </h1>
        <p class="hero-subtitle">{{ mainBio }}</p>
        <p class="hero-sub-desc" v-if="subBio">{{ subBio }}</p>
        <div class="hero-actions">
          <router-link to="/notes" class="btn-primary">
            <n-icon :component="BookOutline" :size="18" /> 在线笔记
          </router-link>
          <a v-if="profile?.github" :href="profile.github" target="_blank" rel="noopener" class="btn-outline">
            <n-icon :component="LogoGithub" :size="18" /> GitHub
          </a>
        </div>
      </div>
    </section>

    <!-- About + Stats Section -->
    <section class="section fade-in-up" ref="aboutRef">
      <h2 class="section-title gradient-text">关于我</h2>
      <div class="glass-card about-card">
        <div class="about-layout">
          <div class="about-text-area">
            <h3 class="about-greeting">👋 你好！我是 <strong>{{ profile?.name || '...' }}</strong></h3>
            <p class="about-desc">{{ mainBio }}</p>
            <p class="about-motto" v-if="subBio">🚀 {{ subBio }}</p>
          </div>
          <div class="stats-grid" v-if="skillsData?.stats">
            <div v-for="stat in skillsData.stats" :key="stat.label" class="stat-item" :class="`stat-${stat.color}`">
              <n-icon :component="getStatIcon(stat.icon)" :size="28" class="stat-icon" />
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>
    <!-- Education Section -->
    <section class="section fade-in-up" ref="eduRef" v-if="skillsData?.education?.length">
      <h2 class="section-title gradient-text">
        <n-icon :component="SchoolOutline" :size="28" /> 教育背景
      </h2>
      <div class="timeline">
        <div v-for="(edu, i) in skillsData.education" :key="i" class="glass-card timeline-item" :class="i % 2 === 0 ? 'timeline-cyan' : 'timeline-coral'">
          <div class="timeline-header">
            <h3 class="timeline-title">{{ edu.degree }}</h3>
            <span class="timeline-period">{{ edu.period }}</span>
          </div>
          <p class="timeline-school">{{ edu.school }}</p>
          <p class="timeline-detail">{{ edu.detail }}</p>
        </div>
      </div>
    </section>

    <!-- Experience Section -->
    <section class="section fade-in-up" ref="expRef" v-if="skillsData?.experience?.length">
      <h2 class="section-title gradient-text">
        <n-icon :component="BriefcaseOutline" :size="28" /> 实习经历
      </h2>
      <div class="timeline">
        <div v-for="(exp, i) in skillsData.experience" :key="i" class="glass-card timeline-item timeline-orange">
          <div class="timeline-header">
            <h3 class="timeline-title">{{ exp.company }}</h3>
            <span class="timeline-period">{{ exp.period }}</span>
          </div>
          <p class="timeline-role">{{ exp.role }}</p>
          <p class="timeline-detail">{{ exp.detail }}</p>
        </div>
      </div>
    </section>

    <!-- Skills Section -->
    <section class="section fade-in-up" ref="skillsRef" v-if="skillsData?.categories?.length">
      <h2 class="section-title gradient-text">技术栈</h2>
      <div class="skills-grid">
        <div v-for="cat in skillsData.categories" :key="cat.name" class="glass-card skill-card">
          <n-icon :component="getStatIcon(cat.icon)" :size="28" class="skill-icon" />
          <h3 class="skill-cat-name">{{ cat.name }}</h3>
          <p class="skill-items">{{ cat.items.join(', ') }}</p>
        </div>
      </div>
    </section>

    <!-- Hobbies Section -->
    <section class="section fade-in-up" ref="hobbiesRef" v-if="skillsData?.hobbies?.length">
      <h2 class="section-title gradient-text">🎨 兴趣爱好</h2>
      <div class="hobbies-grid">
        <span v-for="hobby in skillsData.hobbies" :key="hobby" class="glass-card hobby-tag">{{ hobby }}</span>
      </div>
    </section>

    <!-- Contact Section -->
    <section class="section fade-in-up" ref="contactRef">
      <h2 class="section-title gradient-text">联系我</h2>
      <div class="glass-card contact-wrapper">
        <p class="contact-desc">💬 有学术交流或技术合作机会？欢迎随时联系！</p>
        <div class="contact-grid">
          <a v-if="profile?.email" :href="`mailto:${profile.email}`" class="contact-card">
            <n-icon :component="MailOutline" :size="32" class="contact-icon contact-icon-cyan" />
            <span class="contact-value">{{ profile.email }}</span>
          </a>
          <a v-if="profile?.github" :href="profile.github" target="_blank" rel="noopener" class="contact-card">
            <n-icon :component="LogoGithub" :size="32" class="contact-icon contact-icon-coral" />
            <span class="contact-value">@{{ profile.github.split('/').pop() }}</span>
          </a>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { NIcon } from 'naive-ui'
import {
  LogoGithub, MailOutline, GlobeOutline, BookOutline,
  SchoolOutline, BriefcaseOutline, CodeSlashOutline,
  ServerOutline, CloudOutline, StorefrontOutline,
  HeartOutline
} from '@vicons/ionicons5'
import { profileApi } from '@/api/profile'
import type { Profile } from '@/api/types'

interface SkillsData {
  stats?: { label: string; value: string; icon: string; color: string }[]
  categories?: { name: string; icon: string; items: string[] }[]
  education?: { degree: string; school: string; period: string; detail: string }[]
  experience?: { company: string; role: string; period: string; detail: string }[]
  hobbies?: string[]
}

const profile = ref<Profile | null>(null)
const aboutRef = ref<HTMLElement | null>(null)
const eduRef = ref<HTMLElement | null>(null)
const expRef = ref<HTMLElement | null>(null)
const skillsRef = ref<HTMLElement | null>(null)
const hobbiesRef = ref<HTMLElement | null>(null)
const contactRef = ref<HTMLElement | null>(null)

const identityTags = computed(() => {
  if (!profile.value?.title) return []
  return profile.value.title.split(/[,，]/).map(s => s.trim()).filter(Boolean)
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

const iconMap: Record<string, any> = {
  school: SchoolOutline,
  briefcase: BriefcaseOutline,
  code: CodeSlashOutline,
  heart: HeartOutline,
  server: ServerOutline,
  cloud: CloudOutline,
  database: StorefrontOutline,
}

function getStatIcon(name: string) {
  return iconMap[name] || CodeSlashOutline
}

let observer: IntersectionObserver | null = null

onMounted(async () => {
  try {
    const res = await profileApi.get()
    if (res.success && res.data) {
      profile.value = res.data
    }
  } catch (e) {
    console.error('获取个人信息失败:', e)
  }

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
        }
      })
    },
    { threshold: 0.1 }
  )

  ;[aboutRef, eduRef, expRef, skillsRef, hobbiesRef, contactRef].forEach((r) => {
    if (r.value) observer!.observe(r.value)
  })
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<style scoped>
.home-page {
  padding-bottom: 80px;
}

/* Hero */
.hero {
  min-height: calc(100vh - 64px);
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 20px;
}

.hero-content {
  max-width: 700px;
}

.hero-tags {
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.hero-tag {
  padding: 6px 16px;
  font-size: 13px;
  color: var(--accent-cyan);
  font-weight: 500;
}

.hero-title {
  font-size: 52px;
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 20px;
  color: var(--color-text-primary);
}

.name-link {
  text-decoration: none;
  transition: text-shadow 0.3s;
}

.name-link:hover {
  text-shadow: 0 0 12px var(--accent-cyan);
}

.hero-subtitle {
  font-size: 18px;
  color: var(--color-text-secondary);
  line-height: 1.8;
  margin-bottom: 12px;
}

.hero-sub-desc {
  font-size: 15px;
  color: var(--color-text-tertiary);
  margin-bottom: 36px;
}

.hero-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.btn-primary {
  padding: 12px 32px;
  background: linear-gradient(135deg, var(--accent-cyan), var(--accent-teal));
  color: #fff;
  border-radius: 30px;
  font-weight: 600;
  font-size: 15px;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 212, 170, 0.3);
  color: #fff;
}

.btn-outline {
  padding: 12px 32px;
  border: 1px solid var(--glass-border);
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  color: var(--color-text-primary);
  border-radius: 30px;
  font-weight: 600;
  font-size: 15px;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-outline:hover {
  transform: translateY(-2px);
  box-shadow: var(--glass-shadow);
  color: var(--color-text-primary);
}

/* Sections */
.section {
  max-width: 900px;
  margin: 0 auto;
  padding: 60px 20px;
}

.section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 32px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

/* About */
.about-card {
  padding: 32px;
}

.about-layout {
  display: flex;
  align-items: flex-start;
  gap: 3rem;
  flex-wrap: wrap;
}

.about-text-area {
  flex: 1;
  min-width: 280px;
}

.about-greeting {
  color: var(--accent-cyan);
  margin-bottom: 1rem;
  font-size: 1.2rem;
}

.about-desc {
  font-size: 16px;
  line-height: 1.8;
  color: var(--color-text-secondary);
  margin-bottom: 1rem;
}

.about-motto {
  color: var(--color-text-tertiary);
  line-height: 1.5;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  flex: 0 0 auto;
  min-width: 260px;
}

.stat-item {
  text-align: center;
  padding: 1.2rem 0.8rem;
  border-radius: var(--radius-md);
  border: 1px solid transparent;
}

.stat-cyan { background: rgba(0, 229, 191, 0.08); border-color: rgba(0, 229, 191, 0.2); }
.stat-coral { background: rgba(255, 123, 123, 0.08); border-color: rgba(255, 123, 123, 0.2); }
.stat-teal { background: rgba(56, 217, 169, 0.08); border-color: rgba(56, 217, 169, 0.2); }
.stat-orange { background: rgba(255, 169, 77, 0.08); border-color: rgba(255, 169, 77, 0.2); }

.stat-cyan .stat-icon, .stat-cyan .stat-value { color: var(--accent-cyan); }
.stat-coral .stat-icon, .stat-coral .stat-value { color: var(--accent-coral); }
.stat-teal .stat-icon, .stat-teal .stat-value { color: var(--accent-teal); }
.stat-orange .stat-icon, .stat-orange .stat-value { color: var(--accent-orange); }

.stat-icon {
  display: block;
  margin: 0 auto 0.5rem;
}

.stat-value {
  font-size: 1.1rem;
  font-weight: 700;
}

.stat-label {
  font-size: 0.85rem;
  color: var(--color-text-tertiary);
  margin-top: 2px;
}

/* Timeline (Education & Experience) */
.timeline {
  display: grid;
  gap: 1.2rem;
}

.timeline-item {
  padding: 1.5rem;
  border-left: 4px solid var(--accent-cyan);
}

.timeline-cyan { border-left-color: var(--accent-cyan); }
.timeline-coral { border-left-color: var(--accent-coral); }
.timeline-orange { border-left-color: var(--accent-orange); }

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 0.5rem;
}

.timeline-title {
  font-size: 1rem;
  font-weight: 700;
  margin: 0;
}

.timeline-cyan .timeline-title { color: var(--accent-cyan); }
.timeline-coral .timeline-title { color: var(--accent-coral); }
.timeline-orange .timeline-title { color: var(--accent-orange); }

.timeline-period {
  font-size: 0.85rem;
  color: var(--color-text-tertiary);
}

.timeline-school, .timeline-role {
  color: var(--accent-teal);
  margin: 0.3rem 0;
  font-size: 0.95rem;
}

.timeline-detail {
  color: var(--color-text-secondary);
  margin: 0;
  font-size: 0.9rem;
}

/* Skills */
.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.skill-card {
  padding: 24px 16px;
  text-align: center;
  transition: transform 0.2s;
}

.skill-card:hover {
  transform: translateY(-3px);
}

.skill-icon {
  color: var(--accent-cyan);
  margin-bottom: 8px;
}

.skill-cat-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 8px;
}

.skill-items {
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin: 0;
}

/* Hobbies */
.hobbies-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1rem;
}

.hobby-tag {
  padding: 0.6rem 1.2rem;
  font-size: 0.95rem;
  color: var(--accent-cyan);
  font-weight: 500;
}

/* Contact */
.contact-wrapper {
  padding: 32px;
  text-align: center;
  max-width: 600px;
  margin: 0 auto;
}

.contact-desc {
  font-size: 1.05rem;
  color: var(--color-text-secondary);
  margin-bottom: 2rem;
}

.contact-grid {
  display: flex;
  justify-content: center;
  gap: 3rem;
  flex-wrap: wrap;
}

.contact-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  transition: transform 0.2s;
}

.contact-card:hover {
  transform: translateY(-3px);
}

.contact-icon-cyan { color: var(--accent-cyan); }
.contact-icon-coral { color: var(--accent-coral); }

.contact-value {
  font-weight: 500;
  font-size: 0.95rem;
  color: var(--color-text-primary);
}

/* Responsive */
@media (max-width: 768px) {
  .hero-title { font-size: 36px; }
  .hero-subtitle { font-size: 16px; }
  .hero-actions { flex-direction: column; align-items: center; }
  .section-title { font-size: 24px; }
  .about-layout { flex-direction: column; }
  .stats-grid { min-width: unset; width: 100%; }
  .skills-grid { grid-template-columns: repeat(2, 1fr); }
}
</style>
