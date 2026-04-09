<template>
  <div class="card" v-if="isConnected()">
    <div class="card-title">Access Check</div>
    <label>Chain</label>
    <div class="chain-selector" :class="{ disabled: loadingChains }">
      <button
        v-for="chain in chains"
        :key="chain.id"
        class="chain-option"
        :class="{ active: selectedChain === chain.id }"
        :disabled="loadingChains"
        @click="selectedChain = chain.id"
      >
        <span class="chain-dot" :style="{ background: chainColor(chain.id) }"></span>
        <span class="chain-name">{{ chain.name }}</span>
        <span class="chain-id-badge">{{ chain.chain_id }}</span>
      </button>
    </div>
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
import { ref, onMounted } from 'vue'
import { useWallet } from '../composables/useWallet'
import { useApi } from '../composables/useApi'
import ResultDisplay from './ResultDisplay.vue'

const { currentAddress, jwtToken, isConnected } = useWallet()
const { checkAccess, getChains } = useApi()

const chains = ref([])
const selectedChain = ref('ethereum')
const loadingChains = ref(false)
const contractAddress = ref('0xdAC17F958D2ee523a2206206994597C13D831ec7')
const result = ref(null)
const checking = ref(false)

const CHAIN_COLORS = {
  ethereum: '#627EEA',
  'avax-fuji': '#E84142',
}
function chainColor(id) {
  return CHAIN_COLORS[id] || '#8B5CF6'
}

onMounted(async () => {
  loadingChains.value = true
  try {
    chains.value = await getChains()
    if (chains.value.length > 0) {
      selectedChain.value = chains.value[0].id
    }
  } catch (err) {
    console.error('Failed to load chains:', err)
  } finally {
    loadingChains.value = false
  }
})

async function handleCheck() {
  result.value = null
  checking.value = true
  try {
    const hasAccess = await checkAccess(currentAddress.value, contractAddress.value, jwtToken.value, selectedChain.value)
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
