<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useCategoryStore } from '@/stores/category';
import CategoryDialog from '@/components/category/CategoryDialog.vue';
import type { Category, CreateCategoryRequest } from '@/types';

const categoryStore = useCategoryStore();
const { loading, error } = categoryStore;

const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const editingCategory = ref<Category | null>(null);

const incomeCategories = computed(() => categoryStore.getIncomeCategories());
const expenseCategories = computed(() => categoryStore.getExpenseCategories());

const editCategory = (category: Category) => {
  editingCategory.value = category;
  showEditDialog.value = true;
};

const deleteCategory = async (id: number) => {
  if (confirm('このカテゴリを削除しますか？関連する取引がある場合は削除できません。')) {
    try {
      await categoryStore.deleteCategory(id);
    } catch (err) {
      alert('カテゴリの削除に失敗しました。関連する取引が存在する可能性があります。');
    }
  }
};

const closeDialog = () => {
  showCreateDialog.value = false;
  showEditDialog.value = false;
  editingCategory.value = null;
};

const handleSave = async (data: CreateCategoryRequest) => {
  try {
    if (editingCategory.value) {
      await categoryStore.updateCategory(editingCategory.value.id, data);
    } else {
      await categoryStore.createCategory(data);
    }
    closeDialog();
  } catch (err) {
    alert('カテゴリの保存に失敗しました。');
  }
};

onMounted(() => {
  categoryStore.fetchCategories();
});
</script>

<template>
  <div class="categories">
    <h2>カテゴリ管理</h2>

    <div class="category-grid">
      <div class="card header-card">
        <div class="card-header">
          <h3>カテゴリ一覧</h3>
          <button class="btn btn-primary" @click="showCreateDialog = true">新規カテゴリ</button>
        </div>
      </div>

      <div class="card">
        <h3>収入カテゴリ</h3>
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="incomeCategories.length === 0">収入カテゴリがありません</div>
        <div v-else class="category-list">
          <div v-for="category in incomeCategories" :key="category.id" class="category-item">
            <div class="category-color" :style="{ backgroundColor: category.color }"></div>
            <span class="category-name">{{ category.name }}</span>
            <div class="category-actions">
              <button class="btn btn-secondary" @click="editCategory(category)">編集</button>
              <button class="btn btn-danger" @click="deleteCategory(category.id)">削除</button>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <h3>支出カテゴリ</h3>
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="expenseCategories.length === 0">支出カテゴリがありません</div>
        <div v-else class="category-list">
          <div v-for="category in expenseCategories" :key="category.id" class="category-item">
            <div class="category-color" :style="{ backgroundColor: category.color }"></div>
            <span class="category-name">{{ category.name }}</span>
            <div class="category-actions">
              <button class="btn btn-secondary" @click="editCategory(category)">編集</button>
              <button class="btn btn-danger" @click="deleteCategory(category.id)">削除</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <CategoryDialog
      v-if="showCreateDialog || showEditDialog"
      :show="showCreateDialog || showEditDialog"
      :category="editingCategory"
      @close="closeDialog"
      @save="handleSave"
    />
  </div>
</template>

<style scoped>
.categories {
  max-width: 1200px;
  margin: 0 auto;
}

.categories h2 {
  margin-bottom: 2rem;
  color: #333;
}

.category-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.header-card {
  grid-column: 1 / -1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0;
}

.card-header h3 {
  margin: 0;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.category-item {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  border: 1px solid #eee;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.category-item:hover {
  background-color: #f8f9fa;
}

.category-color {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  margin-right: 0.75rem;
  border: 1px solid #ddd;
}

.category-name {
  flex: 1;
  color: #000;
  font-weight: 500;
}

.category-actions {
  display: flex;
  gap: 0.5rem;
}

.category-actions .btn {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
}

@media (max-width: 768px) {
  .category-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
}
</style>
