import { useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'

export function useRequireAuth(redirectTo = '/login') {
  const { user, isLoading } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  useEffect(() => {
    if (!isLoading && !user) {
      navigate(redirectTo, { state: { from: location.pathname }, replace: true })
    }
  }, [user, isLoading, navigate, redirectTo, location])

  return { user, isLoading, isAuthenticated: !!user }
}

