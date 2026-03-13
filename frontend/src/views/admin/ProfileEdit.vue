<template>
  <div class="profile-edit-page">
    <el-card class="edit-card" shadow="never">
      <template #header>
        <div class="edit-card__header">
          <span>基本信息</span>
        </div>
      </template>

      <el-form ref="formRef" :model="form" label-width="100px" :rules="rules" class="profile-form">
        <el-row :gutter="16">
          <el-col :md="12" :sm="24">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="form.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :md="12" :sm="24">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="身份标签" prop="title">
          <el-input v-model="form.title" placeholder="多个标签用逗号分隔，如：🎓 硕士在读,💻 技术爱好者" />
        </el-form-item>

        <el-form-item label="个人简介" prop="bio">
          <el-input v-model="form.bio" type="textarea" :rows="4" placeholder="支持换行，第一行为主简介" />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :md="12" :sm="24">
            <el-form-item label="GitHub">
              <el-input v-model="form.github" placeholder="https://github.com/xxx" />
            </el-form-item>
          </el-col>
          <el-col :md="12" :sm="24">
            <el-form-item label="个人网站">
              <el-input v-model="form.website" placeholder="https://…" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :md="12" :sm="24">
            <el-form-item label="LinkedIn">
              <el-input v-model="form.linkedin" placeholder="LinkedIn 链接" />
            </el-form-item>
          </el-col>
          <el-col :md="12" :sm="24">
            <el-form-item label="Twitter">
              <el-input v-model="form.twitter" placeholder="Twitter 链接" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="头像 URL">
          <el-input v-model="form.avatar" placeholder="头像图片 URL" />
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="edit-card" shadow="never">
      <template #header>
        <div class="edit-card__header">
          <span>统计卡片</span>
          <el-button type="primary" size="small" @click="addStat">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(stat, i) in skillsForm.stats" :key="i" class="list-row">
        <el-input v-model="stat.label" placeholder="标签" class="field-120" />
        <el-input v-model="stat.value" placeholder="值" class="field-100" />
        <el-select v-model="stat.icon" placeholder="图标" class="field-120">
          <el-option v-for="ic in iconOptions" :key="ic" :label="ic" :value="ic" />
        </el-select>
        <el-select v-model="stat.color" placeholder="颜色" class="field-110">
          <el-option v-for="c in colorOptions" :key="c" :label="c" :value="c" />
        </el-select>
        <el-button type="danger" size="small" @click="skillsForm.stats.splice(i, 1)">删除</el-button>
      </div>
      <p v-if="!skillsForm.stats.length" class="empty-tip">暂未添加统计卡片。</p>
    </el-card>

    <el-card class="edit-card" shadow="never">
      <template #header>
        <div class="edit-card__header">
          <span>技能分类</span>
          <el-button type="primary" size="small" @click="addCategory">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(cat, i) in skillsForm.categories" :key="i" class="list-row">
        <el-input v-model="cat.name" placeholder="分类名" class="field-120" />
        <el-select v-model="cat.icon" placeholder="图标" class="field-120">
          <el-option v-for="ic in iconOptions" :key="ic" :label="ic" :value="ic" />
        </el-select>
        <el-input v-model="cat.itemsStr" placeholder="技能项，逗号分隔" class="field-grow" />
        <el-button type="danger" size="small" @click="skillsForm.categories.splice(i, 1)">删除</el-button>
      </div>
      <p v-if="!skillsForm.categories.length" class="empty-tip">暂未添加技能分类。</p>
    </el-card>

    <el-card class="edit-card" shadow="never">
      <template #header>
        <div class="edit-card__header">
          <span>教育背景</span>
          <el-button type="primary" size="small" @click="addEducation">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(edu, i) in skillsForm.education" :key="i" class="list-row list-row--stack">
        <el-row :gutter="12" class="row-full">
          <el-col :md="8" :sm="24"><el-input v-model="edu.degree" placeholder="学位" /></el-col>
          <el-col :md="8" :sm="24"><el-input v-model="edu.school" placeholder="学校" /></el-col>
          <el-col :md="6" :sm="24"><el-input v-model="edu.period" placeholder="时间段" /></el-col>
          <el-col :md="2" :sm="24"><el-button type="danger" size="small" @click="skillsForm.education.splice(i, 1)">删除</el-button></el-col>
        </el-row>
        <el-input v-model="edu.detail" placeholder="详情描述" />
      </div>
      <p v-if="!skillsForm.education.length" class="empty-tip">暂未添加教育背景。</p>
    </el-card>

    <el-card class="edit-card" shadow="never">
      <template #header>
        <div class="edit-card__header">
          <span>工作经历</span>
          <el-button type="primary" size="small" @click="addExperience">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(exp, i) in skillsForm.experience" :key="i" class="list-row list-row--stack">
        <el-row :gutter="12" class="row-full">
          <el-col :md="8" :sm="24"><el-input v-model="exp.company" placeholder="公司" /></el-col>
          <el-col :md="8" :sm="24"><el-input v-model="exp.role" placeholder="角色" /></el-col>
          <el-col :md="6" :sm="24"><el-input v-model="exp.period" placeholder="时间段" /></el-col>
          <el-col :md="2" :sm="24"><el-button type="danger" size="small" @click="skillsForm.experience.splice(i, 1)">删除</el-button></el-col>
        </el-row>
        <el-input v-model="exp.detail" placeholder="详情描述" />
      </div>
      <p v-if="!skillsForm.experience.length" class="empty-tip">暂未添加工作经历。</p>
    </el-card>

    <el-card class="edit-card" shadow="never">
      <template #header>
        <div class="edit-card__header">
          <span>兴趣爱好</span>
          <el-button type="primary" size="small" @click="skillsForm.hobbies.push('')">+ 添加</el-button>
        </div>
      </template>
      <div class="hobbies-row">
        <div v-for="(_, i) in skillsForm.hobbies" :key="i" class="hobby-item">
          <el-input v-model="skillsForm.hobbies[i]" placeholder="爱好" />
          <el-button type="danger" size="small" @click="skillsForm.hobbies.splice(i, 1)">×</el-button>
        </div>
      </div>
      <p v-if="!skillsForm.hobbies.length" class="empty-tip">暂未添加兴趣爱好。</p>
    </el-card>

    <div class="actions">
      <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">保存所有更改</el-button>
      <el-button size="large" @click="handleReset">重置</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  ElButton,
  ElCard,
  ElCol,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElRow,
  ElSelect
} from 'element-plus'
import { profileApi } from '@/api/profile'

