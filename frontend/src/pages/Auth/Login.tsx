// src/pages/Auth/Login.tsx
import React from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Card, Form, Input, Button, message } from 'antd';
import { useUserStore } from '../../store/userStore';
import apiClient from '../../services/api';
import type { AuthCredentials } from '../../types';

const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(false);
  const setToken = useUserStore((state) => state.setToken);

  const onFinish = async (values: AuthCredentials) => {
    setLoading(true);
    try {
      const response = await apiClient.post('/public/login', values);
      if (response.data && response.data.token) {
        setToken(response.data.token); // 使用 store action 来更新状态
        message.success('登录成功！');
        navigate('/'); // 跳转到首页
      } else {
        message.error(response.data.error || '登录失败，请稍后重试');
      }
    } catch (error) {
      const errorMsg = (error as { response?: { data?: { error?: string } } })?.response?.data?.error || '登录时发生错误，请检查您的凭据';
      message.error(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', paddingTop: '50px' }}>
      <Card title="用户登录" style={{ width: 400 }}>
        <Form
          name="login"
          onFinish={onFinish}
          layout="vertical"
          requiredMark="optional"
        >
          <Form.Item
            name="user_name"
            label="用户名"
            rules={[{ required: true, message: '请输入您的用户名!' }]}
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入您的密码!' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} style={{ width: '100%' }}>
              登录
            </Button>
          </Form.Item>
        </Form>
        <div style={{ textAlign: 'center' }}>
          还没有账户？ <Link to="/register">立即注册</Link>
        </div>
      </Card>
    </div>
  );
};

export default LoginPage;
