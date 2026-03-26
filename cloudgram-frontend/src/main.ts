import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { setupStore } from '@/store'

// 验证环境变量（仅在开发环境显示）
if (import.meta.env.DEV) {
  console.log('Environment Variables:', {
    VITE_APP_TITLE: import.meta.env.VITE_APP_TITLE,
    VITE_API_BASE_URL: import.meta.env.VITE_API_BASE_URL,
    VITE_APP_ENV: import.meta.env.VITE_APP_ENV,
    BASE_URL: import.meta.env.BASE_URL
  })
}

const app = createApp(App)
setupStore(app)
app.use(router)

app.mount('#app')
