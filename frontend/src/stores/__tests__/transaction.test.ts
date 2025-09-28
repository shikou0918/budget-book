import { describe, test, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useTransactionStore } from '../transaction';
import { transactionApi } from '@/services/api';
import type { Transaction, CreateTransactionRequest } from '@/types';
import type { AxiosResponse } from 'axios';

// APIをモック
vi.mock('@/services/api', () => ({
  transactionApi: {
    getAll: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
}));

const mockTransactionApi = vi.mocked(transactionApi);

const mockTransaction: Transaction = {
  id: 1,
  transaction_date: '2024-01-15',
  type: 'income',
  amount: 50000,
  memo: '給与',
  category: {
    id: 1,
    name: '給与',
    type: 'income',
    color: '#4CAF50',
    created_at: '2024-01-15T00:00:00Z',
    updated_at: '2024-01-15T00:00:00Z',
  },
  category_id: 1,
  created_at: '2024-01-15T00:00:00Z',
  updated_at: '2024-01-15T00:00:00Z',
};

const mockCreateRequest: CreateTransactionRequest = {
  transaction_date: '2024-01-15',
  type: 'income',
  amount: 50000,
  memo: '給与',
  category_id: 1,
};

describe('取引ストア', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  test('初期状態が空で初期化される', () => {
    const store = useTransactionStore();

    expect(store.transactions).toEqual([]);
    expect(store.loading).toBe(false);
    expect(store.error).toBe(null);
  });

  describe('取引データ取得', () => {
    test('取引データを正常に取得する', async () => {
      const mockTransactions = [mockTransaction];
      mockTransactionApi.getAll.mockResolvedValue({
        data: mockTransactions,
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });

      const store = useTransactionStore();
      await store.fetchTransactions();

      expect(store.loading).toBe(false);
      expect(store.error).toBe(null);
      expect(store.transactions).toEqual(mockTransactions);
      expect(mockTransactionApi.getAll).toHaveBeenCalledOnce();
    });

    test('取得エラーを処理する', async () => {
      mockTransactionApi.getAll.mockRejectedValue(new Error('Network error'));

      const store = useTransactionStore();
      await store.fetchTransactions();

      expect(store.loading).toBe(false);
      expect(store.error).toBe('Failed to fetch transactions');
      expect(store.transactions).toEqual([]);
    });

    test('取得中にローディング状態を設定する', async () => {
      let resolvePromise: (value: AxiosResponse<Transaction[]>) => void;
      const promise = new Promise<AxiosResponse<Transaction[]>>(resolve => {
        resolvePromise = resolve;
      });
      mockTransactionApi.getAll.mockReturnValue(promise);

      const store = useTransactionStore();
      const fetchPromise = store.fetchTransactions();

      expect(store.loading).toBe(true);

      resolvePromise!({
        data: [mockTransaction],
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });
      await fetchPromise;

      expect(store.loading).toBe(false);
    });
  });

  describe('取引作成', () => {
    test('取引を正常に作成する', async () => {
      mockTransactionApi.create.mockResolvedValue({
        data: mockTransaction,
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });

      const store = useTransactionStore();
      const result = await store.createTransaction(mockCreateRequest);

      expect(store.loading).toBe(false);
      expect(store.error).toBe(null);
      expect(store.transactions).toEqual([mockTransaction]);
      expect(result).toEqual(mockTransaction);
      expect(mockTransactionApi.create).toHaveBeenCalledWith(mockCreateRequest);
    });

    test('作成エラーを処理する', async () => {
      mockTransactionApi.create.mockRejectedValue(new Error('Create failed'));

      const store = useTransactionStore();

      await expect(store.createTransaction(mockCreateRequest)).rejects.toThrow();
      expect(store.loading).toBe(false);
      expect(store.error).toBe('Failed to create transaction');
      expect(store.transactions).toEqual([]);
    });

    test('新しい取引をリストの先頭に追加する', async () => {
      const existingTransaction = { ...mockTransaction, id: 2 };
      const newTransaction = { ...mockTransaction, id: 3 };

      mockTransactionApi.create.mockResolvedValue({
        data: newTransaction,
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });

      const store = useTransactionStore();
      store.transactions = [existingTransaction];

      await store.createTransaction(mockCreateRequest);

      expect(store.transactions).toEqual([newTransaction, existingTransaction]);
    });
  });

  describe('取引更新', () => {
    test('取引を正常に更新する', async () => {
      const updatedTransaction = { ...mockTransaction, amount: 60000 };
      mockTransactionApi.update.mockResolvedValue({
        data: updatedTransaction,
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });

      const store = useTransactionStore();
      store.transactions = [mockTransaction];

      const result = await store.updateTransaction(1, mockCreateRequest);

      expect(store.loading).toBe(false);
      expect(store.error).toBe(null);
      expect(store.transactions[0]).toEqual(updatedTransaction);
      expect(result).toEqual(updatedTransaction);
      expect(mockTransactionApi.update).toHaveBeenCalledWith(1, mockCreateRequest);
    });

    test('更新エラーを処理する', async () => {
      mockTransactionApi.update.mockRejectedValue(new Error('Update failed'));

      const store = useTransactionStore();
      store.transactions = [mockTransaction];

      await expect(store.updateTransaction(1, mockCreateRequest)).rejects.toThrow();
      expect(store.loading).toBe(false);
      expect(store.error).toBe('Failed to update transaction');
      expect(store.transactions).toEqual([mockTransaction]); // unchanged
    });

    test('取引が見つからない場合は更新しない', async () => {
      const updatedTransaction = { ...mockTransaction, id: 999 };
      mockTransactionApi.update.mockResolvedValue({
        data: updatedTransaction,
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });

      const store = useTransactionStore();
      store.transactions = [mockTransaction];

      await store.updateTransaction(999, mockCreateRequest);

      expect(store.transactions).toEqual([mockTransaction]); // unchanged
    });
  });

  describe('取引削除', () => {
    test('取引を正常に削除する', async () => {
      mockTransactionApi.delete.mockResolvedValue({
        data: {},
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });

      const store = useTransactionStore();
      store.transactions = [mockTransaction];

      await store.deleteTransaction(1);

      expect(store.loading).toBe(false);
      expect(store.error).toBe(null);
      expect(store.transactions).toEqual([]);
      expect(mockTransactionApi.delete).toHaveBeenCalledWith(1);
    });

    test('削除エラーを処理する', async () => {
      mockTransactionApi.delete.mockRejectedValue(new Error('Delete failed'));

      const store = useTransactionStore();
      store.transactions = [mockTransaction];

      // TODO(human): Fix the error handling test
      // The deleteTransaction function should catch errors internally and set the error state
      // rather than throwing the error. Update this test to properly await the function
      // without expecting it to throw
      await store.deleteTransaction(1);

      expect(store.loading).toBe(false);
      expect(store.error).toBe('Failed to delete transaction');
      expect(store.transactions).toEqual([mockTransaction]); // unchanged
    });

    test('指定された取引のみを削除する', async () => {
      const transaction1 = { ...mockTransaction, id: 1 };
      const transaction2 = { ...mockTransaction, id: 2 };
      mockTransactionApi.delete.mockResolvedValue({
        data: {},
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      });

      const store = useTransactionStore();
      store.transactions = [transaction1, transaction2];

      await store.deleteTransaction(1);

      expect(store.transactions).toEqual([transaction2]);
    });
  });
});
