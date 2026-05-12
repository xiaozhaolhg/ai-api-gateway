import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/admin': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  // Ant Design CSS-in-JS compatibility
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
      },
    },
  },
  // Ensure JSON files are properly imported
  json: {
    stringify: true,
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          if (id.includes('react') || id.includes('react-dom')) {
            return 'vendor';
          }
          if (id.includes('antd')) {
            return 'antd';
          }
          if (id.includes('i18next') || id.includes('react-i18next')) {
            return 'i18n';
          }
        },
      },
    },
  },
})
