import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      // Проксируем запросы с /api на бэкенд
      '/api': {
        target: 'http://localhost:8080', // В режиме разработки используем локальный порт
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '') // Убираем /api перед отправкой на бэкенд
      }
    }
  },
  // Настройки для продакшн сборки
  build: {
    outDir: 'dist',
    // Режим поддержки при CORS ошибках и других проблемах
    sourcemap: false,
    // Минимизация для продакшна
    minify: 'terser',
    // Настройки делают сборку более эффективной
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    }
  }
})
