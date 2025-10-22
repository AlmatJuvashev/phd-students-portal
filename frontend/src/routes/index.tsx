import React, { lazy, Suspense } from 'react'
import { createBrowserRouter } from 'react-router-dom'
import { AppLayout } from '@/pages/layout'
import { ProtectedRoute } from '@/components/auth/ProtectedRoute'
import { useAuth } from '@/contexts/AuthContext'

const LoginPage = lazy(() => import('@/pages/login').then(m => ({ default: m.LoginPage })))
const DoctoralJourney = lazy(() => import('@/pages/doctoral.journey').then(m => ({ default: m.DoctoralJourney })))
const HomePage = lazy(() => import('@/pages/home').then(m => ({ default: m.HomePage })))
const ContactsPage = lazy(() => import('@/pages/contacts').then(m => ({ default: m.ContactsPage })))
const AdminUsers = lazy(() => import('@/pages/admin.users').then(m => ({ default: m.AdminUsers })))
const AdvisorInbox = lazy(() => import('@/pages/advisor.inbox').then(m => ({ default: m.AdvisorInbox })))
const Dashboard = lazy(() => import('@/pages/dashboard').then(m => ({ default: m.Dashboard })))
const ForgotPassword = lazy(() => import('@/pages/forgot').then(m => ({ default: m.ForgotPassword })))
const ResetPassword = lazy(() => import('@/pages/reset').then(m => ({ default: m.ResetPassword })))

const WithSuspense = (el: React.ReactNode) => (
  <Suspense fallback={<div className="p-4 text-sm">Loading…</div>}>{el}</Suspense>
)

function PublicOnly({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth()
  if (isLoading) return <div className="p-4 text-sm">Loading…</div>
  if (user) return WithSuspense(<DoctoralJourney />)
  return <>{children}</>
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      { index: true, element: WithSuspense(<HomePage />) },
      {
        path: 'journey',
        element: (
          <ProtectedRoute>
            {WithSuspense(<DoctoralJourney />)}
          </ProtectedRoute>
        ),
      },
      { path: 'contacts', element: WithSuspense(<ContactsPage />) },
      {
        path: 'login',
        element: <PublicOnly>{WithSuspense(<LoginPage />)}</PublicOnly>,
      },
      { path: 'forgot-password', element: WithSuspense(<ForgotPassword />) },
      { path: 'reset-password', element: WithSuspense(<ResetPassword />) },
      {
        path: 'admin/users',
        element: (
          <ProtectedRoute requiredRole="admin">
            {WithSuspense(<AdminUsers />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'advisor/inbox',
        element: (
          <ProtectedRoute requiredRole="advisor">
            {WithSuspense(<AdvisorInbox />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'dashboard',
        element: (
          <ProtectedRoute>
            {WithSuspense(<Dashboard />)}
          </ProtectedRoute>
        ),
      },
    ],
  },
])
