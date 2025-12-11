import { Link, Outlet, useLocation } from "react-router-dom";
import React, { useState } from "react";
import { APP_NAME } from "../config";
import { useTranslation } from "react-i18next";
import { Menu, GraduationCap } from "lucide-react";
import { Button } from "@/components/ui/button";
import { NotificationCenter } from "@/components/NotificationCenter";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { GlobalSearch } from "@/components/GlobalSearch";
import { LanguageSwitcher } from "@/components/layout/LanguageSwitcher";
import { UserMenu } from "@/components/layout/UserMenu";
import { TenantSwitcher } from "@/components/TenantSwitcher";
import { useAuth } from "@/contexts/AuthContext";
import { useServiceEnabled } from "@/contexts/TenantServicesContext";
import { FloatingMenu } from "@/components/map/FloatingMenu";

export function AppLayout({ children }: { children?: React.ReactNode }) {
  const { t: T, i18n } = useTranslation("common");
  const { user, isLoading } = useAuth();
  const authed = !!user && !isLoading;
  const role = user?.role;
  const { pathname } = useLocation();
  
  // Check which optional services are enabled for the tenant
  const chatEnabled = useServiceEnabled('chat');
  const calendarEnabled = useServiceEnabled('calendar');
  
  const active = (p: string) =>
    pathname === p
      ? "text-primary font-medium bg-primary/10"
      : "text-muted-foreground hover:text-foreground hover:bg-muted/50";

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const onMobileLinkClick = () => setIsMobileMenuOpen(false);

  const NavLink = ({
    to,
    children,
    mobile,
  }: {
    to: string;
    children: React.ReactNode;
    mobile?: boolean;
  }) => (
    <Link
      to={to}
      className={`px-3 py-2 rounded-md transition-colors duration-200 ${active(
        to
      )} ${mobile ? "block w-full text-lg" : "text-sm"}`}
      onClick={mobile ? onMobileLinkClick : undefined}
    >
      {children}
    </Link>
  );

  const NavLinks = ({ mobile = false }: { mobile?: boolean }) => (
    <>
      {authed && (
          <NavLink to="/journey" mobile={mobile}>
            {T("nav.journey")}
          </NavLink>
        )}
      {authed && chatEnabled && (
        <NavLink to="/chat" mobile={mobile}>
          {T("nav.chat", { defaultValue: "Messages" })}
        </NavLink>
      )}
      {authed && calendarEnabled && (
        <NavLink to="/calendar" mobile={mobile}>
          {T("nav.calendar", { defaultValue: "Calendar" })}
        </NavLink>
      )}

      {authed && (role === "admin" || role === "superadmin") && (
        <NavLink to="/admin" mobile={mobile}>
          {T("nav.admin")}
        </NavLink>
      )}
      {!authed && pathname !== "/login" && (
        <NavLink to="/login" mobile={mobile}>
          {T("nav.login")}
        </NavLink>
      )}
    </>
  );

  return (
    <div className="min-h-screen flex flex-col">
      {(pathname !== "/" || authed) && (
        <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between gap-4">
          <Link to="/" className="font-bold text-xl shrink-0 flex items-center gap-2">
            <div className="bg-primary/10 p-1.5 rounded-lg">
              <GraduationCap className="h-5 w-5 text-primary" />
            </div>
            <span className="bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent">
              {APP_NAME}
            </span>
          </Link>

          {/* Global Search (Desktop) */}
          <div className="hidden md:block flex-1 max-w-md mx-4">
            {authed && <GlobalSearch />}
          </div>

          <nav className="flex items-center gap-3 text-sm shrink-0">
            {/* Desktop Navigation */}
            <div className="hidden md:flex items-center gap-1 mr-2">
              <NavLinks />
            </div>

            {/* Common Items (TenantSwitcher, Notification & Language) */}
            {authed && <TenantSwitcher className="hidden sm:flex" />}
            {authed && <NotificationCenter />}

            <LanguageSwitcher />

            {/* User Menu */}
            {authed && <UserMenu />}

            {/* Mobile Navigation Trigger */}
            <div className="md:hidden ml-2">
              <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
                <SheetTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <Menu className="h-5 w-5" />
                    <span className="sr-only">Toggle menu</span>
                  </Button>
                </SheetTrigger>
                <SheetContent side="right">
                  <SheetHeader className="text-left mb-6">
                    <SheetTitle>{APP_NAME}</SheetTitle>
                  </SheetHeader>
                  <div className="flex flex-col gap-2">
                    <NavLinks mobile />
                  </div>
                </SheetContent>
              </Sheet>
            </div>
          </nav>
        </div>
      </header>
      )}
      <main className="flex-1">{children ?? <Outlet />}</main>
      
      {authed && <FloatingMenu />}
    </div>
  );
}
