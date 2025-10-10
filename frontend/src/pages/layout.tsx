import { Link, Outlet, useLocation } from "react-router-dom";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { APP_NAME } from "../config";
import { useTranslation } from "react-i18next";

export function AppLayout({ children }: { children?: React.ReactNode }) {
  const { t: T, i18n } = useTranslation("common");
  const { data: me } = useQuery({
    queryKey: ["me"],
    queryFn: () => api("/me"),
  });
  const authed = !!me;
  const role = me?.role;
  const { pathname } = useLocation();
  const active = (p: string) =>
    pathname === p ? "font-semibold" : "text-muted-foreground hover:underline";

  return (
    <div className="max-w-4xl mx-auto p-4">
      <header className="flex items-center justify-between py-2">
        <h1 className="font-semibold">{APP_NAME}</h1>
        <nav className="flex gap-3 text-sm">
          {authed && (
            <Link to="/" className={active("/")}>
              {T("nav.home")}
            </Link>
          )}
          {/* Checklist removed; Journey is the primary view */}
          {authed && (
            <Link to="/journey" className={active("/journey")}>
              {T("nav.journey")}
            </Link>
          )}
          {authed &&
            (role === "advisor" ||
              role === "chair" ||
              role === "admin" ||
              role === "superadmin") && (
              <Link to="/advisor/inbox" className={active("/advisor/inbox")}>
                {T("nav.inbox")}
              </Link>
            )}
          {authed && (role === "admin" || role === "superadmin") && (
            <Link to="/admin/users" className={active("/admin/users")}>
              {T("nav.admin")}
            </Link>
          )}
          {authed ? (
            <button
              className="text-muted-foreground hover:underline"
              onClick={() => {
                localStorage.removeItem("token");
                location.href = "/login";
              }}
            >
              {T("nav.logout")}
            </button>
          ) : (
            pathname !== "/login" ? (
              <Link to="/login" className="text-muted-foreground hover:underline">
                {T("nav.login")}
              </Link>
            ) : null
          )}
          <select
            className="ml-2 border rounded px-1 py-0.5 text-xs"
            value={i18n.language}
            onChange={(e) => i18n.changeLanguage(e.target.value)}
          >
            <option value="ru">RU</option>
            <option value="kz">KZ</option>
            <option value="en">EN</option>
          </select>
        </nav>
      </header>
      <main>{children ?? <Outlet />}</main>
    </div>
  );
}
