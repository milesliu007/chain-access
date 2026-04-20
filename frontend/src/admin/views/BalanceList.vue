<template>
  <div class="admin-layout">
    <Sidebar />
    <main class="main-content">
      <div class="page-header">
        <h2>User Balances</h2>
        <div class="search-bar">
          <input
            v-model="searchAddr"
            placeholder="Search by wallet address..."
            @keyup.enter="loadData"
          />
          <button class="btn-search" @click="loadData">Search</button>
        </div>
      </div>

      <div v-if="error" class="error-msg">{{ error }}</div>

      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Address</th>
              <th>Chain</th>
              <th>Token</th>
              <th>Balance</th>
              <th>Updated At</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="loading-cell">Loading...</td>
            </tr>
            <tr v-else-if="rows.length === 0">
              <td colspan="6" class="empty-cell">No records found</td>
            </tr>
            <tr v-for="row in rows" :key="row.id">
              <td>{{ row.id }}</td>
              <td class="addr-cell" :title="row.address">{{ shortAddr(row.address) }}</td>
              <td>{{ row.chain }}</td>
              <td class="addr-cell" :title="row.token">{{ shortAddr(row.token) }}</td>
              <td class="balance-cell">{{ row.balance }}</td>
              <td>{{ row.updated_at }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <button :disabled="page <= 1" @click="changePage(page - 1)">← Prev</button>
        <span>Page {{ page }} / {{ totalPages }}</span>
        <button :disabled="page >= totalPages" @click="changePage(page + 1)">Next →</button>
        <span class="total-count">Total: {{ total }}</span>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import { useAdminAuth } from '../composables/useAdminAuth.js'

const { fetchBalances } = useAdminAuth()

const rows = ref([])
const total = ref(0)
const page = ref(1)
const size = 20
const searchAddr = ref('')
const loading = ref(false)
const error = ref('')

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size)))

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const res = await fetchBalances(page.value, size, searchAddr.value)
    rows.value = res.data
    total.value = res.total
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function changePage(p) {
  page.value = p
  loadData()
}

function shortAddr(addr) {
  if (!addr) return ''
  return addr.length > 14 ? `${addr.slice(0, 8)}...${addr.slice(-6)}` : addr
}

onMounted(loadData)
</script>
