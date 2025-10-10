<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTransactionStore } from '@/stores/transaction';
import TransactionDialog from '@/components/transaction/TransactionDialog.vue';
import type { Transaction, CreateTransactionRequest } from '@/types';
import { storeToRefs } from 'pinia';
import TransactionTable from '@/components/transaction/TransactionTable.vue';

const transactionStore = useTransactionStore();
const { transactions, loading, error } = storeToRefs(transactionStore);

const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const editingTransaction = ref<Transaction | null>(null);
const editTransaction = (transaction: Transaction) => {
  editingTransaction.value = transaction;
  showEditDialog.value = true;
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

const closeDialog = () => {
  showCreateDialog.value = false;
  showEditDialog.value = false;
  editingTransaction.value = null;
};

const handleSave = async (data: CreateTransactionRequest) => {
  try {
    if (editingTransaction.value) {
      await transactionStore.updateTransaction(editingTransaction.value.id, data);
    } else {
      await transactionStore.createTransaction(data);
    }
    closeDialog();
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
    <h2>取引管理</h2>

    <div class="transactions-grid">
      <div class="card">
        <div class="card-header">
          <h3>取引一覧</h3>
          <button class="btn btn-primary" @click="showCreateDialog = true">新規取引</button>
        </div>
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="transactions.length === 0">取引がありません</div>
        <div v-else>
          <TransactionTable
            :transactions="transactions"
            :loading="loading"
            :items-per-page="10"
            @edit="editTransaction"
            @delete="deleteTransaction"
          />
        </div>
      </div>
    </div>

    <TransactionDialog
      v-if="showCreateDialog || showEditDialog"
      :show="showCreateDialog || showEditDialog"
      :transaction="editingTransaction"
      @close="closeDialog"
      @save="handleSave"
    />
  </div>
</template>

<style scoped>
.transactions {
}

.transactions h2 {
  margin-bottom: 2rem;
  color: #333;
}

.transactions-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.transactions-grid .card {
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

@media (max-width: 768px) {
  .transactions-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
}
</style>
