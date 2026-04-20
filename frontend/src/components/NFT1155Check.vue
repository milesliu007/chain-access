<template>
  <div class="card" v-if="isConnected()">
    <div class="card-title">NFT 1155 Check</div>
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
    <label for="nft1155ContractInput">ERC-1155 Contract Address</label>
    <input
      id="nft1155ContractInput"
      v-model="contractAddress"
      type="text"
      placeholder="0x..."
    />
    <label for="nft1155TokenIdInput">Token ID</label>
    <input
      id="nft1155TokenIdInput"
      v-model="tokenId"
      type="text"
      placeholder="Token ID (e.g. 1)"
    />
    <button class="btn-secondary" :disabled="!contractAddress || !tokenId || checking" @click="handleCheck">
      {{ checking ? 'Checking...' : 'Check NFT' }}
    </button>
    <ResultDisplay :result="result" />
    <div v-if="result && result.type === 'access'" class="token-list">
      <div class="token-list-title">Balance: {{ balance }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useWallet } from '../composables/useWallet'
import { useApi } from '../composables/useApi'
import ResultDisplay from './ResultDisplay.vue'

const { currentAddress, jwtToken, isConnected } = useWallet()
const { checkNFT1155, getChains } = useApi()

const chains = ref([])
const selectedChain = ref('ethereum')
const loadingChains = ref(false)
const contractAddress = ref('')
const tokenId = ref('')
const result = ref(null)
const balance = ref(0)
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
  balance.value = 0
  checking.value = true
  try {
    const data = await checkNFT1155(currentAddress.value, contractAddress.value, tokenId.value, jwtToken.value, selectedChain.value)
    balance.value = data.balance || 0
    result.value = {
      type: data.has_nft ? 'access' : 'no-access',
      message: data.has_nft
        ? `NFT Holder — Balance: ${balance.value}`
        : 'Not a Holder — Wallet does not own this NFT',
    }
  } catch (err) {
    result.value = { type: 'error', message: err.message }
  } finally {
    checking.value = false
  }
}
</script>
