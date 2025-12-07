import React from "react";
import { Outlet, NavLink, useLocation, Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { useServiceEnabled } from "@/contexts/TenantServicesContext";
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
  LayoutDashboard,
  Users,
  UserPlus,
  UserCog,
  Monitor,
  Bell,
  MessageCircle,
  BookOpen,
  BarChart3,
  PhoneCall,
} from "lucide-react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/client";
import { Badge } from "@/components/ui/badge";
import { NotificationCenter } from "@/components/NotificationCenter";
import { LanguageSwitcher } from "@/components/layout/LanguageSwitcher";
import { UserMenu } from "@/components/layout/UserMenu";

function SidebarNav({ collapsed }: { collapsed?: boolean }) {
  const { t } = useTranslation("common");
  const { user } = useAuth();
  const location = useLocation();
  const isActive = (to: string) => location.pathname === to;
  const canSeeAdmins = user?.role === "superadmin";
  const isAdmin = user?.role === "admin" || user?.role === "superadmin";

  // Fetch unread notifications count
  const { data: unreadData } = useQuery<{ count: number }>({
    queryKey: ["admin", "notifications", "unread-count"],
    queryFn: () => api("/admin/notifications/unread-count"),
    refetchInterval: 30000, // Refetch every 30 seconds
  });
  const unreadCount = unreadData?.count || 0;
  
  // Check if optional services are enabled
  const chatEnabled = useServiceEnabled('chat');

  return (
    <nav className={cn("p-4 space-y-2", collapsed && "px-2")}>
      {!collapsed && (
        <div className="text-xs uppercase text-muted-foreground px-2">
          {t("admin.sidebar.core", "Core")}
        </div>
      )}
      <NavLink
        to="/admin"
        className={cn(
          "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
          isActive("/admin") && "bg-muted font-medium",
          collapsed && "justify-center px-2"
        )}
        title={t("admin.sidebar.dashboard", "Dashboard")}
      >
        <LayoutDashboard className="h-4 w-4" />
        {!collapsed && <span>{t("admin.sidebar.dashboard", "Dashboard")}</span>}
      </NavLink>
      <NavLink
        to="/admin/analytics"
        className={cn(
          "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
          isActive("/admin/analytics") && "bg-muted font-medium",
          collapsed && "justify-center px-2"
        )}
        title={t("admin.sidebar.analytics", "Analytics")}
      >
        <BarChart3 className="h-4 w-4" />
        {!collapsed && <span>{t("admin.sidebar.analytics", "Analytics")}</span>}
      </NavLink>
      <NavLink
        to="/admin/students-monitor"
        className={cn(
          "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
          isActive("/admin/students-monitor") && "bg-muted font-medium",
          collapsed && "justify-center px-2"
        )}
        title={t("admin.sidebar.students_monitor", "Students Monitor")}
      >
        <Monitor className="h-4 w-4" />
        {!collapsed && (
          <span>{t("admin.sidebar.students_monitor", "Students Monitor")}</span>
        )}
      </NavLink>
      <NavLink
        to="/admin/notifications"
        className={cn(
          "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted relative",
          isActive("/admin/notifications") && "bg-muted font-medium",
          collapsed && "justify-center px-2"
        )}
        title={t("admin.sidebar.notifications", "Notifications")}
      >
        <div className="relative">
          <Bell className="h-4 w-4" />
          {unreadCount > 0 && (
            <Badge className="absolute -top-2 -right-2 h-4 min-w-4 flex items-center justify-center p-0 text-[10px] bg-red-500 text-white">
              {unreadCount > 99 ? "99+" : unreadCount}
            </Badge>
          )}
        </div>
        {!collapsed && (
          <div className="flex items-center justify-between flex-1">
            <span>{t("admin.sidebar.notifications", "Notifications")}</span>
            {unreadCount > 0 && (
              <Badge className="bg-red-500 text-white text-[10px] px-1 py-0.5 rounded-full min-w-[1.5rem] flex items-center justify-center">
                {unreadCount}
              </Badge>
            )}
          </div>
        )}
      </NavLink>
      {isAdmin && (
        <>
          {!collapsed && (
            <div className="text-xs uppercase text-muted-foreground mt-4 px-2">
              {t("admin.sidebar.management", "Management")}
            </div>
          )}
          
          {/* User Management */}
          <NavLink
            to="/admin/users"
            className={cn(
              "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
              isActive("/admin/users") && "bg-muted font-medium",
              collapsed && "justify-center px-2"
            )}
            title={t("admin.sidebar.users", "Users")}
          >
            <Users className="h-4 w-4" />
            {!collapsed && (
              <span>
                {t("admin.sidebar.users", "Users")}
              </span>
            )}
          </NavLink>

          {/* Chat Rooms - Only if chat service is enabled */}
          {chatEnabled && (
            <NavLink
              to="/admin/chat-rooms"
              className={cn(
                "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
                isActive("/admin/chat-rooms") && "bg-muted font-medium",
                collapsed && "justify-center px-2"
              )}
              title={t("admin.sidebar.chat_rooms", "Chat rooms")}
            >
              <MessageCircle className="h-4 w-4" />
              {!collapsed && <span>{t("admin.sidebar.chat_rooms", "Chat rooms")}</span>}
            </NavLink>
          )}
          <NavLink
            to="/admin/dictionaries"
            className={cn(
              "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
              isActive("/admin/dictionaries") && "bg-muted font-medium",
              collapsed && "justify-center px-2"
            )}
            title={t("admin.sidebar.dictionaries", "Dictionaries")}
          >
            <BookOpen className="h-4 w-4" />
            {!collapsed && (
              <span>
                {t("admin.sidebar.dictionaries", "Dictionaries")}
              </span>
            )}
          </NavLink>
          <NavLink
            to="/admin/contacts"
            className={cn(
              "flex items-center gap-2 rounded px-3 py-2 hover:bg-muted",
              isActive("/admin/contacts") && "bg-muted font-medium",
              collapsed && "justify-center px-2"
            )}
            title={t("admin.sidebar.contacts", "Contacts")}
          >
            <PhoneCall className="h-4 w-4" />
            {!collapsed && (
              <span>{t("admin.sidebar.contacts", "Contacts")}</span>
            )}
          </NavLink>
        </>
      )}
    </nav>
  );
}

