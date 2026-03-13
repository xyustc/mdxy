<template>
  <div class="login-page">
    <div class="login-page__ornament login-page__ornament--left"></div>
    <div class="login-page__ornament login-page__ornament--right"></div>

    <div class="login-page__shell">
      <section class="login-intro">
        <span class="section-kicker">Editor Access</span>
        <h1>进入内容编辑台</h1>
        <p>
          在这里维护个人资料、整理工具箱、查看站点数据，并把公开页面持续打磨成一份有节奏感的个人出版系统。
        </p>
        <router-link to="/" class="login-intro__link">返回前台</router-link>
      </section>

      <section class="paper-sheet login-panel">
        <div class="login-panel__head">
          <span class="meta-label">Admin Login</span>
          <h2>管理后台登录</h2>
        </div>

        <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin">
          <el-form-item prop="username">
            <label class="sr-only" for="login-username">用户名</label>
            <el-input
              id="login-username"
              v-model="form.username"
              name="username"
              autocomplete="username"
              aria-label="用户名"
              spellcheck="false"
              placeholder="用户名"
              size="large"
              :prefix-icon="ElIconUser"
            />
          </el-form-item>

          <el-form-item prop="password">
            <label class="sr-only" for="login-password">密码</label>
            <el-input
              id="login-password"
              v-model="form.password"
              type="password"
              name="password"
              autocomplete="current-password"
              aria-label="密码"
              placeholder="密码"
              size="large"
              :prefix-icon="ElIconLock"
              @keyup.enter="handleLogin"
            />
          </el-form-item>

          <el-form-item class="login-panel__submit">
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleLogin"
              style="width: 100%"
            >
              登录工作台
            </el-button>
          </el-form-item>
        </el-form>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus'
import { User as ElIconUser, Lock as ElIconLock } from '@element-plus/icons-vue'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminStore()

const formRef = ref()
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  try {
    await formRef.value.validate()
    loading.value = true

    const success = await adminStore.login(form)

    if (success) {
      ElMessage.success('登录成功')
      const redirect = (route.query.redirect as string) || '/admin/dashboard'
      router.push(redirect)
    } else {
      ElMessage.error('用户名或密码错误')
    }
  } catch (error: any) {
    if (error.response) {
      ElMessage.error(error.response.data?.error || '登录失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 32px;
}

.login-page__ornament {
  position: absolute;
  width: 28rem;
  height: 28rem;
  border-radius: 50%;
  filter: blur(72px);
  opacity: 0.18;
  pointer-events: none;
}

.login-page__ornament--left {
  left: -10rem;
  top: -10rem;
  background: rgba(139, 94, 60, 0.4);
}

.login-page__ornament--right {
  right: -10rem;
  bottom: -10rem;
  background: rgba(41, 70, 58, 0.42);
}

.login-page__shell {
  position: relative;
  z-index: 1;
  width: min(1080px, 100%);
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(360px, 420px);
  gap: var(--space-2xl);
  align-items: center;
}

.login-intro {
  padding: clamp(1rem, 3vw, 2rem);
}

.login-intro h1 {
  margin-top: var(--space-md);
  font-family: var(--font-display);
  font-size: clamp(3rem, 7vw, 5rem);
  line-height: 0.96;
  letter-spacing: -0.05em;
}

.login-intro p {
  max-width: 34rem;
  margin-top: var(--space-lg);
  color: var(--text-secondary);
  font-size: 1.05rem;
  line-height: 1.85;
}

.login-intro__link {
  display: inline-flex;
  margin-top: var(--space-xl);
  font-weight: 600;
}

.login-panel {
  padding: clamp(1.35rem, 3vw, 2.2rem);
}

.login-panel__head h2 {
  margin-top: var(--space-sm);
  font-family: var(--font-display);
  font-size: 2rem;
  letter-spacing: -0.04em;
}

.login-panel :deep(.el-form) {
  margin-top: var(--space-xl);
}

.login-panel :deep(.el-form-item) {
  margin-bottom: 1rem;
}

.login-panel :deep(.el-input__wrapper) {
  min-height: 3.25rem;
  border-radius: var(--radius-lg);
  background: var(--bg-panel);
  box-shadow: inset 0 0 0 1px var(--border-primary);
}

.login-panel :deep(.el-input__wrapper.is-focus) {
  box-shadow: inset 0 0 0 1px var(--accent-primary);
}

.login-panel__submit {
  padding-top: var(--space-sm);
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

@media (max-width: 900px) {
  .login-page__shell {
    grid-template-columns: 1fr;
  }

  .login-intro {
    padding: 0;
  }
}

@media (max-width: 640px) {
  .login-page {
    padding: 18px;
  }
}
</style>
