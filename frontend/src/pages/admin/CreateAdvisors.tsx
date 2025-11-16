import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { api } from "@/api/client";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Copy, Loader2, RefreshCw, Search, ChevronUp, ChevronDown, ChevronLeft, ChevronRight, Trash2, CheckCircle } from "lucide-react";
import { ConfirmModal } from "@/features/forms/ConfirmModal";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

const Schema = z.object({
  first_name: z.string().min(1, "Required"),
  last_name: z.string().min(1, "Required"),
  email: z.string().email("Invalid email"),
});

type Form = z.infer<typeof Schema>;
type UserRow = {
  id: string;
  name: string;
  email: string;
  username?: string;
  created_at?: string;
  role?: string;
};
type PaginatedResponse = {
  data: UserRow[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
};

export function CreateAdvisors() {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const [created, setCreated] = React.useState<{ username: string; temp_password: string } | null>(null);
  const [resetInfo, setResetInfo] = React.useState<{ username: string; temp_password: string } | null>(null);
  const [searchTerm, setSearchTerm] = React.useState("");
  const [sortField, setSortField] = React.useState<"name" | "username" | "email" | "created_at">("name");
  const [sortDirection, setSortDirection] = React.useState<"asc" | "desc">("asc");
  const [page, setPage] = React.useState(1);
  const [pendingResetId, setPendingResetId] = React.useState<string | null>(null);
  const [pendingActiveId, setPendingActiveId] = React.useState<string | null>(null);
  const [activeFilter, setActiveFilter] = React.useState<"all" | "active" | "inactive">("all");
  const [confirmState, setConfirmState] = React.useState<{ open: boolean; kind: "reset" | "deactivate" | "activate" | null; advisor: UserRow | null }>({ open: false, kind: null, advisor: null });
  const PAGE_SIZE = 10;
  const { register, handleSubmit, formState: { errors }, reset } = useForm<Form>({ resolver: zodResolver(Schema) });

  const { data: usersResponse, isLoading, isError, refetch } = useQuery<PaginatedResponse>({
    queryKey: ["admin", "users", "advisors"],
    queryFn: () => api("/admin/users?role=advisor&limit=200&active=all"),
  });

  const advisors = React.useMemo(() => usersResponse?.data || [], [usersResponse]);

  const createAdvisorMutation = useMutation({
    mutationFn: (payload: Form) =>
      api("/admin/users", {
        method: "POST",
        body: JSON.stringify({ ...payload, role: "advisor" }),
      }),
    onSuccess: (res: { username: string; temp_password: string }) => {
      setCreated(res);
      reset();
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    onError: (err: any) => alert(err?.message || "Failed to create advisor"),
  });

  const resetPasswordMutation = useMutation({
    mutationFn: (userId: string) =>
      api(`/admin/users/${userId}/reset-password`, { method: "POST" }),
    onSuccess: (res: { username: string; temp_password: string }) => {
      setResetInfo(res);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    onError: (err: any) => alert(err?.message || "Failed to reset password"),
  });

  const setActiveMutation = useMutation({
    mutationFn: (payload: { id: string; active: boolean }) =>
      api(`/admin/users/${payload.id}/active`, {
        method: "PATCH",
        body: JSON.stringify({ active: payload.active }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      refetch();
    },
    onError: (err: any) => alert(err?.message || "Failed to update status"),
  });

  const onSubmit = (data: Form) => createAdvisorMutation.mutate(data);

  const copyCredentials = (creds: { username: string; temp_password: string }) => {
    navigator.clipboard.writeText(`Username: ${creds.username}\nPassword: ${creds.temp_password}`);
  };

  const handleReset = (advisor: UserRow) => {
    setConfirmState({ open: true, kind: "reset", advisor });
  };

  const handleDeactivate = (advisor: UserRow) => {
    setConfirmState({ open: true, kind: "deactivate", advisor });
  };

  const handleActivate = (advisor: UserRow) => {
    setConfirmState({ open: true, kind: "activate", advisor });
  };

  const formatDate = (value?: string) => {
    if (!value) return "—";
    const d = new Date(value);
    return Number.isNaN(d.getTime()) ? value : d.toLocaleDateString();
  };

  const filteredAdvisors = React.useMemo(() => {
    const term = searchTerm.trim().toLowerCase();
    const filtered = term
      ? advisors.filter((advisor) =>
          [advisor.name, advisor.email, advisor.username]
            .filter(Boolean)
            .join(" ")
            .toLowerCase()
            .includes(term)
        )
      : advisors;
    const filteredByStatus = filtered.filter((advisor) => {
      if (activeFilter === "active" && advisor.is_active === false) return false;
      if (activeFilter === "inactive" && advisor.is_active !== false) return false;
      return true;
    });
    return [...filteredByStatus].sort((a, b) => {
      const aVal = (a[sortField] || "").toString().toLowerCase();
      const bVal = (b[sortField] || "").toString().toLowerCase();
      if (aVal < bVal) return sortDirection === "asc" ? -1 : 1;
      if (aVal > bVal) return sortDirection === "asc" ? 1 : -1;
      return 0;
    });
  }, [advisors, searchTerm, sortField, sortDirection, activeFilter]);

  const totalPages = Math.max(1, Math.ceil(filteredAdvisors.length / PAGE_SIZE));
  const currentPage = Math.min(page, totalPages);
  const paginated = React.useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE;
    return filteredAdvisors.slice(start, start + PAGE_SIZE);
  }, [filteredAdvisors, currentPage]);

  React.useEffect(() => setPage(1), [searchTerm, sortField, sortDirection, advisors.length]);
  React.useEffect(() => setPage(1), [activeFilter]);

  const handleSort = (field: typeof sortField) => {
    if (sortField === field) {
      setSortDirection((prev) => (prev === "asc" ? "desc" : "asc"));
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold">{t('admin.forms.create_advisor.title','Create Advisor')}</h2>
        <p className="text-muted-foreground">{t('admin.forms.create_advisor.subtitle','Add a new academic advisor.')}</p>
      </div>
      <Card>
        <CardHeader><CardTitle>{t('admin.forms.create_advisor.heading','Advisor Details')}</CardTitle></CardHeader>
        <CardContent>
          <form className="grid grid-cols-1 gap-4 md:grid-cols-2" onSubmit={handleSubmit(onSubmit)}>
            <div>
              <Input placeholder={t('admin.forms.first_name','First name')} {...register("first_name")} />
              {errors.first_name && <div className="text-xs text-red-600 mt-1">{errors.first_name.message}</div>}
            </div>
            <div>
              <Input placeholder={t('admin.forms.last_name','Last name')} {...register("last_name")} />
              {errors.last_name && <div className="text-xs text-red-600 mt-1">{errors.last_name.message}</div>}
            </div>
            <div className="md:col-span-2">
              <Input type="email" placeholder={t('admin.forms.email','Email')} {...register("email")} />
              {errors.email && <div className="text-xs text-red-600 mt-1">{errors.email.message}</div>}
            </div>
            <div className="md:col-span-2 pt-2">
              <Button type="submit" className="w-full md:w-auto" disabled={createAdvisorMutation.isPending}>
                {createAdvisorMutation.isPending ? t('admin.forms.creating','Creating…') : t('admin.forms.create_advisor.submit','Create Advisor')}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader><CardTitle className="text-green-800">{t('admin.forms.create_advisor.success','Advisor Created')}</CardTitle></CardHeader>
          <CardContent className="flex flex-wrap items-center gap-4">
            <div>
              <div className="text-sm">{t('admin.forms.username','Username')}: <span className="font-mono">{created.username}</span></div>
              <div className="text-sm">{t('admin.forms.temp_password','Temp password')}: <span className="font-mono">{created.temp_password}</span></div>
            </div>
            <Button variant="outline" size="sm" className="gap-2" onClick={() => copyCredentials(created)}>
              <Copy className="h-4 w-4" />
              {t('admin.forms.copy_credentials','Copy credentials')}
            </Button>
          </CardContent>
        </Card>
      )}

      {resetInfo && (
        <Card className="border-blue-200 bg-blue-50">
          <CardHeader><CardTitle className="text-blue-900">{t('admin.review.password_reset','Password reset')}</CardTitle></CardHeader>
          <CardContent className="flex flex-wrap items-center gap-4">
            <div>
              <div className="text-sm">{t('admin.forms.username','Username')}: <span className="font-mono">{resetInfo.username}</span></div>
              <div className="text-sm">{t('admin.forms.temp_password','Temp password')}: <span className="font-mono">{resetInfo.temp_password}</span></div>
            </div>
            <Button variant="outline" size="sm" className="gap-2" onClick={() => copyCredentials(resetInfo)}>
              <Copy className="h-4 w-4" />
              {t('admin.forms.copy_credentials','Copy credentials')}
            </Button>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader className="space-y-4">
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div>
              <CardTitle>{t('admin.forms.advisors_table','Advisors')} · {advisors.length}</CardTitle>
              <p className="text-sm text-muted-foreground">
                {t('admin.forms.advisors_summary','{{count}} advisors · page {{page}} of {{pages}}')
                  .replace('{{count}}', filteredAdvisors.length.toString())
                  .replace('{{page}}', currentPage.toString())
                  .replace('{{pages}}', totalPages.toString())}
              </p>
            </div>
            <div className="relative w-full sm:w-64">
              <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                value={searchTerm}
                onChange={(event) => setSearchTerm(event.target.value)}
                placeholder={t('admin.forms.search_advisors_input','Search advisors…')}
                className="pl-9"
              />
            </div>
          </div>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <label className="mb-1 block text-xs text-muted-foreground">
                {t("admin.forms.active_state", { defaultValue: "Status" })}
              </label>
              <Select value={activeFilter} onValueChange={(v: any) => setActiveFilter(v)}>
                <SelectTrigger>
                  <SelectValue placeholder={t("admin.forms.status_all", { defaultValue: "All" })} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("admin.forms.status_all", { defaultValue: "All" })}</SelectItem>
                  <SelectItem value="active">{t("admin.forms.status_active", { defaultValue: "Active" })}</SelectItem>
                  <SelectItem value="inactive">{t("admin.forms.status_inactive", { defaultValue: "Inactive" })}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="max-h-[60vh] overflow-auto rounded-md border border-border/50">
            <table className="min-w-full text-sm">
              <thead className="sticky top-0 z-20 bg-card/95 backdrop-blur text-left text-muted-foreground">
                <tr>
                  <th className="py-2 pr-4 font-medium">#</th>
                  <th className="py-2 pr-4 font-medium cursor-pointer select-none" onClick={() => handleSort('name')}>
                    <div className="flex items-center">
                      {t('table.name','Name')}
                      {sortField === 'name' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2 pr-4 font-medium cursor-pointer select-none" onClick={() => handleSort('username')}>
                    <div className="flex items-center">
                      Username
                      {sortField === 'username' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2 pr-4 font-medium cursor-pointer select-none" onClick={() => handleSort('email')}>
                    <div className="flex items-center">
                      Email
                      {sortField === 'email' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2 pr-4 font-medium cursor-pointer select-none" onClick={() => handleSort('created_at')}>
                    <div className="flex items-center">
                      {t('admin.forms.registration_date','Registered')}
                      {sortField === 'created_at' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2 text-right font-medium">{t('table.actions','Actions')}</th>
                </tr>
              </thead>
              <tbody>
                {isLoading && (
                  <tr>
                    <td colSpan={6} className="py-6 text-center text-muted-foreground">
                      <Loader2 className="mx-auto mb-2 h-5 w-5 animate-spin" />
                      {t('common.loading','Loading…')}
                    </td>
                  </tr>
                )}
                {isError && (
                  <tr>
                    <td colSpan={6} className="py-6 text-center text-red-600">
                      {t('common.error','Error')} · {t('admin.forms.advisors_error','Unable to load advisors.')}
                      <div>
                        <Button variant="outline" size="sm" className="mt-2" onClick={() => refetch()}>
                          {t('common.retry','Retry')}
                        </Button>
                      </div>
                    </td>
                  </tr>
                )}
                {!isLoading && !isError && filteredAdvisors.length === 0 && (
                  <tr>
                    <td colSpan={6} className="py-6 text-center text-muted-foreground">
                      {t('admin.forms.advisors_empty','No advisors yet.')}
                    </td>
                  </tr>
                )}
                {!isLoading && !isError &&
                  paginated.map((advisor, idx) => (
                    <tr key={advisor.id} className="border-t border-border/60">
                      <td className="py-3 pr-4 text-muted-foreground">
                        {(currentPage - 1) * PAGE_SIZE + idx + 1}
                      </td>
                      <td className="py-3 pr-4">
                        <div className="font-medium">{advisor.name}</div>
                      </td>
                      <td className="py-3 pr-4 font-mono">{advisor.username || '—'}</td>
                      <td className="py-3 pr-4">{advisor.email}</td>
                      <td className="py-3 pr-4">{formatDate(advisor.created_at)}</td>
                      <td className="py-3 text-right">
                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button variant="ghost" size="icon" onClick={() => handleReset(advisor)}>
                                {pendingResetId === advisor.id && resetPasswordMutation.isPending ? (
                                  <Loader2 className="h-4 w-4 animate-spin" />
                                ) : (
                                  <RefreshCw className="h-4 w-4" />
                                )}
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>{t('admin.review.reset_password','Reset password')}</TooltipContent>
                          </Tooltip>
                          {advisor.is_active !== false ? (
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button variant="ghost" size="icon" onClick={() => handleDeactivate(advisor)}>
                                  {pendingActiveId === advisor.id && setActiveMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <Trash2 className="h-4 w-4" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>{t("admin.forms.delete_student", { defaultValue: "Deactivate" })}</TooltipContent>
                            </Tooltip>
                          ) : (
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button variant="ghost" size="icon" onClick={() => handleActivate(advisor)}>
                                  {pendingActiveId === advisor.id && setActiveMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <CheckCircle className="h-4 w-4 text-emerald-600" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>{t("admin.forms.mark_active", { defaultValue: "Activate" })}</TooltipContent>
                            </Tooltip>
                          )}
                        </TooltipProvider>
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
          <div className="mt-4 flex flex-col items-center gap-3 sm:flex-row sm:justify-between">
            <div className="text-sm text-muted-foreground">
              {t('admin.forms.pagination_label','Page {{page}} of {{pages}}')
                .replace('{{page}}', currentPage.toString())
                .replace('{{pages}}', totalPages.toString())}
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
                {t('admin.forms.prev_page','Prev')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={currentPage === totalPages}
                className="gap-2"
              >
                {t('admin.forms.next_page','Next')}
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <ConfirmModal
        open={confirmState.open}
        onOpenChange={(open) => setConfirmState((s) => ({ ...s, open }))}
        message={(() => {
          const adv = confirmState.advisor;
          if (!adv) return "";
          const details = `${adv.username || "—"} · ${adv.email || "—"}`;
          if (confirmState.kind === "reset") {
            return `${t('admin.review.confirm_reset',{defaultValue:\"Reset this user's password?\"})}\n\n${adv.name} — ${details}`;
          }
          if (confirmState.kind === "deactivate") {
            return `${t('admin.forms.confirm_deactivate_named',{defaultValue:\"Deactivate this user?\"})}\n\n${adv.name} — ${details}`;
          }
          if (confirmState.kind === "activate") {
            return `${t('admin.forms.confirm_activate_named',{defaultValue:\"Activate this user?\"})}\n\n${adv.name} — ${details}`;
          }
          return "";
        })()}
        confirmLabel={t('common.confirm','Confirm')}
        cancelLabel={t('common.cancel','Cancel')}
        busy={
          (confirmState.kind === "reset" && resetPasswordMutation.isPending) ||
          (confirmState.kind !== "reset" && setActiveMutation.isPending)
        }
        onConfirm={() => {
          const adv = confirmState.advisor;
          if (!adv) return;
          if (confirmState.kind === "reset") {
            setPendingResetId(adv.id);
            resetPasswordMutation.mutate(adv.id, {
              onSettled: () => {
                setPendingResetId(null);
                setConfirmState({ open: false, kind: null, advisor: null });
              },
            });
            return;
          }
          if (confirmState.kind === "deactivate") {
            setPendingActiveId(adv.id);
            setActiveMutation.mutate(
              { id: adv.id, active: false },
              {
                onSettled: () => {
                  setPendingActiveId(null);
                  setConfirmState({ open: false, kind: null, advisor: null });
                },
              }
            );
            return;
          }
          if (confirmState.kind === "activate") {
            setPendingActiveId(adv.id);
            setActiveMutation.mutate(
              { id: adv.id, active: true },
              {
                onSettled: () => {
                  setPendingActiveId(null);
                  setConfirmState({ open: false, kind: null, advisor: null });
                },
              }
            );
          }
        }}
      />
    </div>
  );
}

export default CreateAdvisors;
