import React from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { User, Settings, LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownItem } from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useAuth } from "@/contexts/AuthContext";

export function UserMenu() {
  const { t } = useTranslation("common");
  const { user: me } = useAuth();

  if (!me) return null;

  return (
    <DropdownMenu
      trigger={
        <Button
          variant="ghost"
          className="relative h-9 w-9 rounded-full"
          data-testid="user-menu-button"
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
            <span>{t("nav.profile", { defaultValue: "Profile" })}</span>
          </div>
        </DropdownItem>
      </Link>
      <Link to="/profile">
        <DropdownItem>
          <div className="flex items-center gap-2">
            <Settings className="h-4 w-4" />
            <span>{t("nav.settings")}</span>
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
          <span>{t("nav.logout")}</span>
        </div>
      </DropdownItem>
    </DropdownMenu>
  );
}
