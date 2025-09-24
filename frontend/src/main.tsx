import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider, createRootRoute, createRouter } from '@tanstack/react-router'
import './index.css'
import { ToastProvider } from './components/toast'
import { AppLayout } from './pages/layout'
import { requireAuth, requireRole } from './auth/auth'
import { LoginPage } from './pages/login'
import { Dashboard } from './pages/dashboard'
import { ForgotPassword } from './pages/forgot'
import { ResetPassword } from './pages/reset'
import { AdminUsers } from './pages/admin.users'
import { StudentChecklist } from './pages/checklist'
import { AdvisorInbox } from './pages/advisor.inbox'
import { DocumentDetail } from './pages/document.detail'

const rootRoute = createRootRoute({
  component: AppLayout,
})

const routeTree = rootRoute.addChildren([
  { path: '/', component: Dashboard },
  { path: '/login', component: LoginPage },
  { path: '/forgot-password', component: ForgotPassword },
  { path: '/reset-password', component: ResetPassword },
  { path: '/admin/users', component: AdminUsers, beforeLoad: () => { requireAuth(); requireRole('admin','superadmin') }, errorComponent: GuardErrorBoundary },
  { path: '/checklist', component: StudentChecklist, beforeLoad: () => { requireAuth(); requireRole('student','advisor','chair','admin','superadmin') }, errorComponent: GuardErrorBoundary },
  { path: '/advisor/inbox', component: AdvisorInbox, beforeLoad: () => { requireAuth(); requireRole('advisor','chair','admin','superadmin') }, errorComponent: GuardErrorBoundary },
  , { path: '/documents/$docId', component: () => {
      const params = router.useMatch({ from: '/documents/$docId' }).params as any
      return <DocumentDetail docId={params.docId} />
    }, beforeLoad: () => { requireAuth(); requireRole('student','advisor','chair','admin','superadmin') }, errorComponent: GuardErrorBoundary }
])


function GuardErrorBoundary({ error }: { error: any }) {
  if (String(error?.message).includes('UNAUTHENTICATED')) {
    location.href = '/login'
    return null
  }
  if (String(error?.message).includes('FORBIDDEN')) {
    return <div className="max-w-lg mx-auto mt-10"><h2 className="text-xl font-semibold">403 — Forbidden</h2><p className="text-sm text-gray-600">You don’t have access to this page.</p></div>
  }
  return <div>Something went wrong.</div>
}

const router = createRouter({ routeTree, defaultPreload: 'intent' })
const qc = new QueryClient()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={qc}>
            <ToastProvider>
        <RouterProvider router={router} />
      </ToastProvider>
    </QueryClientProvider>
  </React.StrictMode>
)
