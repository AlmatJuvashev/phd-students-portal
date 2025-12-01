import React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { listCohorts, createCohort, updateCohort, deleteCohort, Cohort } from "./api";
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

export const CohortsList = () => {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = React.useState(false);
  const [editingCohort, setEditingCohort] = React.useState<Cohort | null>(null);

  const { data: cohorts, isLoading } = useQuery({
    queryKey: ["cohorts"],
    queryFn: () => listCohorts(false),
  });

  const createMutation = useMutation({
    mutationFn: createCohort,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cohorts"] });
      setIsCreateOpen(false);
      toast({ title: "Cohort created" });
    },
    onError: (err: any) => {
      toast({ title: "Error", description: err.response?.data?.error || "Failed to create cohort", variant: "destructive" });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Cohort> }) => updateCohort(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cohorts"] });
      setEditingCohort(null);
      toast({ title: "Cohort updated" });
    },
    onError: (err: any) => {
      toast({ title: "Error", description: err.response?.data?.error || "Failed to update cohort", variant: "destructive" });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteCohort,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cohorts"] });
      toast({ title: "Cohort deleted (soft)" });
    },
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, isActive }: { id: string; isActive: boolean }) =>
      updateCohort(id, { is_active: isActive }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cohorts"] });
      toast({ title: "Cohort status updated" });
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
      start_date: formData.get("start_date") as string,
      end_date: formData.get("end_date") as string,
    });
  };

  const handleSubmitUpdate = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!editingCohort) return;
    const formData = new FormData(e.currentTarget);
    updateMutation.mutate({
      id: editingCohort.id,
      data: {
        name: formData.get("name") as string,
        start_date: formData.get("start_date") as string,
        end_date: formData.get("end_date") as string,
        is_active: formData.get("is_active") === "on",
      },
    });
  };

  if (isLoading) return <Loader2 className="h-8 w-8 animate-spin mx-auto" />;

  const filteredCohorts = cohorts || [];

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">Cohorts</h2>
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2">
              <Plus className="h-4 w-4" /> Add Cohort
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add New Cohort</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmitCreate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Name</Label>
                <Input id="name" name="name" placeholder="e.g. 2024-2025" required />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="start_date">Start Date</Label>
                  <Input id="start_date" name="start_date" type="date" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="end_date">End Date</Label>
                  <Input id="end_date" name="end_date" type="date" />
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
              <TableHead>Start Date</TableHead>
              <TableHead>End Date</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredCohorts.map((cohort: Cohort) => (
              <TableRow key={cohort.id} className={!cohort.is_active ? "opacity-50" : ""}>
                <TableCell className="font-medium">{cohort.name}</TableCell>
                <TableCell>{cohort.start_date || "—"}</TableCell>
                <TableCell>{cohort.end_date || "—"}</TableCell>
                <TableCell>
                  <Switch
                    checked={cohort.is_active}
                    onCheckedChange={(checked) =>
                      toggleActiveMutation.mutate({ id: cohort.id, isActive: checked })
                    }
                    disabled={toggleActiveMutation.isPending}
                  />
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setEditingCohort(cohort)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => deleteMutation.mutate(cohort.id)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
            {filteredCohorts.length === 0 && (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                  No cohorts found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog open={!!editingCohort} onOpenChange={(open) => !open && setEditingCohort(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Cohort</DialogTitle>
          </DialogHeader>
          {editingCohort && (
            <form onSubmit={handleSubmitUpdate} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="edit-name">Name</Label>
                <Input id="edit-name" name="name" defaultValue={editingCohort.name} required />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="edit-start_date">Start Date</Label>
                  <Input id="edit-start_date" name="start_date" type="date" defaultValue={editingCohort.start_date} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="edit-end_date">End Date</Label>
                  <Input id="edit-end_date" name="end_date" type="date" defaultValue={editingCohort.end_date} />
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <Switch id="edit-active" name="is_active" defaultChecked={editingCohort.is_active} />
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
