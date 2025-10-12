<script setup lang="ts">
import { ref } from 'vue';
import { RouterView } from 'vue-router';
import { useNotification } from '@/composables/useNotification';

const drawer = ref(false);

const navItems = [
  { title: 'ダッシュボード', icon: 'mdi-view-dashboard', to: '/' },
  { title: '取引', icon: 'mdi-cash-multiple', to: '/transactions' },
  { title: 'カテゴリ', icon: 'mdi-tag-multiple', to: '/categories' },
  { title: '予算', icon: 'mdi-calculator', to: '/budgets' },
  { title: 'サマリー', icon: 'mdi-chart-bar', to: '/summary' },
];

// グローバル通知システム
const notification = useNotification();
</script>

<template>
  <v-app>
    <!-- App Bar -->
    <v-app-bar color="primary" prominent>
      <v-app-bar-nav-icon variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>

      <v-toolbar-title class="text-h5 font-weight-bold"> Budget Book </v-toolbar-title>

      <v-spacer></v-spacer>

      <!-- Desktop Navigation -->
      <template v-for="item in navItems" :key="item.to">
        <v-btn :to="item.to" variant="text" class="d-none d-md-flex">
          <v-icon start>{{ item.icon }}</v-icon>
          {{ item.title }}
        </v-btn>
      </template>
    </v-app-bar>

    <!-- Mobile Navigation Drawer -->
    <v-navigation-drawer v-model="drawer" temporary>
      <v-list>
        <v-list-item
          prepend-icon="mdi-wallet"
          title="Budget Book"
          subtitle="家計簿アプリ"
        ></v-list-item>
      </v-list>

      <v-divider></v-divider>

      <v-list density="compact" nav>
        <v-list-item
          v-for="item in navItems"
          :key="item.to"
          :prepend-icon="item.icon"
          :title="item.title"
          :to="item.to"
          @click="drawer = false"
        ></v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- Main Content -->
    <v-main>
      <v-container fluid>
        <RouterView />
      </v-container>
    </v-main>

    <!-- Global Snackbar for Notifications -->
    <v-snackbar
      v-model="notification.show.value"
      :color="notification.color.value"
      :timeout="notification.timeout.value"
      location="top"
    >
      {{ notification.message.value }}

      <template #actions>
        <v-btn variant="text" @click="notification.show.value = false"> 閉じる </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<style>
/* グローバルスタイルは最小限に - Vuetifyが大部分を処理 */
</style>
