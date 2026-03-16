// src/pages/Admin/BookManagement.tsx
import React, { useState, useEffect, useCallback } from 'react';
import { Table, Button, Space, Popconfirm, message, Alert, Modal, Form, Input, InputNumber, Upload } from 'antd';
import type { TablePaginationConfig } from 'antd';
import apiClient from '../../services/api';
import type { Book } from '../../types';

// 定义表单和模态框的状态
type ModalState = {
  visible: boolean;
  mode: 'add' | 'edit';
  book: Book | null;
};

const BookManagementPage: React.FC = () => {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState<TablePaginationConfig>({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [modalState, setModalState] = useState<ModalState>({ visible: false, mode: 'add', book: null });
  const [form] = Form.useForm();

  const fetchBooks = useCallback(async (params: { page: number; pageSize: number }) => {
    setLoading(true);
    setError(null);
    try {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const response = await apiClient.get<any>('/public/books', {
        params: {
          page: params.page,
          page_size: params.pageSize,
        },
      });
      // 兼容后端返回的大写字段名
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const normalizedBooks: Book[] = (response.data.books || []).map((b: any) => ({
        id: b.ID ?? b.id,
        createdAt: b.CreatedAt ?? b.createdAt,
        updatedAt: b.UpdatedAt ?? b.updatedAt,
        isbn: b.ISBN ?? b.isbn,
        title: b.Title ?? b.title,
        author: b.Author ?? b.author,
        coverURL: b.CoverURL ?? b.coverURL,
        available: b.Available ?? b.available,
        price: b.Price ?? b.price,
      }));
      setBooks(normalizedBooks);
      setPagination(p => ({
        ...p,
        total: response.data.total,
        current: response.data.page,
      }));
    } catch (err) {
      const errorMsg = (err as { response?: { data?: { error: string } } }).response?.data?.error || '获取书籍列表失败';
      setError(errorMsg);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchBooks({ page: pagination.current!, pageSize: pagination.pageSize! });
  }, [fetchBooks]);

  // --- CRUD 操作 ---

  const handleDelete = async (isbn: string) => {
    try {
      await apiClient.delete(`/admin/book/del/${isbn}`);
      message.success('删除成功');
      // 如果当前页只剩一项，删除后请求前一页
      if (books.length === 1 && pagination.current! > 1) {
        await fetchBooks({ page: pagination.current! - 1, pageSize: pagination.pageSize! });
      } else {
        await fetchBooks({ page: pagination.current!, pageSize: pagination.pageSize! });
      }
    } catch (err) {
      const errorMsg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '删除失败';
      message.error(errorMsg);
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      
      // 处理图片转 base64
      let coverPicBase64: string | undefined;
      const fileList = values.cover_pic_base64;
      if (fileList && fileList.length > 0) {
        const file = fileList[0].originFileObj;
        if (file) {
          const base64DataUrl = await new Promise<string>((resolve, reject) => {
            const reader = new FileReader();
            reader.readAsDataURL(file as File);
            reader.onload = () => resolve(reader.result as string);
            reader.onerror = (error) => reject(error);
          });
          coverPicBase64 = base64DataUrl.split(',')[1];
        }
      }

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const payload: any = {
        isbn: values.isbn,
        title: values.title,
        author: values.author,
        available: values.available,
        price: values.price,
      };
      
      if (coverPicBase64) {
        payload.cover_pic_base64 = Array.from(new Uint8Array(
          atob(coverPicBase64).split('').map(char => char.charCodeAt(0))
        ));
      }

      if (modalState.mode === 'add') {
        await apiClient.post('/admin/book', payload);
        message.success('添加成功');
      } else if (modalState.mode === 'edit' && modalState.book) {
        // 后端更新接口需要ISBN在body中，我们确保它存在
        await apiClient.put('/admin/book', { ...payload, isbn: modalState.book.isbn });
        message.success('更新成功');
      }
      
      setModalState({ visible: false, mode: 'add', book: null });
      // 刷新当前页
      fetchBooks({ page: pagination.current!, pageSize: pagination.pageSize! });
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


  const handleTableChange = (newPagination: TablePaginationConfig) => {
    fetchBooks({ page: newPagination.current!, pageSize: newPagination.pageSize! });
  };

  const columns = [
    { title: 'ISBN', dataIndex: 'isbn', key: 'isbn' },
    { title: '书名', dataIndex: 'title', key: 'title' },
    { title: '作者', dataIndex: 'author', key: 'author' },
    { title: '库存', dataIndex: 'available', key: 'available' },
    { title: '赔偿金额', dataIndex: 'price', key: 'price', render: (val: number) => val ? (val / 100).toFixed(2) : '0.00' },
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

  if (error && !loading) return <Alert message="Error" description={error} type="error" showIcon />;

  return (
    <div>
      <Button type="primary" onClick={showAddModal} style={{ marginBottom: 16 }}>
        添加新书
      </Button>
      <Table
        columns={columns}
        dataSource={books}
        rowKey="isbn"
        pagination={pagination}
        loading={loading}
        onChange={handleTableChange}
      />

      <Modal
        title={modalState.mode === 'add' ? '添加新书' : '编辑书籍'}
        open={modalState.visible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={false} // 可以根据API请求状态来设置
        forceRender // 确保 Form 在 Modal 第一次打开前就渲染并与 useForm() 实例绑定
      >
        <Form form={form} layout="vertical" name="book_form" initialValues={{ price: 0, available: 0 }}>
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
          <Form.Item name="price" label="赔偿金额(分)" rules={[{ type: 'number', min: 0 }]}>
            <InputNumber style={{ width: '100%' }} placeholder="金额" />
          </Form.Item>
          <Form.Item name="cover_pic_base64" label="封面图片" valuePropName="fileList" getValueFromEvent={(e) => {
            if (Array.isArray(e)) return e;
            return e?.fileList;
          }}>
            <Upload beforeUpload={() => false} maxCount={1} listType="picture" accept="image/*">
              <Button>点击选择图片</Button>
            </Upload>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BookManagementPage;
