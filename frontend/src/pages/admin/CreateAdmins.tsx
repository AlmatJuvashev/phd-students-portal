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
import { Copy, Loader2, RefreshCw, Search, ChevronUp, ChevronDown, ChevronLeft, ChevronRight } from "lucide-react";

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
  const Schema = React.useMemo(() => z.object({
    first_name: z.string().min(1, t('validation.required','Required')),
    last_name: z.string().min(1, t('validation.required','Required')),
    email: z.string().email(t('validation.invalid_email','Invalid email')),
  }), [t]);
  const [created, setCreated] = React.useState<{ username: string; temp_password: string } | null>(null);
  const [resetInfo, setResetInfo] = React.useState<{ username: string; temp_password: string } | null>(null);
  const [searchTerm, setSearchTerm] = React.useState("");
  const [sortField, setSortField] = React.useState<"name" | "username" | "email" | "created_at">("name");
  const [sortDirection, setSortDirection] = React.useState<"asc" | "desc">("asc");
  const [pendingResetId, setPendingResetId] = React.useState<string | null>(null);
  const [page, setPage] = React.useState(1);
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

  const handleReset = (admin: UserRow) => {
    if (
      confirm(
        t("admin.review.confirm_reset", {
          defaultValue:
            "Reset this user's password? They will need the new temporary password to login.",
        })
      )
    ) {
      setPendingResetId(admin.id);
      resetPasswordMutation.mutate(admin.id, {
        onSettled: () => setPendingResetId(null),
      });
    }
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
    return [...filtered].sort((a, b) => {
      const aVal = (a[sortField] || "").toString().toLowerCase();
      const bVal = (b[sortField] || "").toString().toLowerCase();
      if (aVal < bVal) return sortDirection === "asc" ? -1 : 1;
      if (aVal > bVal) return sortDirection === "asc" ? 1 : -1;
      return 0;
    });
  }, [admins, searchTerm, sortField, sortDirection]);

  const totalPages = Math.max(1, Math.ceil(filteredAdmins.length / PAGE_SIZE));
  const currentPage = Math.min(page, totalPages);
  const paginated = React.useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE;
    return filteredAdmins.slice(start, start + PAGE_SIZE);
  }, [filteredAdmins, currentPage]);

  React.useEffect(() => {
    setPage(1);
  }, [searchTerm, sortField, sortDirection, admins.length]);

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
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
    onError: (err: any) => alert(err?.message || "Failed to create admin"),
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

  async function onSubmit(data: Form) {
    createAdminMutation.mutate(data);
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold">{t('admin.forms.create_admins.title','Create Admins')}</h2>
        <p className="text-muted-foreground">{t('admin.forms.create_admins.subtitle','Only superadmins can create other admins.')}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t('admin.forms.create_admins.heading','New Admin')}</CardTitle>
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
                  paginated.map((admin, idx) => (
                    <tr key={admin.id} className="border-t border-border/60">
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
                      <td className="py-3 text-right">
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
    </div>
  );
}

export default CreateAdmins;