export function AdminLayout() {
  const { t, i18n } = useTranslation("common");
  const { user, logout } = useAuth();
  const [collapsed, setCollapsed] = React.useState(false);
  const languages = [
    { code: "ru", label: "RU" },
    { code: "kz", label: "KZ" },
    { code: "en", label: "EN" },
  ];

  return (
    <div className="flex min-h-screen">
      {/* Desktop sidebar */}
      <aside
        className={cn(
          "hidden md:block border-r sidebar-gradient transition-all duration-200",
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
              "font-bold truncate text-gradient text-lg",
              collapsed && "text-base"
            )}
            title={t("admin.topbar.admin_panel", "Admin Panel")}
          >
            {collapsed ? "A" : t("admin.topbar.admin_panel", "Admin Panel")}
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
            >
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
            <SheetTitle className="sr-only">
              {t("admin.topbar.menu", "Menu")}
            </SheetTitle>
            <div className="h-14 flex items-center px-4 border-b font-semibold">
              {t("admin.topbar.admin_panel", "Admin Panel")}
            </div>
            <SidebarNav />
          </SheetContent>
        </Sheet>
        <div className="font-semibold">{t("admin.topbar.admin", "Admin")}</div>
        <div className="flex items-center gap-2">
          <span className="text-xs text-muted-foreground">
            {user?.role?.toUpperCase()}
          </span>
          <Link
            to="/"
            aria-label={t("admin.topbar.back_to_main", "Back to main page")}
          >
            <Button variant="ghost" size="sm">
              {t("admin.topbar.main", "Main")}
            </Button>
          </Link>
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
            {t("admin.topbar.signed_in_as", "Signed in as")}{" "}
            <span className="font-medium">{user?.email}</span>
          </div>
          <div className="flex items-center gap-3">
            <LanguageSwitcher />
            <NotificationCenter />
            <span className="text-xs px-2 py-1 rounded bg-muted">
              {user?.role}
            </span>
            <Link
              to="/"
              aria-label={t("admin.topbar.back_to_main", "Back to main page")}
            >
              <Button variant="ghost" size="sm">
                {t("admin.topbar.main", "Main")}
              </Button>
            </Link>
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

export default AdminLayout;
