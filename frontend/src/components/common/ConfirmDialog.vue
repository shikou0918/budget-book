<script setup lang="ts">
import { ref, watch } from 'vue';

/**
 * 確認ダイアログコンポーネントのProps定義
 */
interface Props {
  /** ダイアログの表示状態（v-model） */
  modelValue: boolean;
  /** ダイアログのタイトル（任意） */
  title?: string;
  /** 確認メッセージ */
  message: string;
  /** 確認ボタンのテキスト（任意） */
  confirmText?: string;
  /** キャンセルボタンのテキスト（任意） */
  cancelText?: string;
  /** 確認ボタンの色（任意） */
  confirmColor?: string;
}

/**
 * 確認ダイアログコンポーネントのEmits定義
 */
interface Emits {
  /** ダイアログ表示状態の更新イベント */
  (e: 'update:modelValue', value: boolean): void;
  /** 確認ボタンクリックイベント */
  (e: 'confirm'): void;
  /** キャンセルボタンクリックイベント */
  (e: 'cancel'): void;
}

const props = withDefaults(defineProps<Props>(), {
  title: '確認',
  confirmText: '確認',
  cancelText: 'キャンセル',
  confirmColor: 'error',
});

const emit = defineEmits<Emits>();

const dialog = ref(props.modelValue);

watch(
  () => props.modelValue,
  newVal => {
    dialog.value = newVal;
  }
);

watch(dialog, newVal => {
  emit('update:modelValue', newVal);
});

const handleConfirm = () => {
  emit('confirm');
  dialog.value = false;
};

const handleCancel = () => {
  emit('cancel');
  dialog.value = false;
};
</script>

<template>
  <v-dialog v-model="dialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        {{ props.title }}
      </v-card-title>

      <v-card-text>
        {{ props.message }}
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="handleCancel">
          {{ props.cancelText }}
        </v-btn>
        <v-btn :color="props.confirmColor" variant="elevated" @click="handleConfirm">
          {{ props.confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
/* Vuetifyコンポーネントを使用 */
</style>
