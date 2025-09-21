<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import { useTransactionStore } from '@/stores/transaction';
  import { summaryApi } from '@/services/api';
  import type { MonthlySummary } from '@/types';

  const transactionStore = useTransactionStore();
  const { transactions, loading: transactionLoading, error: transactionError } = transactionStore;

  const summary = ref<MonthlySummary | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const recentTransactions = computed(() => transactions.slice(0, 5));

  const formatNumber = (num: number) => {
    return new Intl.NumberFormat('ja-JP').format(num);
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('ja-JP');
  };

  const fetchMonthlySummary = async () => {
    loading.value = true;
    error.value = null;
    try {
      const now = new Date();
      const response = await summaryApi.getMonthly(now.getFullYear(), now.getMonth() + 1);
      summary.value = response.data;
    } catch (err) {
      error.value = '月次サマリーの取得に失敗しました';
    } finally {
      loading.value = false;
    }
  };

  onMounted(async () => {
    await Promise.all([transactionStore.fetchTransactions(), fetchMonthlySummary()]);
  });
</script>

<template>
  <div class="dashboard">
    <h2>ダッシュボード</h2>

    <div class="dashboard-grid">
      <div class="card summary-card">
        <h3>今月のサマリー</h3>
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="summary">
          <div class="summary-stats">
            <div class="stat">
              <span class="stat-label">収入</span>
              <span class="stat-value income">¥{{ formatNumber(summary.total_income) }}</span>
            </div>
            <div class="stat">
              <span class="stat-label">支出</span>
              <span class="stat-value expense">¥{{ formatNumber(summary.total_expense) }}</span>
            </div>
            <div class="stat">
              <span class="stat-label">残高</span>
              <span
                class="stat-value"
                :class="{ income: summary.balance >= 0, expense: summary.balance < 0 }"
              >
                ¥{{ formatNumber(summary.balance) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <h3>最近の取引</h3>
        <div v-if="transactionLoading" class="loading">読み込み中...</div>
        <div v-else-if="transactionError" class="error">{{ transactionError }}</div>
        <div v-else-if="recentTransactions.length === 0">取引がありません</div>
        <div v-else>
          <table class="table">
            <thead>
              <tr>
                <th>日付</th>
                <th>カテゴリ</th>
                <th>金額</th>
                <th>メモ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="transaction in recentTransactions" :key="transaction.id">
                <td>{{ formatDate(transaction.transaction_date) }}</td>
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
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .dashboard {
    max-width: 1200px;
    margin: 0 auto;
  }

  .dashboard h2 {
    margin-bottom: 2rem;
    color: #333;
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    margin-bottom: 2rem;
  }

  .summary-card {
    grid-column: 1 / -1;
  }

  .summary-stats {
    display: flex;
    justify-content: space-around;
    margin-top: 1rem;
  }

  .stat {
    text-align: center;
  }

  .stat-label {
    display: block;
    font-size: 0.875rem;
    color: #666;
    margin-bottom: 0.5rem;
  }

  .stat-value {
    display: block;
    font-size: 1.5rem;
    font-weight: 700;
  }

  .stat-value.income {
    color: #28a745;
  }

  .stat-value.expense {
    color: #dc3545;
  }

  .income {
    color: #28a745;
  }

  .expense {
    color: #dc3545;
  }

  @media (max-width: 768px) {
    .dashboard-grid {
      grid-template-columns: 1fr;
    }

    .summary-stats {
      flex-direction: column;
      gap: 1rem;
    }
  }
</style>
