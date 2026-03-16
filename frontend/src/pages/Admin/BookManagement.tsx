// src/pages/Admin/BookManagement.tsx
import React, { useState, useEffect, useCallback } from 'react';
import { Table, Button, Space, Popconfirm, message, Spin, Alert, Modal, Form, Input, InputNumber } from 'antd';
import apiClient from '../../services/api';
import type { Book, NewBookInfo } from '../../types';

// 定义表单和模态框的状态
type ModalState = {
  visible: boolean;
  mode: 'add' | 'edit';
  book: Book | null;
};

const BookManagementPage: React.FC = () => {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalState, setModalState] = useState<ModalState>({ visible: false, mode: 'add', book: null });
  const [form] = Form.useForm();

  // 获取所有书籍
  const fetchBooks = useCallback(async () => {
    setLoading(true);
    try {
      // 注意：后端分页接口可能不适合一次性获取所有书籍
      // 这里暂时假设有一个能返回所有书籍的接口或调整分页参数
      const response = await apiClient.get<{ books: Book[] }>('/public/books?page_size=999');
      setBooks(response.data.books || []);
    } catch (err) {
      const errorMsg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '获取书籍列表失败';
      setError(errorMsg);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchBooks();
  }, [fetchBooks]);

  // --- CRUD 操作 ---

  const handleDelete = async (isbn: string) => {
    try {
      await apiClient.delete(`/admin/book/del/${isbn}`);
      message.success('删除成功');
      fetchBooks(); // 重新获取数据
    } catch (err) {
      const errorMsg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '删除失败';
      message.error(errorMsg);
    }
  };

  const handleModalOk = async () => {
    try {
      const values: NewBookInfo = await form.validateFields();
      
      if (modalState.mode === 'add') {
        await apiClient.post('/admin/book', values);
        message.success('添加成功');
      } else if (modalState.mode === 'edit' && modalState.book) {
        // 后端更新接口需要ISBN在body中，我们确保它存在
        await apiClient.put('/admin/book', { ...values, isbn: modalState.book.isbn });
        message.success('更新成功');
      }
      
      setModalState({ visible: false, mode: 'add', book: null });
      fetchBooks(); // 重新获取数据
    } catch (info) {
      console.log('Validate Failed:', info);
      message.error('请检查表单输入！');
    }
  };

  // --- 模态框和表单处理 ---

  const showAddModal = () => {
    form.resetFields();
    setModalState({ visible: true, mode: 'add', book: null });
  };

  const showEditModal = (book: Book) => {
    form.setFieldsValue(book);
    setModalState({ visible: true, mode: 'edit', book });
  };

  const handleModalCancel = () => {
    setModalState({ visible: false, mode: 'add', book: null });
  };


  const columns = [
    { title: 'ISBN', dataIndex: 'isbn', key: 'isbn' },
    { title: '书名', dataIndex: 'title', key: 'title' },
    { title: '作者', dataIndex: 'author', key: 'author' },
    { title: '库存', dataIndex: 'available', key: 'available' },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Book) => (
        <Space size="middle">
          <Button type="link" onClick={() => showEditModal(record)}>编辑</Button>
          <Popconfirm
            title="确定要删除这本书吗？"
            onConfirm={() => handleDelete(record.isbn)}
            okText="是"
            cancelText="否"
          >
            <Button type="link" danger>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  if (loading) return <Spin />;
  if (error) return <Alert message="Error" description={error} type="error" showIcon />;

  return (
    <div>
      <Button type="primary" onClick={showAddModal} style={{ marginBottom: 16 }}>
        添加新书
      </Button>
      <Table columns={columns} dataSource={books} rowKey="isbn" />

      <Modal
        title={modalState.mode === 'add' ? '添加新书' : '编辑书籍'}
        visible={modalState.visible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={false} // 可以根据API请求状态来设置
      >
        <Form form={form} layout="vertical" name="book_form">
          <Form.Item name="isbn" label="ISBN" rules={[{ required: true, message: '请输入ISBN' }]}>
            <Input disabled={modalState.mode === 'edit'} />
          </Form.Item>
          <Form.Item name="title" label="书名" rules={[{ required: true, message: '请输入书名' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="author" label="作者">
            <Input />
          </Form.Item>
          <Form.Item name="available" label="库存" rules={[{ type: 'number', min: 0, message: '库存不能为负' }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BookManagementPage;
