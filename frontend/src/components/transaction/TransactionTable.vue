<script setup lang="ts">
import { computed, ref } from 'vue';
import type { Transaction } from '@/types';

interface Props {
  transactions: Transaction[];
  loading?: boolean;
  showSearch?: boolean;
  itemsPerPage?: number;
  height?: string | number;
}

interface Emits {
  (e: 'edit', transaction: Transaction): void;
  (e: 'delete', id: number): void;
}

// データテーブルのヘッダー定義
const headers = computed(() => [
  { title: '日付', value: 'transaction_date', sortable: true },
  { title: '種別', value: 'type', sortable: true },
  { title: 'カテゴリ', value: 'category.name', sortable: true },
  { title: '金額', value: 'amount', sortable: true },
  { title: 'メモ', value: 'memo', sortable: false },
  { title: '操作', value: 'actions', sortable: false, width: '120px' },
]);

withDefaults(defineProps<Props>(), {
  loading: false,
  showSearch: true,
  itemsPerPage: 10,
});

const emit = defineEmits<Emits>();

const search = ref('');

const formatNumber = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(num);
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ja-JP');
};

const handleEdit = (transaction: Transaction) => {
  emit('edit', transaction);
};

const handleDelete = (id: number) => {
  emit('delete', id);
};
</script>

<template>
  <div class="transaction-table">
    <v-data-table
      :headers="headers"
      :items="transactions"
      :search="showSearch ? search : undefined"
      :loading="loading"
      class="elevation-1"
      :items-per-page="itemsPerPage"
      :items-per-page-options="[5, 10, 25, 50]"
      :height="height"
    >
      <!-- 検索フィールド -->
      <template v-if="showSearch" #top>
        <v-toolbar flat>
          <v-spacer></v-spacer>
          <v-text-field
            v-model="search"
            prepend-inner-icon="mdi-magnify"
            label="検索"
            single-line
            hide-details
            variant="outlined"
            density="compact"
            style="max-width: 300px"
          ></v-text-field>
        </v-toolbar>
      </template>

      <!-- 日付フォーマット -->
      <template #[`item.transaction_date`]="{ item }">
        {{ formatDate(item.transaction_date) }}
      </template>

      <!-- 種別カスタム表示 -->
      <template #[`item.type`]="{ item }">
        <v-chip
          :color="item.type === 'income' ? 'success' : 'error'"
          variant="outlined"
          size="small"
        >
          {{ item.type === 'income' ? '収入' : '支出' }}
        </v-chip>
      </template>

      <!-- 金額フォーマット -->
      <template #[`item.amount`]="{ item }">
        <span
          :class="{
            'text-success': item.type === 'income',
            'text-error': item.type === 'expense',
          }"
        >
          {{ item.type === 'income' ? '+' : '-' }}¥{{ formatNumber(item.amount) }}
        </span>
      </template>

      <!-- メモ表示 -->
      <template #[`item.memo`]="{ item }">
        {{ item.memo || '-' }}
      </template>

      <!-- 操作ボタン -->
      <template #[`item.actions`]="{ item }">
        <v-btn
          color="primary"
          variant="outlined"
          size="small"
          class="me-2"
          @click="handleEdit(item)"
        >
          編集
        </v-btn>
        <v-btn color="error" variant="outlined" size="small" @click="handleDelete(item.id)">
          削除
        </v-btn>
      </template>

      <!-- データなしの表示 -->
      <template #no-data>
        <div class="text-center pa-4">
          <v-icon size="48" color="grey-lighten-1">mdi-database-off</v-icon>
          <p class="text-h6 mt-2">取引がありません</p>
        </div>
      </template>
    </v-data-table>
  </div>
</template>

<style scoped>
.transaction-table {
  width: 100%;
}
</style>
