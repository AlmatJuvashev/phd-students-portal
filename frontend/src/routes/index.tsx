import React, { lazy, Suspense } from 'react'
import { createBrowserRouter } from 'react-router-dom'
import { AppLayout } from '@/pages/layout'
import { ProtectedRoute } from '@/components/auth/ProtectedRoute'
import { useAuth } from '@/contexts/AuthContext'
import { AdminLayout } from '@/layouts/AdminLayout'

const LoginPage = lazy(() => import('@/pages/login').then(m => ({ default: m.LoginPage })))
const DoctoralJourney = lazy(() => import('@/pages/doctoral.journey').then(m => ({ default: m.DoctoralJourney })))
const HomePage = lazy(() => import('@/pages/home').then(m => ({ default: m.HomePage })))
const ContactsPage = lazy(() => import('@/pages/contacts').then(m => ({ default: m.ContactsPage })))
const AdminUsers = lazy(() => import('@/pages/admin.users').then(m => ({ default: m.AdminUsers })))
const AdminDashboard = lazy(() => import('@/pages/dashboard').then(m => ({ default: m.Dashboard })))
const CreateAdmins = lazy(() => import('@/pages/admin/CreateAdmins').then(m => ({ default: m.CreateAdmins })))
const CreateUsers = lazy(() => import('@/pages/admin/CreateUsers').then(m => ({ default: m.CreateUsers })))
const CreateStudents = lazy(() => import('@/pages/admin/CreateStudents').then(m => ({ default: m.CreateStudents })))
const CreateAdvisors = lazy(() => import('@/pages/admin/CreateAdvisors').then(m => ({ default: m.CreateAdvisors })))
const StudentProgress = lazy(() => import('@/pages/admin/StudentProgress').then(m => ({ default: m.StudentProgress })))
const StudentsMonitorPage = lazy(() => import('@/features/students-monitor/StudentsMonitorPage').then(m => ({ default: m.StudentsMonitorPage })))
const StudentDetailPage = lazy(() => import('@/features/students-monitor/pages/StudentDetailPage').then(m => ({ default: m.StudentDetailPage })))
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
  // App routes (constrained width via AppLayout)
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
  // Admin routes (full-width layout)
  {
    path: '/admin',
    element: (
      <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
        {WithSuspense(<AdminLayout />)}
      </ProtectedRoute>
    ),
    children: [
      { index: true, element: WithSuspense(<AdminDashboard />) },
      {
        path: 'create-admins',
        element: (
          <ProtectedRoute requiredAnyRole={["superadmin"]}>
            {WithSuspense(<CreateAdmins />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'create-students',
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<CreateStudents />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'create-advisors',
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<CreateAdvisors />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'create-users',
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<CreateUsers />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'student-progress',
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<StudentProgress />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'students-monitor',
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
            {WithSuspense(<StudentsMonitorPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'students-monitor/:id',
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
            {WithSuspense(<StudentDetailPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: 'users',
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<AdminUsers />)}
          </ProtectedRoute>
        ),
      },
    ],
  },
])
