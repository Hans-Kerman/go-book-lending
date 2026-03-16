// src/pages/Admin/BookManagement.tsx
import React from 'react';
import { Typography } from 'antd';

const BookManagementPage: React.FC = () => {
  return (
    <div>
      <Typography.Title level={2}>图书管理 (管理员)</Typography.Title>
      <Typography.Paragraph>
        欢迎来到图书管理后台。在这里，您可以添加、编辑和删除系统中的所有书籍。
      </Typography.Paragraph>
      {/* 后续将在这里实现图书的增删改查功能 */}
    </div>
  );
};

export default BookManagementPage;
