// src/pages/Book/BookDetail.tsx
import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Spin, Alert, Row, Col, Image, Typography, Button, message, Descriptions } from 'antd';
import apiClient from '../../services/api';
import { useUserStore } from '../../store/userStore';
import type { Book, LendRecordResponse } from '../../types';

const { Title, Paragraph } = Typography;

const BookDetailPage: React.FC = () => {
  const { isbn } = useParams<{ isbn: string }>();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useUserStore();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [book, setBook] = useState<Book | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [isBorrowed, setIsBorrowed] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      if (!isbn) return;
      setLoading(true);
      setError(null);
      try {
        // 1. 获取书籍详情
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const bookResponse = await apiClient.get<any>(`/public/book/${isbn}`);
        // 兼容后端返回的大写字段名
        const b = bookResponse.data.book || bookResponse.data.Book || bookResponse.data;
        const fetchedBook: Book = {
          id: b.ID ?? b.id,
          createdAt: b.CreatedAt ?? b.createdAt,
          updatedAt: b.UpdatedAt ?? b.updatedAt,
          isbn: b.ISBN ?? b.isbn,
          title: b.Title ?? b.title,
          author: b.Author ?? b.author,
          coverURL: b.CoverURL ?? b.coverURL,
          available: b.Available ?? b.available,
          price: b.Price ?? b.price,
        };
        setBook(fetchedBook);

        // 2. 如果用户已登录，获取其借阅记录以判断是否已借阅此书
        if (isAuthenticated && user) {
          try {
            const borrowsResponse = await apiClient.get<{ data: LendRecordResponse[] }>(`/user/borrows`);
            const userRecords = borrowsResponse.data.data || [];
            const hasBorrowed = userRecords.some(
              record => record.book_id === isbn && !record.return_time
            );
            setIsBorrowed(hasBorrowed);
          } catch (borrowErr) {
            console.error('获取借阅记录失败:', borrowErr);
          }
        }
      } catch (err) {
        console.error('获取书籍详情失败:', err);
        setError('无法加载书籍详情，请检查书籍ISBN是否正确或稍后再试。');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [isbn, isAuthenticated, user]);

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
      setIsBorrowed(true);
      // 乐观更新UI
      setBook(prevBook => prevBook ? { ...prevBook, available: prevBook.available - 1 } : null);
    } catch (err) {
      const errorMsg = (err as { response?: { data?: { error: string } } }).response?.data?.error || '借阅失败';
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
      setIsBorrowed(false);
      // 乐观更新UI
      setBook(prevBook => prevBook ? { ...prevBook, available: prevBook.available + 1 } : null);
    } catch (err) {
      const errorMsg = (err as { response?: { data?: { error: string } } }).response?.data?.error || '还书失败';
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
        {isAuthenticated && (
          <div style={{ marginTop: '24px' }}>
            {isBorrowed ? (
              <Button
                type="primary"
                onClick={handleReturn}
                loading={actionLoading}
              >
                归还本书
              </Button>
            ) : (
              <Button
                type="primary"
                onClick={handleBorrow}
                loading={actionLoading}
                disabled={book.available <= 0}
              >
                {book.available > 0 ? '立即借阅' : '已无库存'}
              </Button>
            )}
          </div>
        )}
      </Col>
    </Row>
  );
};

export default BookDetailPage;
