import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue')
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import('@/views/TransactionView.vue')
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import('@/views/CategoryView.vue')
    },
    {
      path: '/budgets',
      name: 'budgets',
      component: () => import('@/views/BudgetView.vue')
    },
    {
      path: '/summary',
      name: 'summary',
      component: () => import('@/views/SummaryView.vue')
    }
  ]
})

export default router