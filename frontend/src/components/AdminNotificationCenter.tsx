import React from "react";
import { Bell } from "lucide-react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/api/client";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu-radix";
import { Badge } from "@/components/ui/badge";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { useTranslation } from "react-i18next";

interface Notification {
  id: string;
  title: string;
  message: string;
  link?: string;
  is_read: boolean;
  created_at: string;
}

export const AdminNotificationCenter = () => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { user, isLoading: authLoading } = useAuth();
  const { t } = useTranslation("common");

  // Check if token exists in localStorage
  const token = localStorage.getItem("token");

  const { data: notificationsData } = useQuery<Notification[] | null>({
    queryKey: ["admin", "notifications"],
    queryFn: () => api("/admin/notifications"),
    refetchInterval: 30000,
    enabled: !!token && !!user && !authLoading, // Only fetch if token exists, user is loaded, and auth is not loading
  });

  // Ensure notifications is always an array
  const notifications = notificationsData ?? [];

  const markAsRead = useMutation({
    mutationFn: (id: string) =>
      api(`/admin/notifications/${id}/read`, { method: "PATCH" }),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["admin", "notifications"] }),
  });

  const markAllAsRead = useMutation({
    mutationFn: () => api("/admin/notifications/read-all", { method: "POST" }),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["admin", "notifications"] }),
  });

  const handleClick = (notif: Notification) => {
    if (!notif.is_read) {
      markAsRead.mutate(notif.id);
    }
    if (notif.link) {
      navigate(notif.link);
    }
  };

  const unreadCount = notifications.filter((n) => !n.is_read).length;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="relative">
          <Bell className="h-5 w-5" />
          {unreadCount > 0 && (
            <Badge className="absolute -top-1 -right-1 h-4 w-4 p-0 flex items-center justify-center bg-red-500 text-[10px]">
              {unreadCount}
            </Badge>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-80">
        <DropdownMenuLabel className="flex justify-between items-center">
          <span>{t("notifications.title")}</span>
          {unreadCount > 0 && (
            <Button
              variant="ghost"
              size="sm"
              className="text-xs h-auto py-1"
              onClick={() => markAllAsRead.mutate()}
            >
              {t("notifications.mark_all_read")}
            </Button>
          )}
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        {notifications.length === 0 ? (
          <div className="p-4 text-center text-sm text-muted-foreground">
            {t("notifications.no_notifications")}
          </div>
        ) : (
          <div className="max-h-[300px] overflow-y-auto">
            {notifications.map((notif) => (
              <DropdownMenuItem
                key={notif.id}
                className="flex flex-col items-start p-3 cursor-pointer"
                onClick={() => handleClick(notif)}
              >
                <div className="flex justify-between w-full">
                  <span
                    className={`font-medium text-sm ${
                      !notif.is_read
                        ? "text-foreground"
                        : "text-muted-foreground"
                    }`}
                  >
                    {notif.title}
                  </span>
                  {!notif.is_read && (
                    <span className="h-2 w-2 rounded-full bg-blue-500 mt-1" />
                  )}
                </div>
                <span className="text-xs text-muted-foreground mt-1 line-clamp-2">
                  {notif.message}
                </span>
                <span className="text-[10px] text-muted-foreground mt-2">
                  {new Date(notif.created_at).toLocaleDateString()}
                </span>
              </DropdownMenuItem>
            ))}
          </div>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
