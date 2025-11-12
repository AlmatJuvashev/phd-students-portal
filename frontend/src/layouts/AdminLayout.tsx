import React from "react";
import { Outlet, NavLink, useLocation, Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Sheet, SheetTrigger, SheetContent, SheetTitle } from "@/components/ui/sheet";
import { cn } from "@/lib/utils";
import { Menu } from "lucide-react";

function SidebarNav() {
  const { user } = useAuth();
  const location = useLocation();
  const isActive = (to: string) => location.pathname === to;
  const canSeeAdmins = user?.role === "superadmin";
  const isAdmin = user?.role === "admin" || user?.role === "superadmin";

  return (
    <nav className="p-4 space-y-2">
      <div className="text-xs uppercase text-muted-foreground px-2">Core</div>
      <NavLink
        to="/admin"
        className={cn(
          "block rounded px-3 py-2 hover:bg-muted",
          isActive("/admin") && "bg-muted font-medium"
        )}
      >
        Dashboard
      </NavLink>
      <NavLink
        to="/admin/student-progress"
        className={cn(
          "block rounded px-3 py-2 hover:bg-muted",
          isActive("/admin/student-progress") && "bg-muted font-medium"
        )}
      >
        Student Progress
      </NavLink>
      <NavLink
        to="/admin/students-monitor"
        className={cn(
          "block rounded px-3 py-2 hover:bg-muted",
          isActive("/admin/students-monitor") && "bg-muted font-medium"
        )}
      >
        Students Monitor
      </NavLink>
      {isAdmin && (
        <>
          <div className="text-xs uppercase text-muted-foreground mt-4 px-2">Management</div>
          <NavLink
            to="/admin/create-students"
            className={cn(
              "block rounded px-3 py-2 hover:bg-muted",
              isActive("/admin/create-students") && "bg-muted font-medium"
            )}
          >
            Create Students
          </NavLink>
          <NavLink
            to="/admin/create-advisors"
            className={cn(
              "block rounded px-3 py-2 hover:bg-muted",
              isActive("/admin/create-advisors") && "bg-muted font-medium"
            )}
          >
            Create Advisors
          </NavLink>
        </>
      )}
      {canSeeAdmins && (
        <NavLink
          to="/admin/create-admins"
          className={cn(
            "block rounded px-3 py-2 hover:bg-muted",
            isActive("/admin/create-admins") && "bg-muted font-medium"
          )}
        >
          Create Admins
        </NavLink>
      )}
    </nav>
  );
}

export function AdminLayout() {
  const { user, logout } = useAuth();

  return (
    <div className="flex min-h-screen">
      {/* Desktop sidebar */}
      <aside className="hidden md:block w-64 border-r bg-background">
        <div className="h-14 flex items-center px-4 border-b font-semibold">
          Admin Panel
        </div>
        <SidebarNav />
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
            <SheetTitle className="sr-only">Menu</SheetTitle>
            <div className="h-14 flex items-center px-4 border-b font-semibold">
              Admin Panel
            </div>
            <SidebarNav />
          </SheetContent>
        </Sheet>
        <div className="font-semibold">Admin</div>
        <div className="flex items-center gap-2">
          <span className="text-xs text-muted-foreground">
            {user?.role?.toUpperCase()}
          </span>
          <Link to="/" aria-label="Back to main page">
            <Button variant="ghost" size="sm">Main</Button>
          </Link>
          <Button variant="ghost" size="sm" onClick={logout}>
            Logout
          </Button>
        </div>
      </div>

      {/* Main content */}
      <section className="flex-1">
        {/* Desktop top bar */}
        <div className="hidden md:flex h-14 items-center justify-between border-b px-6">
          <div className="text-sm text-muted-foreground">
            Signed in as <span className="font-medium">{user?.email}</span>
          </div>
          <div className="flex items-center gap-3">
            <span className="text-xs px-2 py-1 rounded bg-muted">
              {user?.role}
            </span>
            <Link to="/" aria-label="Back to main page">
              <Button variant="outline" size="sm">Back to main</Button>
            </Link>
            <Button variant="outline" size="sm" onClick={logout}>
              Logout
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
