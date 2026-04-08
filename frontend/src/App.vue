<template>
  <div class="container">
    <h1>Chain Access</h1>
    <p class="subtitle">Web3 权限验证 Demo</p>
    <WalletConnect />
    <AccessCheck />
  </div>
</template>

<script setup>
import WalletConnect from './components/WalletConnect.vue'
import AccessCheck from './components/AccessCheck.vue'
</script>
