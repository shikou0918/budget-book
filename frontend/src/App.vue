<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <div id="app">
    <header class="header">
      <nav class="nav">
        <div class="nav-brand">
          <h1>Budget Book</h1>
        </div>
        <ul class="nav-links">
          <li><RouterLink to="/">ダッシュボード</RouterLink></li>
          <li><RouterLink to="/transactions">取引</RouterLink></li>
          <li><RouterLink to="/categories">カテゴリ</RouterLink></li>
          <li><RouterLink to="/budgets">予算</RouterLink></li>
          <li><RouterLink to="/summary">サマリー</RouterLink></li>
        </ul>
      </nav>
    </header>
    <main class="main">
      <RouterView />
    </main>
  </div>
</template>

<style>
#app {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  min-height: 100vh;
  width: 100%;
  background-color: #f5f5f5;
}

.header {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
}

.nav-brand h1 {
  color: #007bff;
  margin: 0;
  font-size: 1.5rem;
}

.nav-links {
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;
  gap: 2rem;
}

.nav-links a {
  text-decoration: none;
  color: #666;
  font-weight: 500;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: all 0.2s;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  color: #007bff;
  background-color: #f8f9fa;
}

.main {
  margin: 2rem 0;
  padding: 0 2rem;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-primary:hover {
  background-color: #0056b3;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #545b62;
}

.btn-danger {
  background-color: #dc3545;
  color: white;
}

.btn-danger:hover {
  background-color: #c82333;
}

.card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
}

.form-input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.error {
  color: #dc3545;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
}

.table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

.table th,
.table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.table th {
  font-weight: 600;
  background-color: #f8f9fa;
}

.table tbody tr:hover {
  background-color: #f8f9fa;
}

@media (max-width: 768px) {
  .nav {
    flex-direction: column;
    gap: 1rem;
  }

  .nav-links {
    flex-wrap: wrap;
    justify-content: center;
    gap: 1rem;
  }

  .main {
    padding: 0 1rem;
  }
}
</style>
