// src/App.tsx
import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Outlet } from 'react-router-dom';
import { Layout, Menu, theme } from 'antd';

const { Header, Content, Footer } = Layout;

// 简单的占位符组件
const HomePage = () => <div><h1>欢迎来到图书借阅系统</h1></div>;
const LoginPage = () => <div>登录页</div>;

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
          {/* 路由匹配的组件将在这里渲染 */}
          <Outlet />
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
    <Routes>
      <Route path="/" element={<AppLayout />}>
        {/* 嵌套路由，所有在 AppLayout 内部显示的页面都在这里定义 */}
        <Route index element={<HomePage />} />
        <Route path="login" element={<LoginPage />} />
        {/* 后续页面将在这里添加 */}
      </Route>
    </Routes>
  </Router>
);

export default App;
