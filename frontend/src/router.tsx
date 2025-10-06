import React from "react";
import { useTranslation } from "react-i18next";
import { createBrowserRouter, redirect, useParams } from "react-router-dom";
import { AppLayout } from "./pages/layout";
import { LoginPage } from "./pages/login";
import { Dashboard } from "./pages/dashboard";
import { ForgotPassword } from "./pages/forgot";
import { ResetPassword } from "./pages/reset";
import { AdminUsers } from "./pages/admin.users";
import { AdvisorInbox } from "./pages/advisor.inbox";
import { DocumentDetail } from "./pages/document.detail";
import { DoctoralJourney } from "./pages/doctoral.journey";
import type { Role } from "./auth/auth";

function getAuth() {
  const token = localStorage.getItem("token");
  let role: Role | null = null;
  try {
    if (token) {
      const p = JSON.parse(
        atob(token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/"))
      );
      role = p.role as Role;
    }
  } catch {}
  return { token, role };
}

function guardAuth() {
  const { token } = getAuth();
  if (!token) return redirect("/login");
  return null;
}

function guardRole(allowed: Role[]) {
  const { token, role } = getAuth();
  if (!token) return redirect("/login");
  if (!role || !allowed.includes(role)) {
    throw new Response("Forbidden", { status: 403, statusText: "FORBIDDEN" });
  }
  return null;
}

function Forbidden() {
  const { t: T } = useTranslation("common");
  return (
    <div className="max-w-lg mx-auto mt-10">
      <h2 className="text-xl font-semibold">403 — Forbidden</h2>
      <p className="text-sm text-gray-600">You don’t have access to this page.</p>
    </div>
  );
}

function DocumentDetailPage() {
  const { docId = "" } = useParams() as { docId?: string };
  return <DocumentDetail docId={docId} />;
}

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppLayout />,
    children: [
      { index: true, element: <DoctoralJourney />, loader: guardAuth },
      { path: "login", element: <LoginPage /> },
      { path: "forgot-password", element: <ForgotPassword /> },
      { path: "reset-password", element: <ResetPassword /> },
      {
        path: "admin/users",
        element: <AdminUsers />,
        loader: () => guardRole(["admin", "superadmin"]),
        errorElement: <Forbidden />,
      },
      {
        path: "advisor/inbox",
        element: <AdvisorInbox />,
        loader: () => guardRole(["advisor", "chair", "admin", "superadmin"]),
        errorElement: <Forbidden />,
      },
      {
        path: "journey",
        element: <DoctoralJourney />,
        loader: () =>
          guardRole(["student", "advisor", "chair", "admin", "superadmin"]),
        errorElement: <Forbidden />,
      },
      {
        path: "documents/:docId",
        loader: () =>
          guardRole(["student", "advisor", "chair", "admin", "superadmin"]),
        element: <DocumentDetailPage />,
        errorElement: <Forbidden />,
      },
    ],
  },
]);
