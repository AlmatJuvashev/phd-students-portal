import React, { createContext, useContext, useMemo } from 'react'
import { useQueryClient, useQuery } from '@tanstack/react-query'
import { api } from '@/api/client'

export type Role = 'student' | 'advisor' | 'secretary' | 'chair' | 'admin' | 'superadmin'

export interface User {
  id: string
  full_name?: string
  email?: string
  role: Role
  progress?: Record<string, any>
  completedNodes?: string[]
}

interface AuthContextType {
  user: User | null
  isLoading: boolean
  login: (credentials: { username?: string; email?: string; password: string }) => Promise<void>
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
    login: async (credentials) => {
      const res = await api('/auth/login', {
        method: 'POST',
        body: JSON.stringify({
          username: credentials.username ?? credentials.email,
          password: credentials.password,
        }),
      })
      localStorage.setItem('token', res.token)
      // Refresh user info
      await qc.invalidateQueries({ queryKey: ['me'] })
    },
    logout: () => {
      localStorage.removeItem('token')
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

