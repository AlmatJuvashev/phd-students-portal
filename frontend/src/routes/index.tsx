import React, { lazy, Suspense } from "react";
import { createBrowserRouter, Navigate } from "react-router-dom";
import RouteErrorBoundary from "@/pages/errors/RouteErrorBoundary";
import NotFound from "@/pages/errors/NotFound";
import { AppLayout } from "@/pages/layout";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { useAuth } from "@/contexts/AuthContext";
import { useTenantServices, OptionalService } from "@/contexts/TenantServicesContext";
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
const StudentDashboard = lazy(() =>
  import("@/features/student-portal/StudentDashboard").then((m) => ({ default: m.StudentDashboard }))
);
const ForgotPassword = lazy(() =>
  import("@/pages/ForgotPasswordPage").then((m) => ({ default: m.ForgotPasswordPage }))
);
const ResetPassword = lazy(() =>
  import("@/pages/ResetPasswordPage").then((m) => ({ default: m.ResetPasswordPage }))
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
const CalendarPage = lazy(() =>
  import("@/features/calendar").then((m) => ({ default: m.CalendarPage }))
);
const AnalyticsDashboard = lazy(() =>
  import("@/features/analytics/AnalyticsDashboard").then((m) => ({ default: m.AnalyticsDashboard }))
);
const SchedulerPage = lazy(() =>
  import("@/features/scheduler/SchedulerPage").then((m) => ({ default: m.default }))
);
const ProgramsPage = lazy(() =>
  import("@/features/curriculum/ProgramsPage").then((m) => ({ default: m.ProgramsPage }))
);
const CoursesPage = lazy(() =>
  import("@/features/curriculum/CoursesPage").then((m) => ({ default: m.CoursesPage }))
);
const EnrollmentsPage = lazy(() =>
  import("@/features/enrollments/EnrollmentsPage").then((m) => ({ default: m.EnrollmentsPage }))
);
const ItemBanksPage = lazy(() =>
  import("@/features/item-bank/BanksPage").then((m) => ({ default: m.BanksPage }))
);
const ItemBankItemsPage = lazy(() =>
  import("@/features/item-bank/BankItemsPage").then((m) => ({ default: m.BankItemsPage }))
);
const CourseBuilder = lazy(() =>
  import("@/features/studio/CourseBuilder").then((m) => ({ default: m.CourseBuilder }))
);
const ProgramBuilderPage = lazy(() =>
  import("@/features/studio/ProgramBuilderPage").then((m) => ({ default: m.ProgramBuilderPage }))
);
const TeacherDashboard = lazy(() =>
  import("@/features/teacher/TeacherDashboard").then((m) => ({ default: m.TeacherDashboard }))
);
const TeacherCoursesPage = lazy(() =>
  import("@/features/teacher/TeacherCoursesPage").then((m) => ({ default: m.TeacherCoursesPage }))
);
const TeacherCourseDetail = lazy(() =>
  import("@/features/teacher/TeacherCourseDetail").then((m) => ({ default: m.TeacherCourseDetail }))
);
const TeacherGradingPage = lazy(() =>
  import("@/features/teacher/TeacherGradingPage").then((m) => ({ default: m.TeacherGradingPage }))
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
  if (user) {
    // Superadmins should go to /superadmin, not DoctoralJourney
    if (user.is_superadmin || user.role === 'superadmin') {
      return <Navigate to="/superadmin" replace />;
    }
    return WithSuspense(<DoctoralJourney />);
  }
  return <>{children}</>;
}

// Service-gated route - shows "service not available" for disabled services
function ServiceProtectedRoute({ 
  children, 
  service 
}: { 
  children: React.ReactNode; 
  service: OptionalService;
}) {
  const { isServiceEnabled, isLoading } = useTenantServices();
  
  if (isLoading) return <div className="p-4 text-sm">Loading…</div>;
  
  if (!isServiceEnabled(service)) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] p-8 text-center">
        <h2 className="text-xl font-semibold mb-2">Service Not Available</h2>
        <p className="text-muted-foreground">
          This feature is not enabled for your institution.
        </p>
      </div>
    );
  }
  
  return <>{children}</>;
}


function IndexRoute() {
  const { user } = useAuth();
  if (user) {
    return WithSuspense(<HomePage />);
  }
  return WithSuspense(<LandingPage />);
}

export const router = createBrowserRouter([
  // App routes (constrained width via AppLayout)
  {
    path: "/",
    element: <AppLayout />,
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: <IndexRoute /> },
      {
        path: "journey",
        element: (
          <ProtectedRoute>{WithSuspense(<DoctoralJourney />)}</ProtectedRoute>
        ),
      },
      {
        path: "chat",
        element: (
          <ProtectedRoute>
            <ServiceProtectedRoute service="chat">
              {WithSuspense(<ChatPage />)}
            </ServiceProtectedRoute>
          </ProtectedRoute>
        ),
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
        path: "student/dashboard",
        element: (
          <ProtectedRoute requiredRole="student">
            {WithSuspense(<StudentDashboard />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "profile",
        element: <ProtectedRoute>{WithSuspense(<ProfilePage />)}</ProtectedRoute>,
      },
      {
        path: "calendar",
        element: (
          <ProtectedRoute>
            <ServiceProtectedRoute service="calendar">
              {WithSuspense(<CalendarPage />)}
            </ServiceProtectedRoute>
          </ProtectedRoute>
        ),
      },
    ],
  },
  // Admin routes (full-width layout)
  {
    path: "/admin",
    element: (
      <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
        {WithSuspense(<AdminLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<AdminDashboard />) },
      {
        path: "create-admins",
        element: (
          <ProtectedRoute requiredAnyRole={[]}>
            {WithSuspense(<CreateAdmins />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "create-students",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<CreateStudents />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "create-advisors",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<CreateAdvisors />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "create-users",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<CreateUsers />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "contacts",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<ContactsAdminPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "students-monitor",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            {WithSuspense(<StudentsMonitorPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "students-monitor/:id",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            {WithSuspense(<StudentDetailPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "notifications",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            {WithSuspense(<NotificationsPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "users",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<AdminUsersPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "chat-rooms",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            <ServiceProtectedRoute service="chat">
              {WithSuspense(<ChatRoomsAdminPage />)}
            </ServiceProtectedRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: "dictionaries",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<DictionariesPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "calendar",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            <ServiceProtectedRoute service="calendar">
              {WithSuspense(<CalendarPage />)}
            </ServiceProtectedRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: "analytics",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "chair"]}>
            {WithSuspense(<AnalyticsDashboard />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "scheduler",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
             {WithSuspense(<SchedulerPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "programs",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<ProgramsPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "courses",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<CoursesPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "enrollments",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<EnrollmentsPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "item-banks",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<ItemBanksPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "item-banks/:bankId",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<ItemBankItemsPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "studio/courses/:courseId/builder",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<CourseBuilder />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "studio/programs/:programId/builder",
        element: (
          <ProtectedRoute requiredAnyRole={["admin"]}>
            {WithSuspense(<ProgramBuilderPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "teacher/dashboard",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            {WithSuspense(<TeacherDashboard />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "teacher/courses",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            {WithSuspense(<TeacherCoursesPage />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "teacher/courses/:courseId",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            {WithSuspense(<TeacherCourseDetail />)}
          </ProtectedRoute>
        ),
      },
      {
        path: "teacher/grading",
        element: (
          <ProtectedRoute requiredAnyRole={["admin", "advisor"]}>
            {WithSuspense(<TeacherGradingPage />)}
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
