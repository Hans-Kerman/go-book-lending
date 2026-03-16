// src/services/api.ts
import axios from 'axios';
import { API_BASE_URL } from '../config';

// 创建一个 Axios 实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // 设置请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
});

/*
  后续可以根据需要添加请求和响应拦截器
  例如，在请求头中自动添加 JWT Token
*/

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 假设 token 存储在 localStorage
    const token = localStorage.getItem('jwt_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 (可选，用于统一处理错误)
apiClient.interceptors.response.use(
  (response) => {
    // 对响应数据做点什么
    return response;
  },
  (error) => {
    // 对响应错误做点什么
    // 例如，如果收到 401 Unauthorized，可以重定向到登录页
    if (error.response && error.response.status === 401) {
      // window.location.href = '/login';
      console.error('Unauthorized request, redirecting to login may be needed.');
    }
    return Promise.reject(error);
  }
);

export default apiClient;
