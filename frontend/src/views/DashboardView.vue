<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useTransactionStore } from '@/stores/transaction';
import { summaryApi } from '@/services/api';
import { formatNumber } from '@/utils/formatters';
import type { MonthlySummary } from '@/types';
import { storeToRefs } from 'pinia';
import TransactionTable from '@/components/transaction/TransactionTable.vue';
import PieChart from '@/components/common/PieChart.vue';

const transactionStore = useTransactionStore();
const {
  transactions,
  loading: transactionLoading,
  error: transactionError,
} = storeToRefs(transactionStore);

const summary = ref<MonthlySummary | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);
const chartType = ref<'income' | 'expense'>('income');

const pieChartData = computed(() => {
  if (!summary.value?.category_summary) {
    return { labels: [], data: [] };
  }

  const categories = Object.values(summary.value.category_summary)
    .filter(detail => detail.category_type === chartType.value && detail.total > 0)
    .sort((a, b) => b.total - a.total);

  return {
    labels: categories.map(detail => detail.category_name),
    data: categories.map(detail => detail.total),
  };
});

// summaryが読み込まれた時、適切なchartTypeを設定
watch(
  summary,
  newSummary => {
    if (!newSummary?.category_summary) return;

    const hasIncome = Object.values(newSummary.category_summary).some(
      detail => detail.category_type === 'income' && detail.total > 0
    );
    const hasExpense = Object.values(newSummary.category_summary).some(
      detail => detail.category_type === 'expense' && detail.total > 0
    );

    // 初期表示時、incomeがあればincome、なければexpenseを表示
    if (hasIncome) {
      chartType.value = 'income';
    } else if (hasExpense) {
      chartType.value = 'expense';
    }
  },
  { immediate: true }
);

const fetchMonthlySummary = async () => {
  loading.value = true;
  error.value = null;
  try {
    const now = new Date();
    const response = await summaryApi.getMonthly(now.getFullYear(), now.getMonth() + 1);
    if (!response) {
      throw new Error('月次サマリーの取得に失敗しました');
    }
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
  <div>
    <h2 class="mb-2">ダッシュボード</h2>

    <div class="dashboard-grid">
      <div class="card summary-card">
        <h3>今月のサマリー</h3>
        <div v-if="loading" class="loading">読み込み中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="summary">
          <div class="summary-stats">
            <div class="stat">
              <span>収入</span>
              <span class="stat-value income">¥{{ formatNumber(summary.total_income) }}</span>
            </div>
            <div class="stat">
              <span>支出</span>
              <span class="stat-value expense">¥{{ formatNumber(summary.total_expense) }}</span>
            </div>
            <div class="stat">
              <span>残高</span>
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

      <v-card v-if="pieChartData.data.length > 0">
        <div class="chart-header">
          <v-card-title
            >今月のカテゴリ別{{ chartType === 'expense' ? '支出' : '収入' }}</v-card-title
          >
          <v-card-actions>
            <v-btn @click="chartType = 'expense'">支出</v-btn>
            <v-btn @click="chartType = 'income'">収入</v-btn>
          </v-card-actions>
        </div>
        <PieChart
          :labels="pieChartData.labels"
          :data="pieChartData.data"
          :title="''"
          :height="250"
        />
      </v-card>
      <v-card>
        <v-card-title>最近の取引</v-card-title>
        <v-card-text>
          <TransactionTable
            :transactions="transactions"
            :loading="transactionLoading"
            :items-per-page="10"
            :error="transactionError"
          />
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<style scoped>
.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
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

.stat-value {
  display: block;
  font-size: 20px;
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

.chart-header {
  display: flex;
  justify-content: space-between;
}
</style>
