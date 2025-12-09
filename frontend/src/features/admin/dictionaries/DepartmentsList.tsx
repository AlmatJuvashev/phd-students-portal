import React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { listDepartments, createDepartment, updateDepartment, deleteDepartment, Department } from "./api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Plus, Pencil, Trash2, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/use-toast";
import { useTranslation } from "react-i18next";

export const DepartmentsList = () => {
  const { t } = useTranslation("common");
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = React.useState(false);
  const [editingDepartment, setEditingDepartment] = React.useState<Department | null>(null);

  const { data: departments, isLoading } = useQuery({
    queryKey: ["departments"],
    queryFn: () => listDepartments(false),
  });

  const createMutation = useMutation({
    mutationFn: createDepartment,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["departments"] });
      setIsCreateOpen(false);
      toast({
        title: t("admin.dictionaries.departments.created", {
          defaultValue: "Department created",
        }),
      });
    },
    onError: (err: any) => {
      toast({
        title: t("common.error", { defaultValue: "Error" }),
        description:
          err.response?.data?.error ||
          t("admin.dictionaries.departments.create_error", {
            defaultValue: "Failed to create department",
          }),
        variant: "destructive",
      });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Department> }) =>
      updateDepartment(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["departments"] });
      setEditingDepartment(null);
      toast({
        title: t("admin.dictionaries.departments.updated", {
          defaultValue: "Department updated",
        }),
      });
    },
    onError: (err: any) => {
      toast({
        title: t("common.error", { defaultValue: "Error" }),
        description:
          err.response?.data?.error ||
          t("admin.dictionaries.departments.update_error", {
            defaultValue: "Failed to update department",
          }),
        variant: "destructive",
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteDepartment,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["departments"] });
      toast({
        title: t("admin.dictionaries.departments.deleted", {
          defaultValue: "Department deleted",
        }),
      });
    },
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, isActive }: { id: string; isActive: boolean }) =>
      updateDepartment(id, { is_active: isActive }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["departments"] });
      toast({
        title: t("admin.dictionaries.departments.status_updated", {
          defaultValue: "Department status updated",
        }),
      });
    },
    onError: (err: any) => {
      toast({
        title: t("common.error", { defaultValue: "Error" }),
        description:
          err.response?.data?.error ||
          t("admin.dictionaries.actions.status_error", {
            defaultValue: "Failed to update status",
          }),
        variant: "destructive",
      });
    },
  });

  const handleSubmitCreate = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    createMutation.mutate({
      name: formData.get("name") as string,
      code: formData.get("code") as string,
    });
  };

  const handleSubmitUpdate = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!editingDepartment) return;
    const formData = new FormData(e.currentTarget);
    updateMutation.mutate({
      id: editingDepartment.id,
      data: {
        name: formData.get("name") as string,
        code: formData.get("code") as string,
        is_active: formData.get("is_active") === "on",
      },
    });
  };

  if (isLoading)
    return (
      <div className="flex justify-center py-6">
        <Loader2 className="h-8 w-8 animate-spin" aria-hidden />
      </div>
    );

  const filteredDepartments = departments || [];

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">
          {t("admin.dictionaries.departments.title", { defaultValue: "Departments" })}
        </h2>
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />{" "}
              {t("admin.dictionaries.departments.add", { defaultValue: "Add Department" })}
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {t("admin.dictionaries.departments.add_title", {
                  defaultValue: "Add New Department",
                })}
              </DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmitCreate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">
                  {t("admin.dictionaries.fields.name", { defaultValue: "Name" })}
                </Label>
                <Input
                  id="name"
                  name="name"
                  placeholder={t("admin.dictionaries.departments.name_placeholder", {
                    defaultValue: "e.g. Department of Surgery",
                  })}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="code">
                  {t("admin.dictionaries.fields.code", { defaultValue: "Code" })}
                </Label>
                <Input
                  id="code"
                  name="code"
                  placeholder={t("admin.dictionaries.departments.code_placeholder", {
                    defaultValue: "e.g. SURG",
                  })}
                />
              </div>
              <DialogFooter>
                <Button type="submit" disabled={createMutation.isPending}>
                  {createMutation.isPending
                    ? t("admin.dictionaries.actions.creating", { defaultValue: "Creating..." })
                    : t("admin.dictionaries.actions.create", { defaultValue: "Create" })}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="border rounded-md">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t("table.name", { defaultValue: "Name" })}</TableHead>
              <TableHead>
                {t("admin.dictionaries.fields.code", { defaultValue: "Code" })}
              </TableHead>
              <TableHead>
                {t("admin.forms.active_state", { defaultValue: "Status" })}
              </TableHead>
              <TableHead className="text-right">
                {t("table.actions", { defaultValue: "Actions" })}
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredDepartments.map((dept: Department) => (
              <TableRow key={dept.id} className={!dept.is_active ? "opacity-50" : ""}>
                <TableCell className="font-medium">{dept.name}</TableCell>
                <TableCell>{dept.code || "—"}</TableCell>
                <TableCell>
                  <Switch
                    checked={dept.is_active}
                    onCheckedChange={(checked) =>
                      toggleActiveMutation.mutate({ id: dept.id, isActive: checked })
                    }
                    disabled={toggleActiveMutation.isPending}
                  />
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setEditingDepartment(dept)}
                      aria-label={t("common.edit", { defaultValue: "Edit" })}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => deleteMutation.mutate(dept.id)}
                      disabled={deleteMutation.isPending}
                      aria-label={t("common.remove", { defaultValue: "Remove" })}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
            {filteredDepartments.length === 0 && (
              <TableRow>
                <TableCell colSpan={4} className="text-center py-8 text-muted-foreground">
                  {t("admin.dictionaries.departments.empty", {
                    defaultValue: "No departments found.",
                  })}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog open={!!editingDepartment} onOpenChange={(open) => !open && setEditingDepartment(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {t("admin.dictionaries.departments.edit_title", { defaultValue: "Edit Department" })}
            </DialogTitle>
          </DialogHeader>
          {editingDepartment && (
            <form onSubmit={handleSubmitUpdate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="edit-name">
                  {t("admin.dictionaries.fields.name", { defaultValue: "Name" })}
                </Label>
                <Input id="edit-name" name="name" defaultValue={editingDepartment.name} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-code">
                  {t("admin.dictionaries.fields.code", { defaultValue: "Code" })}
                </Label>
                <Input id="edit-code" name="code" defaultValue={editingDepartment.code} />
              </div>
              <div className="flex items-center space-x-2">
                <Switch id="edit-active" name="is_active" defaultChecked={editingDepartment.is_active} />
                <Label htmlFor="edit-active">
                  {t("admin.forms.status_active", { defaultValue: "Active" })}
                </Label>
              </div>
              <DialogFooter>
                <Button type="submit" disabled={updateMutation.isPending}>
                  {updateMutation.isPending
                    ? t("admin.dictionaries.actions.saving", { defaultValue: "Saving..." })
                    : t("admin.dictionaries.actions.save_changes", {
                        defaultValue: "Save Changes",
                      })}
                </Button>
              </DialogFooter>
            </form>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};
