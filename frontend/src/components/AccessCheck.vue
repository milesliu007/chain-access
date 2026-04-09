<template>
  <div class="card" v-if="isConnected()">
    <div class="card-title">Access Check</div>
    <label for="contractInput">ERC-20 Contract Address</label>
    <input
      id="contractInput"
      v-model="contractAddress"
      type="text"
      placeholder="0x..."
    />
    <button class="btn-secondary" :disabled="!contractAddress || checking" @click="handleCheck">
      {{ checking ? 'Checking...' : 'Check Access' }}
    </button>
    <ResultDisplay :result="result" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useWallet } from '../composables/useWallet'
import { useApi } from '../composables/useApi'
import ResultDisplay from './ResultDisplay.vue'

const { currentAddress, jwtToken, isConnected } = useWallet()
const { checkAccess } = useApi()

const contractAddress = ref('0xdAC17F958D2ee523a2206206994597C13D831ec7')
const result = ref(null)
const checking = ref(false)

async function handleCheck() {
  result.value = null
  checking.value = true
  try {
    const hasAccess = await checkAccess(currentAddress.value, contractAddress.value, jwtToken.value)
    result.value = {
      type: hasAccess ? 'access' : 'no-access',
      message: hasAccess ? 'Access Granted — Wallet holds this token' : 'Access Denied — Wallet does not hold this token',
    }
  } catch (err) {
    result.value = { type: 'error', message: err.message }
  } finally {
    checking.value = false
  }
}
</script>
