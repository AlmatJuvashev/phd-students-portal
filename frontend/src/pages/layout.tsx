import { Link, Outlet, useLocation } from "react-router-dom";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { APP_NAME } from "../config";
import { useTranslation } from "react-i18next";
import { Globe } from "lucide-react";
import { DropdownMenu, DropdownItem } from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";

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
          <DropdownMenu
            trigger={
              <Button
                variant="ghost"
                size="sm"
                className="gap-2 h-8 px-3 text-sm font-medium hover:bg-primary/10 transition-all duration-200"
              >
                <Globe className="h-4 w-4" />
                <span className="hidden sm:inline">{i18n.language.toUpperCase()}</span>
              </Button>
            }
          >
            <DropdownItem
              onClick={() => i18n.changeLanguage("ru")}
            >
              <div
                className={`flex items-center gap-3 ${
                  i18n.language === "ru" ? "font-semibold text-primary" : ""
                }`}
              >
                <span className="text-lg">🇷🇺</span>
                <span>Русский</span>
              </div>
            </DropdownItem>
            <DropdownItem
              onClick={() => i18n.changeLanguage("kz")}
            >
              <div
                className={`flex items-center gap-3 ${
                  i18n.language === "kz" ? "font-semibold text-primary" : ""
                }`}
              >
                <span className="text-lg">🇰🇿</span>
                <span>Қазақша</span>
              </div>
            </DropdownItem>
            <DropdownItem
              onClick={() => i18n.changeLanguage("en")}
            >
              <div
                className={`flex items-center gap-3 ${
                  i18n.language === "en" ? "font-semibold text-primary" : ""
                }`}
              >
                <span className="text-lg">🇬🇧</span>
                <span>English</span>
              </div>
            </DropdownItem>
          </DropdownMenu>
        </nav>
      </header>
      <main>{children ?? <Outlet />}</main>
    </div>
  );
}
