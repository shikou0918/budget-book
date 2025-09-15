import { createRouter, createWebHistory } from 'vue-router';
import DashboardView from '../views/DashboardView.vue';
import CategoryView from '../views/CategoryView.vue';
import BudgetView from '../views/BudgetView.vue';
import SummaryView from '../views/SummaryView.vue';
import TransactionView from '../views/TransactionView.vue';

const routes = [
  {
    path: '/',
    name: 'dashboard',
    component: DashboardView,
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: TransactionView,
  },
  {
    path: '/categories',
    name: 'categories',
    component: () => CategoryView,
  },
  {
    path: '/budgets',
    name: 'budgets',
    component: () => BudgetView,
  },
  {
    path: '/summary',
    name: 'summary',
    component: () => SummaryView,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
