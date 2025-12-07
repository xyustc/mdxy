import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import './styles/main.css'
import { initCopyProtection, addWatermark } from './utils/copyProtection'
import { ensureVisitorId } from './utils/visitorId'

// 确保用户标识符存在
ensureVisitorId();

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

// 注册所有Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(ElementPlus)

// 初始化防复制保护
initCopyProtection()

// 添加水印（可选，取消注释启用）
addWatermark('仅供个人学习(XingYu)')

