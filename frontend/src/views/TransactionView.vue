<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useTransactionStore } from '@/stores/transaction';
import TransactionModal from '@/components/transaction/TransactionModal.vue';
import type { Transaction, CreateTransactionRequest } from '@/types';
import { storeToRefs } from 'pinia';

const transactionStore = useTransactionStore();
const { transactions, loading, error } = storeToRefs(transactionStore);

const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingTransaction = ref<Transaction | null>(null);
const search = ref('');

const formatNumber = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(num);
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ja-JP');
};

// データテーブルのヘッダー定義
const headers = computed(() => [
  { title: '日付', value: 'transaction_date', sortable: true },
  { title: '種別', value: 'type', sortable: true },
  { title: 'カテゴリ', value: 'category.name', sortable: true },
  { title: '金額', value: 'amount', sortable: true },
  { title: 'メモ', value: 'memo', sortable: false },
  { title: '操作', value: 'actions', sortable: false, width: '120px' },
]);

const editTransaction = (transaction: Transaction) => {
  editingTransaction.value = transaction;
  showEditModal.value = true;
};

const deleteTransaction = async (id: number) => {
  if (confirm('この取引を削除しますか？')) {
    try {
      await transactionStore.deleteTransaction(id);
    } catch (err) {
      alert('取引の削除に失敗しました。');
    }
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  editingTransaction.value = null;
};

const handleSave = async (data: CreateTransactionRequest) => {
  try {
    if (editingTransaction.value) {
      await transactionStore.updateTransaction(editingTransaction.value.id, data);
    } else {
      await transactionStore.createTransaction(data);
    }
    closeModal();
  } catch (err) {
    alert('取引の保存に失敗しました。');
  }
};

onMounted(() => {
  transactionStore.fetchTransactions();
});
</script>

<template>
  <div class="transactions">
    <div class="page-header">
      <h2>取引管理</h2>
      <button class="btn btn-primary" @click="showCreateModal = true">新規取引</button>
    </div>

    <div class="card">
      <div v-if="loading" class="loading">読み込み中...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="transactions.length === 0">取引がありません</div>
      <div v-else>
        <v-data-table
          :headers="headers"
          :items="transactions"
          :search="search"
          :loading="loading"
          class="elevation-1"
          :items-per-page="10"
          :items-per-page-options="[5, 10, 25, 50]"
        >
          <!-- 検索フィールド -->
          <template #top>
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
              @click="editTransaction(item)"
            >
              編集
            </v-btn>
            <v-btn
              color="error"
              variant="outlined"
              size="small"
              class="me-2"
              @click="deleteTransaction(item.id)"
            >
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
    </div>

    <TransactionModal
      v-if="showCreateModal || showEditModal"
      :show="showCreateModal || showEditModal"
      :transaction="editingTransaction"
      @close="closeModal"
      @save="handleSave"
    />
  </div>
</template>

<style scoped>
.transactions {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h2 {
  margin: 0;
  color: #333;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
}
</style>
