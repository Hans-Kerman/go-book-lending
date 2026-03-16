// src/App.tsx
import React, { Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Outlet, useNavigate } from 'react-router-dom';
import { Layout, Menu, theme, Spin, Button, Typography, Dropdown } from 'antd';
import { useUserStore } from './store/userStore';

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
  const { isAuthenticated, user, logout } = useUserStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const userMenu = (
    <Menu>
      <Menu.Item key="username" disabled>
        <Typography.Text>你好, {user?.name}</Typography.Text>
      </Menu.Item>
      <Menu.Divider />
      {user?.role === 'admin' && (
        <Menu.Item key="admin">
          <Link to="/admin/books">图书管理</Link>
        </Menu.Item>
      )}
      <Menu.Item key="my-borrows">
        <Link to="/user/borrows">我的借阅</Link>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item key="logout" onClick={handleLogout}>
        退出登录
      </Menu.Item>
    </Menu>
  );

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <div className="logo" style={{ color: 'white', marginRight: '24px' }}>图书系统</div>
          <Menu theme="dark" mode="horizontal" defaultSelectedKeys={['1']} items={[
            { key: '1', label: <Link to="/">首页</Link> },
          ]} />
        </div>
        <div>
          {isAuthenticated ? (
            <Dropdown menu={{ items: userMenu.props.items }} placement="bottomRight">
              <Button type="primary">欢迎, {user?.name}</Button>
            </Dropdown>
          ) : (
            <>
              <Button type="primary" style={{ marginRight: '10px' }} onClick={() => navigate('/login')}>
                登录
              </Button>
              <Button onClick={() => navigate('/register')}>
                注册
              </Button>
            </>
          )}
        </div>
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
