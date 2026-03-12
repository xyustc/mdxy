import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    rollupOptions: {
      input: {
        main: fileURLToPath(new URL('./index.html', import.meta.url))
      },
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) return

          if (id.includes('/echarts/') || id.includes('/vue-echarts/')) {
            return 'chart-vendor'
          }

          if (id.includes('/markdown-it/') || id.includes('/highlight.js/')) {
            return 'markdown-vendor'
          }

          if (id.includes('/element-plus/')) {
            return 'element-vendor'
          }

          if (id.includes('/naive-ui/')) {
            return 'naive-vendor'
          }

          if (id.includes('/@element-plus/icons-vue/') || id.includes('/@vicons/')) {
            return 'icon-vendor'
          }

          if (id.includes('/vue/') || id.includes('/vue-router/') || id.includes('/pinia/')) {
            return 'vue-vendor'
          }

          if (id.includes('/axios/') || id.includes('/async-validator/') || id.includes('/dayjs/')) {
            return 'utility-vendor'
          }

          return 'vendor'
        }
      }
    }
  }
})
