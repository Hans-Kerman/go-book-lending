// src/store/userStore.ts
import { create } from 'zustand';
import { jwtDecode } from 'jwt-decode';

/**
 * 后端 JWT payload 中编码的用户信息结构
 * @property {number} id - 用户ID
 * @property {string} name - 用户名
 * @property {'admin' | 'user'} role - 用户角色
 * @property {number} exp - Token 过期时间戳
 * @property {number} iat - Token 签发时间戳
 */
interface JwtPayload {
  user_id: number;
  user_name: string;
  role: 'admin' | 'user';
  exp: number;
  iat: number;
}


/**
 * 应用中使用的用户信息结构
 */
interface User {
  id: number;
  name: string;
  role: 'admin' | 'user';
}

/**
 * Zustand store 的 state 和 actions 类型定义
 */
interface UserState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  isAdmin: boolean;
  setToken: (token: string | null) => void;
  logout: () => void;
  initializeAuth: () => void;
}

/**
 * 安全地解析 JWT 的辅助函数
 * @param token JWT 字符串
 * @returns 解析后的用户信息或 null
 */
const decodeToken = (token: string): User | null => {
  try {
    const decoded = jwtDecode<JwtPayload>(token);
    // 检查 token 是否过期
    if (decoded.exp * 1000 < Date.now()) {
      console.warn("Token has expired.");
      return null;
    }
    return {
      id: decoded.user_id,
      name: decoded.user_name,
      role: decoded.role,
    };
  } catch (error) {
    console.error("Failed to decode token:", error);
    return null;
  }
};

export const useUserStore = create<UserState>((set) => ({
  // --- STATE ---
  token: null,
  user: null,
  isAuthenticated: false,
  isAdmin: false,

  // --- ACTIONS ---
  
  /**
   * 在应用启动时初始化认证状态
   */
  initializeAuth: () => {
    const token = localStorage.getItem('jwt_token');
    if (token) {
      const decodedUser = decodeToken(token);
      if (decodedUser) {
        set({
          token,
          user: decodedUser,
          isAuthenticated: true,
          isAdmin: decodedUser.role === 'admin',
        });
      } else {
        // Token 无效或已过期，执行登出逻辑
        localStorage.removeItem('jwt_token');
      }
    }
  },

  /**
   * 设置新的 token（通常在登录或注册后调用）
   * @param token 新的 JWT
   */
  setToken: (token) => {
    if (token) {
      localStorage.setItem('jwt_token', token);
      const decodedUser = decodeToken(token);
      set({
        token,
        user: decodedUser,
        isAuthenticated: !!decodedUser,
        isAdmin: decodedUser?.role === 'admin',
      });
    } else {
      // 如果传入 null，则视为登出
      localStorage.removeItem('jwt_token');
      set({ token: null, user: null, isAuthenticated: false, isAdmin: false });
    }
  },

  /**
   * 退出登录
   */
  logout: () => {
    localStorage.removeItem('jwt_token');
    set({ token: null, user: null, isAuthenticated: false, isAdmin: false });
  },
}));

// 在模块加载时立即初始化认证状态
useUserStore.getState().initializeAuth();
