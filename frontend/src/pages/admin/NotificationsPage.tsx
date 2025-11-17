import React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { api } from "@/api/client";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Bell, BellOff, CheckCheck, Clock } from "lucide-react";
import { useTranslation } from "react-i18next";
import { cn } from "@/lib/utils";

type Notification = {
  id: string;
  student_id: string;
  student_name: string;
  student_email: string;
  node_id: string;
  node_instance_id: string;
  event_type: string;
  is_read: boolean;
  message: string;
  metadata: string;
  created_at: string;
};

function getEventIcon(eventType: string) {
  switch (eventType) {
    case "document_submitted":
      return "📄";
    case "document_uploaded":
      return "📎";
    case "form_submitted":
      return "📝";
    default:
      return "🔔";
  }
}

function getEventBadgeColor(eventType: string) {
  switch (eventType) {
    case "document_submitted":
      return "bg-blue-100 text-blue-800";
    case "document_uploaded":
      return "bg-green-100 text-green-800";
    case "form_submitted":
      return "bg-purple-100 text-purple-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
}

function getRelativeTime(
  dateString: string,
  t: (key: string, options?: any) => string,
  locale = "en"
): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffMins < 1)
    return t("admin.notifications.relative.just_now", {
      defaultValue: "Just now",
    });
  if (diffMins < 60)
    return t("admin.notifications.relative.minutes", {
      defaultValue: "{{count}} min ago",
      count: diffMins,
    });
  if (diffHours < 24)
    return t("admin.notifications.relative.hours", {
      defaultValue: "{{count}} h ago",
      count: diffHours,
    });
  if (diffDays === 1)
    return t("admin.notifications.relative.yesterday", {
      defaultValue: "Yesterday",
    });
  if (diffDays < 7)
    return t("admin.notifications.relative.days", {
      defaultValue: "{{count}} d ago",
      count: diffDays,
    });
  const shortDate = date.toLocaleDateString(locale, {
    day: "numeric",
    month: "short",
  });
  return t("admin.notifications.relative.short_date", {
    defaultValue: shortDate,
    date: shortDate,
  });
}

function getInitials(name: string): string {
  const parts = name.split(" ");
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return name.slice(0, 2).toUpperCase();
}

function getEventLabel(
  eventType: string,
  t: (key: string, options?: any) => string
) {
  switch (eventType) {
    case "document_submitted":
      return t("admin.notifications.event_labels.document_submitted", {
        defaultValue: "Document submitted",
      });
    case "document_uploaded":
      return t("admin.notifications.event_labels.document_uploaded", {
        defaultValue: "Document uploaded",
      });
    case "form_submitted":
      return t("admin.notifications.event_labels.form_submitted", {
        defaultValue: "Form submitted",
      });
    default:
      return t("admin.notifications.event_labels.default", {
        defaultValue: "Notification",
      });
  }
}

function parseMetadata(raw?: string | null) {
  if (!raw) return {};
  if (typeof raw !== "string") return raw;
  try {
    return JSON.parse(raw);
  } catch {
    return {};
  }
}

function getNotificationMessage(
  notification: Notification,
  nodeTitle: string,
  metadata: Record<string, any>,
  t: (key: string, options?: any) => string
) {
  const context = {
    node: nodeTitle,
    slot: metadata?.slot_key || "",
    file: metadata?.filename || "",
  };
  switch (notification.event_type) {
    case "document_submitted":
      return t("admin.notifications.messages.document_submitted", {
        defaultValue: "Documents submitted for {{node}}",
        ...context,
      });
    case "document_uploaded":
      return t("admin.notifications.messages.document_uploaded", {
        defaultValue: "New file uploaded in {{node}}",
        ...context,
      });
    case "form_submitted":
      return t("admin.notifications.messages.form_submitted", {
        defaultValue: "Form submitted in {{node}}",
        ...context,
      });
    default:
      return t("admin.notifications.messages.default", {
        defaultValue: "Update in {{node}}",
        ...context,
      });
  }
}

