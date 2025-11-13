import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { api } from "@/api/client";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Modal } from "@/components/ui/modal";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  Plus,
  RefreshCw,
  Copy,
  Loader2,
  X,
  Search,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  ChevronDown,
} from "lucide-react";

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
type StudentRow = UserLite & {
  username?: string;
  program?: string;
  department?: string;
  cohort?: string;
  created_at?: string;
};
type Creds = { username: string; temp_password: string };

export function CreateStudents() {
  const { t } = useTranslation("common");
  const queryClient = useQueryClient();

  const [advisorSearch, setAdvisorSearch] = React.useState("");
  const [selectedAdvisors, setSelectedAdvisors] = React.useState<UserLite[]>([]);
  const [showModal, setShowModal] = React.useState(false);
  const [created, setCreated] = React.useState<Creds | null>(null);
  const [resetInfo, setResetInfo] = React.useState<Creds | null>(null);
  const [searchTerm, setSearchTerm] = React.useState("");
  const [sortField, setSortField] = React.useState<
    "name" | "username" | "program" | "department" | "cohort" | "created_at"
  >("name");
  const [sortDirection, setSortDirection] = React.useState<"asc" | "desc">("asc");
  const [page, setPage] = React.useState(1);

  const { data: advisors = [] } = useQuery<UserLite[]>({
    queryKey: ["admin", "advisors", advisorSearch],
    queryFn: () =>
      api(`/admin/users?role=advisor&q=${encodeURIComponent(advisorSearch)}`),
  });

  const {
    data: allUsers = [],
    isLoading: studentsLoading,
    isError: studentsError,
    refetch: refetchStudents,
  } = useQuery<StudentRow[]>({
    queryKey: ["admin", "users"],
    queryFn: () => api(`/admin/users`),
  });

  const students = React.useMemo(
    () => allUsers.filter((user) => user.role === "student"),
    [allUsers]
  );
  const normalizedStudents = React.useMemo(() => students ?? [], [students]);

  const {
    register,
    handleSubmit,
    setValue,
    reset,
    formState: { errors },
  } = useForm<Form>({
    resolver: zodResolver(Schema),
    defaultValues: {
      email: "",
      advisor_ids: [],
      phone: "",
      program: "",
      department: "",
      cohort: "",
    },
  });

  React.useEffect(() => {
    setValue(
      "advisor_ids",
      selectedAdvisors.map((a) => a.id)
    );
  }, [selectedAdvisors, setValue]);

  const createStudentMutation = useMutation({
    mutationFn: (payload: Form) =>
      api("/admin/users", {
        method: "POST",
        body: JSON.stringify({ ...payload, role: "student" }),
      }),
    onSuccess: (result: Creds) => {
      setCreated(result);
      setShowModal(false);
      reset();
      setAdvisorSearch("");
      setSelectedAdvisors([]);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    onError: (err: any) => {
      alert(err?.message || "Failed to create student");
    },
  });

  const resetPasswordMutation = useMutation({
    mutationFn: (userId: string) =>
      api(`/admin/users/${userId}/reset-password`, { method: "POST" }),
    onSuccess: (result: Creds) => {
      setResetInfo(result);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    onError: (err: any) => alert(err?.message || "Failed to reset password"),
  });

  const [pendingResetId, setPendingResetId] = React.useState<string | null>(
    null
  );

  const onSubmit = (data: Form) => createStudentMutation.mutate(data);

  const addAdvisor = (advisor: UserLite) => {
    if (selectedAdvisors.find((a) => a.id === advisor.id)) return;
    setSelectedAdvisors((prev) => [...prev, advisor]);
    setAdvisorSearch("");
  };
  const removeAdvisor = (advisorId: string) => {
    setSelectedAdvisors((prev) => prev.filter((a) => a.id !== advisorId));
  };

  const copyCredentials = (creds: Creds) => {
    const message = `Username: ${creds.username}\nPassword: ${creds.temp_password}`;
    navigator.clipboard.writeText(message);
  };

  const handleResetPassword = (student: StudentRow) => {
    if (
      confirm(
        t("admin.review.confirm_reset", {
          defaultValue:
            "Reset this student's password? They will need the new temporary password to login.",
        })
      )
    ) {
      setPendingResetId(student.id);
      resetPasswordMutation.mutate(student.id, {
        onSettled: () => setPendingResetId(null),
      });
    }
  };

  const formatDate = (value?: string) => {
    if (!value) return "—";
    const d = new Date(value);
    if (Number.isNaN(d.getTime())) return value;
    return d.toLocaleDateString();
  };

  const PAGE_SIZE = 10;
  const filteredStudents = React.useMemo(() => {
    const term = searchTerm.trim().toLowerCase();
    const base = normalizedStudents;
    const filtered = term
      ? base.filter((student) => {
          const haystack = [
            student.name,
            student.email,
            student.username,
            student.program,
            student.department,
            student.cohort,
          ]
            .filter(Boolean)
            .join(" ")
            .toLowerCase();
          return haystack.includes(term);
        })
      : base;
    const sorted = [...filtered].sort((a, b) => {
      const aVal = (a[sortField] || "").toString().toLowerCase();
      const bVal = (b[sortField] || "").toString().toLowerCase();
      if (aVal < bVal) return sortDirection === "asc" ? -1 : 1;
      if (aVal > bVal) return sortDirection === "asc" ? 1 : -1;
      return 0;
    });
    return sorted;
  }, [normalizedStudents, searchTerm, sortField, sortDirection]);

  const totalPages = Math.max(1, Math.ceil(filteredStudents.length / PAGE_SIZE));
  const currentPage = Math.min(page, totalPages);
  const paginatedStudents = React.useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE;
    return filteredStudents.slice(start, start + PAGE_SIZE);
  }, [filteredStudents, currentPage]);

  React.useEffect(() => {
    setPage(1);
  }, [searchTerm, sortField, sortDirection, normalizedStudents.length]);

  const handleSort = (
    field: "name" | "username" | "program" | "department" | "cohort" | "created_at"
  ) => {
    if (sortField === field) {
      setSortDirection((prev) => (prev === "asc" ? "desc" : "asc"));
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
  };

  const renderSortIcon = (field: typeof sortField) => {
    if (sortField !== field) return null;
    return sortDirection === "asc" ? (
      <ChevronUp className="ml-1 h-3 w-3" />
    ) : (
      <ChevronDown className="ml-1 h-3 w-3" />
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold">
            {t("admin.forms.create_students_title", {
              defaultValue: "Create Students",
            })}
          </h2>
          <p className="text-sm text-muted-foreground">
            {t("admin.forms.create_students_subtitle", {
              defaultValue:
                "Invite new students, assign advisors, and manage existing accounts.",
            })}
          </p>
        </div>
        <Button onClick={() => setShowModal(true)} className="w-full sm:w-auto">
          <Plus className="h-4 w-4 mr-2" />
          {t("admin.forms.create_student.submit", { defaultValue: "Create Student" })}
        </Button>
      </div>

      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader>
            <CardTitle className="text-green-800">
              {t("admin.forms.create_student.success", {
                defaultValue: "Student Created",
              })}
            </CardTitle>
          </CardHeader>
          <CardContent className="flex flex-wrap items-center gap-4">
            <div>
              <div className="text-sm">
                {t("admin.forms.username", { defaultValue: "Username" })}:&nbsp;
                <span className="font-mono">{created.username}</span>
              </div>
              <div className="text-sm">
                {t("admin.forms.temp_password", { defaultValue: "Temp password" })}
                :&nbsp;
                <span className="font-mono">{created.temp_password}</span>
              </div>
            </div>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => copyCredentials(created)}
              className="gap-2"
            >
              <Copy className="h-4 w-4" />
              {t("admin.forms.copy_credentials", { defaultValue: "Copy credentials" })}
            </Button>
          </CardContent>
        </Card>
      )}

      {resetInfo && (
        <Card className="border-blue-200 bg-blue-50">
          <CardHeader>
            <CardTitle className="text-blue-900">
              {t("admin.review.password_reset", { defaultValue: "Password reset" })}
            </CardTitle>
          </CardHeader>
          <CardContent className="flex flex-wrap items-center gap-4">
            <div>
              <div className="text-sm">
                {t("admin.forms.username", { defaultValue: "Username" })}:&nbsp;
                <span className="font-mono">{resetInfo.username}</span>
              </div>
              <div className="text-sm">
                {t("admin.forms.temp_password", { defaultValue: "Temp password" })}
                :&nbsp;
                <span className="font-mono">{resetInfo.temp_password}</span>
              </div>
            </div>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => copyCredentials(resetInfo)}
              className="gap-2"
            >
              <Copy className="h-4 w-4" />
              {t("admin.forms.copy_credentials", { defaultValue: "Copy credentials" })}
            </Button>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader className="space-y-4">
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div>
              <CardTitle>
                {t("admin.forms.students_table", { defaultValue: "Students" })} · {normalizedStudents.length}
              </CardTitle>
              <p className="text-sm text-muted-foreground">
                {t("admin.forms.students_summary", {
                  defaultValue: "{{count}} students · page {{page}} of {{pages}}",
                })
                  .replace("{{count}}", filteredStudents.length.toString())
                  .replace("{{page}}", currentPage.toString())
                  .replace("{{pages}}", totalPages.toString())}
              </p>
            </div>
            <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
              <div className="relative w-full sm:w-64">
                <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  value={searchTerm}
                  onChange={(event) => setSearchTerm(event.target.value)}
                  placeholder={t("admin.forms.search_students", {
                    defaultValue: "Search students…",
                  })}
                  className="pl-9"
                />
              </div>
              <Button onClick={() => setShowModal(true)} className="w-full sm:w-auto">
                <Plus className="h-4 w-4 mr-2" />
                {t("admin.forms.create_student.submit", { defaultValue: "Create Student" })}
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="min-w-full text-sm">
              <thead className="sticky top-0 z-10 bg-background text-left text-muted-foreground">
                <tr>
                  <th className="py-2 pr-4 font-medium">#</th>
                  <th
                    className="py-2 pr-4 font-medium cursor-pointer select-none"
                    onClick={() => handleSort("name")}
                  >
                    <div className="flex items-center">
                      {t("table.name", { defaultValue: "Name" })}
                      {renderSortIcon("name")}
                    </div>
                  </th>
                  <th
                    className="py-2 pr-4 font-medium cursor-pointer select-none"
                    onClick={() => handleSort("username")}
                  >
                    <div className="flex items-center">
                      Username
                      {renderSortIcon("username")}
                    </div>
                  </th>
                  <th
                    className="py-2 pr-4 font-medium cursor-pointer select-none"
                    onClick={() => handleSort("program")}
                  >
                    <div className="flex items-center">
                      {t("admin.forms.program", { defaultValue: "Program" })}
                      {renderSortIcon("program")}
                    </div>
                  </th>
                  <th
                    className="py-2 pr-4 font-medium cursor-pointer select-none"
                    onClick={() => handleSort("department")}
                  >
                    <div className="flex items-center">
                      {t("admin.forms.department", { defaultValue: "Department" })}
                      {renderSortIcon("department")}
                    </div>
                  </th>
                  <th
                    className="py-2 pr-4 font-medium cursor-pointer select-none"
                    onClick={() => handleSort("cohort")}
                  >
                    <div className="flex items-center">
                      {t("admin.forms.cohort", { defaultValue: "Cohort" })}
                      {renderSortIcon("cohort")}
                    </div>
                  </th>
                  <th
                    className="py-2 pr-4 font-medium cursor-pointer select-none"
                    onClick={() => handleSort("created_at")}
                  >
                    <div className="flex items-center">
                      {t("admin.forms.registration_date", { defaultValue: "Registered" })}
                      {renderSortIcon("created_at")}
                    </div>
                  </th>
                  <th className="py-2 text-right font-medium">
                    {t("table.actions", { defaultValue: "Actions" })}
                  </th>
                </tr>
              </thead>
              <tbody>
                {studentsLoading && (
                  <tr>
                    <td colSpan={8} className="py-6 text-center text-muted-foreground">
                      <Loader2 className="mx-auto mb-2 h-5 w-5 animate-spin" />
                      {t("common.loading", { defaultValue: "Loading…" })}
                    </td>
                  </tr>
                )}
                {studentsError && (
                  <tr>
                    <td colSpan={8} className="py-6 text-center text-red-600">
                      {t("common.error", { defaultValue: "Error" })}: {" "}
                      {t("admin.forms.students_error", {
                        defaultValue: "Unable to load students.",
                      })}
                      <div>
                        <Button
                          variant="outline"
                          size="sm"
                          className="mt-2"
                          onClick={() => refetchStudents()}
                        >
                          {t("common.retry", { defaultValue: "Retry" })}
                        </Button>
                      </div>
                    </td>
                  </tr>
                )}
                {!studentsLoading && !studentsError && filteredStudents.length === 0 && (
                  <tr>
                    <td colSpan={8} className="py-6 text-center text-muted-foreground">
                      {t("admin.forms.students_empty", {
                        defaultValue: "No students yet.",
                      })}
                    </td>
                  </tr>
                )}
                {!studentsLoading &&
                  !studentsError &&
                  paginatedStudents.map((student, idx) => (
                    <tr
                      key={student.id}
                      className="border-t border-border/60 text-foreground"
                    >
                      <td className="py-3 pr-4 align-top text-muted-foreground">
                        {(currentPage - 1) * PAGE_SIZE + idx + 1}
                      </td>
                      <td className="py-3 pr-4 align-top">
                        <div className="font-medium">{student.name}</div>
                        <div className="text-xs text-muted-foreground">
                          {student.email}
                        </div>
                      </td>
                      <td className="py-3 pr-4 align-top font-mono">
                        {student.username || "—"}
                      </td>
                      <td className="py-3 pr-4 align-top">
                        {student.program || "—"}
                      </td>
                      <td className="py-3 pr-4 align-top">
                        {student.department || "—"}
                      </td>
                      <td className="py-3 pr-4 align-top">
                        {student.cohort || "—"}
                      </td>
                      <td className="py-3 pr-4 align-top">
                        {formatDate(student.created_at)}
                      </td>
                      <td className="py-3 text-right align-top">
                        <TooltipProvider delayDuration={100}>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleResetPassword(student)}
                                aria-label={t("admin.review.reset_password", {
                                  defaultValue: "Reset password",
                                })}
                              >
                                {pendingResetId === student.id &&
                                resetPasswordMutation.isPending ? (
                                  <Loader2 className="h-4 w-4 animate-spin" />
                                ) : (
                                  <RefreshCw className="h-4 w-4" />
                                )}
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                              {t("admin.review.reset_password", {
                                defaultValue: "Reset password",
                              })}
                            </TooltipContent>
                          </Tooltip>
                        </TooltipProvider>
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
          <div className="mt-4 flex flex-col items-center gap-3 sm:flex-row sm:justify-between">
            <div className="text-sm text-muted-foreground">
              {t("admin.forms.pagination_label", {
                defaultValue: "Page {{page}} of {{pages}}",
              })
                .replace("{{page}}", currentPage.toString())
                .replace("{{pages}}", totalPages.toString())}
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="gap-2"
              >
                <ChevronLeft className="h-4 w-4" />
                {t("admin.forms.prev_page", { defaultValue: "Prev" })}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={currentPage === totalPages}
                className="gap-2"
              >
                {t("admin.forms.next_page", { defaultValue: "Next" })}
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Modal open={showModal} onClose={() => setShowModal(false)}>
        <div className="max-w-3xl max-h-[85vh] overflow-y-auto p-1">
          <Card>
            <CardHeader className="flex flex-row items-start justify-between gap-2">
              <div>
                <CardTitle>
                  {t("admin.forms.create_student.heading", {
                    defaultValue: "Student Details",
                  })}
                </CardTitle>
                <p className="text-sm text-muted-foreground">
                  {t("admin.forms.create_student.subtitle", {
                    defaultValue: "Add a new student with program details and advisors.",
                  })}
                </p>
              </div>
              <Button variant="ghost" size="icon" onClick={() => setShowModal(false)}>
                <X className="h-4 w-4" />
              </Button>
            </CardHeader>
            <CardContent>
              <form
                className="grid grid-cols-1 gap-4 md:grid-cols-2"
                onSubmit={handleSubmit(onSubmit)}
              >
                <div>
                  <Input
                    placeholder={t("admin.forms.first_name", { defaultValue: "First name" })}
                    {...register("first_name")}
                  />
                  {errors.first_name && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.first_name.message}
                    </p>
                  )}
                </div>
                <div>
                  <Input
                    placeholder={t("admin.forms.last_name", { defaultValue: "Last name" })}
                    {...register("last_name")}
                  />
                  {errors.last_name && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.last_name.message}
                    </p>
                  )}
                </div>

                <div>
                  <Input
                    placeholder={t("admin.forms.phone_optional", {
                      defaultValue: "Phone (optional)",
                    })}
                    {...register("phone")}
                  />
                </div>
                <div>
                  <Input
                    type="email"
                    placeholder={t("admin.forms.email_optional", {
                      defaultValue: "Email (optional)",
                    })}
                    {...register("email")}
                  />
                  {errors.email && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.email.message as string}
                    </p>
                  )}
                </div>

                <div>
                  <Input
                    placeholder={t("admin.forms.program", { defaultValue: "Program" })}
                    {...register("program")}
                  />
                  {errors.program && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.program.message}
                    </p>
                  )}
                </div>
                <div>
                  <Input
                    placeholder={t("admin.forms.department", { defaultValue: "Department" })}
                    {...register("department")}
                  />
                  {errors.department && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.department.message}
                    </p>
                  )}
                </div>

                <div>
                  <Input
                    placeholder={t("admin.forms.cohort", { defaultValue: "Cohort" })}
                    {...register("cohort")}
                  />
                  {errors.cohort && (
                    <p className="mt-1 text-xs text-red-600">
                      {errors.cohort.message}
                    </p>
                  )}
                </div>

                <div className="md:col-span-2 space-y-2">
                  <label className="text-sm font-medium">
                    {t("admin.forms.advisors", { defaultValue: "Advisors" })}
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {selectedAdvisors.map((advisor) => (
                      <Badge key={advisor.id} className="gap-1">
                        {advisor.name}
                        <button
                          type="button"
                          className="ml-1 text-xs"
                          onClick={() => removeAdvisor(advisor.id)}
                          aria-label={t("common.remove", { defaultValue: "Remove" })}
                        >
                          ×
                        </button>
                      </Badge>
                    ))}
                  </div>
                  <div className="relative">
                    <Input
                      placeholder={t("admin.forms.search_advisors", {
                        defaultValue: "Search advisors…",
                      })}
                      value={advisorSearch}
                      onChange={(event) => setAdvisorSearch(event.target.value)}
                    />
                    {advisorSearch && advisors.length > 0 && (
                      <div className="absolute z-50 mt-1 max-h-56 w-full overflow-auto rounded border bg-card shadow">
                        {advisors.map((advisor) => (
                          <button
                            type="button"
                            key={advisor.id}
                            className="w-full px-3 py-2 text-left hover:bg-muted"
                            onClick={() => addAdvisor(advisor)}
                          >
                            <div className="font-medium">{advisor.name}</div>
                            <div className="text-xs text-muted-foreground">
                              {advisor.email}
                            </div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                <div className="md:col-span-2 flex gap-2 pt-2">
                  <Button type="submit" className="w-full" disabled={createStudentMutation.isPending}>
                    {createStudentMutation.isPending ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      t("admin.forms.create_student.submit", { defaultValue: "Create Student" })
                    )}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </Modal>
    </div>
  );
}

export default CreateStudents;
