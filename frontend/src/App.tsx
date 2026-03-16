// src/App.tsx
import React, { Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Outlet } from 'react-router-dom';
import { Layout, Menu, theme, Spin } from 'antd';

const { Header, Content, Footer } = Layout;

// 懒加载页面组件
// 注意：这里的路径现在指向目录，而不是具体文件
const HomePage = React.lazy(() => import('./pages/Home')); 
const LoginPage = React.lazy(() => import('./pages/Auth/Login'));
const RegisterPage = React.lazy(() => import('./pages/Auth/Register'));

// 全屏加载指示器
const FullScreenSpinner: React.FC = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
    <Spin size="large" />
  </div>
);

const AppLayout: React.FC = () => {
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center' }}>
        <div className="logo" style={{ color: 'white', marginRight: '24px' }}>图书系统</div>
        <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['1']}>
          <Menu.Item key="1"><Link to="/">首页</Link></Menu.Item>
          <Menu.Item key="2"><Link to="/login">登录</Link></Menu.Item>
        </Menu>
      </Header>
      <Content style={{ padding: '48px' }}>
        <div
          style={{
            background: colorBgContainer,
            minHeight: 280,
            padding: 24,
            borderRadius: borderRadiusLG,
          }}
        >
          {/* Suspense 的 fallback 会在懒加载组件下载和解析时显示 */}
          <Suspense fallback={<Spin />}>
            <Outlet />
          </Suspense>
        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        Go Book Lending ©{new Date().getFullYear()} Created by Roo
      </Footer>
    </Layout>
  );
};

const App: React.FC = () => (
  <Router>
    {/* Suspense 需要放在路由定义的外层，以便在切换路由、加载新组件时显示 fallback */}
    <Suspense fallback={<FullScreenSpinner />}>
      <Routes>
        <Route path="/" element={<AppLayout />}>
          <Route index element={<HomePage />} />
          <Route path="login" element={<LoginPage />} />
          <Route path="register" element={<RegisterPage />} />
          {/* 后续页面将在这里添加 */}
        </Route>
      </Routes>
    </Suspense>
  </Router>
);

export default App;
