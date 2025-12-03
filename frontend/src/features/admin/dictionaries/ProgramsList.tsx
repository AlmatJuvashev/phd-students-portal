import React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { listPrograms, createProgram, updateProgram, deleteProgram, Program } from "./api";
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

export const ProgramsList = () => {
  const { t } = useTranslation("common");
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = React.useState(false);
  const [editingProgram, setEditingProgram] = React.useState<Program | null>(null);

  const { data: programs, isLoading } = useQuery({
    queryKey: ["programs"],
    queryFn: () => listPrograms(false), // Fetch all, including inactive
  });

  const createMutation = useMutation({
    mutationFn: createProgram,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["programs"] });
      setIsCreateOpen(false);
      toast({
        title: t("admin.dictionaries.programs.created", {
          defaultValue: "Program created",
        }),
      });
    },
    onError: (err: any) => {
      toast({
        title: t("common.error", { defaultValue: "Error" }),
        description:
          err.response?.data?.error ||
          t("admin.dictionaries.programs.create_error", {
            defaultValue: "Failed to create program",
          }),
        variant: "destructive",
      });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Program> }) => updateProgram(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["programs"] });
      setEditingProgram(null);
      toast({
        title: t("admin.dictionaries.programs.updated", {
          defaultValue: "Program updated",
        }),
      });
    },
    onError: (err: any) => {
      toast({
        title: t("common.error", { defaultValue: "Error" }),
        description:
          err.response?.data?.error ||
          t("admin.dictionaries.programs.update_error", {
            defaultValue: "Failed to update program",
          }),
        variant: "destructive",
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteProgram,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["programs"] });
      toast({
        title: t("admin.dictionaries.programs.deleted", {
          defaultValue: "Program deleted",
        }),
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
    if (!editingProgram) return;
    const formData = new FormData(e.currentTarget);
    updateMutation.mutate({
      id: editingProgram.id,
      data: {
        name: formData.get("name") as string,
        code: formData.get("code") as string,
        is_active: formData.get("is_active") === "on",
      },
    });
  };

  if (isLoading) return <Loader2 className="h-8 w-8 animate-spin mx-auto" />;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">
          {t("admin.dictionaries.programs.title", { defaultValue: "Programs" })}
        </h2>
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />{" "}
              {t("admin.dictionaries.programs.add", { defaultValue: "Add Program" })}
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {t("admin.dictionaries.programs.add_title", {
                  defaultValue: "Add New Program",
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
                  required
                  placeholder={t("admin.dictionaries.programs.name_placeholder", {
                    defaultValue: "Program name",
                  })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="code">
                  {t("admin.dictionaries.fields.code_optional", {
                    defaultValue: "Code (optional)",
                  })}
                </Label>
                <Input id="code" name="code" />
              </div>
              <DialogFooter>
                <Button type="submit" disabled={createMutation.isPending}>
                  {createMutation.isPending
                    ? t("admin.dictionaries.actions.creating", {
                        defaultValue: "Creating...",
                      })
                    : t("admin.dictionaries.actions.create", {
                        defaultValue: "Create",
                      })}
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
              <TableHead>
                {t("table.name", { defaultValue: "Name" })}
              </TableHead>
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
            {programs?.map((program: Program) => (
              <TableRow key={program.id} className={!program.is_active ? "opacity-50" : ""}>
                <TableCell className="font-medium">{program.name}</TableCell>
                <TableCell>{program.code}</TableCell>
                <TableCell>
                  {program.is_active ? (
                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      {t("admin.forms.status_active", { defaultValue: "Active" })}
                    </span>
                  ) : (
                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                      {t("admin.forms.status_inactive", { defaultValue: "Inactive" })}
                    </span>
                  )}
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setEditingProgram(program)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    {program.is_active && (
                      <Button
                        variant="ghost"
                        size="icon"
                        className="text-destructive hover:text-destructive"
                        onClick={() => {
                          if (
                            confirm(
                              t("admin.dictionaries.programs.delete_confirm", {
                                defaultValue: "Are you sure you want to delete this program?",
                              })
                            )
                          ) {
                            deleteMutation.mutate(program.id);
                          }
                        }}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                </TableCell>
              </TableRow>
            ))}
            {programs?.length === 0 && (
              <TableRow>
                <TableCell colSpan={4} className="text-center py-8 text-muted-foreground">
                  {t("admin.dictionaries.programs.empty", {
                    defaultValue: "No programs found.",
                  })}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog open={!!editingProgram} onOpenChange={(open: boolean) => !open && setEditingProgram(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {t("admin.dictionaries.programs.edit_title", {
                defaultValue: "Edit Program",
              })}
            </DialogTitle>
          </DialogHeader>
          {editingProgram && (
            <form onSubmit={handleSubmitUpdate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="edit-name">
                  {t("admin.dictionaries.fields.name", { defaultValue: "Name" })}
                </Label>
                <Input id="edit-name" name="name" defaultValue={editingProgram.name} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-code">
                  {t("admin.dictionaries.fields.code", { defaultValue: "Code" })}
                </Label>
                <Input id="edit-code" name="code" defaultValue={editingProgram.code} />
              </div>
              <div className="flex items-center space-x-2">
                <Switch id="edit-active" name="is_active" defaultChecked={editingProgram.is_active} />
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
