import axios, { AxiosError } from 'axios';
import type { AxiosInstance } from 'axios';
import type {
  Transaction,
  Category,
  Budget,
  MonthlySummary,
  CreateTransactionRequest,
  CreateCategoryRequest,
  CreateBudgetRequest,
} from '@/types';
import { ApplicationError as AppError } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Timeout and retry configuration
const REQUEST_TIMEOUT = 10000; // 10 seconds
const MAX_RETRIES = 3;
const RETRY_DELAY = 1000; // 1 second

const api: AxiosInstance = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  timeout: REQUEST_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Exponential backoff retry logic
const wait = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// Error handler utility
const handleApiError = (error: unknown): never => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError;
    const statusCode = axiosError.response?.status;
    const message = axiosError.response?.data
      ? JSON.stringify(axiosError.response.data)
      : axiosError.message;

    // Log error for debugging (in production, send to monitoring service)
    console.error('[API Error]', {
      statusCode,
      message,
      url: axiosError.config?.url,
      method: axiosError.config?.method,
      timestamp: new Date().toISOString(),
    });

    throw new AppError(message || 'ネットワークエラーが発生しました', statusCode, axiosError);
  }

  // Unknown error
  console.error('[Unknown Error]', error);
  throw new AppError('予期しないエラーが発生しました', undefined, error);
};

// Retry interceptor with exponential backoff
api.interceptors.response.use(
  response => response,
  async error => {
    const config = error.config;

    // Don't retry if not configured or max retries exceeded
    if (!config || !config.retry) {
      return Promise.reject(error);
    }

    config.retryCount = config.retryCount || 0;

    // Check if we should retry
    const shouldRetry =
      config.retryCount < MAX_RETRIES &&
      // Retry on network errors or 5xx server errors
      (!error.response || (error.response.status >= 500 && error.response.status < 600));

    if (!shouldRetry) {
      return Promise.reject(error);
    }

    config.retryCount += 1;

    // Exponential backoff
    const delayMs = RETRY_DELAY * Math.pow(2, config.retryCount - 1);
    console.log(`[Retry] Attempt ${config.retryCount}/${MAX_RETRIES} after ${delayMs}ms`);

    await wait(delayMs);
    return api(config);
  }
);

// Input validation utilities
const validatePositiveNumber = (value: number, fieldName: string): void => {
  if (value <= 0) {
    throw new AppError(`${fieldName}は正の数である必要があります`);
  }
};

const validateId = (id: number): void => {
  if (!Number.isInteger(id) || id <= 0) {
    throw new AppError('無効なIDです');
  }
};

const validateTransactionRequest = (data: CreateTransactionRequest): void => {
  if (!data.transaction_date) {
    throw new AppError('取引日は必須です');
  }
  if (!data.type || !['income', 'expense'].includes(data.type)) {
    throw new AppError('取引タイプは収入または支出である必要があります');
  }
  validatePositiveNumber(data.amount, '金額');
  validateId(data.category_id);
};

const validateCategoryRequest = (data: CreateCategoryRequest): void => {
  if (!data.name || data.name.trim().length === 0) {
    throw new AppError('カテゴリ名は必須です');
  }
  if (!data.type || !['income', 'expense'].includes(data.type)) {
    throw new AppError('カテゴリタイプは収入または支出である必要があります');
  }
};

const validateBudgetRequest = (data: CreateBudgetRequest): void => {
  validateId(data.category_id);
  validatePositiveNumber(data.amount, '予算金額');
  if (data.target_year < 2000 || data.target_year > 2100) {
    throw new AppError('対象年が無効です');
  }
  if (data.target_month < 1 || data.target_month > 12) {
    throw new AppError('対象月は1〜12の範囲である必要があります');
  }
};

// API methods with error handling and retry support
export const transactionApi = {
  getAll: async () => {
    try {
      return await api.get<Transaction[]>('/transactions', { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  getById: async (id: number) => {
    try {
      validateId(id);
      return await api.get<Transaction>(`/transactions/${id}`, { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  create: async (data: CreateTransactionRequest) => {
    try {
      validateTransactionRequest(data);
      return await api.post<Transaction>('/transactions', data, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  update: async (id: number, data: CreateTransactionRequest) => {
    try {
      validateId(id);
      validateTransactionRequest(data);
      return await api.put<Transaction>(`/transactions/${id}`, data, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  delete: async (id: number) => {
    try {
      validateId(id);
      return await api.delete(`/transactions/${id}`, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
};

export const categoryApi = {
  getAll: async () => {
    try {
      return await api.get<Category[]>('/categories', { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  getByType: async (type: 'income' | 'expense') => {
    try {
      if (!['income', 'expense'].includes(type)) {
        throw new AppError('カテゴリタイプは収入または支出である必要があります');
      }
      return await api.get<Category[]>(`/categories?type=${type}`, { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  getById: async (id: number) => {
    try {
      validateId(id);
      return await api.get<Category>(`/categories/${id}`, { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  create: async (data: CreateCategoryRequest) => {
    try {
      validateCategoryRequest(data);
      return await api.post<Category>('/categories', data, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  update: async (id: number, data: CreateCategoryRequest) => {
    try {
      validateId(id);
      validateCategoryRequest(data);
      return await api.put<Category>(`/categories/${id}`, data, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  delete: async (id: number) => {
    try {
      validateId(id);
      return await api.delete(`/categories/${id}`, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
};

export const budgetApi = {
  getAll: async () => {
    try {
      return await api.get<Budget[]>('/budgets', { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  getByMonth: async (year: number, month: number) => {
    try {
      if (year < 2000 || year > 2100) {
        throw new AppError('対象年が無効です');
      }
      if (month < 1 || month > 12) {
        throw new AppError('対象月は1〜12の範囲である必要があります');
      }
      return await api.get<Budget[]>(`/budgets?year=${year}&month=${month}`, {
        retry: true,
      } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  getById: async (id: number) => {
    try {
      validateId(id);
      return await api.get<Budget>(`/budgets/${id}`, { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  create: async (data: CreateBudgetRequest) => {
    try {
      validateBudgetRequest(data);
      return await api.post<Budget>('/budgets', data, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  update: async (id: number, data: CreateBudgetRequest) => {
    try {
      validateId(id);
      validateBudgetRequest(data);
      return await api.put<Budget>(`/budgets/${id}`, data, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
  delete: async (id: number) => {
    try {
      validateId(id);
      return await api.delete(`/budgets/${id}`, { retry: false } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
};

export const summaryApi = {
  getMonthly: async (year: number, month: number) => {
    try {
      if (year < 2000 || year > 2100) {
        throw new AppError('対象年が無効です');
      }
      if (month < 1 || month > 12) {
        throw new AppError('対象月は1〜12の範囲である必要があります');
      }
      return await api.get<MonthlySummary>(`/summary/${year}/${month}`, { retry: true } as any);
    } catch (error) {
      handleApiError(error);
    }
  },
};

export default api;
