<template>
  <div class="home-page">
    <div class="container">
      <div class="hero">
        <div v-if="profile" class="profile-section">
          <div v-if="profile.avatar" class="avatar">
            <img :src="profile.avatar" :alt="profile.name" />
          </div>
          <h1 class="name">{{ profile.name }}</h1>
          <p class="title">{{ profile.title }}</p>
          <p class="bio">{{ profile.bio }}</p>

          <div v-if="hasLinks" class="links">
            <a v-if="profile.github" :href="profile.github" target="_blank" rel="noopener">
              <n-button text>GitHub</n-button>
            </a>
            <a v-if="profile.email" :href="`mailto:${profile.email}`">
              <n-button text>Email</n-button>
            </a>
            <a v-if="profile.website" :href="profile.website" target="_blank" rel="noopener">
              <n-button text>Website</n-button>
            </a>
          </div>
        </div>

        <n-spin v-else :show="loading">
          <div style="height: 300px"></div>
        </n-spin>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { NButton, NSpin } from 'naive-ui'
import { profileApi } from '@/api/profile'
import type { Profile } from '@/api/types'

const profile = ref<Profile | null>(null)
const loading = ref(true)

const hasLinks = computed(() => {
  return profile.value && (profile.value.github || profile.value.email || profile.value.website)
})

const fetchProfile = async () => {
  try {
    const response = await profileApi.get()
    if (response.success && response.data) {
      profile.value = response.data
    }
  } catch (error) {
    console.error('获取个人信息失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>
.home-page {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero {
  text-align: center;
  padding: 60px 20px;
}

.profile-section {
  max-width: 600px;
  margin: 0 auto;
}

.avatar {
  width: 120px;
  height: 120px;
  margin: 0 auto 24px;
  border-radius: 50%;
  overflow: hidden;
  border: 4px solid var(--color-border);
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.name {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 12px;
  color: var(--color-text-primary);
}

.title {
  font-size: 24px;
  color: var(--color-text-secondary);
  margin-bottom: 24px;
}

.bio {
  font-size: 16px;
  line-height: 1.8;
  color: var(--color-text-secondary);
  margin-bottom: 32px;
}

.links {
  display: flex;
  gap: 16px;
  justify-content: center;
}
</style>
