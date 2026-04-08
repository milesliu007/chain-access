<template>
  <div class="card" v-if="isConnected()">
    <label for="contractInput">ERC-20 合约地址</label>
    <input
      id="contractInput"
      v-model="contractAddress"
      type="text"
      placeholder="0x..."
    />
    <button class="btn-secondary" :disabled="!contractAddress" @click="handleCheck">
      查询权限
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

async function handleCheck() {
  result.value = null
  try {
    const hasAccess = await checkAccess(currentAddress.value, contractAddress.value, jwtToken.value)
    result.value = {
      type: hasAccess ? 'access' : 'no-access',
      message: hasAccess ? '有权限 — 该钱包持有此 Token' : '无权限 — 该钱包未持有此 Token',
    }
  } catch (err) {
    result.value = { type: 'error', message: err.message }
  }
}
</script>
