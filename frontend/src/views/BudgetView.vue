<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { budgetApi } from '@/services/api';
import BudgetDialog from '@/components/budget/BudgetDialog.vue';
import type { Budget, CreateBudgetRequest } from '@/types';

const budgets = ref<Budget[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const editingBudget = ref<Budget | null>(null);

const formatNumber = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(num);
};

const fetchBudgets = async () => {
  loading.value = true;
  error.value = null;
  try {
    const response = await budgetApi.getAll();
    budgets.value = response.data;
  } catch (err) {
    error.value = '予算の取得に失敗しました';
  } finally {
    loading.value = false;
  }
};

const editBudget = (budget: Budget) => {
  editingBudget.value = budget;
  showEditDialog.value = true;
};

const deleteBudget = async (id: number) => {
  if (confirm('この予算を削除しますか？')) {
    try {
      await budgetApi.delete(id);
      budgets.value = budgets.value.filter(b => b.id !== id);
    } catch (err) {
      alert('予算の削除に失敗しました');
    }
  }
};

const closeDialog = () => {
  showCreateDialog.value = false;
  showEditDialog.value = false;
  editingBudget.value = null;
};

const handleSave = async (data: CreateBudgetRequest) => {
  try {
    if (editingBudget.value) {
      const response = await budgetApi.update(editingBudget.value.id, data);
      const index = budgets.value.findIndex(b => b.id === editingBudget.value!.id);
      if (index !== -1) {
        budgets.value[index] = response.data;
      }
    } else {
      const response = await budgetApi.create(data);
      budgets.value.unshift(response.data);
    }
    closeDialog();
  } catch (err) {
    alert('予算の保存に失敗しました');
  }
};

onMounted(() => {
  fetchBudgets();
});
</script>

<template>
  <div class="budgets">
    <h2>予算管理</h2>

    <div class="budgets-grid">
      <div class="card">
        <div class="card-header">
          <h3>予算一覧</h3>
          <button class="btn btn-primary" @click="showCreateDialog = true">新規予算</button>
        </div>
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="budgets.length === 0">予算が設定されていません</div>
        <div v-else>
          <table class="table">
            <thead>
              <tr>
                <th>期間</th>
                <th>カテゴリ</th>
                <th>予算額</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="budget in budgets" :key="budget.id">
                <td>{{ budget.target_year }}年{{ budget.target_month }}月</td>
                <td>{{ budget.category?.name }}</td>
                <td>¥{{ formatNumber(budget.amount) }}</td>
                <td>
                  <button class="btn btn-secondary" @click="editBudget(budget)">編集</button>
                  <button class="btn btn-danger" @click="deleteBudget(budget.id)">削除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <BudgetDialog
      v-if="showCreateDialog || showEditDialog"
      :show="showCreateDialog || showEditDialog"
      :budget="editingBudget"
      @close="closeDialog"
      @save="handleSave"
    />
  </div>
</template>

<style scoped>
.budgets {
  max-width: 1200px;
  margin: 0 auto;
}

.budgets h2 {
  margin-bottom: 2rem;
  color: #333;
}

.budgets-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.budgets-grid .card {
  grid-column: 1 / -1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.card-header h3 {
  margin: 0;
}

.table td .btn {
  margin-right: 0.5rem;
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
}

@media (max-width: 768px) {
  .budgets-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
}
</style>
