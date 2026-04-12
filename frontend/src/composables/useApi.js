const API_BASE = window.location.origin

export function useApi() {
  async function requestChallenge(address) {
    const res = await fetch(`${API_BASE}/auth/challenge`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ address }),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || '获取 challenge 失败')
    }
    const data = await res.json()
    return data.challenge
  }

  async function verifySignature(address, signature) {
    const res = await fetch(`${API_BASE}/auth/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ address, signature }),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || '签名验证失败')
    }
    const data = await res.json()
    return data.token
  }

  async function getChains() {
    const res = await fetch(`${API_BASE}/chains`)
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || 'Failed to fetch chains')
    }
    const data = await res.json()
    return data.chains
  }

  async function checkAccess(address, contractAddress, token, chainId = 'ethereum') {
    const res = await fetch(`${API_BASE}/check-access`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        chain_id: chainId,
        address,
        contract_address: contractAddress,
      }),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || '查询失败')
    }
    const data = await res.json()
    return data.has_access
  }

  async function checkNFT(address, contractAddress, token, chainId = 'ethereum') {
    const res = await fetch(`${API_BASE}/check-nft`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        chain_id: chainId,
        address,
        contract_address: contractAddress,
      }),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || 'NFT query failed')
    }
    return await res.json()
  }

  return { requestChallenge, verifySignature, checkAccess, checkNFT, getChains }
}
