import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { api } from "@/api/client";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

type Form = {
  first_name: string;
  last_name: string;
  email: string;
  role: "student" | "advisor" | "chair";
};

export function CreateUsers() {
  const { t } = useTranslation('common');
  const Schema = React.useMemo(() => z.object({
    first_name: z.string().min(1, t('validation.required','Required')),
    last_name: z.string().min(1, t('validation.required','Required')),
    email: z.string().email(t('validation.invalid_email','Invalid email')),
    role: z.enum(["student", "advisor", "chair"]),
  }), [t]);
  const [created, setCreated] = React.useState<{ username: string; temp_password: string } | null>(null);
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<Form>({ resolver: zodResolver(Schema), defaultValues: { role: "student" } });

  async function onSubmit(data: Form) {
    try {
      const res = await api("/admin/users", {
        method: "POST",
        body: JSON.stringify(data),
      });
      setCreated(res);
      reset({ first_name: "", last_name: "", email: "", role: data.role });
    } catch (e: any) {
      alert(e.message || "Failed to create user");
    }
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div>
        <h2 className="text-2xl font-bold">{t('admin.forms.create_users.title','Create Users')}</h2>
        <p className="text-muted-foreground">{t('admin.forms.create_users.subtitle','Create students, advisors and department chairs.')}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t('admin.forms.create_users.heading','New User')}</CardTitle>
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
            <div>
              <select className="w-full border border-gray-300 p-2 rounded-md" {...register("role")}>
                <option value="student">{t('admin.forms.roles.student','Student')}</option>
                <option value="advisor">{t('admin.forms.roles.advisor','Advisor')}</option>
                <option value="chair">{t('admin.forms.roles.chair','Department Chair')}</option>
              </select>
            </div>
            <div className="pt-2">
              <Button type="submit" disabled={isSubmitting}>{isSubmitting ? t('admin.forms.creating','Creating…') : t('admin.forms.create_users.submit','Create User')}</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader>
            <CardTitle className="text-green-800">{t('admin.forms.create_users.success','User Created')}</CardTitle>
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

export default CreateUsers;
