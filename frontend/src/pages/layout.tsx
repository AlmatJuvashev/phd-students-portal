import { Link, useRouterState } from '@tanstack/react-router'
import { Breadcrumbs } from '../components/ui/breadcrumbs'
import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../api/client'
import { APP_NAME } from '../config'

export function AppLayout({ children }: { children?: React.ReactNode }) {
  const { data: me } = useQuery({ queryKey: ['me'], queryFn: ()=> api('/me') })
  const authed = !!me
  const role = me?.role
  const pathname = location.pathname
  const active = (p:string) => pathname===p ? 'font-semibold underline' : 'underline'

  return (
    <div className="max-w-4xl mx-auto p-4">
      <header className="flex items-center justify-between py-2">
        <h1 className="font-semibold">{APP_NAME}</h1>
        <nav className="flex gap-3 text-sm">
          {authed && <Link to="/" className={active("/")}>Home</Link>}
          {authed && <Link to="/checklist" className={active("/")}>Checklist</Link>}
          {authed && (role==='advisor' || role==='chair' || role==='admin' || role==='superadmin') && <Link to="/advisor/inbox" className={active("/")}>Inbox</Link>}
          {authed && (role==='admin' || role==='superadmin') && <Link to="/admin/users" className={active("/")}>Admin</Link>}
          {authed ? (
            <button
              className={active("/")}
              onClick={() => { localStorage.removeItem('token'); location.href = '/login' }}
            >Logout</button>
          ) : (
            <Link to="/login" className={active("/")}>Login</Link>
          )}
        </nav>
      </header>
      <main>
        <Breadcrumbs />
        {children}
      </main>
    </div>
  )
}
