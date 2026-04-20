<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">⛓ Chain Access</div>
        <h1>Admin Panel</h1>
        <p>Connect wallet and verify ERC721 NFT ownership to login</p>
      </div>

      <div v-if="error" class="error-msg">{{ error }}</div>

      <div class="form-group">
        <label>Chain</label>
        <select v-model="chainId">
          <option value="ethereum">Ethereum Mainnet</option>
          <option value="avax-fuji">Avalanche Fuji Testnet</option>
        </select>
      </div>

      <div class="form-group">
        <label>ERC721 Contract Address</label>
        <input
          v-model="nftContract"
          placeholder="0x..."
          spellcheck="false"
        />
      </div>

      <button class="btn-connect" :disabled="loading || !nftContract" @click="handleLogin">
        <span v-if="loading" class="spinner"></span>
        <span v-else>🦊 Connect &amp; Verify NFT</span>
      </button>

      <div v-if="connectedAddr" class="address-display">
        Connected: {{ shortAddr(connectedAddr) }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminAuth } from '../composables/useAdminAuth.js'

const router = useRouter()
const { login, address } = useAdminAuth()

const loading = ref(false)
const error = ref('')
const chainId = ref('ethereum')
const nftContract = ref('')
const connectedAddr = address

async function handleLogin() {
  if (!nftContract.value.startsWith('0x')) {
    error.value = 'Invalid contract address'
    return
  }
  loading.value = true
  error.value = ''
  try {
    await login(chainId.value, nftContract.value)
    router.push('/balances')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function shortAddr(addr) {
  return addr.value ? `${addr.value.slice(0, 6)}...${addr.value.slice(-4)}` : ''
}
</script>
