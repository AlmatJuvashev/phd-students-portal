import React from "react";
import { useTranslation } from "react-i18next";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CreateStudents } from "./admin/CreateStudents";
import { CreateAdvisors } from "./admin/CreateAdvisors";
import { CreateAdmins } from "./admin/CreateAdmins";
import { useAuth } from "@/contexts/AuthContext";

export function AdminUsersPage() {
  const { t } = useTranslation("common");
  const { user } = useAuth();
  const isSuperAdmin = user?.role === "superadmin";

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">
          {t("admin.users.title", { defaultValue: "User Management" })}
        </h1>
      </div>

      <Tabs defaultValue="students" className="space-y-4">
        <TabsList>
          <TabsTrigger value="students">
            {t("admin.users.students", { defaultValue: "Students" })}
          </TabsTrigger>
          <TabsTrigger value="advisors">
            {t("admin.users.advisors", { defaultValue: "Advisors" })}
          </TabsTrigger>
          {isSuperAdmin && (
            <TabsTrigger value="admins">
              {t("admin.users.admins", { defaultValue: "Admins" })}
            </TabsTrigger>
          )}
        </TabsList>
        <TabsContent value="students" className="space-y-4">
          <CreateStudents />
        </TabsContent>
        <TabsContent value="advisors" className="space-y-4">
          <CreateAdvisors />
        </TabsContent>
        {isSuperAdmin && (
          <TabsContent value="admins" className="space-y-4">
            <CreateAdmins />
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}
