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

  async function checkAccess(address, contractAddress, token) {
    const res = await fetch(
      `${API_BASE}/check-access?address=${encodeURIComponent(address)}&contract_address=${encodeURIComponent(contractAddress)}`,
      { headers: { Authorization: `Bearer ${token}` } }
    )
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || '查询失败')
    }
    const data = await res.json()
    return data.has_access
  }

  return { requestChallenge, verifySignature, checkAccess }
}
