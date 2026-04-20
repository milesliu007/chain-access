<template>
  <div class="container">
    <header class="app-header">
      <div class="app-logo">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
            <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
          </svg>
        </div>
        <h1>Chain Access</h1>
      </div>
      <p class="subtitle">Web3 Token-Gated Access Verification</p>
    </header>
    <WalletConnect />
    <div class="tab-bar">
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'erc20' }"
        @click="activeTab = 'erc20'"
      >ERC-20</button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'nft' }"
        @click="activeTab = 'nft'"
      >NFT (ERC-721)</button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'nft1155' }"
        @click="activeTab = 'nft1155'"
      >ERC-1155</button>
    </div>
    <AccessCheck v-if="activeTab === 'erc20'" />
    <NFTCheck v-if="activeTab === 'nft'" />
    <NFT1155Check v-if="activeTab === 'nft1155'" />
    <footer class="app-footer">
      Powered by Multi-Chain EVM &middot; ERC-20 &amp; ERC-721 Token Gating
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import WalletConnect from './components/WalletConnect.vue'
import AccessCheck from './components/AccessCheck.vue'
import NFTCheck from './components/NFTCheck.vue'
import NFT1155Check from './components/NFT1155Check.vue'

const activeTab = ref('erc20')
</script>
