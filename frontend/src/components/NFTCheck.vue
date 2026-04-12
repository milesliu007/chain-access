<template>
  <div class="card" v-if="isConnected()">
    <div class="card-title">NFT Check</div>
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
    <label for="nftContractInput">ERC-721 Contract Address</label>
    <input
      id="nftContractInput"
      v-model="contractAddress"
      type="text"
      placeholder="0x..."
    />
    <button class="btn-secondary" :disabled="!contractAddress || checking" @click="handleCheck">
      {{ checking ? 'Checking...' : 'Check NFT' }}
    </button>
    <ResultDisplay :result="result" />
    <div v-if="tokenIds.length > 0" class="token-list">
      <div class="token-list-title">Token IDs</div>
      <div class="token-id-grid">
        <span v-for="id in tokenIds" :key="id" class="token-id-badge">
          #{{ id }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useWallet } from '../composables/useWallet'
import { useApi } from '../composables/useApi'
import ResultDisplay from './ResultDisplay.vue'

const { currentAddress, jwtToken, isConnected } = useWallet()
const { checkNFT, getChains } = useApi()

const chains = ref([])
const selectedChain = ref('ethereum')
const loadingChains = ref(false)
const contractAddress = ref('')
const result = ref(null)
const tokenIds = ref([])
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
  tokenIds.value = []
  checking.value = true
  try {
    const data = await checkNFT(currentAddress.value, contractAddress.value, jwtToken.value, selectedChain.value)
    tokenIds.value = data.token_ids || []
    result.value = {
      type: data.has_nft ? 'access' : 'no-access',
      message: data.has_nft
        ? `NFT Holder — ${tokenIds.value.length} token(s) found`
        : 'Not a Holder — Wallet does not own this NFT',
    }
  } catch (err) {
    result.value = { type: 'error', message: err.message }
  } finally {
    checking.value = false
  }
}
</script>
