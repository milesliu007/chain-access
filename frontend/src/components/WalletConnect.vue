<template>
  <div class="card">
    <div class="status">
      <span class="dot" :class="{ connected: isConnected() }"></span>
      <span>{{ isConnected() ? '已连接' : '未连接钱包' }}</span>
    </div>
    <p v-if="currentAddress" class="addr">{{ currentAddress }}</p>
    <button
      class="btn-primary"
      :disabled="isConnecting || isConnected()"
      @click="handleConnect"
    >
      {{ isConnecting ? '连接中...' : isConnected() ? '已连接' : '连接钱包' }}
    </button>
    <p v-if="error" class="result error">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useWallet } from '../composables/useWallet'

const { currentAddress, isConnecting, isConnected, connect } = useWallet()
const error = ref(null)

async function handleConnect() {
  error.value = null
  try {
    await connect()
  } catch (err) {
    error.value = err.message
  }
}
</script>
