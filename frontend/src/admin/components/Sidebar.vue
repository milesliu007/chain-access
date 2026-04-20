<template>
  <nav class="sidebar">
    <div class="sidebar-logo">⛓ Admin</div>
    <ul class="sidebar-menu">
      <li>
        <router-link to="/dashboard" class="menu-item">
          <span class="icon">📊</span> Dashboard
        </router-link>
      </li>
      <li>
        <router-link to="/balances" class="menu-item">
          <span class="icon">💰</span> User Balances
        </router-link>
      </li>
    </ul>
    <div class="sidebar-footer">
      <div class="addr-badge">{{ shortAddr(address) }}</div>
      <button class="btn-logout" @click="handleLogout">Logout</button>
    </div>
  </nav>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAdminAuth } from '../composables/useAdminAuth.js'

const router = useRouter()
const { address, logout } = useAdminAuth()

function shortAddr(addr) {
  return addr.value ? `${addr.value.slice(0, 6)}...${addr.value.slice(-4)}` : ''
}

function handleLogout() {
  logout()
  router.push('/login')
}
</script>
