import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { Map, Clock, CheckCircle2, TrendingUp } from "lucide-react";

export function Dashboard() {
  return (
    <div className="max-w-6xl mx-auto px-4 py-6 sm:py-8 space-y-6">
      {/* Welcome Header */}
      <div className="bg-gradient-to-r from-primary/10 via-primary/5 to-transparent rounded-2xl p-6 sm:p-8 border-l-4 border-primary">
        <h1 className="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent mb-2">
          Welcome to Your Dashboard
        </h1>
        <p className="text-sm sm:text-base text-muted-foreground">
          Track your doctoral journey progress and manage your research
          milestones
        </p>
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
              In Progress
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
              Completed
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
              Progress
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
              Total Tasks
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card className="bg-gradient-to-br from-card to-card/50">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Map className="w-5 h-5 text-primary" />
            Quick Actions
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <p className="text-sm text-muted-foreground mb-4">
            Navigate to your journey map to track progress, upload documents,
            and complete required tasks.
          </p>
          <Link to="/journey">
            <Button size="lg" className="w-full sm:w-auto">
              <Map className="w-4 h-4 mr-2" />
              Open Journey Map
            </Button>
          </Link>
        </CardContent>
      </Card>

      {/* Info Card */}
      <Card className="border-l-4 border-primary">
        <CardContent className="p-6">
          <h3 className="font-semibold mb-2">📝 Coming Soon</h3>
          <p className="text-sm text-muted-foreground">
            This dashboard will soon display your journey progress, recent
            uploads, advisor feedback, and upcoming deadlines.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
