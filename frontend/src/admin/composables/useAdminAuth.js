import { ref } from 'vue'

const token = ref(localStorage.getItem('admin_token') || '')
const address = ref(localStorage.getItem('admin_address') || '')

export function useAdminAuth() {
  const BASE = window.location.origin

  async function getChallenge(addr) {
    const res = await fetch(`${BASE}/auth/challenge`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ address: addr })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Failed to get challenge')
    return data.challenge
  }

  async function login(chainId, nftContract) {
    if (!window.ethereum) throw new Error('MetaMask not found')
    const { ethers } = await import('ethers')
    const provider = new ethers.BrowserProvider(window.ethereum)
    const accounts = await provider.send('eth_requestAccounts', [])
    const walletAddr = accounts[0]

    const message = await getChallenge(walletAddr)
    const signer = await provider.getSigner()
    const sig = await signer.signMessage(message)

    const res = await fetch(`${BASE}/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ address: walletAddr, signature: sig, chain_id: chainId, nft_contract: nftContract })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Login failed')

    token.value = data.token
    address.value = walletAddr
    localStorage.setItem('admin_token', data.token)
    localStorage.setItem('admin_address', walletAddr)
    return walletAddr
  }

  function logout() {
    token.value = ''
    address.value = ''
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_address')
  }

  async function fetchBalances(page, size, addrFilter) {
    const params = new URLSearchParams({ page, size })
    if (addrFilter) params.set('address', addrFilter)
    const res = await fetch(`${BASE}/admin/balances?${params}`, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Failed to fetch balances')
    return data
  }

  return { token, address, login, logout, fetchBalances }
}
