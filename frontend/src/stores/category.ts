import { defineStore } from 'pinia';
import { ref } from 'vue';
import { categoryApi } from '../services/api';
import type { Category, CreateCategoryRequest } from '../types';

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchCategories = async () => {
    loading.value = true;
    error.value = null;
    try {
      const response = await categoryApi.getAll();
      categories.value = response.data;
    } catch (err) {
      error.value = 'Failed to fetch categories';
      console.error('Error fetching categories:', err);
    } finally {
      loading.value = false;
    }
  };

  const fetchCategoriesByType = async (type: 'income' | 'expense') => {
    loading.value = true;
    error.value = null;
    try {
      const response = await categoryApi.getByType(type);
      return response.data;
    } catch (err) {
      error.value = 'Failed to fetch categories';
      console.error('Error fetching categories by type:', err);
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
      categories.value.push(response.data);
      return response.data;
    } catch (err) {
      error.value = 'Failed to create category';
      console.error('Error creating category:', err);
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
      const index = categories.value.findIndex(c => c.id === id);
      if (index !== -1) {
        categories.value[index] = response.data;
      }
      return response.data;
    } catch (err) {
      error.value = 'Failed to update category';
      console.error('Error updating category:', err);
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
      categories.value = categories.value.filter(c => c.id !== id);
    } catch (err) {
      error.value = 'Failed to delete category';
      console.error('Error deleting category:', err);
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
