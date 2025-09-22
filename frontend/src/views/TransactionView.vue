<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTransactionStore } from '@/stores/transaction';
import TransactionModal from '@/components/transaction/TransactionModal.vue';
import type { Transaction, CreateTransactionRequest } from '@/types';
import { storeToRefs } from 'pinia';
import TransactionTable from '@/components/transaction/TransactionTable.vue';

const transactionStore = useTransactionStore();
const { transactions, loading, error } = storeToRefs(transactionStore);

const showCreateModal = ref(false);
const showEditModal = ref(false);
const editingTransaction = ref<Transaction | null>(null);
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
        <TransactionTable
          :transactions="transactions"
          :loading="loading"
          :items-per-page="10"
          @edit="editTransaction"
          @delete="deleteTransaction"
        />
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
