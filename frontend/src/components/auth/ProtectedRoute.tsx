import React from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth, Role } from '@/contexts/AuthContext'

interface ProtectedRouteProps {
  children: React.ReactNode
  fallback?: React.ReactNode
  requiredRole?: 'student' | 'advisor' | 'secretary' | 'chair' | 'admin'
  requiredAnyRole?: Role[]
}

export function ProtectedRoute({ children, fallback, requiredRole, requiredAnyRole }: ProtectedRouteProps) {
  const { user, isLoading } = useAuth()
  const location = useLocation()

  if (isLoading) {
    return fallback || (
      <div className="flex items-center justify-center min-h-[40vh] text-sm text-muted-foreground">Loading…</div>
    )
  }

  if (!user) {
    return (
      <Navigate to="/login" state={{ from: location.pathname }} replace />
    )
  }

  // Backward-compat: requiredRole OR new requiredAnyRole
  if (requiredRole || (requiredAnyRole && requiredAnyRole.length > 0)) {
    const allowed = new Set<Role>([ 'superadmin', 'admin' ])
    if (requiredRole) allowed.add(requiredRole as Role)
    if (requiredAnyRole) requiredAnyRole.forEach(r => allowed.add(r))
    const roleOk = allowed.has(user.role)
    if (!roleOk) {
      return (
        <div className="flex flex-col items-center justify-center min-h-[40vh] gap-2">
          <h1 className="text-xl font-semibold">Доступ запрещён</h1>
          <p className="text-sm text-muted-foreground">У вас нет прав для просмотра этой страницы</p>
        </div>
      )
    }
  }

  return <>{children}</>
}
