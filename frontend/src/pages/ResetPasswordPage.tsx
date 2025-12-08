import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { api } from "../api/client";
import { Input } from "../components/ui/input";
import { Button } from "../components/ui/button";
import { Card } from "../components/ui/card";
import { Label } from "../components/ui/label";
import { Link, useSearchParams, useNavigate } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";

const Schema = z.object({
  password: z.string().min(6, "Password must be at least 6 characters"),
  confirm: z.string()
}).refine((data) => data.password === data.confirm, {
  message: "Passwords don't match",
  path: ["confirm"],
});

type Form = z.infer<typeof Schema>;

export function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token");
  const navigate = useNavigate();
  
  const [showPassword, setShowPassword] = React.useState(false);
  const [success, setSuccess] = React.useState(false);

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<Form>({ resolver: zodResolver(Schema) });

  React.useEffect(() => {
    if (!token) {
        // Redirect or show error if no token
    }
  }, [token]);

  const onSubmit = async (data: Form) => {
    if (!token) return;
    try {
      await api.post("/auth/reset-password", {
        token,
        new_password: data.password
      });
      setSuccess(true);
      setTimeout(() => {
          navigate("/login");
      }, 3000);
    } catch (e: any) {
      console.error("Reset password failed", e);
      let msg = "Failed to reset password. Token may be invalid or expired.";
      try {
          const parsed = JSON.parse(e.message);
          if (parsed.error) msg = parsed.error;
      } catch {}
      setError("root", { message: msg });
    }
  };

  if (!token) {
      return (
        <div className="min-h-screen flex items-center justify-center p-4">
            <Card className="p-6 text-center text-rose-600">
                Invalid reset link.
            </Card>
        </div>
      )
  }

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <Card className="w-full max-w-sm p-6 shadow-md text-center">
          <h2 className="text-xl font-semibold mb-4 text-green-600">Password Reset!</h2>
          <p className="text-muted-foreground mb-6">
            Your password has been successfully updated. Redirecting to login...
          </p>
          <Button asChild className="w-full">
            <Link to="/login">Go to Login</Link>
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <Card className="w-full max-w-sm p-6 shadow-md">
        <div className="mb-4">
          <h2 className="text-xl font-semibold">Set New Password</h2>
        </div>

        {errors.root && (
          <div className="mb-4 p-3 rounded-md bg-rose-50 border border-rose-200 text-rose-700 text-sm">
            {errors.root.message}
          </div>
        )}

        <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-1">
            <Label htmlFor="password">New Password</Label>
            <div className="relative">
                <Input
                id="password"
                type={showPassword ? "text" : "password"}
                {...register("password")}
                />
                <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-0 top-0 h-full px-3 text-muted-foreground hover:text-foreground transition-colors"
                >
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
            </div>
            {errors.password && (
              <div className="text-xs text-rose-600">
                {errors.password.message}
              </div>
            )}
          </div>

          <div className="grid gap-1">
            <Label htmlFor="confirm">Confirm Password</Label>
            <Input
              id="confirm"
              type="password"
              {...register("confirm")}
            />
            {errors.confirm && (
              <div className="text-xs text-rose-600">
                {errors.confirm.message}
              </div>
            )}
          </div>
          
          <Button
            className="w-full h-11"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Resetting..." : "Reset Password"}
          </Button>
        </form>
      </Card>
    </div>
  );
}
