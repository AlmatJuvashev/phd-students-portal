import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { api } from "../api/client";
import { Input } from "../components/ui/input";
import { Button } from "../components/ui/button";
import { Card } from "../components/ui/card";
import { Label } from "../components/ui/label";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/contexts/AuthContext";
import { useLocation, useNavigate } from "react-router-dom";

const Schema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
});
type Form = z.infer<typeof Schema>;

export function LoginPage() {
  const { t: T } = useTranslation("common");
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [errorMessage, setErrorMessage] = React.useState<string>("");
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<Form>({ resolver: zodResolver(Schema) });
  const onSubmit = async (data: Form) => {
    try {
      setErrorMessage("");
      // Use AuthContext for login to refresh user state
      await login({ username: data.username, password: data.password });
      const from = (location.state as any)?.from || "/journey";
      navigate(from, { replace: true });
    } catch (e: any) {
      console.error("Login failed", e);
      const message = e?.message || T("auth.failed");
      setErrorMessage(message);
    }
  };
  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <Card className="w-full max-w-sm p-6 shadow-md">
        <div className="mb-4">
          <h2 className="text-xl font-semibold">{T("auth.login_title")}</h2>
          <p className="text-xs text-muted-foreground mt-1">
            {T("app.name", { defaultValue: "PhD Portal" })}
          </p>
        </div>
        {errorMessage && (
          <div className="mb-4 p-3 rounded-md bg-rose-50 border border-rose-200 text-rose-700 text-sm">
            {errorMessage}
          </div>
        )}
        <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-1">
            <Label htmlFor="username">{T("auth.username")}</Label>
            <Input
              id="username"
              placeholder={T("auth.username")}
              aria-invalid={!!errors.username}
              aria-describedby={errors.username ? "username-error" : undefined}
              {...register("username")}
            />
            {errors.username && (
              <div id="username-error" className="text-xs text-rose-600">
                {errors.username.message}
              </div>
            )}
          </div>
          <div className="grid gap-1">
            <Label htmlFor="password">{T("auth.password")}</Label>
            <Input
              id="password"
              placeholder={T("auth.password")}
              type="password"
              aria-invalid={!!errors.password}
              aria-describedby={errors.password ? "password-error" : undefined}
              {...register("password")}
            />
            {errors.password && (
              <div id="password-error" className="text-xs text-rose-600">
                {errors.password.message}
              </div>
            )}
          </div>
          <Button
            className="w-full h-11"
            disabled={isSubmitting}
            aria-busy={isSubmitting}
          >
            {T("auth.signin")}
          </Button>
        </form>
        <div className="mt-4 text-sm text-center">
          <a
            href="/forgot-password"
            className="text-muted-foreground hover:underline"
          >
            {T("auth.forgot")}
          </a>
        </div>
      </Card>
    </div>
  );
}
