import React, { lazy, Suspense } from 'react'
import { createBrowserRouter } from 'react-router-dom'
import { AppLayout } from '@/pages/layout'

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

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      { index: true, element: WithSuspense(<HomePage />) },
      { path: 'journey', element: WithSuspense(<DoctoralJourney />) },
      { path: 'contacts', element: WithSuspense(<ContactsPage />) },
      { path: 'login', element: WithSuspense(<LoginPage />) },
      { path: 'forgot-password', element: WithSuspense(<ForgotPassword />) },
      { path: 'reset-password', element: WithSuspense(<ResetPassword />) },
      { path: 'admin/users', element: WithSuspense(<AdminUsers />) },
      { path: 'advisor/inbox', element: WithSuspense(<AdvisorInbox />) },
      { path: 'dashboard', element: WithSuspense(<Dashboard />) },
    ],
  },
])
