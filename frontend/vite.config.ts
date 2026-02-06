import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/todo.v1.TodoService': 'http://localhost:8080',
      '/event.v1.EventService': 'http://localhost:8080',
    },
  },
})
