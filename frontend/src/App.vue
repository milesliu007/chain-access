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
    <AccessCheck />
    <footer class="app-footer">
      Powered by Ethereum &middot; ERC-20 Token Gating
    </footer>
  </div>
</template>

<script setup>
import WalletConnect from './components/WalletConnect.vue'
import AccessCheck from './components/AccessCheck.vue'
</script>
