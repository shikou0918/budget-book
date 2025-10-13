<script setup lang="ts">
/**
 * ベースダイアログコンポーネントのProps定義
 */
interface Props {
  /** ダイアログの表示状態 */
  show: boolean;
  /** ダイアログのタイトル */
  title: string;
  /** ダイアログの最大幅（任意） */
  maxWidth?: string;
}

/**
 * ベースダイアログコンポーネントのEmits定義
 */
interface Emits {
  /** ダイアログを閉じるイベント */
  close: [];
}

withDefaults(defineProps<Props>(), {
  maxWidth: '500px',
});

const emit = defineEmits<Emits>();

const handleOverlayClick = (e: Event) => {
  if (e.target === e.currentTarget) {
    emit('close');
  }
};
</script>

<template>
  <div v-if="show" class="modal-overlay" @click="handleOverlayClick">
    <div class="modal" :style="{ maxWidth }" @click.stop>
      <div class="modal-header">
        <h3>{{ title }}</h3>
        <button class="modal-close" @click="emit('close')">&times;</button>
      </div>

      <div class="modal-body">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

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
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.25rem;
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
  transition: color 0.2s;
}

.modal-close:hover {
  color: #333;
}

.modal-body {
  padding: 1.5rem;
}
</style>
