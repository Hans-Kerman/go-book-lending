// src/pages/Auth/Register.tsx
import React from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Card, Form, Input, Button, message } from 'antd';
import { useUserStore } from '../../store/userStore';
import apiClient from '../../services/api';
import type { AuthCredentials } from '../../types';

const RegisterPage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(false);
  const setToken = useUserStore((state) => state.setToken);

  const onFinish = async (values: AuthCredentials) => {
    setLoading(true);
    try {
      const response = await apiClient.post('/public/register', values);
      if (response.data && response.data.token) {
        // 注册成功，使用 store action 更新状态
        setToken(response.data.token);
        message.success('注册成功！已自动为您登录。');
        navigate('/');
      } else {
        message.error(response.data.error || '注册失败，请稍后重试');
      }
    } catch (error) {
      const errorMsg = (error as { response?: { data?: { error?: string } } })?.response?.data?.error || '注册时发生错误';
      message.error(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', paddingTop: '50px' }}>
      <Card title="用户注册" style={{ width: 400 }}>
        <Form
          name="register"
          onFinish={onFinish}
          layout="vertical"
          requiredMark="optional"
        >
          <Form.Item
            name="user_name"
            label="用户名"
            rules={[
              { required: true, message: '请输入您的用户名!' },
              { min: 2, message: '用户名至少2个字符' },
              { max: 50, message: '用户名最多50个字符' },
            ]}
          >
            <Input placeholder="设置您的用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入您的密码!' },
              { min: 6, message: '密码至少6个字符' },
            ]}
            hasFeedback
          >
            <Input.Password placeholder="设置您的密码" />
          </Form.Item>

          <Form.Item
            name="confirm"
            label="确认密码"
            dependencies={['password']}
            hasFeedback
            rules={[
              { required: true, message: '请确认您的密码!' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不匹配!'));
                },
              }),
            ]}
          >
            <Input.Password placeholder="再次输入密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} style={{ width: '100%' }}>
              注册
            </Button>
          </Form.Item>
        </Form>
        <div style={{ textAlign: 'center' }}>
          已经有账户了？ <Link to="/login">直接登录</Link>
        </div>
      </Card>
    </div>
  );
};

export default RegisterPage;
