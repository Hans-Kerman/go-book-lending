import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173, // 前端开发服务器端口
    proxy: {
      // '/api' 是我们在 src/config/index.ts 中定义的基础路径
      '/api': {
        target: 'http://localhost:8080', // 后端服务地址
        changeOrigin: true, // 改变源，解决跨域问题
        // 后端的基础路径已经包含了 /api，所以这里不需要 rewrite
        configure: (proxy) => {
          proxy.on('proxyReq', (proxyReq) => {
            // 重写 Origin 请求头，防止由于浏览器访问 127.0.0.1 导致后端 CORS 403 拦截
            proxyReq.setHeader('Origin', 'http://localhost:5173');
          });
        },
      },
    },
  },
})
