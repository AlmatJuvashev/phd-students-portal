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
import { Eye, EyeOff } from "lucide-react";

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
  const [showPassword, setShowPassword] = React.useState(false);
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<Form>({ resolver: zodResolver(Schema) });
  const onSubmit = async (data: Form) => {
    try {
      // Use AuthContext for login to refresh user state
      const res = await login({ username: data.username, password: data.password });
      if (res.is_superadmin) {
        navigate("/superadmin", { replace: true });
        return;
      }
      const from = (location.state as any)?.from || "/";
      navigate(from, { replace: true });
    } catch (e: any) {
      console.error("Login failed", e);
      let message = e?.message || T("auth.failed");
      try {
        // Try to parse JSON error message from backend
        const parsed = JSON.parse(message);
        if (parsed.error) message = parsed.error;
      } catch (err) {
        // Not JSON, use original string
      }
      setError("password", { type: "manual", message });
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
            <div className="relative">
              <Input
                id="password"
                placeholder={T("auth.password")}
                type={showPassword ? "text" : "password"}
                aria-invalid={!!errors.password}
                aria-describedby={
                  errors.password ? "password-error" : undefined
                }
                className="pr-10"
                {...register("password")}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-0 top-0 h-full px-3 text-muted-foreground hover:text-foreground transition-colors"
                aria-label={
                  showPassword
                    ? T("auth.hide_password")
                    : T("auth.show_password")
                }
              >
                {showPassword ? (
                  <EyeOff className="h-4 w-4" />
                ) : (
                  <Eye className="h-4 w-4" />
                )}
              </button>
            </div>
            {errors.password && (
              <div id="password-error" className="text-xs text-rose-600">
                {errors.password.message}
              </div>
            )}
          </div>
          <Button
            type="submit"
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
