<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { COLOR_PRESETS, DEFAULT_COLOR } from '@/config/colors';
import BaseDialog from '@/components/common/BaseDialog.vue';
import type { Category, CreateCategoryRequest } from '@/types';

interface Props {
  show: boolean;
  category?: Category | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  close: [];
  save: [data: CreateCategoryRequest];
}>();

const form = ref<CreateCategoryRequest>({
  name: '',
  type: 'expense',
  color: DEFAULT_COLOR,
});

const dialogTitle = computed(() => (props.category ? 'カテゴリ編集' : '新規カテゴリ'));

const isFormValid = computed(() => {
  return form.value.name.trim() && form.value.type && form.value.color;
});

const handleSubmit = () => {
  if (isFormValid.value) {
    emit('save', { ...form.value });
  }
};

watch(
  () => props.category,
  newCategory => {
    if (newCategory) {
      form.value = {
        name: newCategory.name,
        type: newCategory.type,
        color: newCategory.color || DEFAULT_COLOR,
      };
    } else {
      form.value = {
        name: '',
        type: 'expense',
        color: DEFAULT_COLOR,
      };
    }
  },
  { immediate: true }
);
</script>

<template>
  <BaseDialog :show="show" :title="dialogTitle" @close="emit('close')">
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label class="form-label">カテゴリ名</label>
        <input
          v-model="form.name"
          type="text"
          class="form-input"
          maxlength="50"
          required
          placeholder="カテゴリ名を入力"
        />
      </div>

      <div class="form-group">
        <label class="form-label">種別</label>
        <select v-model="form.type" class="form-input" required>
          <option value="">選択してください</option>
          <option value="income">収入</option>
          <option value="expense">支出</option>
        </select>
      </div>

      <div class="form-group">
        <label class="form-label">色</label>
        <div class="color-input-group">
          <input v-model="form.color" type="color" class="form-color" />
          <input
            v-model="form.color"
            type="text"
            class="form-input color-text"
            pattern="#[0-9A-Fa-f]{6}"
            placeholder="#007BFF"
          />
        </div>
        <div class="color-presets">
          <button
            v-for="color in COLOR_PRESETS"
            :key="color"
            type="button"
            class="color-preset"
            :style="{ backgroundColor: color }"
            @click="form.color = color"
          ></button>
        </div>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" @click="emit('close')">キャンセル</button>
        <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
          {{ category ? '更新' : '作成' }}
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

.color-input-group {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.form-color {
  width: 60px;
  height: 40px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
}

.color-text {
  flex: 1;
}

.color-presets {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
  flex-wrap: wrap;
}

.color-preset {
  width: 30px;
  height: 30px;
  border: 2px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  transition: transform 0.2s;
}

.color-preset:hover {
  transform: scale(1.1);
  border-color: #007bff;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}
</style>
