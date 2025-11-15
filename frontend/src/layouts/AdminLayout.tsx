import React from "react";
import { Outlet, NavLink, useLocation, Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Sheet, SheetTrigger, SheetContent, SheetTitle } from "@/components/ui/sheet";
import { cn } from "@/lib/utils";
import { Menu, LayoutDashboard, Users, UserPlus, UserCog, Monitor } from "lucide-react";
import { useTranslation } from "react-i18next";

function SidebarNav({ collapsed }: { collapsed?: boolean }) {
  const { t } = useTranslation('common');
  const { user } = useAuth();
  const location = useLocation();
  const isActive = (to: string) => location.pathname === to;
  const canSeeAdmins = user?.role === "superadmin";
  const isAdmin = user?.role === "admin" || user?.role === "superadmin";

  return (
    <nav className={cn("p-4 space-y-2", collapsed && "px-2") }>
      {!collapsed && (
        <div className="text-xs uppercase text-muted-foreground px-2">{t('admin.sidebar.core','Core')}</div>
      )}
      <NavLink
        to="/admin"
        className={cn(
          "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
          isActive("/admin") && "bg-muted font-medium",
          collapsed && "justify-center px-2"
        )}
        title={t('admin.sidebar.dashboard','Dashboard')}
      >
        <LayoutDashboard className="h-4 w-4" />
        {!collapsed && <span>{t('admin.sidebar.dashboard','Dashboard')}</span>}
      </NavLink>
      <NavLink
        to="/admin/students-monitor"
        className={cn(
          "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
          isActive("/admin/students-monitor") && "bg-muted font-medium",
          collapsed && "justify-center px-2"
        )}
        title={t('admin.sidebar.students_monitor','Students Monitor')}
      >
        <Monitor className="h-4 w-4" />
        {!collapsed && <span>{t('admin.sidebar.students_monitor','Students Monitor')}</span>}
      </NavLink>
      {isAdmin && (
        <>
          {!collapsed && (
            <div className="text-xs uppercase text-muted-foreground mt-4 px-2">{t('admin.sidebar.management','Management')}</div>
          )}
          <NavLink
            to="/admin/create-students"
            className={cn(
              "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
              isActive("/admin/create-students") && "bg-muted font-medium",
              collapsed && "justify-center px-2"
            )}
            title={t('admin.sidebar.create_students','Create Students')}
          >
            <Users className="h-4 w-4" />
            {!collapsed && <span>{t('admin.sidebar.create_students','Create Students')}</span>}
          </NavLink>
          <NavLink
            to="/admin/create-advisors"
            className={cn(
              "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
              isActive("/admin/create-advisors") && "bg-muted font-medium",
              collapsed && "justify-center px-2"
            )}
            title={t('admin.sidebar.create_advisors','Create Advisors')}
          >
            <UserPlus className="h-4 w-4" />
            {!collapsed && <span>{t('admin.sidebar.create_advisors','Create Advisors')}</span>}
          </NavLink>
        </>
      )}
      {canSeeAdmins && (
        <NavLink
          to="/admin/create-admins"
          className={cn(
            "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
            isActive("/admin/create-admins") && "bg-muted font-medium",
            collapsed && "justify-center px-2"
          )}
          title={t('admin.sidebar.create_admins','Create Admins')}
        >
          <UserCog className="h-4 w-4" />
          {!collapsed && <span>{t('admin.sidebar.create_admins','Create Admins')}</span>}
        </NavLink>
      )}
    </nav>
  );
}

export function AdminLayout() {
  const { t } = useTranslation('common');
  const { user, logout } = useAuth();
  const [collapsed, setCollapsed] = React.useState(false);

  return (
    <div className="flex min-h-screen">
      {/* Desktop sidebar */}
      <aside className={cn("hidden md:block border-r bg-background transition-all duration-200", collapsed ? "w-16" : "w-64") }>
        <div className={cn("h-14 flex items-center border-b", collapsed ? "justify-center" : "px-4 justify-between") }>
          <div className="font-semibold truncate" title={t('admin.topbar.admin_panel','Admin Panel')}>
            {collapsed ? "A" : t('admin.topbar.admin_panel','Admin Panel')}
          </div>
          {!collapsed && (
            <Button variant="ghost" size="icon" onClick={() => setCollapsed(true)} aria-label={t('common.collapse','Collapse')}>
              <span className="text-lg leading-none">«</span>
            </Button>
          )}
          {collapsed && (
            <Button variant="ghost" size="icon" onClick={() => setCollapsed(false)} aria-label={t('common.expand','Expand')}>
              <span className="text-lg leading-none">»</span>
            </Button>
          )}
        </div>
        <SidebarNav collapsed={collapsed} />
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
            <SheetTitle className="sr-only">{t('admin.topbar.menu','Menu')}</SheetTitle>
            <div className="h-14 flex items-center px-4 border-b font-semibold">
              {t('admin.topbar.admin_panel','Admin Panel')}
            </div>
            <SidebarNav />
          </SheetContent>
        </Sheet>
        <div className="font-semibold">{t('admin.topbar.admin','Admin')}</div>
        <div className="flex items-center gap-2">
          <span className="text-xs text-muted-foreground">
            {user?.role?.toUpperCase()}
          </span>
          <Link to="/" aria-label={t('admin.topbar.back_to_main','Back to main page')}>
            <Button variant="ghost" size="sm">{t('admin.topbar.main','Main')}</Button>
          </Link>
          <Button variant="ghost" size="sm" onClick={logout}>
            {t('nav.logout','Logout')}
          </Button>
        </div>
      </div>

      {/* Main content */}
      <section className="flex-1">
        {/* Desktop top bar */}
        <div className="hidden md:flex h-14 items-center justify-between border-b px-6">
          <div className="text-sm text-muted-foreground">{t('admin.topbar.signed_in_as','Signed in as')} <span className="font-medium">{user?.email}</span></div>
          <div className="flex items-center gap-3">
            <span className="text-xs px-2 py-1 rounded bg-muted">
              {user?.role}
            </span>
            <Link to="/" aria-label={t('admin.topbar.back_to_main','Back to main page')}>
              <Button variant="outline" size="sm">{t('admin.topbar.back_to_main_button','Back to main')}</Button>
            </Link>
            <Button variant="outline" size="sm" onClick={logout}>
              {t('nav.logout','Logout')}
            </Button>
          </div>
        </div>
        <div className="p-4 md:p-6 pt-16 md:pt-4">
          <Outlet />
        </div>
      </section>
    </div>
  );
}

export default AdminLayout;
