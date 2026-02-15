<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-content fade-in-up visible">
        <div class="hero-tags">
          <span v-for="tag in identityTags" :key="tag" class="hero-tag glass-card">{{ tag }}</span>
        </div>
        <h1 class="hero-title">
          Hi，我是 <span class="gradient-text">{{ profile?.name || '...' }}</span>
        </h1>
        <p class="hero-subtitle">{{ profile?.bio || '' }}</p>
        <div class="hero-actions">
          <router-link to="/notes" class="btn-primary">在线笔记</router-link>
          <a v-if="profile?.github" :href="profile.github" target="_blank" rel="noopener" class="btn-outline">
            <n-icon :component="LogoGithub" :size="18" />
            GitHub
          </a>
        </div>
      </div>
    </section>

    <!-- About Section -->
    <section class="section fade-in-up" ref="aboutRef">
      <h2 class="section-title gradient-text">关于我</h2>
      <div class="about-grid">
        <div class="glass-card about-card">
          <p class="about-text">{{ profile?.bio || '加载中...' }}</p>
          <div class="stats-grid" v-if="profile">
            <div class="stat-item" v-if="profile.title">
              <span class="stat-value">{{ profile.title }}</span>
              <span class="stat-label">职位</span>
            </div>
            <div class="stat-item" v-if="profile.email">
              <span class="stat-value">{{ profile.email }}</span>
              <span class="stat-label">邮箱</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Skills Section -->
    <section class="section fade-in-up" ref="skillsRef" v-if="skills.length">
      <h2 class="section-title gradient-text">技能</h2>
      <div class="skills-grid">
        <div v-for="skill in skills" :key="skill" class="glass-card skill-card">
          <span class="skill-name">{{ skill }}</span>
        </div>
      </div>
    </section>

    <!-- Contact Section -->
    <section class="section fade-in-up" ref="contactRef">
      <h2 class="section-title gradient-text">联系我</h2>
      <div class="contact-grid">
        <a v-if="profile?.github" :href="profile.github" target="_blank" rel="noopener" class="glass-card contact-card">
          <n-icon :component="LogoGithub" :size="32" class="contact-icon" />
          <span class="contact-label">GitHub</span>
        </a>
        <a v-if="profile?.email" :href="`mailto:${profile.email}`" class="glass-card contact-card">
          <n-icon :component="MailOutline" :size="32" class="contact-icon" />
          <span class="contact-label">Email</span>
        </a>
        <a v-if="profile?.website" :href="profile.website" target="_blank" rel="noopener" class="glass-card contact-card">
          <n-icon :component="GlobeOutline" :size="32" class="contact-icon" />
          <span class="contact-label">Website</span>
        </a>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { NIcon } from 'naive-ui'
import { LogoGithub, MailOutline, GlobeOutline } from '@vicons/ionicons5'
import { profileApi } from '@/api/profile'
import type { Profile } from '@/api/types'

const profile = ref<Profile | null>(null)
const aboutRef = ref<HTMLElement | null>(null)
const skillsRef = ref<HTMLElement | null>(null)
const contactRef = ref<HTMLElement | null>(null)

const identityTags = computed(() => {
  if (!profile.value?.title) return []
  return profile.value.title.split(/[,，/|]/).map(s => s.trim()).filter(Boolean)
})

const skills = computed(() => {
  if (!profile.value?.skills) return []
  return profile.value.skills.split(/[,，]/).map(s => s.trim()).filter(Boolean)
})

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

  ;[aboutRef, skillsRef, contactRef].forEach((r) => {
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

.hero-subtitle {
  font-size: 18px;
  color: var(--color-text-secondary);
  line-height: 1.8;
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
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 212, 170, 0.3);
  color: #fff;
}
/* STYLE_CHUNK_2 */

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
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 32px;
  text-align: center;
}

/* About */
.about-card {
  padding: 32px;
}

.about-text {
  font-size: 16px;
  line-height: 1.8;
  color: var(--color-text-secondary);
  margin-bottom: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-item {
  text-align: center;
  padding: 16px;
  border-radius: var(--radius-md);
  background: var(--glass-bg);
}

.stat-value {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--accent-cyan);
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-tertiary);
}
/* STYLE_CHUNK_3 */

/* Skills */
.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}

.skill-card {
  padding: 16px;
  text-align: center;
  transition: transform 0.2s;
}

.skill-card:hover {
  transform: translateY(-2px);
}

.skill-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

/* Contact */
.contact-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
}

.contact-card {
  padding: 32px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  transition: transform 0.2s;
  cursor: pointer;
  color: var(--color-text-primary);
}

.contact-card:hover {
  transform: translateY(-4px);
  color: var(--color-text-primary);
}

.contact-icon {
  color: var(--accent-cyan);
}

.contact-label {
  font-size: 15px;
  font-weight: 500;
}

/* Responsive */
@media (max-width: 768px) {
  .hero-title {
    font-size: 36px;
  }

  .hero-subtitle {
    font-size: 16px;
  }

  .hero-actions {
    flex-direction: column;
    align-items: center;
  }

  .section-title {
    font-size: 26px;
  }
}
</style>