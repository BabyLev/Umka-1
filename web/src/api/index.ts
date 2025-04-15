import axios from 'axios';

const api = axios.create({
  baseURL: '/api', // Базовый URL для всех запросов к API
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
});

// Опционально: можно добавить интерцепторы для централизованной обработки
// ответов и ошибок, логирования, добавления токенов авторизации и т.д.
api.interceptors.response.use(
  (response) => {
    // Любой код состояния, находящийся в диапазоне 2xx, вызывает срабатывание этой функции
    // Можно добавить логирование или модификацию данных ответа здесь
    return response;
  },
  (error) => {
    // Любые коды состояния, выходящие за пределы диапазона 2xx, вызывают срабатывание этой функции
    console.error('API Error:', error.response?.data || error.message || error);
    // Можно добавить более специфичную обработку ошибок (например, показ уведомлений)
    // Важно вернуть Promise.reject, чтобы цепочка .catch() в вызывающем коде работала
    return Promise.reject(error);
  }
);

export default api; 