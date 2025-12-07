import React, { lazy, Suspense } from "react";
import { createBrowserRouter } from "react-router-dom";
import RouteErrorBoundary from "@/pages/errors/RouteErrorBoundary";
import NotFound from "@/pages/errors/NotFound";
import { AppLayout } from "@/pages/layout";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useAuth } from "@/contexts/AuthContext";
import { AdminLayout } from "@/layouts/AdminLayout";

const LoginPage = lazy(() =>
  import("@/pages/login").then((m) => ({ default: m.LoginPage }))
);
const DoctoralJourney = lazy(() =>
  import("@/pages/doctoral.journey").then((m) => ({
    default: m.DoctoralJourney,
  }))
);
const HomePage = lazy(() =>
  import("@/pages/home").then((m) => ({ default: m.HomePage }))
);
const LandingPage = lazy(() =>
  import("@/pages/Landing").then((m) => ({ default: m.Landing }))
);
const ContactsPage = lazy(() =>
  import("@/pages/contacts").then((m) => ({ default: m.ContactsPage }))
);
const AdminUsers = lazy(() =>
  import("@/pages/admin.users").then((m) => ({ default: m.AdminUsers }))
);
const AdminUsersPage = lazy(() =>
  import("@/pages/AdminUsersPage").then((m) => ({ default: m.AdminUsersPage }))
);
const AdminDashboard = lazy(() =>
  import("@/pages/dashboard").then((m) => ({ default: m.Dashboard }))
);
const CreateAdmins = lazy(() =>
  import("@/pages/admin/CreateAdmins").then((m) => ({
    default: m.CreateAdmins,
  }))
);
const ContactsAdminPage = lazy(() =>
  import("@/pages/admin/ContactsAdminPage").then((m) => ({
    default: m.ContactsAdminPage,
  }))
);
const CreateUsers = lazy(() =>
  import("@/pages/admin/CreateUsers").then((m) => ({ default: m.CreateUsers }))
);
const CreateStudents = lazy(() =>
  import("@/pages/admin/CreateStudents").then((m) => ({
    default: m.CreateStudents,
  }))
);
const CreateAdvisors = lazy(() =>
  import("@/pages/admin/CreateAdvisors").then((m) => ({
    default: m.CreateAdvisors,
  }))
);
const StudentsMonitorPage = lazy(() =>
  import("@/features/students-monitor/StudentsMonitorPage").then((m) => ({
    default: m.StudentsMonitorPage,
  }))
);
const StudentDetailPage = lazy(() =>
  import("@/features/students-monitor/pages/StudentDetailPage").then((m) => ({
    default: m.StudentDetailPage,
  }))
);
const NotificationsPage = lazy(() =>
  import("@/pages/admin/NotificationsPage").then((m) => ({
    default: m.default,
  }))
);
const AdvisorInbox = lazy(() =>
  import("@/pages/advisor.inbox").then((m) => ({ default: m.AdvisorInbox }))
);
const Dashboard = lazy(() =>
  import("@/pages/dashboard").then((m) => ({ default: m.Dashboard }))
);
const ForgotPassword = lazy(() =>
  import("@/pages/forgot").then((m) => ({ default: m.ForgotPassword }))
);
const ResetPassword = lazy(() =>
  import("@/pages/reset").then((m) => ({ default: m.ResetPassword }))
);
const ChatPage = lazy(() =>
  import("@/pages/chat").then((m) => ({ default: m.ChatPage }))
);
const ChatRoomsAdminPage = lazy(() =>
  import("@/features/chat-admin/ChatRoomsAdminPage").then((m) => ({
    default: m.ChatRoomsAdminPage,
  }))
);
const DictionariesPage = lazy(() =>
  import("@/features/admin/dictionaries/DictionariesPage").then((m) => ({
    default: m.DictionariesPage,
  }))
);
const ProfilePage = lazy(() =>
  import("@/pages/profile").then((m) => ({ default: m.ProfilePage }))
);
const VerifyEmailPage = lazy(() =>
  import("@/pages/verify-email").then((m) => ({ default: m.VerifyEmailPage }))
);
const CalendarView = lazy(() =>
  import("@/features/calendar").then((m) => ({ default: m.CalendarView }))
);
const AnalyticsDashboard = lazy(() =>
  import("@/features/analytics/AnalyticsDashboard").then((m) => ({ default: m.AnalyticsDashboard }))
);

