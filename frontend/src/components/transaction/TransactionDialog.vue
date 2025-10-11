<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useCategoryStore } from '@/stores/category';
import BaseDialog from '@/components/common/BaseDialog.vue';
import type { Transaction, CreateTransactionRequest } from '@/types';

interface Props {
  show: boolean;
  transaction?: Transaction | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  close: [];
  save: [data: CreateTransactionRequest];
}>();

const categoryStore = useCategoryStore();
const { categories } = categoryStore;

const form = ref<CreateTransactionRequest>({
  type: 'expense',
  amount: 0,
  category_id: 0,
  transaction_date: new Date().toISOString().split('T')[0],
  memo: '',
});

const dialogTitle = computed(() => (props.transaction ? '取引編集' : '新規取引'));

const filteredCategories = computed(() => {
  return categories.filter(c => c.type === form.value.type);
});

const isFormValid = computed(() => {
  return (
    form.value.type &&
    form.value.amount > 0 &&
    form.value.category_id > 0 &&
    form.value.transaction_date
  );
});

const handleSubmit = () => {
  if (isFormValid.value) {
    emit('save', { ...form.value });
  }
};

watch(
  () => props.transaction,
  newTransaction => {
    if (newTransaction) {
      form.value = {
        type: newTransaction.type,
        amount: newTransaction.amount,
        category_id: newTransaction.category_id,
        transaction_date: newTransaction.transaction_date.split('T')[0],
        memo: newTransaction.memo || '',
      };
    } else {
      form.value = {
        type: 'expense',
        amount: 0,
        category_id: 0,
        transaction_date: new Date().toISOString().split('T')[0],
        memo: '',
      };
    }
  },
  { immediate: true }
);

watch(
  () => form.value.type,
  () => {
    form.value.category_id = 0;
  }
);

onMounted(() => {
  categoryStore.fetchCategories();
});
</script>

<template>
  <BaseDialog :show="show" :title="dialogTitle" @close="emit('close')">
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label class="form-label">種別</label>
        <select v-model="form.type" class="form-input">
          <option value="">選択してください</option>
          <option value="income">収入</option>
          <option value="expense">支出</option>
        </select>
      </div>

      <div class="form-group">
        <label class="form-label">カテゴリ</label>
        <select v-model="form.category_id" class="form-input">
          <option :value="0">選択してください</option>
          <option v-for="category in filteredCategories" :key="category.id" :value="category.id">
            {{ category.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label class="form-label">金額</label>
        <input v-model.number="form.amount" type="number" class="form-input" min="0" step="1" />
      </div>

      <div class="form-group">
        <label class="form-label">日付</label>
        <input v-model="form.transaction_date" type="date" class="form-input" />
      </div>

      <div class="form-group">
        <label class="form-label">メモ</label>
        <textarea
          v-model="form.memo"
          class="form-input"
          rows="3"
          placeholder="メモ（任意）"
        ></textarea>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" @click="emit('close')">キャンセル</button>
        <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
          {{ transaction ? '更新' : '作成' }}
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
  gap: 0.75rem;
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}
</style>
