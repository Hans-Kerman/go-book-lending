# 图书借阅系统前端开发计划

本文档详细说明了图书借阅系统前端部分的开发计划，旨在确保项目按计划、高质量地完成。

## 1. 技术选型

*   **构建工具**: [Vite](https://vitejs.dev/)
*   **框架**: [React](https://react.dev/)
*   **语言**: [TypeScript](https://www.typescriptlang.org/)
*   **UI 库**: [Ant Design](https://ant.design/)
*   **路由管理**: [React Router DOM](https://reactrouter.com/)
*   **API 请求**: [Axios](https://axios-http.com/)
*   **状态管理**: [Zustand](https://github.com/pmndrs/zustand)

## 2. 项目目录结构

我们将采用一个模块化、易于维护的目录结构。

```
/
├── public/
│   └── favicon.ico
├── src/
│   ├── assets/
│   │   └── # 存放图片、字体等静态资源
│   ├── components/
│   │   ├── common/ # 可复用的通用组件 (如按钮, 输入框)
│   │   └── layout/ # 布局组件 (Header, Footer, Sider)
│   ├── config/
│   │   └── index.ts # 存放项目配置 (如 API 地址)
│   ├── hooks/
│   │   └── # 自定义 Hooks
│   ├── pages/
│   │   ├── Admin/
│   │   │   └── BookManagement.tsx # 管理员图书管理页
│   │   ├── Auth/
│   │   │   ├── Login.tsx     # 登录页
│   │   │   └── Register.tsx  # 注册页
│   │   ├── Book/
│   │   │   ├── BookDetail.tsx # 图书详情页
│   │   │   └── BookList.tsx   # 图书列表页
│   │   ├── User/
│   │   │   └── MyBorrows.tsx  # 我的借阅页
│   │   └── Home.tsx
│   ├── routes/
│   │   ├── index.tsx         # 路由配置
│   │   └── PrivateRoute.tsx  # 私有路由/路由守卫
│   ├── services/
│   │   └── api.ts            # API 请求封装 (axios 实例和拦截器)
│   ├── store/
│   │   └── userStore.ts      # 用户相关的全局状态 (Zustand)
│   ├── types/
│   │   └── index.ts          # 全局 TypeScript 类型定义
│   ├── utils/
│   │   └── # 工具函数
│   ├── App.tsx
│   ├── main.tsx
│   └── style.css
├── .eslintrc.cjs
├── .gitignore
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 3. 应用核心流程

下图展示了用户和管理员在应用中的主要交互路径。

```mermaid
graph TD
    A[用户访问网站] --> B{是否已登录?};
    B -- 否 --> C[登录/注册页面];
    B -- 是 --> D[图书列表页 (首页)];

    C --> D;

    D --> E[图书详情页];
    E -- 借阅 --> F{触发借阅操作};
    F -- 成功 --> G[跳转至“我的借阅”];
    F -- 失败 --> H[提示错误信息];

    D --> I[用户中心];
    I --> G;

    B -- 是, 且为管理员 --> J[管理员后台入口];
    J --> K[图书管理页面];
    K --> L[增/删/改图书];

```

## 4. 分阶段实施计划

我们将按照以下四个主要阶段来完成开发：

**阶段一：项目基础搭建**
1.  使用 Vite 初始化 React + TypeScript 项目。
2.  安装核心依赖 (`antd`, `react-router-dom`, `axios`, `zustand`)。
3.  根据规划创建项目目录结构。
4.  创建应用的全局布局 (Header, Content, Footer) 及基础路由配置。

**阶段二：核心认证与浏览功能**
1.  实现用户注册和登录页面及功能。
2.  使用 Zustand 创建全局用户状态（存储 Token 和用户信息）。
3.  创建图书列表页面，从后端获取数据并支持分页展示。
4.  创建图书详情页面，展示单本图书的完整信息。

**阶段三：用户交互功能**
1.  实现路由守卫，保护需要登录才能访问的页面（如用户中心）。
2.  创建“我的借阅”页面，展示用户的借阅记录。
3.  在图书详情页实现借书和还书功能，并更新UI状态。

**阶段四：管理员功能与收尾**
1.  创建管理员专用的图书管理页面（增/删/改/查）。
2.  实现管理员权限的路由保护。
3.  为管理员创建添加和编辑图书的表单及相关逻辑。
4.  进行最终的联调测试和代码优化，确保所有功能正常运行。

---

