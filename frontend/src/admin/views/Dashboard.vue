<template>
  <div class="admin-layout">
    <Sidebar />
    <main class="main-content">
      <div class="page-header">
        <h2>Dashboard</h2>
      </div>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-label">Total Records</div>
          <div class="stat-value">—</div>
        </div>
        <div class="stat-card">
          <div class="stat-label">Status</div>
          <div class="stat-value online">Online</div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import Sidebar from '../components/Sidebar.vue'
</script>
