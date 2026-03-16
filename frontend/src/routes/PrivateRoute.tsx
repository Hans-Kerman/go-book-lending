// src/routes/PrivateRoute.tsx
import React from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useUserStore } from '../store/userStore';

interface PrivateRouteProps {
  /**
   * 允许访问该路由的角色数组
   * @default undefined - 默认只需要登录即可访问
   */
  allowedRoles?: ('admin' | 'user')[];
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ allowedRoles }) => {
  const { isAuthenticated, user } = useUserStore();
  const location = useLocation();

  if (!isAuthenticated) {
    // 1. 用户未登录，重定向到登录页
    //    同时，将他们原本想访问的页面路径记录在 state 中
    //    这样登录后，我们可以将他们导航回原来的页面
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (allowedRoles && user?.role && !allowedRoles.includes(user.role)) {
    // 2. 用户已登录，但角色不匹配 (例如，普通用户尝试访问管理员页面)
    //    这里可以重定向到一个统一的 "403 Forbidden" 页面
    //    为简单起见，我们暂时重定向回首页
    return <Navigate to="/" replace />;
  }

  // 3. 用户已登录且角色匹配 (或路由不要求特定角色)
  //    渲染子路由
  return <Outlet />;
};

export default PrivateRoute;
