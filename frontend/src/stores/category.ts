import { defineStore } from 'pinia';
import { ref } from 'vue';
import { categoryApi } from '@/services/api';
import type { Category, CreateCategoryRequest } from '@/types';
import { ApplicationError } from '@/types';

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([]);
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

  const fetchCategories = async () => {
    loading.value = true;
    error.value = null;
    try {
      const response = await categoryApi.getAll();
      if (!response) {
        throw new ApplicationError('カテゴリデータの取得に失敗しました');
      }
      categories.value = response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `カテゴリデータの取得に失敗しました: ${errorMessage}`;
      console.error('[Store Error] fetchCategories:', err);
    } finally {
      loading.value = false;
    }
  };

  const fetchCategoriesByType = async (type: 'income' | 'expense') => {
    loading.value = true;
    error.value = null;
    try {
      const response = await categoryApi.getByType(type);
      if (!response) {
        throw new ApplicationError('カテゴリデータの取得に失敗しました');
      }
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `カテゴリデータの取得に失敗しました: ${errorMessage}`;
      console.error('[Store Error] fetchCategoriesByType:', err);
      return [];
    } finally {
      loading.value = false;
    }
  };

  const createCategory = async (data: CreateCategoryRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await categoryApi.create(data);
      if (!response) {
        throw new ApplicationError('カテゴリの作成に失敗しました');
      }
      categories.value.push(response.data);
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `カテゴリの作成に失敗しました: ${errorMessage}`;
      console.error('[Store Error] createCategory:', err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateCategory = async (id: number, data: CreateCategoryRequest) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await categoryApi.update(id, data);
      if (!response) {
        throw new ApplicationError('カテゴリの更新に失敗しました');
      }
      const index = categories.value.findIndex(c => c.id === id);
      if (index !== -1) {
        categories.value[index] = response.data;
      } else {
        console.warn(`[Store Warning] Category with id ${id} not found in local state`);
      }
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `カテゴリの更新に失敗しました: ${errorMessage}`;
      console.error('[Store Error] updateCategory:', err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteCategory = async (id: number) => {
    loading.value = true;
    error.value = null;
    try {
      await categoryApi.delete(id);
      const previousLength = categories.value.length;
      categories.value = categories.value.filter(c => c.id !== id);

      // Verify deletion happened in local state
      if (categories.value.length === previousLength) {
        console.warn(`[Store Warning] Category with id ${id} was not found in local state`);
      }
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      error.value = `カテゴリの削除に失敗しました: ${errorMessage}`;
      console.error('[Store Error] deleteCategory:', err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const getIncomeCategories = () => categories.value.filter(c => c.type === 'income');
  const getExpenseCategories = () => categories.value.filter(c => c.type === 'expense');

  return {
    categories,
    loading,
    error,
    fetchCategories,
    fetchCategoriesByType,
    createCategory,
    updateCategory,
    deleteCategory,
    getIncomeCategories,
    getExpenseCategories,
  };
});
