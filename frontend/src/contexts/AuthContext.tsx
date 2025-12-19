import React, { createContext, useContext, useMemo } from 'react'
import { useQueryClient, useQuery } from '@tanstack/react-query'
import { api } from '@/api/client'

export type Role = 'student' | 'advisor' | 'secretary' | 'chair' | 'admin' | 'superadmin'

export interface User {
  id: string
  full_name?: string
  first_name?: string
  last_name?: string
  email?: string
  role: Role
  is_superadmin?: boolean
  progress?: Record<string, any>
  completedNodes?: string[]
  phone?: string
  bio?: string
  address?: string
  date_of_birth?: string
  avatar_url?: string
  program?: string
  specialty?: string
  department?: string
  cohort?: string
}

interface AuthContextType {
  user: User | null
  isLoading: boolean
  token: string | null
  login: (credentials: { username?: string; email?: string; password: string }) => Promise<{ role: string; is_superadmin: boolean }>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const qc = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['me'],
    queryFn: () => api('/me'),
    retry: 0,
  })

  const value = useMemo<AuthContextType>(() => ({
    user: (data as User) ?? null,
    isLoading,
    token: null, // Token is handled via cookies
    login: async (credentials) => {
      const res = await api('/auth/login', {
        method: 'POST',
        body: JSON.stringify({
          username: credentials.username ?? credentials.email,
          password: credentials.password,
        }),
      })
      // Token is set in cookie by server
      // Refresh user info
      await qc.invalidateQueries({ queryKey: ['me'] })
      return { role: res.role, is_superadmin: res.is_superadmin }
    },
    logout: async () => {
      try {
        await api.post('/auth/logout')
      } catch (e) {
        console.error("Logout failed:", e)
      }
      qc.removeQueries({ queryKey: ['me'] })
      // Soft reload to reset app state
      window.location.href = '/login'
    },
  }), [data, isLoading, qc])

  return (
    <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
  )
}

export const useAuth = () => {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
