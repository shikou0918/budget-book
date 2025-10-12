<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { summaryApi } from '@/services/api';
import type { MonthlySummary } from '@/types/index';
import PieChart from '@/components/common/PieChart.vue';

const summary = ref<MonthlySummary | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

const selectedYear = ref(new Date().getFullYear());
const selectedMonth = ref(new Date().getMonth() + 1);
const chartType = ref<'income' | 'expense'>('expense');

const years = computed(() => {
  const currentYear = new Date().getFullYear();
  return Array.from({ length: 5 }, (_, i) => currentYear - 2 + i);
});

const months = [
  { value: 1, label: '1月' },
  { value: 2, label: '2月' },
  { value: 3, label: '3月' },
  { value: 4, label: '4月' },
  { value: 5, label: '5月' },
  { value: 6, label: '6月' },
  { value: 7, label: '7月' },
  { value: 8, label: '8月' },
  { value: 9, label: '9月' },
  { value: 10, label: '10月' },
  { value: 11, label: '11月' },
  { value: 12, label: '12月' },
];

const categoryDetails = computed(() => {
  if (!summary.value?.category_summary) return [];
  return Object.values(summary.value.category_summary)
    .filter(detail => detail.total > 0)
    .sort((a, b) => b.total - a.total);
});

const pieChartData = computed(() => {
  const categories = categoryDetails.value.filter(
    detail => detail.category_type === chartType.value
  );
  return {
    labels: categories.map(detail => detail.category_name),
    data: categories.map(detail => detail.total),
  };
});

const formatNumber = (num: number) => {
  return new Intl.NumberFormat('ja-JP').format(num);
};

const getPercentageClass = (percentage: number) => {
  if (percentage <= 70) return 'percentage-good';
  if (percentage <= 90) return 'percentage-warning';
  return 'percentage-danger';
};

const fetchSummary = async () => {
  loading.value = true;
  error.value = null;
  try {
    const response = await summaryApi.getMonthly(selectedYear.value, selectedMonth.value);
    if (!response) {
      throw new Error('サマリーの取得に失敗しました');
    }
    summary.value = response.data;
  } catch (err) {
    error.value = 'サマリーの取得に失敗しました';
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchSummary();
});
</script>

<template>
  <div class="summary">
    <h2>サマリー</h2>

    <div v-if="loading" class="loading">読み込み中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="summary" class="summary-content">
      <div class="summary-grid">
        <div class="card summary-card">
          <div class="card-header">
            <h3>{{ selectedYear }}年{{ selectedMonth }}月の収支</h3>
            <div class="month-selector">
              <select v-model="selectedYear" class="form-input">
                <option v-for="year in years" :key="year" :value="year">{{ year }}年</option>
              </select>
              <select v-model="selectedMonth" class="form-input">
                <option v-for="month in months" :key="month.value" :value="month.value">
                  {{ month.label }}
                </option>
              </select>
              <button class="btn btn-primary" @click="fetchSummary">表示</button>
            </div>
          </div>
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

        <div class="card" v-if="pieChartData.data.length > 0">
          <div class="chart-header">
            <h3>カテゴリ別{{ chartType === 'expense' ? '支出' : '収入' }}</h3>
            <div class="chart-type-selector">
              <button
                class="btn"
                :class="{
                  'btn-primary': chartType === 'expense',
                  'btn-secondary': chartType !== 'expense',
                }"
                @click="chartType = 'expense'"
              >
                支出
              </button>
              <button
                class="btn"
                :class="{
                  'btn-primary': chartType === 'income',
                  'btn-secondary': chartType !== 'income',
                }"
                @click="chartType = 'income'"
              >
                収入
              </button>
            </div>
          </div>
          <PieChart
            :labels="pieChartData.labels"
            :data="pieChartData.data"
            :title="''"
            :height="250"
          />
        </div>

        <div class="card">
          <h3>カテゴリ別詳細</h3>
          <div v-if="categoryDetails.length === 0" class="no-data">データがありません</div>
          <div v-else>
            <table class="table">
              <thead>
                <tr>
                  <th>カテゴリ</th>
                  <th>実績</th>
                  <th>予算</th>
                  <th>達成率</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="detail in categoryDetails" :key="detail.category_id">
                  <td>{{ detail.category_name }}</td>
                  <td
                    :class="{
                      income: detail.category_type === 'income',
                      expense: detail.category_type === 'expense',
                    }"
                  >
                    ¥{{ formatNumber(detail.total) }}
                  </td>
                  <td>
                    <span v-if="detail.budget > 0">¥{{ formatNumber(detail.budget) }}</span>
                    <span v-else>-</span>
                  </td>
                  <td>
                    <span v-if="detail.budget > 0" :class="getPercentageClass(detail.percentage)">
                      {{ detail.percentage.toFixed(1) }}%
                    </span>
                    <span v-else>-</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.summary h2 {
  margin-bottom: 2rem;
  color: #333;
}

.summary-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.summary-card {
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

.month-selector {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.month-selector .form-input {
  width: auto;
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

.percentage-good {
  color: #28a745;
  font-weight: 600;
}

.percentage-warning {
  color: #ffc107;
  font-weight: 600;
}

.percentage-danger {
  color: #dc3545;
  font-weight: 600;
}

.no-data {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.chart-header h3 {
  margin: 0;
}

.chart-type-selector {
  display: flex;
  gap: 0.5rem;
}

.chart-type-selector .btn {
  font-size: 0.875rem;
  padding: 0.5rem 1rem;
}

@media (max-width: 768px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .month-selector {
    flex-wrap: wrap;
  }

  .summary-stats {
    flex-direction: column;
    gap: 1rem;
  }

  .chart-header {
    flex-direction: column;
    gap: 0.75rem;
    align-items: stretch;
  }
}
</style>
