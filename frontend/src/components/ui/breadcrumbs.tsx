import * as React from 'react'
import { Link, useMatches } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { api } from '../../api/client'

const labels: Record<string,string> = {
  '/': 'Home',
  '/login': 'Login',
  '/checklist': 'Checklist',
  '/advisor/inbox': 'Advisor Inbox',
  '/admin/users': 'Admin • Users',
}

export function Breadcrumbs() {
  const matches = useMatches() as Array<{ pathname: string }>
  const { data: me } = useQuery({ queryKey:['me'], queryFn: ()=> api('/me') })
  const role = me?.role
  const items = matches
    .filter((m) => m.pathname !== '/')
    .map((m) => {
      const p = m.pathname
      let label = labels[p] || p
      if (p.startsWith('/documents/')) label = 'Document'
      if (p.startsWith('/admin') && role!=='admin' && role!=='superadmin') return null
      if (p.startsWith('/advisor') && !['advisor','chair','admin','superadmin'].includes(role)) return null
      return { path: p, label }
    })
    .filter(Boolean) as {path:string,label:string}[]
  if (items.length === 0) return null
  return (
    <nav className="text-sm text-gray-600 py-1">
      {items.map((it, idx) => (
        <span key={it.path}>
          <Link to={it.path} className="underline">{it.label}</Link>
          {idx < items.length - 1 ? ' / ' : ''}
        </span>
      ))}
    </nav>
  )
}
