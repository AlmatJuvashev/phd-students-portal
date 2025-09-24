import { createRootRouteWithContext } from '@tanstack/react-router'

export type Role = 'superadmin'|'admin'|'advisor'|'chair'|'student'
export type AuthState = { token: string|null, role: Role|null }

export function decodeJwtRole(token: string|null): Role|null {
  if (!token) return null
  try {
    const parts = token.split('.')
    const payload = JSON.parse(atob(parts[1].replace(/-/g,'+').replace(/_/g,'/')))
    return payload.role as Role || null
  } catch { return null }
}

export function getAuth(): AuthState {
  const token = localStorage.getItem('token')
  const role = decodeJwtRole(token)
  return { token, role }
}

export function requireAuth() {
  const { token } = getAuth()
  if (!token) throw new Error('UNAUTHENTICATED')
}

export function requireRole(...allowed: Role[]) {
  const { role } = getAuth()
  if (!role || !allowed.includes(role)) throw new Error('FORBIDDEN')
}
