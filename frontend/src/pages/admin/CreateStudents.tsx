import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/client";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

const Schema = z.object({
  first_name: z.string().min(1, "Required"),
  last_name: z.string().min(1, "Required"),
  phone: z.string().optional(),
  email: z.string().email().optional().or(z.literal("")),
  program: z.string().min(1, "Required"),
  department: z.string().min(1, "Required"),
  cohort: z.string().min(1, "Required"),
  advisor_ids: z.array(z.string()).optional(),
});

type Form = z.infer<typeof Schema>;

type UserLite = { id: string; name: string; email: string; role: string };

export function CreateStudents() {
  const [created, setCreated] = React.useState<{ username: string; temp_password: string } | null>(null);
  const [advisorSearch, setAdvisorSearch] = React.useState("");
  const [selectedAdvisors, setSelectedAdvisors] = React.useState<UserLite[]>([]);

  const { data: advisors = [] } = useQuery<UserLite[]>({
    queryKey: ["admin", "advisors", advisorSearch],
    queryFn: () => api(`/admin/users?role=advisor&q=${encodeURIComponent(advisorSearch)}`),
  });

  const { register, handleSubmit, setValue, watch, formState: { errors, isSubmitting }, reset } = useForm<Form>({
    resolver: zodResolver(Schema),
    defaultValues: { email: "", advisor_ids: [] },
  });

  React.useEffect(() => {
    setValue("advisor_ids", selectedAdvisors.map(a => a.id));
  }, [selectedAdvisors, setValue]);

  const onSubmit = async (data: Form) => {
    try {
      const res = await api("/admin/users", {
        method: "POST",
        body: JSON.stringify({ ...data, role: "student" }),
      });
      setCreated(res);
      reset({ first_name: "", last_name: "", phone: "", email: "", program: "", department: "", cohort: "", advisor_ids: [] });
      setSelectedAdvisors([]);
    } catch (e: any) {
      alert(e.message || "Failed to create student");
    }
  };

  const addAdvisor = (u: UserLite) => {
    if (selectedAdvisors.find(a => a.id === u.id)) return;
    setSelectedAdvisors(prev => [...prev, u]);
  };
  const removeAdvisor = (id: string) => setSelectedAdvisors(prev => prev.filter(a => a.id !== id));

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div>
        <h2 className="text-2xl font-bold">Create Student</h2>
        <p className="text-muted-foreground">Add a new student with program details and advisor assignments.</p>
      </div>

      <Card>
        <CardHeader><CardTitle>Student Details</CardTitle></CardHeader>
        <CardContent>
          <form className="grid grid-cols-1 md:grid-cols-2 gap-4" onSubmit={handleSubmit(onSubmit)}>
            <div className="md:col-span-1">
              <Input placeholder="First name" {...register("first_name")} />
              {errors.first_name && <div className="text-xs text-red-600 mt-1">{errors.first_name.message}</div>}
            </div>
            <div className="md:col-span-1">
              <Input placeholder="Last name" {...register("last_name")} />
              {errors.last_name && <div className="text-xs text-red-600 mt-1">{errors.last_name.message}</div>}
            </div>

            <div className="md:col-span-1">
              <Input placeholder="Phone (optional)" {...register("phone")} />
            </div>
            <div className="md:col-span-1">
              <Input type="email" placeholder="Email (optional)" {...register("email")} />
              {errors.email && <div className="text-xs text-red-600 mt-1">{errors.email.message as any}</div>}
            </div>

            <div className="md:col-span-1">
              <Input placeholder="Program" {...register("program")} />
              {errors.program && <div className="text-xs text-red-600 mt-1">{errors.program.message}</div>}
            </div>
            <div className="md:col-span-1">
              <Input placeholder="Department" {...register("department")} />
              {errors.department && <div className="text-xs text-red-600 mt-1">{errors.department.message}</div>}
            </div>

            <div className="md:col-span-1">
              <Input placeholder="Cohort" {...register("cohort")} />
              {errors.cohort && <div className="text-xs text-red-600 mt-1">{errors.cohort.message}</div>}
            </div>

            <div className="md:col-span-2 space-y-2">
              <label className="text-sm font-medium">Advisors</label>
              <div className="flex flex-wrap gap-2">
                {selectedAdvisors.map(a => (
                  <Badge key={a.id} className="gap-2">
                    {a.name}
                    <button type="button" onClick={() => removeAdvisor(a.id)} aria-label={`Remove ${a.name}`}>
                      ×
                    </button>
                  </Badge>
                ))}
              </div>
              <div className="relative">
                <Input placeholder="Search advisors…" value={advisorSearch} onChange={(e) => setAdvisorSearch(e.target.value)} />
                {advisorSearch && advisors.length > 0 && (
                  <div className="absolute z-10 bg-white border rounded mt-1 w-full max-h-56 overflow-auto shadow">
                    {advisors.map(u => (
                      <button type="button" key={u.id} className="w-full text-left px-3 py-2 hover:bg-muted" onClick={() => addAdvisor(u)}>
                        <div className="font-medium">{u.name}</div>
                        <div className="text-xs text-muted-foreground">{u.email}</div>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>

            <div className="md:col-span-2 pt-2 flex gap-2">
              <Button type="submit" disabled={isSubmitting}>{isSubmitting ? "Creating…" : "Create Student"}</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader><CardTitle className="text-green-800">Student Created</CardTitle></CardHeader>
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

export default CreateStudents;

