import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import './styles/main.css'
import { removeCopyProtection, addWatermark } from './utils/copyProtection'
import { ensureVisitorId } from './utils/visitorId'
import router from './router'

const app = createApp({
  template: '<router-view />'
})

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 注册所有Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#admin-app')

// 确保用户标识符存在
ensureVisitorId();

// 注意：我们移除了initCopyProtection()调用，因为管理后台不需要防复制保护
// initCopyProtection()

// 移除可能已存在的防复制保护（例如从客户端路由导航过来的情况）
removeCopyProtection();

// 添加水印（可选，取消注释启用）
addWatermark('仅供个人学习(XingYu)')