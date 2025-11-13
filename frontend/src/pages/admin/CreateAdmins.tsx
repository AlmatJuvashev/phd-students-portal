import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { api } from "@/api/client";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { useTranslation } from "react-i18next";
import { Navigate } from "react-router-dom";

type Form = {
  first_name: string;
  last_name: string;
  email: string;
};

export function CreateAdmins() {
  const { user } = useAuth();
  const { t } = useTranslation('common');
  const Schema = React.useMemo(() => z.object({
    first_name: z.string().min(1, t('validation.required','Required')),
    last_name: z.string().min(1, t('validation.required','Required')),
    email: z.string().email(t('validation.invalid_email','Invalid email')),
  }), [t]);
  const [created, setCreated] = React.useState<{ username: string; temp_password: string } | null>(null);
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<Form>({ resolver: zodResolver(Schema) });

  if (user && user.role !== "superadmin") {
    return <Navigate to="/admin" replace />;
  }

  async function onSubmit(data: Form) {
    try {
      const res = await api("/admin/users", {
        method: "POST",
        body: JSON.stringify({ ...data, role: "admin" }),
      });
      setCreated(res);
      reset();
    } catch (e: any) {
      alert(e.message || "Failed to create admin");
    }
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div>
        <h2 className="text-2xl font-bold">{t('admin.forms.create_admins.title','Create Admins')}</h2>
        <p className="text-muted-foreground">{t('admin.forms.create_admins.subtitle','Only superadmins can create other admins.')}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t('admin.forms.create_admins.heading','New Admin')}</CardTitle>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
            <div>
              <Input placeholder={t('admin.forms.first_name','First name')} {...register("first_name")} />
              {errors.first_name && (
                <div className="text-xs text-red-600 mt-1">{errors.first_name.message}</div>
              )}
            </div>
            <div>
              <Input placeholder={t('admin.forms.last_name','Last name')} {...register("last_name")} />
              {errors.last_name && (
                <div className="text-xs text-red-600 mt-1">{errors.last_name.message}</div>
              )}
            </div>
            <div>
              <Input type="email" placeholder={t('admin.forms.email','Email')} {...register("email")} />
              {errors.email && (
                <div className="text-xs text-red-600 mt-1">{errors.email.message}</div>
              )}
            </div>
            <div className="pt-2">
              <Button type="submit" disabled={isSubmitting}>{isSubmitting ? t('admin.forms.creating','Creating…') : t('admin.forms.create_admins.submit','Create Admin')}</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader>
            <CardTitle className="text-green-800">{t('admin.forms.create_admins.success','Admin Created')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="text-sm">{t('admin.forms.username','Username')}: <span className="font-mono">{created.username}</span></div>
              <div className="text-sm">{t('admin.forms.temp_password','Temp password')}: <span className="font-mono">{created.temp_password}</span></div>
              <div className="text-xs text-muted-foreground">{t('admin.forms.share_credentials','Share these credentials securely with the new user.')}</div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

export default CreateAdmins;
