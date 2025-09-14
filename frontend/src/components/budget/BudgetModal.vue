<template>
  <div v-if="show" class="modal-overlay" @click="handleOverlayClick">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h3>{{ budget ? '予算編集' : '新規予算' }}</h3>
        <button class="modal-close" @click="$emit('close')">&times;</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-group">
          <label class="form-label">カテゴリ</label>
          <select v-model="form.category_id" class="form-input" required>
            <option value="">選択してください</option>
            <option 
              v-for="category in expenseCategories" 
              :key="category.id" 
              :value="category.id"
            >
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
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            キャンセル
          </button>
          <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
            {{ budget ? '更新' : '作成' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useCategoryStore } from '../../stores/category'
import type { Budget, CreateBudgetRequest } from '../../types'

interface Props {
  show: boolean
  budget?: Budget | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  save: [data: CreateBudgetRequest]
}>()

const categoryStore = useCategoryStore()
const { categories } = categoryStore

const form = ref<CreateBudgetRequest>({
  category_id: 0,
  amount: 0,
  target_year: new Date().getFullYear(),
  target_month: new Date().getMonth() + 1
})

const months = [
  { value: 1, label: '1月' }, { value: 2, label: '2月' }, { value: 3, label: '3月' },
  { value: 4, label: '4月' }, { value: 5, label: '5月' }, { value: 6, label: '6月' },
  { value: 7, label: '7月' }, { value: 8, label: '8月' }, { value: 9, label: '9月' },
  { value: 10, label: '10月' }, { value: 11, label: '11月' }, { value: 12, label: '12月' }
]

const expenseCategories = computed(() => categories.filter(c => c.type === 'expense'))

const isFormValid = computed(() => {
  return form.value.category_id > 0 && 
         form.value.amount > 0 && 
         form.value.target_year >= 1900 && 
         form.value.target_year <= 2100 &&
         form.value.target_month >= 1 && 
         form.value.target_month <= 12
})

const handleOverlayClick = (e: Event) => {
  if (e.target === e.currentTarget) {
    emit('close')
  }
}

const handleSubmit = () => {
  if (isFormValid.value) {
    emit('save', { ...form.value })
  }
}

watch(() => props.budget, (newBudget) => {
  if (newBudget) {
    form.value = {
      category_id: newBudget.category_id,
      amount: newBudget.amount,
      target_year: newBudget.target_year,
      target_month: newBudget.target_month
    }
  } else {
    form.value = {
      category_id: 0,
      amount: 0,
      target_year: new Date().getFullYear(),
      target_month: new Date().getMonth() + 1
    }
  }
}, { immediate: true })

onMounted(() => {
  categoryStore.fetchCategories()
})
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