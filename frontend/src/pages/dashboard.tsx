import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import {
  Map,
  Clock,
  CheckCircle2,
  TrendingUp,
  MessageCircle,
  Calendar,
  User,
  Settings,
} from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useTranslation } from "react-i18next";

export function Dashboard() {
  const { user } = useAuth();
  const { t } = useTranslation("common");

  return (
    <div className="max-w-6xl mx-auto px-4 py-6 sm:py-8 space-y-8">
      {/* Welcome Header & Profile Widget */}
      <div className="grid gap-6 md:grid-cols-[1fr_300px]">
        <div className="bg-gradient-to-r from-primary/10 via-primary/5 to-transparent rounded-2xl p-6 sm:p-8 border-l-4 border-primary flex flex-col justify-center">
          <h1 className="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent mb-2">
            {t("dashboard.welcome_back", { name: user?.first_name })}
          </h1>
          <p className="text-sm sm:text-base text-muted-foreground">
            {t("dashboard.subtitle")}
          </p>
        </div>

        {/* Profile Widget */}
        <Card>
          <CardContent className="p-6 flex items-center gap-4">
            <Avatar className="h-16 w-16 border-2 border-primary/10">
              <AvatarImage src={user?.avatar_url} />
              <AvatarFallback className="text-lg">
                {user?.first_name?.[0]}
                {user?.last_name?.[0]}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="font-semibold truncate">
                {user?.first_name} {user?.last_name}
              </p>
              <p className="text-xs text-muted-foreground truncate">
                {user?.role}
              </p>
              <Link to="/profile" className="mt-2 inline-block">
                <Button variant="outline" size="sm" className="h-7 text-xs">
                  <Settings className="w-3 h-3 mr-1.5" />
                  {t("dashboard.settings")}
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="bg-gradient-to-br from-blue-50 to-blue-100/50 dark:from-blue-900/20 dark:to-blue-800/10 border-blue-200 dark:border-blue-800">
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-2">
              <div className="p-2 bg-blue-500 rounded-lg">
                <Clock className="w-5 h-5 text-white" />
              </div>
              <span className="text-2xl font-bold text-blue-700 dark:text-blue-300">
                -
              </span>
            </div>
            <div className="text-sm font-medium text-blue-700 dark:text-blue-300">
              {t("dashboard.in_progress")}
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-green-50 to-green-100/50 dark:from-green-900/20 dark:to-green-800/10 border-green-200 dark:border-green-800">
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-2">
              <div className="p-2 bg-green-500 rounded-lg">
                <CheckCircle2 className="w-5 h-5 text-white" />
              </div>
              <span className="text-2xl font-bold text-green-700 dark:text-green-300">
                -
              </span>
            </div>
            <div className="text-sm font-medium text-green-700 dark:text-green-300">
              {t("dashboard.completed")}
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-purple-50 to-purple-100/50 dark:from-purple-900/20 dark:to-purple-800/10 border-purple-200 dark:border-purple-800">
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-2">
              <div className="p-2 bg-purple-500 rounded-lg">
                <TrendingUp className="w-5 h-5 text-white" />
              </div>
              <span className="text-2xl font-bold text-purple-700 dark:text-purple-300">
                -%
              </span>
            </div>
            <div className="text-sm font-medium text-purple-700 dark:text-purple-300">
              {t("dashboard.progress")}
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-amber-50 to-amber-100/50 dark:from-amber-900/20 dark:to-amber-800/10 border-amber-200 dark:border-amber-800">
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-2">
              <div className="p-2 bg-amber-500 rounded-lg">
                <Map className="w-5 h-5 text-white" />
              </div>
              <span className="text-2xl font-bold text-amber-700 dark:text-amber-300">
                -
              </span>
            </div>
            <div className="text-sm font-medium text-amber-700 dark:text-amber-300">
              {t("dashboard.total_tasks")}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions Grid */}
      <div>
        <h2 className="text-lg font-semibold mb-4">{t("dashboard.quick_actions")}</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <Link to="/journey">
            <Card className="hover:shadow-md transition-shadow cursor-pointer h-full">
              <CardContent className="p-6 flex flex-col items-center text-center gap-3">
                <div className="p-3 bg-primary/10 rounded-full">
                  <Map className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <h3 className="font-medium">{t("dashboard.journey_map")}</h3>
                  <p className="text-xs text-muted-foreground mt-1">
                    {t("dashboard.track_progress")}
                  </p>
                </div>
              </CardContent>
            </Card>
          </Link>

          <Link to="/chat">
            <Card className="hover:shadow-md transition-shadow cursor-pointer h-full">
              <CardContent className="p-6 flex flex-col items-center text-center gap-3">
                <div className="p-3 bg-primary/10 rounded-full">
                  <MessageCircle className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <h3 className="font-medium">{t("dashboard.messages")}</h3>
                  <p className="text-xs text-muted-foreground mt-1">
                    {t("dashboard.chat_advisor")}
                  </p>
                </div>
              </CardContent>
            </Card>
          </Link>

          <Link to="/calendar">
            <Card className="hover:shadow-md transition-shadow cursor-pointer h-full">
              <CardContent className="p-6 flex flex-col items-center text-center gap-3">
                <div className="p-3 bg-primary/10 rounded-full">
                  <Calendar className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <h3 className="font-medium">{t("dashboard.calendar")}</h3>
                  <p className="text-xs text-muted-foreground mt-1">
                    {t("dashboard.view_events")}
                  </p>
                </div>
              </CardContent>
            </Card>
          </Link>

          <Link to="/profile">
            <Card className="hover:shadow-md transition-shadow cursor-pointer h-full">
              <CardContent className="p-6 flex flex-col items-center text-center gap-3">
                <div className="p-3 bg-primary/10 rounded-full">
                  <User className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <h3 className="font-medium">{t("dashboard.my_profile")}</h3>
                  <p className="text-xs text-muted-foreground mt-1">
                    {t("dashboard.manage_account")}
                  </p>
                </div>
              </CardContent>
            </Card>
          </Link>
        </div>
      </div>

      {/* Info Card */}
      <Card className="border-l-4 border-primary bg-muted/20">
        <CardContent className="p-6">
          <h3 className="font-semibold mb-2 flex items-center gap-2">
            <span className="text-xl">📝</span> {t("dashboard.coming_soon")}
          </h3>
          <p className="text-sm text-muted-foreground">
            {t("dashboard.coming_soon_text")}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
