import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { api } from "@/api/client";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Modal } from "@/components/ui/modal";
import { ConfirmModal } from "@/features/forms/ConfirmModal";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
  Pencil,
  Trash2,
  CheckCircle,
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
  is_active?: boolean;
};
type Creds = { username: string; temp_password: string };

type PaginatedResponse = {
  data: StudentRow[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
};

const SERVER_PAGE_SIZE = 50;
const CLIENT_PAGE_SIZE = 10;

export function CreateStudents() {
  const { t } = useTranslation("common");
  const queryClient = useQueryClient();

  const [advisorSearch, setAdvisorSearch] = React.useState("");
  const [selectedAdvisors, setSelectedAdvisors] = React.useState<UserLite[]>(
    []
  );
  const [showModal, setShowModal] = React.useState(false);
  const [created, setCreated] = React.useState<Creds | null>(null);
  const [resetInfo, setResetInfo] = React.useState<Creds | null>(null);
  const [editModal, setEditModal] = React.useState<{
    open: boolean;
    student: StudentRow | null;
  }>({ open: false, student: null });
  const [pendingDeleteId, setPendingDeleteId] = React.useState<string | null>(
    null
  );
  const [pendingActiveId, setPendingActiveId] = React.useState<string | null>(
    null
  );
  const [confirmState, setConfirmState] = React.useState<{
    open: boolean;
    kind: "reset" | "deactivate" | "activate" | null;
    student: StudentRow | null;
  }>({ open: false, kind: null, student: null });
  const [searchTerm, setSearchTerm] = React.useState("");
  const [sortField, setSortField] = React.useState<
    "name" | "username" | "program" | "department" | "cohort" | "created_at"
  >("name");
  const [sortDirection, setSortDirection] = React.useState<"asc" | "desc">(
    "asc"
  );
  const [serverPage, setServerPage] = React.useState(1);
  const [clientPage, setClientPage] = React.useState(1);
  const [activeFilter, setActiveFilter] = React.useState<
    "all" | "active" | "inactive"
  >("all");
  const [filterProgram, setFilterProgram] = React.useState<string>("");
  const [filterDepartment, setFilterDepartment] = React.useState<string>("");
  const [filterCohort, setFilterCohort] = React.useState<string>("");

  const { data: advisorResponse } = useQuery<PaginatedResponse>({
    queryKey: ["admin", "advisors", advisorSearch],
    queryFn: () =>
      api(`/admin/users?role=advisor&q=${encodeURIComponent(advisorSearch)}`),
  });

  const advisors = React.useMemo(
    () => advisorResponse?.data || [],
    [advisorResponse]
  );

  const {
    data: usersResponse,
    isLoading: studentsLoading,
    isError: studentsError,
    refetch: refetchStudents,
  } = useQuery<PaginatedResponse>({
    queryKey: ["admin", "users", serverPage],
    queryFn: async () => {
      const result = await api(
        `/admin/users?page=${serverPage}&limit=${SERVER_PAGE_SIZE}&active=all`
      );
      return result;
    },
    refetchOnMount: true,
    staleTime: 0,
  });

  const allUsers = React.useMemo(
    () => usersResponse?.data || [],
    [usersResponse]
  );
  const programOptions = React.useMemo(
    () =>
      Array.from(
        new Set((allUsers || []).map((u: any) => u.program).filter(Boolean))
      ).sort(),
    [allUsers]
  );
  const departmentOptions = React.useMemo(
    () =>
      Array.from(
        new Set((allUsers || []).map((u: any) => u.department).filter(Boolean))
      ).sort(),
    [allUsers]
  );
  const cohortOptions = React.useMemo(
    () =>
      Array.from(
        new Set((allUsers || []).map((u: any) => u.cohort).filter(Boolean))
      ).sort(),
    [allUsers]
  );

  const students = React.useMemo(() => {
    return allUsers.filter((user) => user.role === "student");
  }, [allUsers]);

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
      console.log("[CreateStudents] Student created successfully:", result);
      setCreated(result);
      setShowModal(false);
      reset();
      setAdvisorSearch("");
      setSelectedAdvisors([]);
      console.log("[CreateStudents] Invalidating queries and refetching...");
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      refetchStudents().then(() => {
        console.log("[CreateStudents] Refetch complete");
      });
    },
    onError: (err: any) => {
      console.error("[CreateStudents] Error creating student:", err);
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

  const updateStudentMutation = useMutation({
    mutationFn: (payload: {
      id: string;
      first_name: string;
      last_name: string;
      email: string;
      program?: string;
      department?: string;
      cohort?: string;
    }) =>
      api(`/admin/users/${payload.id}`, {
        method: "PUT",
        body: JSON.stringify({
          first_name: payload.first_name,
          last_name: payload.last_name,
          email: payload.email,
          role: "student",
          program: payload.program || "",
          department: payload.department || "",
          cohort: payload.cohort || "",
        }),
      }),
    onSuccess: () => {
      setEditModal({ open: false, student: null });
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      refetchStudents();
    },
    onError: (err: any) => alert(err?.message || "Failed to update student"),
  });

  const setActiveMutation = useMutation({
    mutationFn: (payload: { id: string; active: boolean }) =>
      api(`/admin/users/${payload.id}/active`, {
        method: "PATCH",
        body: JSON.stringify({ active: payload.active }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      refetchStudents();
    },
    onError: (err: any) =>
      alert(err?.message || "Failed to update active status"),
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

  const openEdit = (student: StudentRow) => {
    setEditModal({ open: true, student });
  };

  const handleDeactivate = (student: StudentRow) => {
    setConfirmState({ open: true, kind: "deactivate", student });
  };

  const handleResetPassword = (student: StudentRow) => {
    setConfirmState({ open: true, kind: "reset", student });
  };

  const handleActivate = (student: StudentRow) => {
    setConfirmState({ open: true, kind: "activate", student });
  };

  const formatDate = (value?: string) => {
    if (!value) return "—";
    const d = new Date(value);
    if (Number.isNaN(d.getTime())) return value;
    return d.toLocaleDateString();
  };

  const filteredStudents = React.useMemo(() => {
    const term = searchTerm.trim().toLowerCase();
    const base = normalizedStudents;
    const filteredBySearch = term
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
    const filtered = filteredBySearch.filter((s: any) => {
      if (activeFilter === "active" && s.is_active === false) return false;
      if (activeFilter === "inactive" && s.is_active !== false) return false;
      if (filterProgram && (s.program || "") !== filterProgram) return false;
      if (filterDepartment && (s.department || "") !== filterDepartment)
        return false;
      if (filterCohort && (s.cohort || "") !== filterCohort) return false;
      return true;
    });
    const sorted = [...filtered].sort((a: any, b: any) => {
      const aVal = (a[sortField] || "").toString().toLowerCase();
      const bVal = (b[sortField] || "").toString().toLowerCase();
      if (aVal < bVal) return sortDirection === "asc" ? -1 : 1;
      if (aVal > bVal) return sortDirection === "asc" ? 1 : -1;
      return 0;
    });
    return sorted;
  }, [
    normalizedStudents,
    searchTerm,
    sortField,
    sortDirection,
    activeFilter,
    filterProgram,
    filterDepartment,
    filterCohort,
  ]);

  const totalPages = Math.max(
    1,
    Math.ceil(filteredStudents.length / CLIENT_PAGE_SIZE)
  );
  const currentPage = Math.min(clientPage, totalPages);
  const paginatedStudents = React.useMemo(() => {
    const start = (currentPage - 1) * CLIENT_PAGE_SIZE;
    const result = filteredStudents.slice(start, start + CLIENT_PAGE_SIZE);
    return result;
  }, [filteredStudents, currentPage]);

  React.useEffect(() => {
    setClientPage(1);
  }, [searchTerm, sortField, sortDirection, normalizedStudents.length]);

  const handleSort = (
    field:
      | "name"
      | "username"
      | "program"
      | "department"
      | "cohort"
      | "created_at"
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

  console.log(
    "[CreateStudents] Render - Loading:",
    studentsLoading,
    "Error:",
    studentsError,
    "Paginated:",
    paginatedStudents.length
  );

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold">
            {t("admin.forms.create_students_title", {
              defaultValue: "Manage Students",
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
          {t("admin.forms.create_student.submit", {
            defaultValue: "Create Student",
          })}
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
                {t("admin.forms.temp_password", {
                  defaultValue: "Temp password",
                })}
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
              {t("admin.forms.copy_credentials", {
                defaultValue: "Copy credentials",
              })}
            </Button>
          </CardContent>
        </Card>
      )}

      {resetInfo && (
        <Card className="border-blue-200 bg-blue-50">
          <CardHeader>
            <CardTitle className="text-blue-900">
              {t("admin.review.password_reset", {
                defaultValue: "Password reset",
              })}
            </CardTitle>
          </CardHeader>
          <CardContent className="flex flex-wrap items-center gap-4">
            <div>
              <div className="text-sm">
                {t("admin.forms.username", { defaultValue: "Username" })}:&nbsp;
                <span className="font-mono">{resetInfo.username}</span>
              </div>
              <div className="text-sm">
                {t("admin.forms.temp_password", {
                  defaultValue: "Temp password",
                })}
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
              {t("admin.forms.copy_credentials", {
                defaultValue: "Copy credentials",
              })}
            </Button>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader className="space-y-4">
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div>
              <CardTitle>
                {t("admin.forms.students_table", { defaultValue: "Students" })}{" "}
                · {normalizedStudents.length}
              </CardTitle>
              <p className="text-sm text-muted-foreground">
                {t("admin.forms.students_summary", {
                  defaultValue:
                    "{{count}} students · page {{page}} of {{pages}}",
                })
                  .replace("{{count}}", filteredStudents.length.toString())
                  .replace("{{page}}", currentPage.toString())
                  .replace("{{pages}}", totalPages.toString())}
              </p>
            </div>
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
          </div>
        </CardHeader>
        <CardContent>
          {/* Filters */}
          <div className="mb-4 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <label className="mb-1 block text-xs text-muted-foreground">
                {t("admin.forms.active_state", { defaultValue: "Status" })}
              </label>
              <Select
                value={activeFilter}
                onValueChange={(v: any) => setActiveFilter(v)}
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.forms.status_all", {
                      defaultValue: "All",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">
                    {t("admin.forms.status_all", { defaultValue: "All" })}
                  </SelectItem>
                  <SelectItem value="active">
                    {t("admin.forms.status_active", { defaultValue: "Active" })}
                  </SelectItem>
                  <SelectItem value="inactive">
                    {t("admin.forms.status_inactive", {
                      defaultValue: "Inactive",
                    })}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div>
              <label className="mb-1 block text-xs text-muted-foreground">
                {t("admin.forms.program", { defaultValue: "Program" })}
              </label>
              <Select
                value={filterProgram || "__all_programs__"}
                onValueChange={(v: any) =>
                  setFilterProgram(v === "__all_programs__" ? "" : v)
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.forms.all_programs", {
                      defaultValue: "All programs",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__all_programs__">
                    {t("admin.forms.all_programs", {
                      defaultValue: "All programs",
                    })}
                  </SelectItem>
                  {programOptions.map((p) => (
                    <SelectItem key={p} value={p}>
                      {p}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div>
              <label className="mb-1 block text-xs text-muted-foreground">
                {t("admin.forms.department", { defaultValue: "Department" })}
              </label>
              <Select
                value={filterDepartment || "__all_departments__"}
                onValueChange={(v: any) =>
                  setFilterDepartment(v === "__all_departments__" ? "" : v)
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.forms.all_departments", {
                      defaultValue: "All departments",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__all_departments__">
                    {t("admin.forms.all_departments", {
                      defaultValue: "All departments",
                    })}
                  </SelectItem>
                  {departmentOptions.map((d) => (
                    <SelectItem key={d} value={d}>
                      {d}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div>
              <label className="mb-1 block text-xs text-muted-foreground">
                {t("admin.forms.cohort", { defaultValue: "Cohort" })}
              </label>
              <Select
                value={filterCohort || "__all_cohorts__"}
                onValueChange={(v: any) =>
                  setFilterCohort(v === "__all_cohorts__" ? "" : v)
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.forms.all_cohorts", {
                      defaultValue: "All cohorts",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__all_cohorts__">
                    {t("admin.forms.all_cohorts", {
                      defaultValue: "All cohorts",
                    })}
                  </SelectItem>
                  {cohortOptions.map((c) => (
                    <SelectItem key={c} value={c}>
                      {c}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <div className="max-h-[60vh] overflow-auto rounded-md border border-border/50">
            <table className="min-w-full text-sm">
              <thead className="sticky top-0 z-20 bg-card/95 backdrop-blur text-left text-muted-foreground">
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
                      {t("admin.forms.department", {
                        defaultValue: "Department",
                      })}
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
                      {t("admin.forms.registration_date", {
                        defaultValue: "Registered",
                      })}
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
                    <td
                      colSpan={8}
                      className="py-6 text-center text-muted-foreground"
                    >
                      <Loader2 className="mx-auto mb-2 h-5 w-5 animate-spin" />
                      {t("common.loading", { defaultValue: "Loading…" })}
                    </td>
                  </tr>
                )}
                {studentsError && (
                  <tr>
                    <td colSpan={8} className="py-6 text-center text-red-600">
                      {t("common.error", { defaultValue: "Error" })}:{" "}
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
                {!studentsLoading &&
                  !studentsError &&
                  filteredStudents.length === 0 && (
                    <tr>
                      <td
                        colSpan={8}
                        className="py-6 text-center text-muted-foreground"
                      >
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
                        {(currentPage - 1) * CLIENT_PAGE_SIZE + idx + 1}
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
                        <div className="flex justify-end gap-1">
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
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => openEdit(student)}
                                aria-label={t("admin.forms.edit_student", {
                                  defaultValue: "Edit",
                                })}
                              >
                                <Pencil className="h-4 w-4" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                              {t("admin.forms.edit_student", {
                                defaultValue: "Edit",
                              })}
                            </TooltipContent>
                          </Tooltip>
                          {student.is_active ? (
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleDeactivate(student)}
                                  aria-label={t("admin.forms.delete_student", {
                                    defaultValue: "Deactivate",
                                  })}
                                >
                                  {pendingDeleteId === student.id &&
                                  setActiveMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <Trash2 className="h-4 w-4" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                {t("admin.forms.delete_student", {
                                  defaultValue: "Deactivate",
                                })}
                              </TooltipContent>
                            </Tooltip>
                          ) : (
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleActivate(student)}
                                  aria-label={t("admin.forms.mark_active", {
                                    defaultValue: "Mark Active",
                                  })}
                                >
                                  {pendingActiveId === student.id &&
                                  setActiveMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <CheckCircle className="h-4 w-4 text-emerald-600" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                {t("admin.forms.mark_active", {
                                  defaultValue: "Mark Active",
                                })}
                              </TooltipContent>
                            </Tooltip>
                          )}
                        </TooltipProvider>
                        </div>
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
                onClick={() => setClientPage((p) => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="gap-2"
              >
                <ChevronLeft className="h-4 w-4" />
                {t("admin.forms.prev_page", { defaultValue: "Prev" })}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() =>
                  setClientPage((p) => Math.min(totalPages, p + 1))
                }
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
                    defaultValue:
                      "Add a new student with program details and advisors.",
                  })}
                </p>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setShowModal(false)}
              >
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
                    placeholder={t("admin.forms.first_name", {
                      defaultValue: "First name",
                    })}
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
                    placeholder={t("admin.forms.last_name", {
                      defaultValue: "Last name",
                    })}
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
                    placeholder={t("admin.forms.program", {
                      defaultValue: "Program",
                    })}
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
                    placeholder={t("admin.forms.department", {
                      defaultValue: "Department",
                    })}
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
                    placeholder={t("admin.forms.cohort", {
                      defaultValue: "Cohort",
                    })}
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
                          aria-label={t("common.remove", {
                            defaultValue: "Remove",
                          })}
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
                  <Button
                    type="submit"
                    className="w-full"
                    disabled={createStudentMutation.isPending}
                  >
                    {createStudentMutation.isPending ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      t("admin.forms.create_student.submit", {
                        defaultValue: "Create Student",
                      })
                    )}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </Modal>

      {/* Action confirmation modal */}
      <ConfirmModal
        open={confirmState.open}
        onOpenChange={(open) => setConfirmState((s) => ({ ...s, open }))}
        message={(() => {
          const st = confirmState.student;
          const name = st?.name || "";
          const details = st
            ? `${st.username || "—"} · ${st.email || "—"}`
            : "";
          if (confirmState.kind === "reset") {
            const base = t("admin.review.confirm_reset", {
              defaultValue:
                "Reset this student's password? They will need the new temporary password to login.",
            });
            return `${base}\n\n${name} — ${details}`;
          }
          if (confirmState.kind === "deactivate") {
            const base = t("admin.forms.confirm_deactivate_named", {
              defaultValue: "Deactivate this student?",
            });
            return `${base}\n\n${name} — ${details}`;
          }
          if (confirmState.kind === "activate") {
            const base = t("admin.forms.confirm_activate_named", {
              defaultValue: "Mark this student as active?",
            });
            return `${base}\n\n${name} — ${details}`;
          }
          return "";
        })()}
        confirmLabel={t("common.confirm", { defaultValue: "Confirm" })}
        cancelLabel={t("common.cancel", { defaultValue: "Cancel" })}
        busy={
          (confirmState.kind === "reset" && resetPasswordMutation.isPending) ||
          (confirmState.kind !== "reset" && setActiveMutation.isPending)
        }
        onConfirm={() => {
          const st = confirmState;
          if (!st.student) return;
          if (st.kind === "reset") {
            setPendingResetId(st.student.id);
            resetPasswordMutation.mutate(st.student.id, {
              onSettled: () => {
                setPendingResetId(null);
                setConfirmState({ open: false, kind: null, student: null });
              },
            });
            return;
          }
          if (st.kind === "deactivate") {
            setPendingDeleteId(st.student.id);
            setActiveMutation.mutate(
              { id: st.student.id, active: false },
              {
                onSettled: () => {
                  setPendingDeleteId(null);
                  setConfirmState({ open: false, kind: null, student: null });
                },
              }
            );
            return;
          }
          if (st.kind === "activate") {
            setPendingActiveId(st.student.id);
            setActiveMutation.mutate(
              { id: st.student.id, active: true },
              {
                onSettled: () => {
                  setPendingActiveId(null);
                  setConfirmState({ open: false, kind: null, student: null });
                },
              }
            );
          }
        }}
      />

      {/* Edit Student Modal */}
      <Modal
        open={editModal.open}
        onClose={() => setEditModal({ open: false, student: null })}
      >
        {editModal.student && (
          <div className="max-w-2xl max-h-[85vh] overflow-y-auto p-1">
            <Card>
              <CardHeader className="flex flex-row items-start justify-between gap-2">
                <div>
                  <CardTitle>
                    {t("admin.forms.edit_student", {
                      defaultValue: "Edit Student",
                    })}
                  </CardTitle>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setEditModal({ open: false, student: null })}
                >
                  <X className="h-4 w-4" />
                </Button>
              </CardHeader>
              <CardContent>
                <EditStudentForm
                  student={editModal.student}
                  onSubmit={(payload) => updateStudentMutation.mutate(payload)}
                  busy={updateStudentMutation.isPending}
                />
              </CardContent>
            </Card>
          </div>
        )}
      </Modal>
    </div>
  );
}

function splitName(name?: string) {
  const n = (name || "").trim();
  if (!n) return { first: "", last: "" };
  const parts = n.split(/\s+/);
  if (parts.length === 1) return { first: parts[0], last: "" };
  return { first: parts[0], last: parts.slice(1).join(" ") };
}

function EditStudentForm({
  student,
  onSubmit,
  busy,
}: {
  student: any;
  onSubmit: (p: {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    program?: string;
    department?: string;
    cohort?: string;
  }) => void;
  busy?: boolean;
}) {
  const { t } = useTranslation("common");
  const { first, last } = splitName(student.name);
  const [firstName, setFirst] = React.useState(first);
  const [lastName, setLast] = React.useState(last);
  const [email, setEmail] = React.useState(student.email || "");
  const [program, setProgram] = React.useState(student.program || "");
  const [department, setDepartment] = React.useState(student.department || "");
  const [cohort, setCohort] = React.useState(student.cohort || "");

  return (
    <form
      className="grid grid-cols-1 gap-4 md:grid-cols-2"
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit({
          id: student.id,
          first_name: firstName,
          last_name: lastName,
          email,
          program,
          department,
          cohort,
        });
      }}
    >
      <Input
        placeholder={t("admin.forms.first_name", {
          defaultValue: "First name",
        })}
        value={firstName}
        onChange={(e) => setFirst(e.target.value)}
      />
      <Input
        placeholder={t("admin.forms.last_name", { defaultValue: "Last name" })}
        value={lastName}
        onChange={(e) => setLast(e.target.value)}
      />
      <Input
        type="email"
        placeholder={t("admin.forms.email_optional", { defaultValue: "Email" })}
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <Input
        placeholder={t("admin.forms.program", { defaultValue: "Program" })}
        value={program}
        onChange={(e) => setProgram(e.target.value)}
      />
      <Input
        placeholder={t("admin.forms.department", {
          defaultValue: "Department",
        })}
        value={department}
        onChange={(e) => setDepartment(e.target.value)}
      />
      <Input
        placeholder={t("admin.forms.cohort", { defaultValue: "Cohort" })}
        value={cohort}
        onChange={(e) => setCohort(e.target.value)}
      />
      <div className="md:col-span-2 flex gap-2 pt-2">
        <Button type="submit" disabled={busy} className="w-full">
          {busy ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            t("common.save", { defaultValue: "Save" })
          )}
        </Button>
      </div>
    </form>
  );
}

export default CreateStudents;