const formRef = ref()
const loading = ref(false)

const iconOptions = ['school', 'briefcase', 'code', 'heart', 'server', 'cloud', 'database']
const colorOptions = ['cyan', 'coral', 'teal', 'orange']

const form = reactive({
  name: '',
  title: '',
  bio: '',
  avatar: '',
  email: '',
  github: '',
  linkedin: '',
  twitter: '',
  website: '',
  skills: ''
})

const skillsForm = reactive({
  stats: [] as { label: string; value: string; icon: string; color: string }[],
  categories: [] as { name: string; icon: string; itemsStr: string }[],
  education: [] as { degree: string; school: string; period: string; detail: string }[],
  experience: [] as { company: string; role: string; period: string; detail: string }[],
  hobbies: [] as string[]
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }]
}

function parseSkills(json: string) {
  try {
    const data = JSON.parse(json)
    skillsForm.stats = data.stats || []
    skillsForm.categories = (data.categories || []).map((c: any) => ({
      name: c.name,
      icon: c.icon,
      itemsStr: (c.items || []).join(', ')
    }))
    skillsForm.education = data.education || []
    skillsForm.experience = data.experience || []
    skillsForm.hobbies = data.hobbies || []
  } catch {
    skillsForm.stats = []
    skillsForm.categories = []
    skillsForm.education = []
    skillsForm.experience = []
    skillsForm.hobbies = []
  }
}

function buildSkillsJson(): string {
  return JSON.stringify({
    stats: skillsForm.stats,
    categories: skillsForm.categories.map((c) => ({
      name: c.name,
      icon: c.icon,
      items: c.itemsStr
        .split(/[,，]/)
        .map((s) => s.trim())
        .filter(Boolean)
    })),
    education: skillsForm.education,
    experience: skillsForm.experience,
    hobbies: skillsForm.hobbies.filter(Boolean)
  })
}

function addStat() {
  skillsForm.stats.push({ label: '', value: '', icon: 'code', color: 'cyan' })
}

function addCategory() {
  skillsForm.categories.push({ name: '', icon: 'code', itemsStr: '' })
}

function addEducation() {
  skillsForm.education.push({ degree: '', school: '', period: '', detail: '' })
}

function addExperience() {
  skillsForm.experience.push({ company: '', role: '', period: '', detail: '' })
}

async function fetchProfile() {
  try {
    const response = await profileApi.get()
    if (response.success && response.data) {
      Object.assign(form, response.data)
      parseSkills(form.skills)
    }
  } catch {
    ElMessage.error('获取个人信息失败')
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    form.skills = buildSkillsJson()
    const response = await profileApi.update(form)
    if (response.success) {
      ElMessage.success('保存成功')
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '保存失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  fetchProfile()
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>
.profile-edit-page {
  display: grid;
  gap: var(--space-lg);
}

.edit-card {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  background: var(--bg-panel);
}

.edit-card :deep(.el-card__header) {
  border-bottom-color: var(--border-primary);
}

.edit-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.profile-form :deep(.el-input__wrapper),
.edit-card :deep(.el-input__wrapper),
.edit-card :deep(.el-textarea__inner),
.edit-card :deep(.el-select__wrapper) {
  border-radius: var(--radius-md);
  background: var(--bg-panel-strong);
  box-shadow: inset 0 0 0 1px var(--border-primary);
}

.profile-form :deep(.el-input__wrapper.is-focus),
.edit-card :deep(.el-input__wrapper.is-focus),
.edit-card :deep(.el-select__wrapper.is-focused),
.edit-card :deep(.el-textarea__inner:focus) {
  box-shadow: inset 0 0 0 1px var(--accent-primary);
}

.list-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.list-row--stack {
  flex-direction: column;
  align-items: stretch;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: 12px;
}

.row-full {
  width: 100%;
}

.field-120 {
  width: 120px;
}

.field-110 {
  width: 110px;
}

.field-100 {
  width: 100px;
}

.field-grow {
  flex: 1;
}

.hobbies-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.hobby-item {
  display: flex;
  align-items: center;
  gap: 4px;
  width: 220px;
}

.empty-tip {
  margin-top: 4px;
  color: var(--text-muted);
  font-size: 0.84rem;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding: 8px 0 20px;
}

@media (max-width: 900px) {
  .list-row {
    flex-wrap: wrap;
  }

  .field-120,
  .field-110,
  .field-100,
  .field-grow,
  .hobby-item {
    width: 100%;
  }

  .actions {
    flex-direction: column;
  }

  .actions :deep(.el-button) {
    width: 100%;
  }
}
</style>
