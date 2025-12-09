import React from "react";
import { Outlet, NavLink, useLocation, Link, Navigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetTrigger,
  SheetContent,
  SheetTitle,
} from "@/components/ui/sheet";
import { cn } from "@/lib/utils";
import {
  Menu,
  Building2,
  Users,
  ScrollText,
  Settings,
  Shield,
  Home,
} from "lucide-react";
import { useTranslation } from "react-i18next";
import { LanguageSwitcher } from "@/components/layout/LanguageSwitcher";
import { UserMenu } from "@/components/layout/UserMenu";

function SuperadminSidebarNav({ collapsed }: { collapsed?: boolean }) {
  const { t } = useTranslation("common");
  const location = useLocation();
  const isActive = (path: string) => location.pathname.startsWith(path);

  const navItems = [
    {
      to: "/superadmin/tenants",
      icon: Building2,
      label: t("superadmin.sidebar.tenants", "Institutions"),
    },
    {
      to: "/superadmin/admins",
      icon: Users,
      label: t("superadmin.sidebar.admins", "Administrators"),
    },
    {
      to: "/superadmin/logs",
      icon: ScrollText,
      label: t("superadmin.sidebar.logs", "Activity Logs"),
    },
    {
      to: "/superadmin/settings",
      icon: Settings,
      label: t("superadmin.sidebar.settings", "Settings"),
    },
  ];

  return (
    <nav className={cn("p-4 space-y-2", collapsed && "px-2")}>
      {!collapsed && (
        <div className="text-xs uppercase text-muted-foreground px-2 mb-2">
          {t("superadmin.sidebar.platform", "Platform Management")}
        </div>
      )}
      {navItems.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          className={cn(
            "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted transition-colors",
            isActive(item.to) && "bg-muted font-medium",
            collapsed && "justify-center px-2"
          )}
          title={item.label}
        >
          <item.icon className="h-4 w-4" />
          {!collapsed && <span>{item.label}</span>}
        </NavLink>
      ))}
    </nav>
  );
}

export function SuperadminLayout() {
  const { t } = useTranslation("common");
  const { user, logout } = useAuth();
  const [collapsed, setCollapsed] = React.useState(false);

  // Redirect non-superadmins
  if (!user?.is_superadmin) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="flex min-h-screen">
      {/* Desktop sidebar */}
      <aside
        className={cn(
          "hidden md:block border-r transition-all duration-200",
          "bg-gradient-to-b from-violet-950/10 via-purple-900/5 to-background",
          collapsed ? "w-16" : "w-64"
        )}
      >
        <div
          className={cn(
            "h-14 flex items-center border-b",
            collapsed ? "justify-center" : "px-4 justify-between"
          )}
        >
          <div
            className={cn(
              "font-bold truncate text-lg flex items-center gap-2",
              collapsed && "text-base"
            )}
            title={t("superadmin.title", "Superadmin Panel")}
          >
            <Shield className="h-5 w-5 text-violet-500" />
            {!collapsed && (
              <span className="bg-gradient-to-r from-violet-600 to-purple-600 bg-clip-text text-transparent">
                {t("superadmin.title", "Superadmin")}
              </span>
            )}
          </div>
          {!collapsed && (
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setCollapsed(true)}
              aria-label={t("common.collapse", "Collapse")}
            >
              <span className="text-lg leading-none">«</span>
            </Button>
          )}
          {collapsed && (
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setCollapsed(false)}
              aria-label={t("common.expand", "Expand")}
              className="mt-2"
            >
              <span className="text-lg leading-none">»</span>
            </Button>
          )}
        </div>
        <SuperadminSidebarNav collapsed={collapsed} />

        {/* Back to Admin link */}
        <div className={cn("p-4 border-t mt-auto", collapsed && "px-2")}>
          <Link
            to="/admin"
            className={cn(
              "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted text-muted-foreground hover:text-foreground transition-colors",
              collapsed && "justify-center px-2"
            )}
            title={t("superadmin.back_to_admin", "Back to Admin")}
          >
            <Home className="h-4 w-4" />
            {!collapsed && (
              <span className="text-sm">
                {t("superadmin.back_to_admin", "Back to Admin")}
              </span>
            )}
          </Link>
        </div>
      </aside>

      {/* Mobile top bar with sheet menu */}
      <div className="md:hidden fixed top-0 left-0 right-0 h-14 border-b bg-background z-40 flex items-center justify-between px-3">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon" aria-label="Open menu">
              <Menu className="h-5 w-5" />
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="p-0 w-64">
            <SheetTitle className="sr-only">
              {t("superadmin.menu", "Menu")}
            </SheetTitle>
            <div className="h-14 flex items-center px-4 border-b font-semibold gap-2">
              <Shield className="h-5 w-5 text-violet-500" />
              {t("superadmin.title", "Superadmin")}
            </div>
            <SuperadminSidebarNav />
            <div className="p-4 border-t">
              <Link
                to="/admin"
                className="flex items-center gap-2 text-muted-foreground hover:text-foreground"
              >
                <Home className="h-4 w-4" />
                <span className="text-sm">
                  {t("superadmin.back_to_admin", "Back to Admin")}
                </span>
              </Link>
            </div>
          </SheetContent>
        </Sheet>
        <div className="font-semibold flex items-center gap-2">
          <Shield className="h-4 w-4 text-violet-500" />
          {t("superadmin.title_short", "Superadmin")}
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={logout}>
            {t("nav.logout", "Logout")}
          </Button>
        </div>
      </div>

      {/* Main content */}
      <section className="flex-1">
        {/* Desktop top bar */}
        <div className="hidden md:flex h-14 items-center justify-between border-b px-6 bg-background/80 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50">
          <div className="text-sm text-muted-foreground">
            {t("superadmin.topbar.platform_admin", "Platform Administrator")} •{" "}
            <span className="font-medium">{user?.email}</span>
          </div>
          <div className="flex items-center gap-3">
            <LanguageSwitcher />
            <span className="text-xs px-2 py-1 rounded bg-violet-500/10 text-violet-600 font-medium">
              SUPERADMIN
            </span>
            <UserMenu />
          </div>
        </div>
        <div className="p-4 md:p-6 pt-16 md:pt-4">
          <Outlet />
        </div>
      </section>
    </div>
  );
}

export default SuperadminLayout;
