import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { createRoom } from "./api";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { Label } from "@/components/ui/label";

interface CreateRoomDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CreateRoomDialog({ open, onOpenChange }: CreateRoomDialogProps) {
  const { t } = useTranslation("common");
  const [name, setName] = useState("");
  const [type, setType] = useState<"group" | "channel">("group");
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: () => createRoom(name, type),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chat", "rooms"] });
      onOpenChange(false);
      setName("");
      setType("group");
    },
    onError: (err: any) => {
      console.error("Failed to create room:", err);
      // Try to extract error message from response if possible, or just stringify
      const msg = err.message || JSON.stringify(err);
      alert(`Failed to create group. Please ensure backend is running and you have permission. Error: ${msg}`);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    createMutation.mutate();
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle>{t("chat.create_group", "Create New Group")}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="name">
              {t("chat.group_name", "Group Name")}
            </Label>
            <Input 
              id="name" 
              type="text"
              value={name} 
              onChange={(e) => setName(e.target.value)} 
              placeholder={t("chat.group_name_placeholder", "e.g. Thesis Discussion")}
              autoFocus
            />
          </div>
          
          <div className="space-y-2">
            <label htmlFor="type" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              {t("chat.group_type", "Type")}
            </label>
            <select
              id="type"
              value={type}
              onChange={(e) => setType(e.target.value as any)}
              className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="group">Group (Private)</option>
              <option value="channel">Channel (Broadcast)</option>
            </select>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              {t("common.cancel", "Cancel")}
            </Button>
            <Button type="submit" disabled={createMutation.isPending || !name.trim()}>
              {createMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {t("common.create", "Create")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
