// src/pages/User/MyBorrows.tsx
import React, { useState, useEffect } from 'react';
import { Table, Typography, Alert, Spin } from 'antd';
import apiClient from '../../services/api';
import { useUserStore } from '../../store/userStore';
import type { LendRecordResponse } from '../../types';
import { Link } from 'react-router-dom';

const { Title } = Typography;

const columns = [
  {
    title: '借阅记录ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '书籍ISBN',
    dataIndex: 'book_id',
    key: 'book_id',
    render: (isbn: string) => <Link to={`/book/${isbn}`}>{isbn}</Link>,
  },
  {
    title: '借阅时间',
    dataIndex: 'created_at',
    key: 'created_at',
    render: (text: string) => new Date(text).toLocaleString(),
  },
  {
    title: '归还时间',
    dataIndex: 'return_time',
    key: 'return_time',
    render: (text: string) => (text ? new Date(text).toLocaleString() : '尚未归还'),
  },
];

const MyBorrowsPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [records, setRecords] = useState<LendRecordResponse[]>([]);
  const { user } = useUserStore();

  useEffect(() => {
    const fetchBorrows = async () => {
      if (!user) return;
      setLoading(true);
      setError(null);
      try {
        const response = await apiClient.get<{ data: LendRecordResponse[] }>(`/user/borrows`);
        setRecords(response.data.data || []);
      } catch (err) {
        const errorMsg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '获取借阅记录失败';
        setError(errorMsg);
      } finally {
        setLoading(false);
      }
    };

    fetchBorrows();
  }, [user]);

  if (loading) {
    return <div style={{ display: 'flex', justifyContent: 'center', padding: '50px' }}><Spin size="large" /></div>;
  }

  if (error) {
    return <Alert message="错误" description={error} type="error" showIcon />;
  }

  return (
    <div>
      <Title level={2}>我的借阅记录</Title>
      <Table columns={columns} dataSource={records} rowKey="id" />
    </div>
  );
};

export default MyBorrowsPage;
