export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase || 'http://localhost:8080/api/v1'

  const getHeaders = () => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    }
    if (process.client) {
      const token = localStorage.getItem('bn_near_token')
      if (token) {
        headers['Authorization'] = `Bearer ${token}`
      }
    }
    return headers
  }

  const fetchAPI = async <T = any>(endpoint: string, options: any = {}): Promise<{ data?: T; error?: any }> => {
    try {
      const res = await $fetch<{ data: T; error?: any }>(`${baseURL}${endpoint}`, {
        ...options,
        headers: {
          ...getHeaders(),
          ...options.headers
        }
      })
      return { data: res.data }
    } catch (err: any) {
      return { error: err.data?.error || { message: err.message } }
    }
  }

  return {
    baseURL,
    fetchAPI
  }
}
