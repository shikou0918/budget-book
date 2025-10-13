<script setup lang="ts">
import { computed } from 'vue';
import { Pie } from 'vue-chartjs';
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  type ChartOptions,
  type ChartData,
} from 'chart.js';

// Register required Chart.js components
ChartJS.register(ArcElement, Tooltip, Legend);

/**
 * 円グラフコンポーネントのProps定義
 */
interface Props {
  /** グラフのラベル配列 */
  labels: string[];
  /** グラフのデータ配列 */
  data: number[];
  /** グラフのタイトル（任意） */
  title?: string;
  /** 各データの背景色配列（任意） */
  backgroundColor?: string[];
  /** グラフの高さ（ピクセル、任意） */
  height?: number;
}

const props = withDefaults(defineProps<Props>(), {
  title: 'カテゴリ別支出',
  backgroundColor: () => [
    '#FF6384',
    '#36A2EB',
    '#FFCE56',
    '#4BC0C0',
    '#9966FF',
    '#FF9F40',
    '#FF6384',
    '#C9CBCF',
    '#4BC0C0',
    '#FF6384',
  ],
  height: 300,
});

// Compute chart data
const chartData = computed<ChartData<'pie'>>(() => ({
  labels: props.labels,
  datasets: [
    {
      data: props.data,
      backgroundColor: props.backgroundColor,
      borderWidth: 1,
      borderColor: '#fff',
    },
  ],
}));

// Chart options
const chartOptions = computed<ChartOptions<'pie'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'right',
      labels: {
        padding: 15,
        font: {
          size: 12,
        },
      },
    },
    tooltip: {
      callbacks: {
        label: context => {
          const label = context.label || '';
          const value = context.parsed || 0;
          const total = context.dataset.data.reduce((acc: number, val) => acc + (val as number), 0);
          const percentage = ((value / total) * 100).toFixed(1);
          return `${label}: ¥${value.toLocaleString('ja-JP')} (${percentage}%)`;
        },
      },
    },
  },
}));
</script>

<template>
  <div class="pie-chart-container">
    <h3 v-if="title" class="chart-title">{{ title }}</h3>
    <div class="chart-wrapper" :style="{ height: `${height}px` }">
      <Pie :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<style scoped>
.pie-chart-container {
  width: 100%;
}

.chart-title {
  margin-bottom: 1rem;
  color: #333;
  font-size: 1.125rem;
  font-weight: 600;
}

.chart-wrapper {
  position: relative;
  width: 100%;
}
</style>
