// src/types/index.ts

/**
 * 代表一本书的完整信息，对应后端的 models.Book
 */
export interface Book {
  id: number;
  createdAt: string;
  updatedAt: string;
  isbn: string;
  title: string;
  author: string;
  coverURL: string;
  available: number; // 可借阅数量
  price: number; // 赔偿凭据金额（分）
}

/**
 * 代表一条借阅记录，对应后端的 models.LendRecord
 */
export interface LendRecord {
  id: number;
  createdAt: string;
  bookID: string; // ISBN
  borrowReader: number; // 借阅者 ID
  returnTime: string;
}

/**
 * 借阅记录的响应体格式
 */
export interface LendRecordResponse {
  id: number;
  created_at: string;
  book_id: string;
  borrow_reader: number;
  return_time: string;
}

/**
 * 登录或注册时发送给后端的凭据
 */
export interface AuthCredentials {
  user_name: string;
  password?: string;
}

/**
 * 创建一本新书时需要的信息
 */
export interface NewBookInfo {
  isbn: string;
  title: string;
  author?: string;
  available?: number;
  price?: number;
  // Swagger 文档在此处定义为整数数组，可能与后端实现有关
  // 前端通常处理为 base64 字符串，在此暂定为 string
  cover_pic_base64?: string;
}

/**
 * 借书或还书的请求体
 */
export interface BorrowRequest {
  book_id: string;      // ISBN
  borrow_reader: number; // 用户 ID
}