// Superadmin pages
const SuperadminLayout = lazy(() =>
  import("@/layouts/SuperadminLayout").then((m) => ({ default: m.SuperadminLayout }))
);
const TenantsPage = lazy(() =>
  import("@/features/superadmin/tenants/TenantsPage").then((m) => ({ default: m.TenantsPage }))
);
const AdminsPage = lazy(() =>
  import("@/features/superadmin/admins/AdminsPage").then((m) => ({ default: m.AdminsPage }))
);
const LogsPage = lazy(() =>
  import("@/features/superadmin/logs/LogsPage").then((m) => ({ default: m.LogsPage }))
);
const SettingsPage = lazy(() =>
  import("@/features/superadmin/settings/SettingsPage").then((m) => ({ default: m.SettingsPage }))
);


const WithSuspense = (el: React.ReactNode) => (
  <Suspense fallback={<div className="p-4 text-sm">Loading…</div>}>
    {el}
  </Suspense>
);

function PublicOnly({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth();
  if (isLoading) return <div className="p-4 text-sm">Loading…</div>;
  if (user) return WithSuspense(<DoctoralJourney />);
  return <>{children}</>;
}

export const router = createBrowserRouter([
  // App routes (constrained width via AppLayout)
  {
    path: "/",
    element: <AppLayout />,
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<LandingPage />) },
      {
        path: "journey",
        element: (
          <ProtectedRoute>{WithSuspense(<DoctoralJourney />)}</ProtectedRoute>
        ),
      },
      {
        path: "chat",
        element: <ProtectedRoute>{WithSuspense(<ChatPage />)}</ProtectedRoute>,
      },
      { path: "contacts", element: WithSuspense(<ContactsPage />) },
      {
        path: "login",
        element: <PublicOnly>{WithSuspense(<LoginPage />)}</PublicOnly>,
      },
      { path: "forgot-password", element: WithSuspense(<ForgotPassword />) },
      { path: "reset-password", element: WithSuspense(<ResetPassword />) },
      { path: "verify-email", element: WithSuspense(<VerifyEmailPage />) },
      { path: "*", element: <NotFound /> },
      {
        path: "advisor/inbox",
        element: (
          <ProtectedRoute requiredRole="advisor">
            {WithSuspense(<AdvisorInbox />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "dashboard",
        element: <ProtectedRoute>{WithSuspense(<Dashboard />)}</ProtectedRoute>,
      },
      {
        path: "profile",
        element: <ProtectedRoute>{WithSuspense(<ProfilePage />)}</ProtectedRoute>,
      },
      {
        path: "calendar",
        element: <ProtectedRoute>{WithSuspense(<CalendarView />)}</ProtectedRoute>,
      },
    ],
  },
  // Admin routes (full-width layout)
  {
    path: "/admin",
    element: (
      <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
        {WithSuspense(<AdminLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<AdminDashboard />) },
      {
        path: "create-admins",
        element: (
          <ProtectedRoute requiredAnyRole={["superadmin"]}>
            {WithSuspense(<CreateAdmins />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "create-students",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<CreateStudents />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "create-advisors",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<CreateAdvisors />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "create-users",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<CreateUsers />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "contacts",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<ContactsAdminPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "students-monitor",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
            {WithSuspense(<StudentsMonitorPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "students-monitor/:id",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
            {WithSuspense(<StudentDetailPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "notifications",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
            {WithSuspense(<NotificationsPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "users",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<AdminUsersPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "chat-rooms",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<ChatRoomsAdminPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "chat-rooms",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<ChatRoomsAdminPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "dictionaries",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
            {WithSuspense(<DictionariesPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "calendar",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin", "advisor"]}>
            {WithSuspense(<CalendarView />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "analytics",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "superadmin", "chair"]}>
            {WithSuspense(<AnalyticsDashboard />)}
          </ProtectedRoute>
        ),
      },
      { path: "*", element: <NotFound /> },
    ],
  },
  // Superadmin routes (platform administration)
  {
    path: "/superadmin",
    element: WithSuspense(<SuperadminLayout />),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<TenantsPage />) },
      { path: "tenants", element: WithSuspense(<TenantsPage />) },
      { path: "admins", element: WithSuspense(<AdminsPage />) },
      { path: "logs", element: WithSuspense(<LogsPage />) },
      { path: "settings", element: WithSuspense(<SettingsPage />) },
      { path: "*", element: <NotFound /> },
    ],
  },
]);