export function NotificationsPage() {
  const { t, i18n } = useTranslation("common");
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [showUnreadOnly, setShowUnreadOnly] = React.useState(false);
  const [nodeTitles, setNodeTitles] = React.useState<Record<string, string>>(
    {}
  );
  const locale = i18n.language || "en";

  React.useEffect(() => {
    let mounted = true;
    import("@/playbooks/playbook.json")
      .then((mod: any) => {
        if (!mounted) return;
        const pb = mod?.default ?? mod;
        const worlds = pb?.worlds || pb?.Worlds || [];
        const lang = (i18n.language || "en").toLowerCase();
        const pick = (obj: any, key: string) =>
          obj?.[key] ||
          obj?.[key?.toUpperCase?.()] ||
          (key
            ? obj?.[key.charAt(0).toUpperCase() + key.slice(1)]
            : undefined);
        const titles: Record<string, string> = {};
        worlds.forEach((w: any) => {
          const nodes = w?.nodes || w?.Nodes || [];
          nodes.forEach((node: any) => {
            const id = node?.id || node?.ID;
            if (!id) return;
            const titleObj = node?.title || node?.Title || {};
            titles[id] =
              pick(titleObj, lang) ||
              pick(titleObj, "en") ||
              pick(titleObj, "ru") ||
              pick(titleObj, "kz") ||
              id;
          });
        });
        setNodeTitles(titles);
      })
      .catch(() => setNodeTitles({}));
    return () => {
      mounted = false;
    };
  }, [i18n.language]);

  const { data: notificationsData, isLoading } = useQuery<Notification[]>({
    queryKey: ["admin", "notifications", showUnreadOnly],
    queryFn: () =>
      api(`/admin/notifications${showUnreadOnly ? "?unread_only=true" : ""}`),
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  // Ensure notifications is always an array
  const notifications = Array.isArray(notificationsData)
    ? notificationsData
    : [];

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) =>
      api(`/admin/notifications/${id}/read`, { method: "PATCH" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "notifications"] });
      queryClient.invalidateQueries({
        queryKey: ["admin", "notifications", "unread-count"],
      });
    },
  });

  const markAllAsReadMutation = useMutation({
    mutationFn: () => api("/admin/notifications/read-all", { method: "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "notifications"] });
      queryClient.invalidateQueries({
        queryKey: ["admin", "notifications", "unread-count"],
      });
    },
  });

  const handleNotificationClick = (notification: Notification) => {
    // Mark as read if unread
    if (!notification.is_read) {
      markAsReadMutation.mutate(notification.id);
    }

    // Navigate to student detail in StudentsMonitor
    navigate(
      `/admin/students-monitor?student=${notification.student_id}&node=${notification.node_id}`
    );
  };

  const unreadCount = notifications.filter((n) => !n.is_read).length;

  return (
    <div className="max-w-4xl mx-auto space-y-4">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Bell className="h-6 w-6 text-primary" />
              <div>
                <CardTitle>
                  {t("admin.notifications.title", {
                    defaultValue: "Уведомления",
                  })}
                  {unreadCount > 0 && (
                    <Badge className="ml-2 bg-red-500 text-white">
                      {unreadCount}
                    </Badge>
                  )}
                </CardTitle>
                <p className="text-sm text-muted-foreground mt-1">
                  {t("admin.notifications.subtitle", {
                    defaultValue: "Новые документы и обновления от студентов",
                  })}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowUnreadOnly(!showUnreadOnly)}
              >
                {showUnreadOnly ? (
                  <>
                    <Bell className="h-4 w-4 mr-2" />
                    {t("admin.notifications.show_all", {
                      defaultValue: "Показать все",
                    })}
                  </>
                ) : (
                  <>
                    <BellOff className="h-4 w-4 mr-2" />
                    {t("admin.notifications.show_unread", {
                      defaultValue: "Только непрочитанные",
                    })}
                  </>
                )}
              </Button>
              {unreadCount > 0 && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => markAllAsReadMutation.mutate()}
                  disabled={markAllAsReadMutation.isPending}
                >
                  <CheckCheck className="h-4 w-4 mr-2" />
                  {t("admin.notifications.mark_all_read", {
                    defaultValue: "Отметить все",
                  })}
                </Button>
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">
              {t("common.loading", { defaultValue: "Загрузка..." })}
            </div>
          ) : notifications.length === 0 ? (
            <div className="text-center py-12">
              <BellOff className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
              <p className="text-muted-foreground">
                {showUnreadOnly
                  ? t("admin.notifications.no_unread", {
                      defaultValue: "Нет непрочитанных уведомлений",
                    })
                  : t("admin.notifications.empty", {
                      defaultValue: "Уведомлений пока нет",
                    })}
              </p>
            </div>
          ) : (
            <div className="space-y-2">
              {notifications.map((notification) => {
                const metadata = parseMetadata(notification.metadata);
                const nodeTitle =
                  nodeTitles[notification.node_id] ||
                  metadata?.node_title ||
                  notification.node_id;
                const localizedMessage =
                  getNotificationMessage(notification, nodeTitle, metadata, t) ||
                  notification.message;
                return (
                  <div
                    key={notification.id}
                    onClick={() => handleNotificationClick(notification)}
                    className={cn(
                      "p-4 rounded-lg border transition-all cursor-pointer hover:shadow-md",
                      !notification.is_read
                        ? "bg-blue-50 border-blue-200 font-medium"
                        : "bg-white hover:bg-gray-50"
                    )}
                  >
                    <div className="flex items-start gap-3">
                    {/* Avatar */}
                    <div
                      className={cn(
                        "w-10 h-10 rounded-full flex items-center justify-center text-white font-semibold text-sm flex-shrink-0",
                        !notification.is_read
                          ? "bg-primary"
                          : "bg-muted-foreground"
                      )}
                    >
                      {getInitials(notification.student_name)}
                    </div>

                      {/* Content */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-2 mb-1">
                          <div className="flex items-center gap-2 flex-wrap">
                            <span
                              className={cn(
                                "font-medium",
                                !notification.is_read && "text-primary"
                              )}
                            >
                              {notification.student_name}
                            </span>
                            <Badge
                              className={getEventBadgeColor(
                                notification.event_type
                              )}
                            >
                              {getEventIcon(notification.event_type)}{" "}
                              {getEventLabel(notification.event_type, t)}
                            </Badge>
                          </div>
                          <div className="flex items-center gap-2 text-xs text-muted-foreground flex-shrink-0">
                            <Clock className="h-3 w-3" />
                            {getRelativeTime(notification.created_at, t, locale)}
                          </div>
                        </div>

                        <p className="text-sm text-foreground mb-1">
                          {localizedMessage}
                        </p>

                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                          <span
                            className="bg-muted px-2 py-0.5 rounded"
                            title={notification.node_id}
                          >
                            {nodeTitle}
                          </span>
                          <span>•</span>
                          <span>{notification.student_email}</span>
                        </div>
                      </div>

                      {/* Unread indicator */}
                      {!notification.is_read && (
                        <div className="w-2 h-2 bg-primary rounded-full flex-shrink-0 mt-2" />
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export default NotificationsPage;
