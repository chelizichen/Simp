import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/simpserver': {
        target: 'http://localhost:8511/',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/simpserver/, '/simpserver/') // 不可以省略rewrite
      }
    }
  }
})
