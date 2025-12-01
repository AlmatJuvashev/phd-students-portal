export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8280/api'

export async function api(path: string, opts: RequestInit = {}) {
  const token = localStorage.getItem('token')
  const headers = {
    'Content-Type': 'application/json',
    ...(opts.headers || {}),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }
  const res = await fetch(`${API_URL}${path}`, { ...opts, headers })
  if (!res.ok) {
    const t = await res.text()
    throw new Error(t || res.statusText)
  }
  const type = res.headers.get('content-type') || ''
  return type.includes('application/json') ? res.json() : res.text()
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

