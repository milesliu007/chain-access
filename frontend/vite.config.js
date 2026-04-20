import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        admin: resolve(__dirname, 'index-admin.html'),
      },
    },
  },
  server: {
    proxy: {
      '/auth': 'http://localhost:8080',
      '/check-access': 'http://localhost:8080',
      '/check-nft': 'http://localhost:8080',
      '/check-nft1155': 'http://localhost:8080',
      '/chains': 'http://localhost:8080',
      '/health': 'http://localhost:8080',
      '/admin/login': 'http://localhost:8080',
      '/admin/balances': 'http://localhost:8080',
    },
  },
})
