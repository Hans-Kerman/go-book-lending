// src/pages/User/MyBorrows.tsx
import React from 'react';
import { Typography } from 'antd';

const MyBorrowsPage: React.FC = () => {
  return (
    <div>
      <Typography.Title level={2}>我的借阅记录</Typography.Title>
      <Typography.Paragraph>
        这里将显示您当前借阅和历史借阅的所有书籍。
      </Typography.Paragraph>
      {/* 后续将在这里实现从 API 获取并展示借阅列表的功能 */}
    </div>
  );
};

export default MyBorrowsPage;
