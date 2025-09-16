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
        <table class="table">
          <thead>
            <tr>
              <th>日付</th>
              <th>種別</th>
              <th>カテゴリ</th>
              <th>金額</th>
              <th>メモ</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="transaction in transactions" :key="transaction.id">
              <td>{{ formatDate(transaction.transaction_date) }}</td>
              <td>
                <span
                  :class="{
                    'badge-income': transaction.type === 'income',
                    'badge-expense': transaction.type === 'expense',
                  }"
                >
                  {{ transaction.type === 'income' ? '収入' : '支出' }}
                </span>
              </td>
              <td>{{ transaction.category?.name }}</td>
              <td
                :class="{
                  income: transaction.type === 'income',
                  expense: transaction.type === 'expense',
                }"
              >
                {{ transaction.type === 'income' ? '+' : '-' }}¥{{
                  formatNumber(transaction.amount)
                }}
              </td>
              <td>{{ transaction.memo || '-' }}</td>
              <td>
                <button class="btn btn-secondary" @click="editTransaction(transaction)">
                  編集
                </button>
                <button class="btn btn-danger" @click="deleteTransaction(transaction.id)">
                  削除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
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

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import { useTransactionStore } from '@/stores/transaction';
  import TransactionModal from '@/components/transaction/TransactionModal.vue';
  import type { Transaction, CreateTransactionRequest } from '@/types';

  const transactionStore = useTransactionStore();
  const { transactions, loading, error } = transactionStore;

  const showCreateModal = ref(false);
  const showEditModal = ref(false);
  const editingTransaction = ref<Transaction | null>(null);

  const formatNumber = (num: number) => {
    return new Intl.NumberFormat('ja-JP').format(num);
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('ja-JP');
  };

  const editTransaction = (transaction: Transaction) => {
    editingTransaction.value = transaction;
    showEditModal.value = true;
  };

  const deleteTransaction = async (id: number) => {
    if (confirm('この取引を削除しますか？')) {
      try {
        await transactionStore.deleteTransaction(id);
      } catch (err) {
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
    }
  };

  onMounted(() => {
    transactionStore.fetchTransactions();
  });
</script>

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

  .badge-income,
  .badge-expense {
    padding: 0.25rem 0.5rem;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 500;
  }

  .badge-income {
    background-color: #d4edda;
    color: #155724;
  }

  .badge-expense {
    background-color: #f8d7da;
    color: #721c24;
  }

  .income {
    color: #28a745;
  }

  .expense {
    color: #dc3545;
  }

  .table td .btn {
    margin-right: 0.5rem;
    font-size: 0.75rem;
    padding: 0.25rem 0.5rem;
  }

  @media (max-width: 768px) {
    .page-header {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }

    .table {
      font-size: 0.875rem;
    }

    .table td .btn {
      font-size: 0.625rem;
      padding: 0.125rem 0.25rem;
    }
  }
</style>
