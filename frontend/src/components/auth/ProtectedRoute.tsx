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

  // Check against active_role if available, otherwise fallback to role
  if (requiredRole || (requiredAnyRole && requiredAnyRole.length > 0)) {
    const allowed = new Set<string>([ 'superadmin', 'admin' ])
    if (requiredRole) allowed.add(requiredRole)
    if (requiredAnyRole) requiredAnyRole.forEach(r => allowed.add(r))
    
    // Check active_role first
    const currentRole = user.active_role || user.role
    const roleOk = allowed.has(currentRole)
    
    if (!roleOk) {
      return (
        <div className="flex flex-col items-center justify-center min-h-[40vh] gap-2">
          <h1 className="text-xl font-semibold">Доступ запрещён</h1>
          <p className="text-sm text-muted-foreground">
             Ваша текущая роль ({currentRole}) не имеет доступа к этой странице. 
             {user.available_roles && user.available_roles.length > 1 && " Попробуйте переключить роль."}
          </p>
        </div>
      )
    }
  }

  return <>{children}</>
}
