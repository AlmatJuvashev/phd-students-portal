import React from "react";
import { useForm } from "react-hook-form";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { api } from "../api/client";
import { Input } from "../components/ui/input";
import { Button } from "../components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "../components/ui/card";
import { Badge } from "../components/ui/badge";
import {
  Copy,
  Plus,
  X,
  Search,
  ChevronUp,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  Edit,
  RefreshCw,
} from "lucide-react";
import { useTranslation } from "react-i18next";

const Schema = z.object({
  first_name: z.string().min(1),
  last_name: z.string().min(1),
  email: z.string().email(),
  role: z.enum(["student", "advisor", "chair", "admin"]),
});

const EditSchema = z.object({
  first_name: z.string().min(1),
  last_name: z.string().min(1),
  email: z.string().email(),
  role: z.enum(["student", "advisor", "chair", "admin"]),
});

type Form = z.infer<typeof Schema>;
type EditForm = z.infer<typeof EditSchema>;

type User = {
  id: string;
  name: string;
  email: string;
  role: string;
};

export function AdminUsers() {
  const [showModal, setShowModal] = React.useState(false);
  const [showEditModal, setShowEditModal] = React.useState(false);
  const [editingUser, setEditingUser] = React.useState<User | null>(null);
  const [created, setCreated] = React.useState<{
    username: string;
    temp_password: string;
  } | null>(null);
  const [resetPassword, setResetPassword] = React.useState<{
    username: string;
    temp_password: string;
  } | null>(null);
  const [searchQuery, setSearchQuery] = React.useState("");
  const [currentPage, setCurrentPage] = React.useState(1);
  const [sortField, setSortField] = React.useState<keyof User>("name");
  const [sortDirection, setSortDirection] = React.useState<"asc" | "desc">(
    "asc"
  );
  const { t } = useTranslation("common");
  const itemsPerPage = 10;
  const push = (msg: { title: string; description?: string }) => {
    // toasts removed; fallback to console + optional alert on important events
    console.log(`[${msg.title}]`, msg.description || "");
  };
  const queryClient = useQueryClient();

  // Fetch users list - filter out superadmins
  const { data: allUsers = [], isLoading } = useQuery<User[]>({
    queryKey: ["admin", "users"],
    queryFn: () => api("/admin/users"),
  });

  // Filter out superadmins from the list
  const users = React.useMemo(() => {
    return allUsers.filter((user) => user.role !== "superadmin");
  }, [allUsers]);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<Form>({
    resolver: zodResolver(Schema),
    defaultValues: { role: "student" },
  });

  const {
    register: registerEdit,
    handleSubmit: handleSubmitEdit,
    reset: resetEdit,
    setValue: setValueEdit,
    formState: { errors: editErrors, isSubmitting: isEditSubmitting },
  } = useForm<EditForm>({
    resolver: zodResolver(EditSchema),
  });

  // Create user mutation
  const createUserMutation = useMutation({
    mutationFn: (data: Form) =>
      api("/admin/users", { method: "POST", body: JSON.stringify(data) }),
    onSuccess: (result) => {
      setCreated(result);
      reset();
      setShowModal(false);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      push({
        title: t("admin.users.toast.created", { defaultValue: "User created" }),
        description: t("admin.users.toast.credentials_ready", {
          defaultValue: "Credentials generated successfully",
        }),
      });
    },
    onError: (error: any) => {
      console.error(error);
      alert(
        error?.message ||
          t("admin.forms.errors.create_user", {
            defaultValue: "Error creating user",
          })
      );
    },
  });

  // Edit user mutation
  const editUserMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: EditForm }) =>
      api(`/admin/users/${id}`, { method: "PUT", body: JSON.stringify(data) }),
    onSuccess: () => {
      resetEdit();
      setShowEditModal(false);
      setEditingUser(null);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      push({
        title: t("admin.users.toast.updated", { defaultValue: "User updated" }),
      });
    },
    onError: (error: any) => {
      console.error(error);
      alert(
        error?.message ||
          t("admin.forms.errors.update_user", {
            defaultValue: "Error updating user",
          })
      );
    },
  });

  // Reset password mutation
  const resetPasswordMutation = useMutation({
    mutationFn: (userId: string) =>
      api(`/admin/users/${userId}/reset-password`, { method: "POST" }),
    onSuccess: (result) => {
      setResetPassword(result);
      push({
        title: t("admin.users.toast.reset", { defaultValue: "Password reset" }),
      });
    },
    onError: (error: any) => {
      console.error(error);
      alert(
        error?.message ||
          t("admin.forms.errors.reset_password", {
            defaultValue: "Error resetting password",
          })
      );
    },
  });

  const onSubmit = (data: Form) => {
    createUserMutation.mutate(data);
  };

  const onEditSubmit = (data: EditForm) => {
    if (editingUser) {
      editUserMutation.mutate({ id: editingUser.id, data });
    }
  };

  const handleEditUser = (user: User) => {
    setEditingUser(user);
    const [firstName, ...lastNameParts] = user.name.split(" ");
    setValueEdit("first_name", firstName);
    setValueEdit("last_name", lastNameParts.join(" "));
    setValueEdit("email", user.email);
    setValueEdit("role", user.role as any);
    setShowEditModal(true);
  };

  const handleResetPassword = (userId: string) => {
    if (
      confirm(
        t("admin.users.confirm_reset", {
          defaultValue:
            "Are you sure you want to reset this user's password? They will need to use the new temporary password to login.",
        })
      )
    ) {
      resetPasswordMutation.mutate(userId);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    push({ title: t("common.copied", { defaultValue: "Copied" }) });
  };

  const copyUserCredentials = (username: string, password: string) => {
    const loginUrl = window.location.origin + "/login";
    const message = `Your password has been reset. You can login at ${loginUrl}.\n\nUsername: ${username}\nNew Password: ${password}\n\nPlease save these credentials securely and change your password after login.`;
    navigator.clipboard.writeText(message);
    push({
      title: t("admin.users.toast.credentials_copied", {
        defaultValue: "Credentials copied",
      }),
    });
  };

  const copyNewUserCredentials = (username: string, password: string) => {
    const loginUrl = window.location.origin + "/login";
    const message = `Your account has been created. You can login at ${loginUrl}.\n\nUsername: ${username}\nPassword: ${password}\n\nPlease save these credentials securely and change your password after first login.`;
    navigator.clipboard.writeText(message);
    push({
      title: t("admin.users.toast.credentials_copied", {
        defaultValue: "Credentials copied",
      }),
    });
  };

  const handleSort = (field: keyof User) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
    setCurrentPage(1); // Reset to first page when sorting
  };

  // Filter and sort users
  const filteredAndSortedUsers = React.useMemo(() => {
    let filtered = users.filter(
      (user) =>
        user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.role.toLowerCase().includes(searchQuery.toLowerCase())
    );

    filtered.sort((a, b) => {
      const aVal = a[sortField];
      const bVal = b[sortField];
      if (aVal < bVal) return sortDirection === "asc" ? -1 : 1;
      if (aVal > bVal) return sortDirection === "asc" ? 1 : -1;
      return 0;
    });

    return filtered;
  }, [users, searchQuery, sortField, sortDirection]);

  // Pagination
  const totalPages = Math.ceil(filteredAndSortedUsers.length / itemsPerPage);
  const paginatedUsers = filteredAndSortedUsers.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setCurrentPage(1); // Reset to first page when searching
  };

  const SortIcon = ({ field }: { field: keyof User }) => {
    if (sortField !== field) return null;
    return sortDirection === "asc" ? (
      <ChevronUp className="w-4 h-4" />
    ) : (
      <ChevronDown className="w-4 h-4" />
    );
  };

  const getRoleBadgeColor = (role: string) => {
    const colors = {
      superadmin: "bg-purple-100 text-purple-800",
      admin: "bg-red-100 text-red-800",
      chair: "bg-orange-100 text-orange-800",
      advisor: "bg-blue-100 text-blue-800",
      student: "bg-green-100 text-green-800",
    };
    return colors[role as keyof typeof colors] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="max-w-6xl mx-auto mt-8 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">
            {t("admin.users.title", { defaultValue: "User Management" })}
          </h2>
          <p className="text-gray-600">
            {t("admin.users.subtitle", {
              defaultValue: "Manage users and their access permissions",
            })}
          </p>
        </div>
        <Button
          onClick={() => setShowModal(true)}
          className="flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          {t("admin.users.create", { defaultValue: "Create User" })}
        </Button>
      </div>

      {/* Users Table */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>
              {t("admin.users.table_title", { defaultValue: "Users" })} (
              {filteredAndSortedUsers.length})
            </CardTitle>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <Input
                placeholder={t("admin.users.search", {
                  defaultValue: "Search users...",
                })}
                value={searchQuery}
                onChange={(e) => handleSearch(e.target.value)}
                className="pl-10 w-64"
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">
              {t("admin.users.loading", { defaultValue: "Loading users..." })}
            </div>
          ) : paginatedUsers.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              {searchQuery
                ? t("admin.users.empty_search", {
                    defaultValue: 'No users found matching "{{query}}"',
                    query: searchQuery,
                  })
                : t("admin.users.empty", { defaultValue: "No users found" })}
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b">
                      <th
                        className="text-left py-3 px-4 font-medium text-gray-600 cursor-pointer hover:bg-gray-50"
                    onClick={() => handleSort("name")}
                  >
                    <div className="flex items-center gap-2">
                      {t("table.name", { defaultValue: "Name" })}
                      <SortIcon field="name" />
                    </div>
                  </th>
                  <th
                    className="text-left py-3 px-4 font-medium text-gray-600 cursor-pointer hover:bg-gray-50"
                    onClick={() => handleSort("email")}
                  >
                    <div className="flex items-center gap-2">
                      {t("table.email", { defaultValue: "Email" })}
                      <SortIcon field="email" />
                    </div>
                  </th>
                  <th
                    className="text-left py-3 px-4 font-medium text-gray-600 cursor-pointer hover:bg-gray-50"
                    onClick={() => handleSort("role")}
                  >
                    <div className="flex items-center gap-2">
                      {t("admin.users.role", { defaultValue: "Role" })}
                      <SortIcon field="role" />
                    </div>
                  </th>
                  <th className="text-left py-3 px-4 font-medium text-gray-600">
                    {t("table.actions", { defaultValue: "Actions" })}
                  </th>
                </tr>
              </thead>
              <tbody>
                {paginatedUsers.map((user) => (
                  <tr key={user.id} className="border-b hover:bg-gray-50">
                    <td className="py-3 px-4 font-medium">{user.name}</td>
                    <td className="py-3 px-4 text-gray-600">
                      {user.email}
                    </td>
                    <td className="py-3 px-4">
                      <Badge className={getRoleBadgeColor(user.role)}>
                        {t(`admin.forms.roles.${user.role}`, {
                          defaultValue: user.role,
                        })}
                      </Badge>
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleEditUser(user)}
                              className="flex items-center gap-1"
                            >
                              <Edit className="w-3 h-3" />
                              {t("common.edit", { defaultValue: "Edit" })}
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleResetPassword(user.id)}
                              className="flex items-center gap-1"
                            >
                              <RefreshCw className="w-3 h-3" />
                              {t("admin.review.reset_password", {
                                defaultValue: "Reset password",
                              })}
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => copyToClipboard(user.email)}
                              className="flex items-center gap-1"
                            >
                              <Copy className="w-3 h-3" />
                              {t("common.copy", { defaultValue: "Copy" })}
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between mt-4 pt-4 border-t">
                  <div className="text-sm text-gray-500">
                    {t("admin.users.pagination", {
                      defaultValue: "Showing {{from}} to {{to}} of {{total}} users",
                      from: (currentPage - 1) * itemsPerPage + 1,
                      to: Math.min(
                        currentPage * itemsPerPage,
                        filteredAndSortedUsers.length
                      ),
                      total: filteredAndSortedUsers.length,
                    })}
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        setCurrentPage((prev) => Math.max(1, prev - 1))
                      }
                      disabled={currentPage === 1}
                      className="flex items-center gap-1"
                    >
                      <ChevronLeft className="w-4 h-4" />
                      {t("admin.forms.prev_page", { defaultValue: "Prev" })}
                    </Button>

                    <div className="flex items-center gap-1">
                      {Array.from({ length: totalPages }, (_, i) => i + 1)
                        .filter(
                          (page) =>
                            page === 1 ||
                            page === totalPages ||
                            Math.abs(page - currentPage) <= 2
                        )
                        .map((page, index, array) => (
                          <React.Fragment key={page}>
                            {index > 0 && array[index - 1] !== page - 1 && (
                              <span className="px-2 text-gray-400">...</span>
                            )}
                            <Button
                              variant={
                                currentPage === page ? "default" : "outline"
                              }
                              size="sm"
                              onClick={() => setCurrentPage(page)}
                              className="min-w-[32px]"
                            >
                              {page}
                            </Button>
                          </React.Fragment>
                        ))}
                    </div>

                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        setCurrentPage((prev) => Math.min(totalPages, prev + 1))
                      }
                      disabled={currentPage === totalPages}
                      className="flex items-center gap-1"
                    >
                      {t("admin.forms.next_page", { defaultValue: "Next" })}
                      <ChevronRight className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      {/* Edit User Modal */}
      {showEditModal && editingUser && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">
                {t("admin.users.edit_title", { defaultValue: "Edit User" })}
              </h3>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setShowEditModal(false);
                  setEditingUser(null);
                  resetEdit();
                }}
                className="p-1"
              >
                <X className="w-4 h-4" />
              </Button>
            </div>

            <form
              className="space-y-4"
              onSubmit={handleSubmitEdit(onEditSubmit)}
            >
              <div>
                <Input
                  placeholder={t("admin.forms.first_name", {
                    defaultValue: "First name",
                  })}
                  {...registerEdit("first_name")}
                />
                {editErrors.first_name && (
                  <div className="text-xs text-red-600 mt-1">
                    {editErrors.first_name.message}
                  </div>
                )}
              </div>

              <div>
                <Input
                  placeholder={t("admin.forms.last_name", {
                    defaultValue: "Last name",
                  })}
                  {...registerEdit("last_name")}
                />
                {editErrors.last_name && (
                  <div className="text-xs text-red-600 mt-1">
                    {editErrors.last_name.message}
                  </div>
                )}
              </div>

              <div>
                <Input
                  placeholder={t("admin.forms.email", {
                    defaultValue: "Email",
                  })}
                  type="email"
                  {...registerEdit("email")}
                />
                {editErrors.email && (
                  <div className="text-xs text-red-600 mt-1">
                    {editErrors.email.message}
                  </div>
                )}
              </div>

              <div>
                <select
                  className="w-full border border-gray-300 p-2 rounded-md"
                  {...registerEdit("role")}
                >
                  <option value="student">
                    {t("admin.forms.roles.student", { defaultValue: "Student" })}
                  </option>
                  <option value="advisor">
                    {t("admin.forms.roles.advisor", { defaultValue: "Advisor" })}
                  </option>
                  <option value="chair">
                    {t("admin.forms.roles.chair", {
                      defaultValue: "Department Chair",
                    })}
                  </option>
                  <option value="admin">
                    {t("admin.forms.roles.admin", { defaultValue: "Administrator" })}
                  </option>
                </select>
              </div>

              <div className="flex gap-2 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    setShowEditModal(false);
                    setEditingUser(null);
                    resetEdit();
                  }}
                  className="flex-1"
                >
                  {t("common.cancel", { defaultValue: "Cancel" })}
                </Button>
                <Button
                  type="submit"
                  disabled={isEditSubmitting || editUserMutation.isPending}
                  className="flex-1"
                >
                  {editUserMutation.isPending
                    ? t("admin.users.updating", { defaultValue: "Updating..." })
                    : t("admin.users.update_action", { defaultValue: "Update User" })}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Create User Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">
                {t("admin.users.create_title", { defaultValue: "Create New User" })}
              </h3>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowModal(false)}
                className="p-1"
              >
                <X className="w-4 h-4" />
              </Button>
            </div>

            <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
              <div>
                <Input
                  placeholder={t("admin.forms.first_name", {
                    defaultValue: "First name",
                  })}
                  {...register("first_name")}
                />
                {errors.first_name && (
                  <div className="text-xs text-red-600 mt-1">
                    {errors.first_name.message}
                  </div>
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
                  <div className="text-xs text-red-600 mt-1">
                    {errors.last_name.message}
                  </div>
                )}
              </div>

              <div>
                <Input
                  placeholder={t("admin.forms.email", {
                    defaultValue: "Email",
                  })}
                  type="email"
                  {...register("email")}
                />
                {errors.email && (
                  <div className="text-xs text-red-600 mt-1">
                    {errors.email.message}
                  </div>
                )}
              </div>

              <div>
                <select
                  className="w-full border border-gray-300 p-2 rounded-md"
                  {...register("role")}
                >
                  <option value="student">
                    {t("admin.forms.roles.student", { defaultValue: "Student" })}
                  </option>
                  <option value="advisor">
                    {t("admin.forms.roles.advisor", { defaultValue: "Advisor" })}
                  </option>
                  <option value="chair">
                    {t("admin.forms.roles.chair", {
                      defaultValue: "Department Chair",
                    })}
                  </option>
                  <option value="admin">
                    {t("admin.forms.roles.admin", { defaultValue: "Administrator" })}
                  </option>
                </select>
              </div>

              <div className="flex gap-2 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setShowModal(false)}
                  className="flex-1"
                >
                  {t("common.cancel", { defaultValue: "Cancel" })}
                </Button>
                <Button
                  type="submit"
                  disabled={isSubmitting || createUserMutation.isPending}
                  className="flex-1"
                >
                  {createUserMutation.isPending
                    ? t("admin.forms.creating", { defaultValue: "Creating..." })
                    : t("admin.users.create", { defaultValue: "Create User" })}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Reset Password Display */}
      {resetPassword && (
        <Card className="border-orange-200 bg-orange-50">
          <CardHeader>
            <CardTitle className="text-orange-800 flex items-center gap-2">
              🔄{" "}
              {t("admin.users.reset_success", {
                defaultValue: "Password reset successfully",
              })}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-white rounded border">
                <div>
                  <div className="text-sm font-medium text-gray-600">
                    {t("admin.forms.username", { defaultValue: "Username" })}
                  </div>
                  <div className="font-mono text-sm">
                    {resetPassword.username}
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => copyToClipboard(resetPassword.username)}
                  className="flex items-center gap-1"
                >
                  <Copy className="w-3 h-3" />
                  {t("common.copy", { defaultValue: "Copy" })}
                </Button>
              </div>

              <div className="flex items-center justify-between p-3 bg-white rounded border">
                <div>
                  <div className="text-sm font-medium text-gray-600">
                    {t("admin.users.new_password", { defaultValue: "New password" })}
                  </div>
                  <div className="font-mono text-sm">
                    {resetPassword.temp_password}
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => copyToClipboard(resetPassword.temp_password)}
                  className="flex items-center gap-1"
                >
                  <Copy className="w-3 h-3" />
                  {t("common.copy", { defaultValue: "Copy" })}
                </Button>
              </div>

              <div className="flex gap-2">
                <Button
                  variant="outline"
                  onClick={() =>
                    copyUserCredentials(
                      resetPassword.username,
                      resetPassword.temp_password
                    )
                  }
                  className="flex items-center gap-1"
                >
                  <Copy className="w-3 h-3" />
                  {t("admin.users.copy_reset_message", {
                    defaultValue: "Copy reset message",
                  })}
                </Button>
                <Button
                  variant="ghost"
                  onClick={() => setResetPassword(null)}
                  className="flex items-center gap-1"
                >
                  <X className="w-3 h-3" />
                  {t("admin.users.dismiss", { defaultValue: "Dismiss" })}
                </Button>
              </div>

              <div className="text-xs text-orange-600 bg-orange-50 p-2 rounded border">
                ⚠️{" "}
                {t("admin.users.reset_notice", {
                  defaultValue:
                    "The user must use this new password to login. Make sure they receive these credentials securely!",
                })}
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Credentials Display */}
      {created && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader>
            <CardTitle className="text-green-800 flex items-center gap-2">
              ✅{" "}
              {t("admin.users.created_success", {
                defaultValue: "User created successfully",
              })}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-white rounded border">
                <div>
                  <div className="text-sm font-medium text-gray-600">
                    {t("admin.forms.username", { defaultValue: "Username" })}
                  </div>
                  <div className="font-mono text-sm">{created.username}</div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => copyToClipboard(created.username)}
                  className="flex items-center gap-1"
                >
                  <Copy className="w-3 h-3" />
                  {t("common.copy", { defaultValue: "Copy" })}
                </Button>
              </div>

              <div className="flex items-center justify-between p-3 bg-white rounded border">
                <div>
                  <div className="text-sm font-medium text-gray-600">
                    {t("admin.users.password", { defaultValue: "Password" })}
                  </div>
                  <div className="font-mono text-sm">
                    {created.temp_password}
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => copyToClipboard(created.temp_password)}
                  className="flex items-center gap-1"
                >
                  <Copy className="w-3 h-3" />
                  {t("common.copy", { defaultValue: "Copy" })}
                </Button>
              </div>

              <div className="flex gap-2">
                <Button
                  variant="outline"
                  onClick={() =>
                    copyNewUserCredentials(
                      created.username,
                      created.temp_password
                    )
                  }
                  className="flex items-center gap-1"
                >
                  <Copy className="w-3 h-3" />
                  {t("admin.users.copy_login_message", {
                    defaultValue: "Copy login message",
                  })}
                </Button>
                <Button
                  variant="outline"
                  onClick={() =>
                    copyToClipboard(
                      `Username: ${created.username}\nPassword: ${created.temp_password}`
                    )
                  }
                  className="flex items-center gap-1"
                >
                  <Copy className="w-3 h-3" />
                  {t("admin.users.copy_credentials", {
                    defaultValue: "Copy credentials",
                  })}
                </Button>
                <Button
                  variant="ghost"
                  onClick={() => setCreated(null)}
                  className="flex items-center gap-1"
                >
                  <X className="w-3 h-3" />
                  {t("admin.users.dismiss", { defaultValue: "Dismiss" })}
                </Button>
              </div>

              <div className="text-xs text-orange-600 bg-orange-50 p-2 rounded border">
                ⚠️{" "}
                {t("admin.users.credentials_notice", {
                  defaultValue:
                    "Make sure to save these credentials - they won't be shown again!",
                })}
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
