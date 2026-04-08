import { ref } from 'vue'
import { useApi } from './useApi'

const currentAddress = ref(null)
const jwtToken = ref(null)
const isConnecting = ref(false)

export function useWallet() {
  const { requestChallenge, verifySignature } = useApi()

  const isConnected = () => !!jwtToken.value

  async function connect() {
    if (typeof window.ethereum === 'undefined') {
      throw new Error('请安装 MetaMask 钱包')
    }

    isConnecting.value = true
    try {
      const { BrowserProvider } = await import('ethers')
      const provider = new BrowserProvider(window.ethereum)
      const signer = await provider.getSigner()
      const address = await signer.getAddress()

      // 1. 请求 challenge
      const challenge = await requestChallenge(address)

      // 2. MetaMask 签名
      const signature = await signer.signMessage(challenge)

      // 3. 验证签名，获取 JWT
      const token = await verifySignature(address, signature)

      currentAddress.value = address
      jwtToken.value = token
    } finally {
      isConnecting.value = false
    }
  }

  return {
    currentAddress,
    jwtToken,
    isConnecting,
    isConnected,
    connect,
  }
}
