export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8280/api'

export async function api<T = any>(path: string, opts: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('token')
  
  // Debug logging for auth issues
  if (path === '/notifications' || path.includes('/notifications')) {
    console.log('[API Debug] Notifications request:', {
      path,
      hasToken: !!token,
      tokenPreview: token ? `${token.substring(0, 20)}...` : 'null',
    })
  }
  
  const headers: Record<string, string> = {
    ...(opts.headers as Record<string, string> || {}),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }
  
  if (!headers['Content-Type'] && !(opts.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }
  const res = await fetch(`${API_URL}${path}`, { ...opts, headers })
  if (!res.ok) {
    const t = await res.text()
    // Enhanced error logging
    if (res.status === 401) {
      console.error('[API 401 Error]', {
        path,
        hasToken: !!token,
        response: t,
      })
    }
    throw new Error(t || res.statusText)
  }
  const type = res.headers.get('content-type') || ''
  return (type.includes('application/json') ? res.json() : res.text()) as Promise<T>
}

// HTTP method helpers
api.get = (path: string, opts?: RequestInit) => api(path, { ...opts, method: 'GET' })
api.post = (path: string, data?: any, opts?: RequestInit) => 
  api(path, { ...opts, method: 'POST', body: data ? JSON.stringify(data) : undefined })
api.put = (path: string, data?: any, opts?: RequestInit) => 
  api(path, { ...opts, method: 'PUT', body: data ? JSON.stringify(data) : undefined })
api.patch = (path: string, data?: any, opts?: RequestInit) => 
  api(path, { ...opts, method: 'PATCH', body: data ? JSON.stringify(data) : undefined })
api.delete = (path: string, opts?: RequestInit) => api(path, { ...opts, method: 'DELETE' })

