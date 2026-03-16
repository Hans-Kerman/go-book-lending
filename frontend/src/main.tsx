import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
// Import Ant Design's global style reset
import 'antd/dist/reset.css';
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
