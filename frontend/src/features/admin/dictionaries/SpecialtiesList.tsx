import React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { listSpecialties, createSpecialty, updateSpecialty, deleteSpecialty, listPrograms, Specialty } from "./api";
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

export const SpecialtiesList = () => {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = React.useState(false);
  const [editingSpecialty, setEditingSpecialty] = React.useState<Specialty | null>(null);
  const [selectedPrograms, setSelectedPrograms] = React.useState<string[]>([]);

  const { data: specialties, isLoading } = useQuery({
    queryKey: ["specialties"],
    queryFn: () => listSpecialties(false),
  });

  const { data: programs } = useQuery({
    queryKey: ["programs", "all"],
    queryFn: () => listPrograms(false), // Fetch all programs, not just active
  });

  const createMutation = useMutation({
    mutationFn: createSpecialty,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      setIsCreateOpen(false);
      toast({ title: "Specialty created" });
    },
    onError: (err: any) => {
      toast({ title: "Error", description: err.response?.data?.error || "Failed to create specialty", variant: "destructive" });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Specialty> }) => updateSpecialty(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      setEditingSpecialty(null);
      toast({ title: "Specialty updated" });
    },
    onError: (err: any) => {
      toast({ title: "Error", description: err.response?.data?.error || "Failed to update specialty", variant: "destructive" });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteSpecialty,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      toast({ title: "Specialty deleted (soft)" });
    },
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, isActive }: { id: string; isActive: boolean }) =>
      updateSpecialty(id, { is_active: isActive }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["specialties"] });
      toast({ title: "Specialty status updated" });
    },
    onError: (err: any) => {
      toast({
        title: "Error",
        description: err.response?.data?.error || "Failed to update status",
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
    if (!ids || ids.length === 0) return "None";
    return ids.map(id => programs?.find(p => p.id === id)?.name || id).join(", ");
  };

  // Filter specialties (currently showing all, can add filters later)
  const filteredSpecialties = specialties || [];

  if (isLoading) return <Loader2 className="h-8 w-8 animate-spin mx-auto" />;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">Specialties</h2>
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" /> Add Specialty
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add New Specialty</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmitCreate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Name</Label>
                <Input id="name" name="name" required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="code">Code (Optional)</Label>
                <Input id="code" name="code" />
              </div>
              <div className="space-y-2">
                <Label>Programs</Label>
                <div className="border rounded p-3 space-y-2 max-h-40 overflow-auto">
                  {programs && programs.length > 0 ? (
                    programs.map((program) => (
                      <div key={program.id} className="flex items-center space-x-2">
                        <input
                          type="checkbox"
                          id={`create-prog-${program.id}`}
                          checked={selectedPrograms.includes(program.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedPrograms([...selectedPrograms, program.id]);
                            } else {
                              setSelectedPrograms(selectedPrograms.filter(id => id !== program.id));
                            }
                          }}
                          className="h-4 w-4"
                        />
                        <label htmlFor={`create-prog-${program.id}`} className="text-sm cursor-pointer">
                          {program.name}
                        </label>
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground">No active programs available</p>
                  )}
                </div>
              </div>
              <DialogFooter>
                <Button type="submit" disabled={createMutation.isPending}>
                  {createMutation.isPending ? "Creating..." : "Create"}
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
              <TableHead>Name</TableHead>
              <TableHead>Code</TableHead>
              <TableHead>Program</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredSpecialties.map((specialty: Specialty) => (
              <TableRow key={specialty.id} className={!specialty.is_active ? "opacity-50" : ""}>
                <TableCell className="font-medium">{specialty.name}</TableCell>
                <TableCell>{specialty.code || "—"}</TableCell>
                <TableCell>
                  {getProgramNames(specialty.program_ids)}
                </TableCell>
                <TableCell>
                  <Switch
                    checked={specialty.is_active}
                    onCheckedChange={(checked) =>
                      toggleActiveMutation.mutate({ id: specialty.id, isActive: checked })
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
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => deleteMutation.mutate(specialty.id)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
            {specialties?.length === 0 && (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                  No specialties found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog open={!!editingSpecialty} onOpenChange={(open: boolean) => {
        if (!open) {
          setEditingSpecialty(null);
          setSelectedPrograms([]);
        } else if (editingSpecialty) {
          setSelectedPrograms(editingSpecialty.program_ids || []);
        }
      }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Specialty</DialogTitle>
          </DialogHeader>
          {editingSpecialty && (
            <form onSubmit={handleSubmitUpdate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="edit-name">Name</Label>
                <Input id="edit-name" name="name" defaultValue={editingSpecialty.name} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-code">Code</Label>
                <Input id="edit-code" name="code" defaultValue={editingSpecialty.code} />
              </div>
              <div className="space-y-2">
                <Label>Programs</Label>
                <div className="border rounded p-3 space-y-2 max-h-40 overflow-auto">
                  {programs && programs.length > 0 ? (
                    programs.map((program) => (
                      <div key={program.id} className="flex items-center space-x-2">
                        <input
                          type="checkbox"
                          id={`edit-prog-${program.id}`}
                          checked={selectedPrograms.includes(program.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedPrograms([...selectedPrograms, program.id]);
                            } else {
                              setSelectedPrograms(selectedPrograms.filter(id => id !== program.id));
                            }
                          }}
                          className="h-4 w-4"
                        />
                        <label htmlFor={`edit-prog-${program.id}`} className="text-sm cursor-pointer">
                          {program.name}
                        </label>
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground">No active programs available</p>
                  )}
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <Switch id="edit-active" name="is_active" defaultChecked={editingSpecialty.is_active} />
                <Label htmlFor="edit-active">Active</Label>
              </div>
              <DialogFooter>
                <Button type="submit" disabled={updateMutation.isPending}>
                  {updateMutation.isPending ? "Saving..." : "Save Changes"}
                </Button>
              </DialogFooter>
            </form>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};
