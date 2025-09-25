import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { api } from "../api/client";
import { Input } from "../components/ui/input";
import { Button } from "../components/ui/button";
import { useToast } from "../components/toast";

const Schema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
});
type Form = z.infer<typeof Schema>;

export function LoginPage() {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<Form>({ resolver: zodResolver(Schema) });
  const { push } = useToast();
  const onSubmit = async (data: Form) => {
    try {
      const res = await api("/auth/login", {
        method: "POST",
        body: JSON.stringify(data),
      });
      localStorage.setItem("token", res.token);
      location.href = "/";
    } catch (e: any) {
      push({ title: "Login failed", description: e.message });
    }
  };
  return (
    <div className="max-w-sm mx-auto mt-10">
      <h2 className="text-xl font-semibold mb-4">Login</h2>
      <form className="space-y-3" onSubmit={handleSubmit(onSubmit)}>
        <Input placeholder="Username" {...register("username")} />
        {errors.username && (
          <div className="text-xs text-rose-600">{errors.username.message}</div>
        )}
        <Input
          placeholder="Password"
          type="password"
          {...register("password")}
        />
        {errors.password && (
          <div className="text-xs text-rose-600">{errors.password.message}</div>
        )}
        <Button className="w-full" disabled={isSubmitting}>
          Sign in
        </Button>
      </form>
      <div className="mt-3 text-sm">
        <a href="/forgot-password" className="underline">
          Forgot password?
        </a>
      </div>
    </div>
  );
}
