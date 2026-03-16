// src/pages/Book/BookDetail.tsx
import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Spin, Alert, Row, Col, Image, Typography, Button, message, Descriptions } from 'antd';
import apiClient from '../../services/api';
import { useUserStore } from '../../store/userStore';
import type { Book } from '../../types';

const { Title, Paragraph } = Typography;

const BookDetailPage: React.FC = () => {
  const { isbn } = useParams<{ isbn: string }>();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useUserStore();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [book, setBook] = useState<Book | null>(null);
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    const fetchBook = async () => {
      if (!isbn) return;
      setLoading(true);
      setError(null);
      try {
        const response = await apiClient.get<{ book: Book }>(`/public/book/${isbn}`);
        setBook(response.data.book);
      } catch (err) {
        setError('无法加载书籍详情，请检查书籍ISBN是否正确或稍后再试。');
      } finally {
        setLoading(false);
      }
    };
    fetchBook();
  }, [isbn]);

  const handleBorrow = async () => {
    if (!isAuthenticated || !user) {
      message.warning('请先登录后再进行借阅！');
      navigate('/login');
      return;
    }
    setActionLoading(true);
    try {
      await apiClient.post('/borrow', { book_id: isbn, borrow_reader: user.id });
      message.success('借阅成功！');
      // 可以选择跳转到“我的借阅”页面或刷新当前页面信息
      // 为了简单起见，我们暂时只显示消息
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || '借阅失败';
      message.error(errorMsg);
    } finally {
      setActionLoading(false);
    }
  };
  
  // 还书功能 (暂未实现，作为后续步骤)
  const handleReturn = async () => {
     if (!isAuthenticated || !user) {
      message.warning('请先登录！');
      navigate('/login');
      return;
    }
    setActionLoading(true);
    try {
      await apiClient.post('/return', { book_id: isbn, borrow_reader: user.id });
      message.success('还书成功！');
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || '还书失败';
      message.error(errorMsg);
    } finally {
      setActionLoading(false);
    }
  };


  if (loading) {
    return <div style={{ display: 'flex', justifyContent: 'center', padding: '50px' }}><Spin size="large" /></div>;
  }

  if (error) {
    return <Alert message="错误" description={error} type="error" showIcon />;
  }

  if (!book) {
    return <Alert message="未找到书籍" type="warning" showIcon />;
  }

  return (
    <Row gutter={[32, 32]}>
      <Col xs={24} md={8} style={{ textAlign: 'center' }}>
        <Image
          width={250}
          src={book.coverURL || '/placeholder.png'}
          alt={book.title}
        />
      </Col>
      <Col xs={24} md={16}>
        <Title level={2}>{book.title}</Title>
        <Descriptions bordered column={1}>
          <Descriptions.Item label="作者">{book.author}</Descriptions.Item>
          <Descriptions.Item label="ISBN">{book.isbn}</Descriptions.Item>
          <Descriptions.Item label="库存剩余">{book.available}</Descriptions.Item>
        </Descriptions>
        <Paragraph style={{ marginTop: '20px' }}>
          这里是书籍的详细介绍... (暂无)
        </Paragraph>
        <div style={{ marginTop: '24px' }}>
          <Button 
            type="primary" 
            style={{ marginRight: '16px' }}
            onClick={handleBorrow}
            loading={actionLoading}
            disabled={book.available <= 0}
          >
            {book.available > 0 ? '立即借阅' : '已无库存'}
          </Button>
          <Button
            onClick={handleReturn}
            loading={actionLoading}
          >
            归还本书
          </Button>
        </div>
      </Col>
    </Row>
  );
};

export default BookDetailPage;
