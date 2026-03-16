// src/pages/Book/BookList.tsx
import React, { useState, useEffect } from 'react';
import { List, Card, Spin, Alert, Pagination, Typography } from 'antd';
import { Link } from 'react-router-dom';
import apiClient from '../../services/api';
import type { Book } from '../../types';

const { Meta } = Card;

interface BooksResponse {
  total: number;
  page: number;
  page_size: number;
  totalPages: number;
  books: Book[];
}

const BookListPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<BooksResponse | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 12; // 每页显示12本书

  useEffect(() => {
    const fetchBooks = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await apiClient.get<BooksResponse>('/public/books', {
          params: {
            page: currentPage,
            page_size: pageSize,
          },
        });
        setData(response.data);
      } catch (err) {
        console.error('获取图书列表失败:', err);
        setError('获取图书列表失败，请稍后再试。');
      }finally {
        setLoading(false);
      }
    };

    fetchBooks();
  }, [currentPage]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  if (loading && !data) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '300px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return <Alert message="错误" description={error} type="error" showIcon />;
  }

  if (!data || data.books.length === 0) {
    return <Typography.Text>暂无图书</Typography.Text>;
  }

  return (
    <div>
      <List
        grid={{
          gutter: 16,
          xs: 1,
          sm: 2,
          md: 3,
          lg: 4,
          xl: 6,
          xxl: 6,
        }}
        dataSource={data.books}
        renderItem={(book) => (
          <List.Item>
            <Link to={`/book/${book.isbn}`}>
              <Card
                hoverable
                cover={<img alt={book.title} src={book.coverURL || '/placeholder.png'} style={{ height: 300, objectFit: 'cover' }} />}
              >
                <Meta title={book.title} description={book.author} />
              </Card>
            </Link>
          </List.Item>
        )}
      />
      <div style={{ marginTop: '24px', textAlign: 'center' }}>
        <Pagination
          current={currentPage}
          pageSize={pageSize}
          total={data.total}
          onChange={handlePageChange}
          showSizeChanger={false}
        />
      </div>
    </div>
  );
};

export default BookListPage;
