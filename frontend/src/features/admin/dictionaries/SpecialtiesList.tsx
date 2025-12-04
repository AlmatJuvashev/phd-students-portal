import React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  listSpecialties,
  createSpecialty,
  updateSpecialty,
  deleteSpecialty,
  listPrograms,
  Specialty,
} from "./api";
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
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Plus, Pencil, Trash2, Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/use-toast";
import { useTranslation } from "react-i18next";

export const SpecialtiesList = () => {
  const { t } = useTranslation("common");
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = React.useState(false);
  const [editingSpecialty, setEditingSpecialty] =
    React.useState<Specialty | null>(null);
  const [selectedPrograms, setSelectedPrograms] = React.useState<string[]>([]);

  const { data: specialties, isLoading } = useQuery({
    queryKey: ["specialties"],
    queryFn: () => listSpecialties(false),
  });

  const { data: programs, refetch: refetchPrograms } = useQuery({
    queryKey: ["programs", "all"],
    queryFn: () => listPrograms(false), // Fetch all programs, not just active
    staleTime: 0, // Always refetch to get latest programs
  });

  // Refetch programs when dialogs open
  React.useEffect(() => {
    if (isCreateOpen || editingSpecialty) {
      refetchPrograms();
    }
  }, [isCreateOpen, editingSpecialty, refetchPrograms]);

  const createMutation = useMutation({
    mutationFn: createSpecialty,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      setIsCreateOpen(false);
      toast({
        title: t("admin.dictionaries.specialties.created", {
          defaultValue: "Specialty created",
        }),
      });
    },
    onError: (err: any) => {
      toast({
        title: t("common.error", { defaultValue: "Error" }),
        description:
          err.response?.data?.error ||
          t("admin.dictionaries.specialties.create_error", {
            defaultValue: "Failed to create specialty",
          }),
        variant: "destructive",
      });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Specialty> }) =>
      updateSpecialty(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      setEditingSpecialty(null);
      toast({
        title: t("admin.dictionaries.specialties.updated", {
          defaultValue: "Specialty updated",
        }),
      });
    },
    onError: (err: any) => {
      toast({
        title: t("common.error", { defaultValue: "Error" }),
        description:
          err.response?.data?.error ||
          t("admin.dictionaries.specialties.update_error", {
            defaultValue: "Failed to update specialty",
          }),
        variant: "destructive",
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteSpecialty,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      toast({
        title: t("admin.dictionaries.specialties.deleted", {
          defaultValue: "Specialty deleted",
        }),
      });
    },
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, isActive }: { id: string; isActive: boolean }) =>
      updateSpecialty(id, { is_active: isActive }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      toast({
        title: t("admin.dictionaries.specialties.status_updated", {
          defaultValue: "Specialty status updated",
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
      program_ids: selectedPrograms,
    });
    setSelectedPrograms([]);
  };

  const handleSubmitUpdate = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!editingSpecialty) return;
    const formData = new FormData(e.currentTarget);
    updateMutation.mutate({
      id: editingSpecialty.id,
      data: {
        name: formData.get("name") as string,
        code: formData.get("code") as string,
        program_ids: selectedPrograms,
        is_active: formData.get("is_active") === "on",
      },
    });
    setSelectedPrograms([]);
  };

  const getProgramNames = (ids: string[]) => {
    if (!ids || ids.length === 0)
      return t("admin.dictionaries.shared.none", { defaultValue: "None" });
    return ids
      .map((id) => programs?.find((p) => p.id === id)?.name || id)
      .join(", ");
  };

  // Filter specialties (currently showing all, can add filters later)
  const filteredSpecialties = specialties || [];

  if (isLoading)
    return (
      <div className="flex justify-center py-6">
        <Loader2 className="h-8 w-8 animate-spin" aria-hidden />
      </div>
    );

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">
          {t("admin.dictionaries.specialties.title", {
            defaultValue: "Specialties",
          })}
        </h2>
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />{" "}
              {t("admin.dictionaries.specialties.add", {
                defaultValue: "Add Specialty",
              })}
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {t("admin.dictionaries.specialties.add_title", {
                  defaultValue: "Add New Specialty",
                })}
              </DialogTitle>
              <DialogDescription className="sr-only">
                {t("admin.dictionaries.specialties.add_title", {
                  defaultValue: "Add New Specialty",
                })}
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmitCreate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">
                  {t("admin.dictionaries.fields.name", {
                    defaultValue: "Name",
                  })}
                </Label>
                <Input
                  id="name"
                  name="name"
                  required
                  placeholder={t(
                    "admin.dictionaries.specialties.name_placeholder",
                    {
                      defaultValue: "Specialty name",
                    }
                  )}
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
              <div className="space-y-2">
                <Label>
                  {t("admin.dictionaries.fields.programs", {
                    defaultValue: "Programs",
                  })}
                </Label>
                <div className="border rounded p-3 space-y-2 max-h-40 overflow-auto">
                  {programs && programs.length > 0 ? (
                    programs.map((program) => (
                      <div
                        key={program.id}
                        className="flex items-center space-x-2"
                      >
                        <input
                          type="checkbox"
                          id={`create-prog-${program.id}`}
                          checked={selectedPrograms.includes(program.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedPrograms([
                                ...selectedPrograms,
                                program.id,
                              ]);
                            } else {
                              setSelectedPrograms(
                                selectedPrograms.filter(
                                  (id) => id !== program.id
                                )
                              );
                            }
                          }}
                          className="h-4 w-4"
                        />
                        <label
                          htmlFor={`create-prog-${program.id}`}
                          className="text-sm cursor-pointer"
                        >
                          {program.name}
                        </label>
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground">
                      {t("admin.dictionaries.specialties.no_programs", {
                        defaultValue: "No active programs available",
                      })}
                    </p>
                  )}
                </div>
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
              <TableHead>{t("table.name", { defaultValue: "Name" })}</TableHead>
              <TableHead>
                {t("admin.dictionaries.fields.code", { defaultValue: "Code" })}
              </TableHead>
              <TableHead>
                {t("admin.forms.program", { defaultValue: "Program" })}
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
            {filteredSpecialties.map((specialty: Specialty) => (
              <TableRow
                key={specialty.id}
                className={!specialty.is_active ? "opacity-50" : ""}
              >
                <TableCell className="font-medium">{specialty.name}</TableCell>
                <TableCell>{specialty.code || "—"}</TableCell>
                <TableCell>{getProgramNames(specialty.program_ids)}</TableCell>
                <TableCell>
                  <Switch
                    checked={specialty.is_active}
                    onCheckedChange={(checked) =>
                      toggleActiveMutation.mutate({
                        id: specialty.id,
                        isActive: checked,
                      })
                    }
                    disabled={toggleActiveMutation.isPending}
                  />
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        setEditingSpecialty(specialty);
                        setSelectedPrograms(specialty.program_ids || []);
                      }}
                      aria-label={t("common.edit", { defaultValue: "Edit" })}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => deleteMutation.mutate(specialty.id)}
                      disabled={deleteMutation.isPending}
                      aria-label={t("common.remove", {
                        defaultValue: "Remove",
                      })}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
            {specialties?.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="text-center py-8 text-muted-foreground"
                >
                  {t("admin.dictionaries.specialties.empty", {
                    defaultValue: "No specialties found.",
                  })}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog
        open={!!editingSpecialty}
        onOpenChange={(open: boolean) => {
          if (!open) {
            setEditingSpecialty(null);
            setSelectedPrograms([]);
          } else if (editingSpecialty) {
            setSelectedPrograms(editingSpecialty.program_ids || []);
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {t("admin.dictionaries.specialties.edit_title", {
                defaultValue: "Edit Specialty",
              })}
            </DialogTitle>
          </DialogHeader>
          {editingSpecialty && (
            <form onSubmit={handleSubmitUpdate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="edit-name">
                  {t("admin.dictionaries.fields.name", {
                    defaultValue: "Name",
                  })}
                </Label>
                <Input
                  id="edit-name"
                  name="name"
                  defaultValue={editingSpecialty.name}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-code">
                  {t("admin.dictionaries.fields.code", {
                    defaultValue: "Code",
                  })}
                </Label>
                <Input
                  id="edit-code"
                  name="code"
                  defaultValue={editingSpecialty.code}
                />
              </div>
              <div className="space-y-2">
                <Label>
                  {t("admin.dictionaries.fields.programs", {
                    defaultValue: "Programs",
                  })}
                </Label>
                <div className="border rounded p-3 space-y-2 max-h-40 overflow-auto">
                  {programs && programs.length > 0 ? (
                    programs.map((program) => (
                      <div
                        key={program.id}
                        className="flex items-center space-x-2"
                      >
                        <input
                          type="checkbox"
                          id={`edit-prog-${program.id}`}
                          checked={selectedPrograms.includes(program.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedPrograms([
                                ...selectedPrograms,
                                program.id,
                              ]);
                            } else {
                              setSelectedPrograms(
                                selectedPrograms.filter(
                                  (id) => id !== program.id
                                )
                              );
                            }
                          }}
                          className="h-4 w-4"
                        />
                        <label
                          htmlFor={`edit-prog-${program.id}`}
                          className="text-sm cursor-pointer"
                        >
                          {program.name}
                        </label>
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground">
                      {t("admin.dictionaries.specialties.no_programs", {
                        defaultValue: "No active programs available",
                      })}
                    </p>
                  )}
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <Switch
                  id="edit-active"
                  name="is_active"
                  defaultChecked={editingSpecialty.is_active}
                />
                <Label htmlFor="edit-active">
                  {t("admin.forms.status_active", { defaultValue: "Active" })}
                </Label>
              </div>
              <DialogFooter>
                <Button type="submit" disabled={updateMutation.isPending}>
                  {updateMutation.isPending
                    ? t("admin.dictionaries.actions.saving", {
                        defaultValue: "Saving...",
                      })
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
