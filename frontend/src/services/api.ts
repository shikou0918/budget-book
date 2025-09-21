import axios from 'axios';
import type {
  Transaction,
  Category,
  Budget,
  MonthlySummary,
  CreateTransactionRequest,
  CreateCategoryRequest,
  CreateBudgetRequest,
} from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const transactionApi = {
  getAll: () => api.get<Transaction[]>('/transactions'),
  getById: (id: number) => api.get<Transaction>(`/transactions/${id}`),
  create: (data: CreateTransactionRequest) => api.post<Transaction>('/transactions', data),
  update: (id: number, data: CreateTransactionRequest) =>
    api.put<Transaction>(`/transactions/${id}`, data),
  delete: (id: number) => api.delete(`/transactions/${id}`),
};

export const categoryApi = {
  getAll: () => api.get<Category[]>('/categories'),
  getByType: (type: 'income' | 'expense') => api.get<Category[]>(`/categories?type=${type}`),
  getById: (id: number) => api.get<Category>(`/categories/${id}`),
  create: (data: CreateCategoryRequest) => api.post<Category>('/categories', data),
  update: (id: number, data: CreateCategoryRequest) => api.put<Category>(`/categories/${id}`, data),
  delete: (id: number) => api.delete(`/categories/${id}`),
};

export const budgetApi = {
  getAll: () => api.get<Budget[]>('/budgets'),
  getByMonth: (year: number, month: number) =>
    api.get<Budget[]>(`/budgets?year=${year}&month=${month}`),
  getById: (id: number) => api.get<Budget>(`/budgets/${id}`),
  create: (data: CreateBudgetRequest) => api.post<Budget>('/budgets', data),
  update: (id: number, data: CreateBudgetRequest) => api.put<Budget>(`/budgets/${id}`, data),
  delete: (id: number) => api.delete(`/budgets/${id}`),
};

export const summaryApi = {
  getMonthly: (year: number, month: number) => api.get<MonthlySummary>(`/summary/${year}/${month}`),
};

export default api;
