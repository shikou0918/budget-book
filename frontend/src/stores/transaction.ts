import { defineStore } from 'pinia';
import { ref } from 'vue';
import { transactionApi } from '@/services/api';
import type { Transaction, CreateTransactionRequest } from '@/types';
import { ApplicationError } from '@/types';

export const useTransactionStore = defineStore('transaction', () => {
  const transactions = ref<Transaction[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Helper function to extract error message
  const getErrorMessage = (err: unknown): string => {
    if (err instanceof ApplicationError) {
      return err.message;
    }
    if (err instanceof Error) {
      return err.message;
    }
    return '予期しないエラーが発生しました';
  };

  const fetchTransactions = async () => {
    loading.value = true;
    error.value = null;
    try {
      const response = await transactionApi.getAll();
      if (!response) {
        throw new ApplicationError('取引データの取得に失敗しました');
      }
      transactions.value = response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `取引データの取得に失敗しました: ${errorMessage}`;
      console.error('[Store Error] fetchTransactions:', err);
      // Don't throw - allow UI to display error message
    } finally {
      loading.value = false;
    }
  };

  const createTransaction = async (data: CreateTransactionRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await transactionApi.create(data);
      if (!response) {
        throw new ApplicationError('取引の作成に失敗しました');
      }
      transactions.value.unshift(response.data);
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `取引の作成に失敗しました: ${errorMessage}`;
      console.error('[Store Error] createTransaction:', err);
      throw err; // Re-throw to notify caller
    } finally {
      loading.value = false;
    }
  };

  const updateTransaction = async (id: number, data: CreateTransactionRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await transactionApi.update(id, data);
      if (!response) {
        throw new ApplicationError('取引の更新に失敗しました');
      }
      const index = transactions.value.findIndex(t => t.id === id);
      if (index !== -1) {
        transactions.value[index] = response.data;
      } else {
        console.warn(`[Store Warning] Transaction with id ${id} not found in local state`);
      }
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `取引の更新に失敗しました: ${errorMessage}`;
      console.error('[Store Error] updateTransaction:', err);
      throw err; // Re-throw to notify caller
    } finally {
      loading.value = false;
    }
  };

  const deleteTransaction = async (id: number) => {
    loading.value = true;
    error.value = null;
    try {
      await transactionApi.delete(id);
      const previousLength = transactions.value.length;
      transactions.value = transactions.value.filter(t => t.id !== id);

      // Verify deletion happened in local state
      if (transactions.value.length === previousLength) {
        console.warn(`[Store Warning] Transaction with id ${id} was not found in local state`);
      }
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `取引の削除に失敗しました: ${errorMessage}`;
      console.error('[Store Error] deleteTransaction:', err);
      throw err; // Re-throw to notify caller
    } finally {
      loading.value = false;
    }
  };

  // Clear error state
  const clearError = () => {
    error.value = null;
  };

  return {
    transactions,
    loading,
    error,
    fetchTransactions,
    createTransaction,
    updateTransaction,
    deleteTransaction,
    clearError,
  };
});
