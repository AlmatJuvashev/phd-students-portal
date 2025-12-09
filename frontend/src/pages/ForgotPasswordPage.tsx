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
import { Link } from "react-router-dom";

const Schema = z.object({
  email: z.string().email("Invalid email address"),
});
type Form = z.infer<typeof Schema>;

export function ForgotPasswordPage() {
  const { t: T } = useTranslation("common");
  const [success, setSuccess] = React.useState(false);
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<Form>({ resolver: zodResolver(Schema) });

  const onSubmit = async (data: Form) => {
    try {
      await api.post("/auth/forgot-password", data);
      setSuccess(true);
    } catch (e: any) {
      console.error("Forgot password failed", e);
      // For security, we might want to show success even on error, or generic error
      // But typically for UX we show generic error if something broke
      setError("root", { message: "Failed to send reset link. Please try again." });
    }
  };

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <Card className="w-full max-w-sm p-6 shadow-md text-center">
          <h2 className="text-xl font-semibold mb-4">Check your email</h2>
          <p className="text-muted-foreground mb-6">
            If an account exists for that email, we have sent password reset instructions.
          </p>
          <Button asChild className="w-full">
            <Link to="/login">Back to Login</Link>
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <Card className="w-full max-w-sm p-6 shadow-md">
        <div className="mb-4">
          <h2 className="text-xl font-semibold">Reset Password</h2>
          <p className="text-xs text-muted-foreground mt-1">
            Enter your email to receive reset instructions
          </p>
        </div>
        
        {errors.root && (
          <div className="mb-4 p-3 rounded-md bg-rose-50 border border-rose-200 text-rose-700 text-sm">
            {errors.root.message}
          </div>
        )}

        <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-1">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="name@example.com"
              {...register("email")}
            />
            {errors.email && (
              <div className="text-xs text-rose-600">
                {errors.email.message}
              </div>
            )}
          </div>
          
          <Button
            className="w-full h-11"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Sending..." : "Send Reset Link"}
          </Button>
        </form>
        <div className="mt-4 text-sm text-center">
          <Link
            to="/login"
            className="text-muted-foreground hover:underline"
          >
            Back to Login
          </Link>
        </div>
      </Card>
    </div>
  );
}
