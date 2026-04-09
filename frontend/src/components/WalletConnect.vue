<template>
  <div class="card">
    <div class="card-title">Wallet</div>
    <div class="status">
      <span class="dot" :class="{ connected: isConnected() }"></span>
      <span>{{ isConnected() ? 'Connected' : 'Not Connected' }}</span>
    </div>
    <p v-if="currentAddress" class="addr">{{ currentAddress }}</p>
    <button
      class="btn-primary"
      :class="{ connecting: isConnecting }"
      :disabled="isConnecting || isConnected()"
      @click="handleConnect"
    >
      {{ isConnecting ? 'Connecting...' : isConnected() ? 'Connected' : 'Connect Wallet' }}
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
