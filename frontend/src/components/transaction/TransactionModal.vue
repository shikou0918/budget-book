<template>
  <div v-if="show" class="modal-overlay" @click="handleOverlayClick">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h3>{{ transaction ? '取引編集' : '新規取引' }}</h3>
        <button class="modal-close" @click="$emit('close')">&times;</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-group">
          <label class="form-label">種別</label>
          <select v-model="form.type" class="form-input" required>
            <option value="">選択してください</option>
            <option value="income">収入</option>
            <option value="expense">支出</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">カテゴリ</label>
          <select v-model="form.category_id" class="form-input" required>
            <option value="">選択してください</option>
            <option v-for="category in filteredCategories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">金額</label>
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
          <label class="form-label">日付</label>
          <input v-model="form.transaction_date" type="date" class="form-input" required />
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
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            キャンセル
          </button>
          <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
            {{ transaction ? '更新' : '作成' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, watch, onMounted } from 'vue';
  import { useCategoryStore } from '../../stores/category';
  import type { Transaction, CreateTransactionRequest } from '../../types';

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

  const handleOverlayClick = (e: Event) => {
    if (e.target === e.currentTarget) {
      emit('close');
    }
  };

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

<style scoped>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal {
    background: white;
    border-radius: 8px;
    width: 90%;
    max-width: 500px;
    max-height: 90vh;
    overflow-y: auto;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 1.5rem 0 1.5rem;
    border-bottom: 1px solid #eee;
  }

  .modal-header h3 {
    margin: 0;
    color: #333;
  }

  .modal-close {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #999;
    padding: 0;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .modal-close:hover {
    color: #333;
  }

  .modal-body {
    padding: 1.5rem;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1.5rem;
  }

  @media (max-width: 768px) {
    .modal {
      width: 95%;
      margin: 1rem;
    }

    .modal-header,
    .modal-body {
      padding: 1rem;
    }
  }
</style>
