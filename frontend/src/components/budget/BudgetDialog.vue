<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useCategoryStore } from '@/stores/category';
import BaseDialog from '@/components/common/BaseDialog.vue';
import type { Budget, CreateBudgetRequest } from '@/types';

interface Props {
  show: boolean;
  budget?: Budget | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  close: [];
  save: [data: CreateBudgetRequest];
}>();

const categoryStore = useCategoryStore();
const { categories } = categoryStore;

const form = ref<CreateBudgetRequest>({
  category_id: 0,
  amount: 0,
  target_year: new Date().getFullYear(),
  target_month: new Date().getMonth() + 1,
});

const months = [
  { value: 1, label: '1月' },
  { value: 2, label: '2月' },
  { value: 3, label: '3月' },
  { value: 4, label: '4月' },
  { value: 5, label: '5月' },
  { value: 6, label: '6月' },
  { value: 7, label: '7月' },
  { value: 8, label: '8月' },
  { value: 9, label: '9月' },
  { value: 10, label: '10月' },
  { value: 11, label: '11月' },
  { value: 12, label: '12月' },
];

const expenseCategories = computed(() => categories.filter(c => c.type === 'expense'));

const dialogTitle = computed(() => (props.budget ? '予算編集' : '新規予算'));

const isFormValid = computed(() => {
  return (
    form.value.category_id > 0 &&
    form.value.amount > 0 &&
    form.value.target_year >= 1900 &&
    form.value.target_year <= 2100 &&
    form.value.target_month >= 1 &&
    form.value.target_month <= 12
  );
});

const handleSubmit = () => {
  if (isFormValid.value) {
    emit('save', { ...form.value });
  }
};

watch(
  () => props.budget,
  newBudget => {
    if (newBudget) {
      form.value = {
        category_id: newBudget.category_id,
        amount: newBudget.amount,
        target_year: newBudget.target_year,
        target_month: newBudget.target_month,
      };
    } else {
      form.value = {
        category_id: 0,
        amount: 0,
        target_year: new Date().getFullYear(),
        target_month: new Date().getMonth() + 1,
      };
    }
  },
  { immediate: true }
);

onMounted(() => {
  categoryStore.fetchCategories();
});
</script>

<template>
  <BaseDialog :show="show" :title="dialogTitle" @close="emit('close')">
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label class="form-label">カテゴリ</label>
        <select v-model="form.category_id" class="form-input" required>
          <option value="">選択してください</option>
          <option v-for="category in expenseCategories" :key="category.id" :value="category.id">
            {{ category.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label class="form-label">予算額</label>
        <input
          v-model.number="form.amount"
          type="number"
          class="form-input"
          min="0.01"
          step="0.01"
          required
        />
      </div>

      <div class="form-group">
        <label class="form-label">対象年</label>
        <input
          v-model.number="form.target_year"
          type="number"
          class="form-input"
          min="1900"
          max="2100"
          required
        />
      </div>

      <div class="form-group">
        <label class="form-label">対象月</label>
        <select v-model="form.target_month" class="form-input" required>
          <option value="">選択してください</option>
          <option v-for="month in months" :key="month.value" :value="month.value">
            {{ month.label }}
          </option>
        </select>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" @click="emit('close')">キャンセル</button>
        <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
          {{ budget ? '更新' : '作成' }}
        </button>
      </div>
    </form>
  </BaseDialog>
</template>

<style scoped>
.form-group {
  margin-bottom: 1.25rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
  font-size: 0.95rem;
}

.form-input {
  width: 100%;
  padding: 0.625rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
  transition: border-color 0.2s;
  box-sizing: border-box;
  color: #333;
  background-color: #fff;
  line-height: 1.5;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}
</style>
