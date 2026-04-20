import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/auth': 'http://localhost:8080',
      '/check-access': 'http://localhost:8080',
      '/check-nft': 'http://localhost:8080',
      '/check-nft1155': 'http://localhost:8080',
      '/chains': 'http://localhost:8080',
      '/health': 'http://localhost:8080',
    },
  },
})
