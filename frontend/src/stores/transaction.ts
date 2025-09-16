import { defineStore } from 'pinia';
import { ref } from 'vue';
import { transactionApi } from '@/services/api';
import type { Transaction, CreateTransactionRequest } from '@/types';

export const useTransactionStore = defineStore('transaction', () => {
  const transactions = ref<Transaction[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchTransactions = async () => {
    loading.value = true;
    error.value = null;
    try {
      const response = await transactionApi.getAll();
      transactions.value = response.data;
    } catch (err) {
      error.value = 'Failed to fetch transactions';
    } finally {
      loading.value = false;
    }
  };

  const createTransaction = async (data: CreateTransactionRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await transactionApi.create(data);
      transactions.value.unshift(response.data);
      return response.data;
    } catch (err) {
      error.value = 'Failed to create transaction';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateTransaction = async (id: number, data: CreateTransactionRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await transactionApi.update(id, data);
      const index = transactions.value.findIndex(t => t.id === id);
      if (index !== -1) {
        transactions.value[index] = response.data;
      }
      return response.data;
    } catch (err) {
      error.value = 'Failed to update transaction';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteTransaction = async (id: number) => {
    loading.value = true;
    error.value = null;
    try {
      await transactionApi.delete(id);
      transactions.value = transactions.value.filter(t => t.id !== id);
    } catch (err) {
      error.value = 'Failed to delete transaction';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    transactions,
    loading,
    error,
    fetchTransactions,
    createTransaction,
    updateTransaction,
    deleteTransaction,
  };
});
