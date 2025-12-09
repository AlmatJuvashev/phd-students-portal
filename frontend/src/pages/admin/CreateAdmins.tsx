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
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { useTranslation } from "react-i18next";
import { Navigate } from "react-router-dom";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Copy, Loader2, RefreshCw, Search, ChevronUp, ChevronDown, ChevronLeft, ChevronRight, Trash2, CheckCircle, Pencil, X, Plus } from "lucide-react";
import { ConfirmModal } from "@/features/forms/ConfirmModal";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Modal } from "@/components/ui/modal";

type Form = {
  first_name: string;
  last_name: string;
  email: string;
};

type UserRow = {
  id: string;
  name: string;
  email: string;
  username?: string;
  role: string;
  created_at?: string;
  is_active?: boolean;
};
type PaginatedResponse = {
  data: UserRow[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
};

export function CreateAdmins() {
  const { user } = useAuth();
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const Schema = React.useMemo(
    () =>
      z.object({
        first_name: z
          .string()
          .min(1, t("validation.required", { defaultValue: "Required" })),
        last_name: z
          .string()
          .min(1, t("validation.required", { defaultValue: "Required" })),
        email: z
          .string()
          .email(
            t("validation.invalid_email", { defaultValue: "Invalid email" })
          ),
      }),
    [t]
  );
  const [created, setCreated] = React.useState<{ username: string; temp_password: string } | null>(null);
  const [resetInfo, setResetInfo] = React.useState<{ username: string; temp_password: string } | null>(null);
  const [searchTerm, setSearchTerm] = React.useState("");
  const [sortField, setSortField] = React.useState<"name" | "username" | "email" | "created_at">("name");
  const [sortDirection, setSortDirection] = React.useState<"asc" | "desc">("asc");
  const [pendingResetId, setPendingResetId] = React.useState<string | null>(null);
  const [page, setPage] = React.useState(1);
  const [pendingActiveId, setPendingActiveId] = React.useState<string | null>(null);
  const [activeFilter, setActiveFilter] = React.useState<"all" | "active" | "inactive">("all");
  const [confirmState, setConfirmState] = React.useState<{ open: boolean; kind: "reset" | "deactivate" | "activate" | null; admin: UserRow | null }>({ open: false, kind: null, admin: null });
  const [editModal, setEditModal] = React.useState<{ open: boolean; admin: UserRow | null }>({ open: false, admin: null });
  const [showCreateModal, setShowCreateModal] = React.useState(false);
  const PAGE_SIZE = 10;
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<Form>({ resolver: zodResolver(Schema) });

  const {
    data: usersResponse,
    isLoading,
    isError,
    refetch,
  } = useQuery<PaginatedResponse>({
    queryKey: ["admin", "users", "admins"],
    queryFn: () => api("/admin/users?role=admin&limit=200&active=all"),
    enabled: !user || user.role === "superadmin",
  });

  const admins = React.useMemo(
    () => usersResponse?.data || [],
    [usersResponse]
  );

  if (user && user.role !== "superadmin") {
    return <Navigate to="/admin" replace />;
  }

  const copyCredentials = (creds: { username: string; temp_password: string }) => {
    navigator.clipboard.writeText(
      `Username: ${creds.username}\nPassword: ${creds.temp_password}`
    );
  };

  const openEdit = (admin: UserRow) => setEditModal({ open: true, admin });

  const handleReset = (admin: UserRow) => {
    setConfirmState({ open: true, kind: "reset", admin });
  };

  const handleDeactivate = (admin: UserRow) => {
    setConfirmState({ open: true, kind: "deactivate", admin });
  };

  const handleActivate = (admin: UserRow) => {
    setConfirmState({ open: true, kind: "activate", admin });
  };

  const formatDate = (value?: string) => {
    if (!value) return "—";
    const d = new Date(value);
    return Number.isNaN(d.getTime()) ? value : d.toLocaleDateString();
  };

  const filteredAdmins = React.useMemo(() => {
    const term = searchTerm.trim().toLowerCase();
    const filtered = term
      ? admins.filter((admin) =>
          [admin.name, admin.email, admin.username]
            .filter(Boolean)
            .join(" ")
            .toLowerCase()
            .includes(term)
        )
      : admins;
    const filteredByStatus = filtered.filter((admin) => {
      if (activeFilter === "active" && admin.is_active === false) return false;
      if (activeFilter === "inactive" && admin.is_active !== false) return false;
      return true;
    });
    return [...filteredByStatus].sort((a, b) => {
      const aVal = (a[sortField] || "").toString().toLowerCase();
      const bVal = (b[sortField] || "").toString().toLowerCase();
      if (aVal < bVal) return sortDirection === "asc" ? -1 : 1;
      if (aVal > bVal) return sortDirection === "asc" ? 1 : -1;
      return 0;
    });
  }, [admins, searchTerm, sortField, sortDirection, activeFilter]);

  const totalPages = Math.max(1, Math.ceil(filteredAdmins.length / PAGE_SIZE));
  const currentPage = Math.min(page, totalPages);
  const paginated = React.useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE;
    return filteredAdmins.slice(start, start + PAGE_SIZE);
  }, [filteredAdmins, currentPage]);

  React.useEffect(() => {
    setPage(1);
  }, [searchTerm, sortField, sortDirection, admins.length, activeFilter]);

  const handleSort = (field: typeof sortField) => {
    if (sortField === field) {
      setSortDirection((prev) => (prev === "asc" ? "desc" : "asc"));
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
  };

  const createAdminMutation = useMutation({
    mutationFn: (data: Form) =>
      api("/admin/users", {
        method: "POST",
        body: JSON.stringify({ ...data, role: "admin" }),
      }),
    onSuccess: (res: { username: string; temp_password: string }) => {
      setCreated(res);
      reset();
      setShowCreateModal(false);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    onError: (err: any) =>
      alert(
        err?.message ||
          t("admin.forms.errors.create_admin", {
            defaultValue: "Failed to create admin",
          })
      ),
  });

  const resetPasswordMutation = useMutation({
    mutationFn: (userId: string) =>
      api(`/admin/users/${userId}/reset-password`, { method: "POST" }),
    onSuccess: (res: { username: string; temp_password: string }) => {
      setResetInfo(res);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    onError: (err: any) =>
      alert(
        err?.message ||
          t("admin.forms.errors.reset_password", {
            defaultValue: "Failed to reset password",
          })
      ),
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
    onError: (err: any) =>
      alert(
        err?.message ||
          t("admin.forms.errors.status_update", {
            defaultValue: "Failed to update status",
          })
      ),
  });

  const updateAdminMutation = useMutation({
    mutationFn: (payload: { id: string; first_name: string; last_name: string; email: string }) =>
      api(`/admin/users/${payload.id}`, {
        method: "PUT",
        body: JSON.stringify({
          first_name: payload.first_name,
          last_name: payload.last_name,
          email: payload.email,
          role: "admin",
        }),
      }),
    onSuccess: () => {
      setEditModal({ open: false, admin: null });
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      refetch();
    },
    onError: (err: any) =>
      alert(
        err?.message ||
          t("admin.forms.errors.update_admin", {
            defaultValue: "Failed to update admin",
          })
      ),
  });

  async function onSubmit(data: Form) {
    createAdminMutation.mutate(data);
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold">{t('admin.forms.create_admins.title','Manage Admins')}</h2>
          <p className="text-muted-foreground">{t('admin.forms.create_admins.subtitle','Only superadmins can manage other admins.')}</p>
        </div>
        <Button onClick={() => setShowCreateModal(true)} className="w-full sm:w-auto">
          <Plus className="h-4 w-4 mr-2" />
          {t('admin.forms.create_admins.submit','Create Admin')}
        </Button>
      </div>

      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader>
            <CardTitle className="text-green-800">{t('admin.forms.create_admins.success','Admin Created')}</CardTitle>
          </CardHeader>
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
          <CardHeader>
            <CardTitle className="text-blue-900">{t('admin.review.password_reset','Password reset')}</CardTitle>
          </CardHeader>
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
              <CardTitle>
                {t('admin.forms.admins_table','Admins')} · {admins.length}
              </CardTitle>
              <p className="text-sm text-muted-foreground">
                {t('admin.forms.admins_summary','{{count}} admins · page {{page}} of {{pages}}')
                  .replace('{{count}}', filteredAdmins.length.toString())
                  .replace('{{page}}', currentPage.toString())
                  .replace('{{pages}}', totalPages.toString())}
              </p>
            </div>
            <div className="relative w-full sm:w-64">
              <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                value={searchTerm}
                onChange={(event) => setSearchTerm(event.target.value)}
                placeholder={t('admin.forms.search_admins','Search admins…')}
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
              <thead className="sticky top-0 z-20 border-b border-border/60 bg-white/80 backdrop-blur supports-[backdrop-filter]:bg-white/60 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                <tr>
                  <th className="py-2.5 pr-4 text-left">#</th>
                  <th className="py-2.5 pr-4 cursor-pointer select-none" onClick={() => handleSort('name')}>
                    <div className="flex items-center">
                      {t('table.name','Name')}
                      {sortField === 'name' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2.5 pr-4 cursor-pointer select-none" onClick={() => handleSort('username')}>
                    <div className="flex items-center">
                      Username
                      {sortField === 'username' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2.5 pr-4 cursor-pointer select-none" onClick={() => handleSort('email')}>
                    <div className="flex items-center">
                      Email
                      {sortField === 'email' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2.5 pr-4 cursor-pointer select-none" onClick={() => handleSort('created_at')}>
                    <div className="flex items-center">
                      {t('admin.forms.registration_date','Registered')}
                      {sortField === 'created_at' ? (sortDirection === 'asc' ? <ChevronUp className="ml-1 h-3 w-3" /> : <ChevronDown className="ml-1 h-3 w-3" />) : null}
                    </div>
                  </th>
                  <th className="py-2.5 pr-4 text-right">{t('table.actions','Actions')}</th>
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
                      {t('common.error','Error')} · {t('admin.forms.admins_error','Unable to load admins.')}
                      <div>
                        <Button variant="outline" size="sm" className="mt-2" onClick={() => refetch()}>
                          {t('common.retry','Retry')}
                        </Button>
                      </div>
                    </td>
                  </tr>
                )}
                {!isLoading && !isError && filteredAdmins.length === 0 && (
                  <tr>
                    <td colSpan={6} className="py-6 text-center text-muted-foreground">
                      {t('admin.forms.admins_empty','No admins yet.')}
                    </td>
                  </tr>
                )}
                {!isLoading && !isError &&
                  paginated.map((admin, idx) => {
                    const rowClass = idx % 2 === 0 ? "bg-background" : "bg-muted/10";
                    return (
                    <tr key={admin.id} className={`border-t border-border/60 transition-colors hover:bg-muted/30 ${rowClass}`}>
                      <td className="py-3 pr-4 text-muted-foreground">
                        {(currentPage - 1) * PAGE_SIZE + idx + 1}
                      </td>
                      <td className="py-3 pr-4">
                        <div className="font-medium">{admin.name}</div>
                        <div className="text-xs text-muted-foreground">{admin.role}</div>
                      </td>
                      <td className="py-3 pr-4 font-mono">{admin.username || '—'}</td>
                      <td className="py-3 pr-4">{admin.email}</td>
                      <td className="py-3 pr-4">{formatDate(admin.created_at)}</td>
                      <td className="py-3 pr-4 text-right">
                        <div className="flex justify-end gap-1">
                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleReset(admin)}
                              >
                                {pendingResetId === admin.id && resetPasswordMutation.isPending ? (
                                  <Loader2 className="h-4 w-4 animate-spin" />
                                ) : (
                                  <RefreshCw className="h-4 w-4" />
                                )}
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                              {t('admin.review.reset_password','Reset password')}
                            </TooltipContent>
                          </Tooltip>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button variant="ghost" size="icon" onClick={() => openEdit(admin)}>
                                <Pencil className="h-4 w-4" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>{t("admin.forms.edit_admin", { defaultValue: "Edit Admin" })}</TooltipContent>
                          </Tooltip>
                          {admin.is_active !== false ? (
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  className="text-destructive hover:text-destructive/80"
                                  onClick={() => handleDeactivate(admin)}
                                >
                                  {pendingActiveId === admin.id && setActiveMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <Trash2 className="h-4 w-4" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                {t("admin.forms.delete_admin", { defaultValue: "Deactivate Admin" })}
                              </TooltipContent>
                            </Tooltip>
                          ) : (
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleActivate(admin)}
                                >
                                  {pendingActiveId === admin.id && setActiveMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                  ) : (
                                    <CheckCircle className="h-4 w-4 text-emerald-600" />
                                  )}
                                </Button>
                              </TooltipTrigger>
                              <TooltipContent>
                                {t("admin.forms.mark_active", { defaultValue: "Activate" })}
                              </TooltipContent>
                            </Tooltip>
                          )}
                        </TooltipProvider>
                        </div>
                      </td>
                    </tr>
                  )})}
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

      <Modal open={showCreateModal} onClose={() => setShowCreateModal(false)}>
        <div className="max-w-xl w-full p-1">
          <Card>
            <CardHeader className="flex flex-row items-start justify-between gap-2">
              <div>
                <CardTitle>{t('admin.forms.create_admins.heading','New Admin')}</CardTitle>
                <p className="text-sm text-muted-foreground">
                  {t('admin.forms.create_admins.subtitle','Only superadmins can manage other admins.')}
                </p>
              </div>
              <Button variant="ghost" size="icon" onClick={() => setShowCreateModal(false)}>
                <X className="h-4 w-4" />
              </Button>
            </CardHeader>
            <CardContent>
              <form className="grid grid-cols-1 gap-4 md:grid-cols-2" onSubmit={handleSubmit(onSubmit)}>
                <div>
                  <Input placeholder={t('admin.forms.first_name','First name')} {...register("first_name")} />
                  {errors.first_name && (
                    <div className="text-xs text-red-600 mt-1">{errors.first_name.message}</div>
                  )}
                </div>
                <div>
                  <Input placeholder={t('admin.forms.last_name','Last name')} {...register("last_name")} />
                  {errors.last_name && (
                    <div className="text-xs text-red-600 mt-1">{errors.last_name.message}</div>
                  )}
                </div>
                <div className="md:col-span-2">
                  <Input type="email" placeholder={t('admin.forms.email','Email')} {...register("email")} />
                  {errors.email && (
                    <div className="text-xs text-red-600 mt-1">{errors.email.message}</div>
                  )}
                </div>
                <div className="md:col-span-2 pt-2">
                  <Button type="submit" className="w-full md:w-auto" disabled={createAdminMutation.isPending}>
                    {createAdminMutation.isPending ? t('admin.forms.creating','Creating…') : t('admin.forms.create_admins.submit','Create Admin')}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </Modal>

      <ConfirmModal
        open={confirmState.open}
        onOpenChange={(open) => setConfirmState((s) => ({ ...s, open }))}
        message={(() => {
          const adm = confirmState.admin;
          if (!adm) return "";
          const details = `${adm.username || "—"} · ${adm.email || "—"}`;
          if (confirmState.kind === "reset") {
            return `${t("admin.review.confirm_reset", { defaultValue: "Reset this user's password?" })}\n\n${adm.name} — ${details}`;
          }
          if (confirmState.kind === "deactivate") {
            return `${t("admin.forms.confirm_deactivate_named", { defaultValue: "Deactivate this user?" })}\n\n${adm.name} — ${details}`;
          }
          if (confirmState.kind === "activate") {
            return `${t("admin.forms.confirm_activate_named", { defaultValue: "Activate this user?" })}\n\n${adm.name} — ${details}`;
          }
          return "";
        })()}
        confirmLabel={t("common.confirm","Confirm")}
        cancelLabel={t("common.cancel","Cancel")}
        busy={
          (confirmState.kind === "reset" && resetPasswordMutation.isPending) ||
          (confirmState.kind !== "reset" && setActiveMutation.isPending)
        }
        onConfirm={() => {
          const adm = confirmState.admin;
          if (!adm) return;
          if (confirmState.kind === "reset") {
            setPendingResetId(adm.id);
            resetPasswordMutation.mutate(adm.id, {
              onSettled: () => {
                setPendingResetId(null);
                setConfirmState({ open: false, kind: null, admin: null });
              },
            });
            return;
          }
          if (confirmState.kind === "deactivate") {
            setPendingActiveId(adm.id);
            setActiveMutation.mutate(
              { id: adm.id, active: false },
              {
                onSettled: () => {
                  setPendingActiveId(null);
                  setConfirmState({ open: false, kind: null, admin: null });
                },
              }
            );
            return;
          }
          if (confirmState.kind === "activate") {
            setPendingActiveId(adm.id);
            setActiveMutation.mutate(
              { id: adm.id, active: true },
              {
                onSettled: () => {
                  setPendingActiveId(null);
                  setConfirmState({ open: false, kind: null, admin: null });
                },
              }
            );
          }
        }}
      />
      {editModal.admin && (
        <Modal open={editModal.open} onClose={() => setEditModal({ open: false, admin: null })}>
          <div className="max-w-xl rounded-lg bg-card p-1">
            <Card>
              <CardHeader className="flex flex-row items-start justify-between gap-2">
                <CardTitle>{t("admin.forms.edit_admin", { defaultValue: "Edit Admin" })}</CardTitle>
                <Button variant="ghost" size="icon" onClick={() => setEditModal({ open: false, admin: null })}>
                  <X className="h-4 w-4" />
                </Button>
              </CardHeader>
              <CardContent>
                <EditAdminForm
                  admin={editModal.admin}
                  busy={updateAdminMutation.isPending}
                  onSubmit={(values) => updateAdminMutation.mutate(values)}
                />
              </CardContent>
            </Card>
          </div>
        </Modal>
      )}
    </div>
  );
}

export default CreateAdmins;

function splitName(name?: string) {
  const parts = (name || "").trim().split(/\s+/);
  if (parts.length === 0) return { first: "", last: "" };
  if (parts.length === 1) return { first: parts[0], last: "" };
  return { first: parts[0], last: parts.slice(1).join(" ") };
}

function EditAdminForm({
  admin,
  onSubmit,
  busy,
}: {
  admin: UserRow;
  onSubmit: (payload: { id: string; first_name: string; last_name: string; email: string }) => void;
  busy?: boolean;
}) {
  const { t } = useTranslation("common");
  const { first, last } = splitName(admin.name);
  const [firstName, setFirst] = React.useState(first);
  const [lastName, setLast] = React.useState(last);
  const [email, setEmail] = React.useState(admin.email || "");

  return (
    <form
      className="grid grid-cols-1 gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit({ id: admin.id, first_name: firstName, last_name: lastName, email });
      }}
    >
      <Input
        value={firstName}
        onChange={(e) => setFirst(e.target.value)}
        placeholder={t("admin.forms.first_name", { defaultValue: "First name" })}
      />
      <Input
        value={lastName}
        onChange={(e) => setLast(e.target.value)}
        placeholder={t("admin.forms.last_name", { defaultValue: "Last name" })}
      />
      <Input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder={t("admin.forms.email", { defaultValue: "Email" })}
      />
      <div className="pt-2">
        <Button type="submit" disabled={busy} className="w-full">
          {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : t("common.save", { defaultValue: "Save" })}
        </Button>
      </div>
    </form>
  );
}
