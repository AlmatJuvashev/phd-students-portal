import { Link, Outlet, useLocation } from "react-router-dom";
import React, { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { APP_NAME } from "../config";
import { useTranslation } from "react-i18next";
import { Globe, Menu, LogOut, User, Settings } from "lucide-react";
import { DropdownMenu, DropdownItem } from "@/components/ui/dropdown-menu";
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

      {authed &&
        (!["admin", "superadmin"].includes(role || "") ||
          !import.meta.env.PROD) && (
          <NavLink to="/journey" mobile={mobile}>
            {T("nav.journey")}
          </NavLink>
        )}
      {authed && (
        <NavLink to="/chat" mobile={mobile}>
          {T("nav.chat", { defaultValue: "Messages" })}
        </NavLink>
      )}
      {authed && (
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
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between gap-4">
          <Link to="/" className="font-semibold text-lg shrink-0">
            {APP_NAME}
          </Link>

          {/* Global Search (Desktop) */}
          <div className="hidden md:block flex-1 max-w-md mx-4">
            {authed && <GlobalSearch />}
          </div>

          <nav className="flex items-center gap-2 text-sm shrink-0">
            {/* Desktop Navigation */}
            <div className="hidden md:flex items-center gap-1 mr-2">
              <NavLinks />
            </div>

            {/* Common Items (Notification & Language) */}
            {authed && <NotificationCenter />}
            <DropdownMenu
              trigger={
                <Button
                  variant="ghost"
                  size="sm"
                  className="gap-2 h-9 px-2 text-muted-foreground hover:text-foreground"
                >
                  <Globe className="h-4 w-4" />
                  <span className="hidden sm:inline">
                    {i18n.language.toUpperCase()}
                  </span>
                </Button>
              }
            >
              <DropdownItem onClick={() => i18n.changeLanguage("ru")}>
                <div
                  className={`flex items-center gap-3 ${
                    i18n.language === "ru" ? "font-semibold text-primary" : ""
                  }`}
                >
                  <span className="text-lg">🇷🇺</span>
                  <span>Русский</span>
                </div>
              </DropdownItem>
              <DropdownItem onClick={() => i18n.changeLanguage("kz")}>
                <div
                  className={`flex items-center gap-3 ${
                    i18n.language === "kz" ? "font-semibold text-primary" : ""
                  }`}
                >
                  <span className="text-lg">🇰🇿</span>
                  <span>Қазақша</span>
                </div>
              </DropdownItem>
              <DropdownItem onClick={() => i18n.changeLanguage("en")}>
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

            {/* User Menu */}
            {authed && (
              <DropdownMenu
                trigger={
                  <Button
                    variant="ghost"
                    className="relative h-9 w-9 rounded-full ml-2"
                  >
                    <Avatar className="h-9 w-9 border border-border">
                      <AvatarImage src={me?.avatar_url} />
                      <AvatarFallback>
                        {me?.first_name?.[0]}
                        {me?.last_name?.[0]}
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                }
              >
                <div className="px-3 py-2 border-b border-border/50 mb-1">
                  <p className="text-sm font-medium leading-none">
                    {me?.first_name} {me?.last_name}
                  </p>
                  <p className="text-xs text-muted-foreground mt-1 truncate max-w-[180px]">
                    {me?.email}
                  </p>
                </div>
                <Link to="/profile">
                  <DropdownItem>
                    <div className="flex items-center gap-2">
                      <User className="h-4 w-4" />
                      <span>{T("nav.profile", { defaultValue: "Profile" })}</span>
                    </div>
                  </DropdownItem>
                </Link>
                <Link to="/profile">
                  <DropdownItem>
                    <div className="flex items-center gap-2">
                      <Settings className="h-4 w-4" />
                      <span>Settings</span>
                    </div>
                  </DropdownItem>
                </Link>
                <DropdownItem
                  onClick={() => {
                    localStorage.removeItem("token");
                    location.href = "/login";
                  }}
                >
                  <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
                    <LogOut className="h-4 w-4" />
                    <span>{T("nav.logout")}</span>
                  </div>
                </DropdownItem>
              </DropdownMenu>
            )}

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
      <main className="flex-1">{children ?? <Outlet />}</main>
    </div>
  );
}
