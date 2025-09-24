import { useQuery } from '@tanstack/react-query'
import { API_URL } from '../api/client'

export function useMe() {
  return useQuery({
    queryKey: ['me'],
    queryFn: async () => {
      const token = localStorage.getItem('token')
      const r = await fetch(`${API_URL}/me`, { headers: { Authorization: `Bearer ${token}` } })
      if (r.status === 401) throw new Error('UNAUTHENTICATED')
      return r.json()
    },
    staleTime: 5 * 60 * 1000,
  })
}
