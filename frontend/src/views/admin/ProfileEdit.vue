<template>
  <div class="profile-edit-page">
    <h1 class="page-title">编辑个人信息</h1>

    <el-card header="基本信息" style="margin-bottom: 20px">
      <el-form ref="formRef" :model="form" label-width="100px" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="form.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="身份标签" prop="title">
          <el-input v-model="form.title" placeholder="多个标签用逗号分隔，如：🎓 硕士在读,💻 技术爱好者" />
        </el-form-item>
        <el-form-item label="个人简介" prop="bio">
          <el-input v-model="form.bio" type="textarea" :rows="3" placeholder="支持换行，第一行为主简介" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="GitHub">
              <el-input v-model="form.github" placeholder="https://github.com/xxx" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="个人网站">
              <el-input v-model="form.website" placeholder="https://..." />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="LinkedIn">
              <el-input v-model="form.linkedin" placeholder="LinkedIn链接" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Twitter">
              <el-input v-model="form.twitter" placeholder="Twitter链接" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="头像URL">
          <el-input v-model="form.avatar" placeholder="头像图片URL" />
        </el-form-item>
      </el-form>
    </el-card>
    <el-card style="margin-bottom: 20px">
      <template #header>
        <div class="card-header">
          <span>统计卡片</span>
          <el-button type="primary" size="small" @click="addStat">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(stat, i) in skillsForm.stats" :key="i" class="list-row">
        <el-input v-model="stat.label" placeholder="标签" style="width: 120px" />
        <el-input v-model="stat.value" placeholder="值" style="width: 100px" />
        <el-select v-model="stat.icon" placeholder="图标" style="width: 120px">
          <el-option v-for="ic in iconOptions" :key="ic" :label="ic" :value="ic" />
        </el-select>
        <el-select v-model="stat.color" placeholder="颜色" style="width: 110px">
          <el-option v-for="c in colorOptions" :key="c" :label="c" :value="c" />
        </el-select>
        <el-button type="danger" size="small" @click="skillsForm.stats.splice(i, 1)">删除</el-button>
      </div>
    </el-card>

    <el-card style="margin-bottom: 20px">
      <template #header>
        <div class="card-header">
          <span>技能分类</span>
          <el-button type="primary" size="small" @click="addCategory">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(cat, i) in skillsForm.categories" :key="i" class="list-row">
        <el-input v-model="cat.name" placeholder="分类名" style="width: 120px" />
        <el-select v-model="cat.icon" placeholder="图标" style="width: 120px">
          <el-option v-for="ic in iconOptions" :key="ic" :label="ic" :value="ic" />
        </el-select>
        <el-input v-model="cat.itemsStr" placeholder="技能项，逗号分隔" style="flex: 1" />
        <el-button type="danger" size="small" @click="skillsForm.categories.splice(i, 1)">删除</el-button>
      </div>
    </el-card>

    <el-card style="margin-bottom: 20px">
      <template #header>
        <div class="card-header">
          <span>教育背景</span>
          <el-button type="primary" size="small" @click="addEducation">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(edu, i) in skillsForm.education" :key="i" class="list-row vertical">
        <el-row :gutter="12" style="width: 100%">
          <el-col :span="8"><el-input v-model="edu.degree" placeholder="学位" /></el-col>
          <el-col :span="8"><el-input v-model="edu.school" placeholder="学校" /></el-col>
          <el-col :span="6"><el-input v-model="edu.period" placeholder="时间段" /></el-col>
          <el-col :span="2"><el-button type="danger" size="small" @click="skillsForm.education.splice(i, 1)">删除</el-button></el-col>
        </el-row>
        <el-input v-model="edu.detail" placeholder="详情描述" style="margin-top: 8px" />
      </div>
    </el-card>
    <el-card style="margin-bottom: 20px">
      <template #header>
        <div class="card-header">
          <span>工作经历</span>
          <el-button type="primary" size="small" @click="addExperience">+ 添加</el-button>
        </div>
      </template>
      <div v-for="(exp, i) in skillsForm.experience" :key="i" class="list-row vertical">
        <el-row :gutter="12" style="width: 100%">
          <el-col :span="8"><el-input v-model="exp.company" placeholder="公司" /></el-col>
          <el-col :span="8"><el-input v-model="exp.role" placeholder="角色" /></el-col>
          <el-col :span="6"><el-input v-model="exp.period" placeholder="时间段" /></el-col>
          <el-col :span="2"><el-button type="danger" size="small" @click="skillsForm.experience.splice(i, 1)">删除</el-button></el-col>
        </el-row>
        <el-input v-model="exp.detail" placeholder="详情描述" style="margin-top: 8px" />
      </div>
    </el-card>

    <el-card style="margin-bottom: 20px">
      <template #header>
        <div class="card-header">
          <span>兴趣爱好</span>
          <el-button type="primary" size="small" @click="skillsForm.hobbies.push('')">+ 添加</el-button>
        </div>
      </template>
      <div class="hobbies-row">
        <div v-for="(_, i) in skillsForm.hobbies" :key="i" class="hobby-item">
          <el-input v-model="skillsForm.hobbies[i]" placeholder="爱好" style="width: 160px" />
          <el-button type="danger" size="small" @click="skillsForm.hobbies.splice(i, 1)">×</el-button>
        </div>
      </div>
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
  ElCard, ElForm, ElFormItem, ElInput, ElButton, ElMessage,
  ElRow, ElCol, ElSelect, ElOption
} from 'element-plus'
import { profileApi } from '@/api/profile'

const formRef = ref()
const loading = ref(false)

const iconOptions = ['school', 'briefcase', 'code', 'heart', 'server', 'cloud', 'database']
const colorOptions = ['cyan', 'coral', 'teal', 'orange']

const form = reactive({
  name: '', title: '', bio: '', avatar: '',
  email: '', github: '', linkedin: '', twitter: '', website: '', skills: ''
})

const skillsForm = reactive({
  stats: [] as { label: string; value: string; icon: string; color: string }[],
  categories: [] as { name: string; icon: string; itemsStr: string }[],
  education: [] as { degree: string; school: string; period: string; detail: string }[],
  experience: [] as { company: string; role: string; period: string; detail: string }[],
  hobbies: [] as string[]
})

const rules = { name: [{ required: true, message: '请输入姓名', trigger: 'blur' }] }

function parseSkills(json: string) {
  try {
    const data = JSON.parse(json)
    skillsForm.stats = data.stats || []
    skillsForm.categories = (data.categories || []).map((c: any) => ({
      name: c.name, icon: c.icon, itemsStr: (c.items || []).join(', ')
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
    categories: skillsForm.categories.map(c => ({
      name: c.name, icon: c.icon,
      items: c.itemsStr.split(/[,，]/).map(s => s.trim()).filter(Boolean)
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

const fetchProfile = async () => {
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

const handleReset = () => { fetchProfile() }

onMounted(() => { fetchProfile() })
</script>

<style scoped>
.profile-edit-page {
  padding: 20px;
  max-width: 1100px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.list-row.vertical {
  flex-direction: column;
  align-items: stretch;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
  margin-bottom: 12px;
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
}

.actions {
  text-align: center;
  padding: 20px 0;
}
</style>