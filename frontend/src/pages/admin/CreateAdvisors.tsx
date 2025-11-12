import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { api } from "@/api/client";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const Schema = z.object({
  first_name: z.string().min(1, "Required"),
  last_name: z.string().min(1, "Required"),
  email: z.string().email("Invalid email"),
});

type Form = z.infer<typeof Schema>;

export function CreateAdvisors() {
  const [created, setCreated] = React.useState<{ username: string; temp_password: string } | null>(null);
  const { register, handleSubmit, formState: { errors, isSubmitting }, reset } = useForm<Form>({ resolver: zodResolver(Schema) });

  const onSubmit = async (data: Form) => {
    try {
      const res = await api("/admin/users", {
        method: "POST",
        body: JSON.stringify({ ...data, role: "advisor" }),
      });
      setCreated(res);
      reset();
    } catch (e: any) {
      alert(e.message || "Failed to create advisor");
    }
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div>
        <h2 className="text-2xl font-bold">Create Advisor</h2>
        <p className="text-muted-foreground">Add a new academic advisor.</p>
      </div>
      <Card>
        <CardHeader><CardTitle>Advisor Details</CardTitle></CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
            <div>
              <Input placeholder="First name" {...register("first_name")} />
              {errors.first_name && <div className="text-xs text-red-600 mt-1">{errors.first_name.message}</div>}
            </div>
            <div>
              <Input placeholder="Last name" {...register("last_name")} />
              {errors.last_name && <div className="text-xs text-red-600 mt-1">{errors.last_name.message}</div>}
            </div>
            <div>
              <Input type="email" placeholder="Email" {...register("email")} />
              {errors.email && <div className="text-xs text-red-600 mt-1">{errors.email.message}</div>}
            </div>
            <div className="pt-2">
              <Button type="submit" disabled={isSubmitting}>{isSubmitting ? "Creating…" : "Create Advisor"}</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader><CardTitle className="text-green-800">Advisor Created</CardTitle></CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="text-sm">Username: <span className="font-mono">{created.username}</span></div>
              <div className="text-sm">Temp password: <span className="font-mono">{created.temp_password}</span></div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

export default CreateAdvisors;

