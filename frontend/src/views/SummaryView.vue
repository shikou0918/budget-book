<template>
  <div class="summary">
    <div class="page-header">
      <h2>サマリー</h2>
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

    <div v-if="loading" class="loading">読み込み中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="summary" class="summary-content">
      <div class="summary-grid">
        <div class="card summary-card">
          <h3>{{ selectedYear }}年{{ selectedMonth }}月の収支</h3>
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

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue';
  import { summaryApi } from '@/services/api';
  import type { MonthlySummary } from '@/types/index';

  const summary = ref<MonthlySummary | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const selectedYear = ref(new Date().getFullYear());
  const selectedMonth = ref(new Date().getMonth() + 1);

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
      summary.value = response.data;
    } catch (err) {
      error.value = 'サマリーの取得に失敗しました';
      console.error('Error fetching summary:', err);
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    fetchSummary();
  });
</script>

<style scoped>
  .summary {
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

  .month-selector {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .month-selector .form-input {
    width: auto;
  }

  .summary-grid {
    display: grid;
    grid-template-columns: 1fr 2fr;
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

  @media (max-width: 768px) {
    .page-header {
      flex-direction: column;
      gap: 1rem;
      align-items: stretch;
    }

    .month-selector {
      flex-wrap: wrap;
    }

    .summary-grid {
      grid-template-columns: 1fr;
    }

    .summary-stats {
      flex-direction: column;
      gap: 1rem;
    }
  }
</style>
