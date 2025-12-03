import React, { useState, useEffect } from "react";
import { Command } from "cmdk";
import { Search, User, Loader2, Check } from "lucide-react";
import { useDebounce } from "use-debounce";
import { api } from "@/api/client";
import { useTranslation } from "react-i18next";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { addMember } from "./api";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { cn } from "@/lib/utils";

type UserResult = {
  id: string;
  name: string;
  email: string;
  role: string;
};

interface AddMemberDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  roomId: string;
}

export function AddMemberDialog({ open, onOpenChange, roomId }: AddMemberDialogProps) {
  const [query, setQuery] = useState("");
  const [debouncedQuery] = useDebounce(query, 300);
  const [results, setResults] = useState<UserResult[]>([]);
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation("common");
  const queryClient = useQueryClient();

  // Search users
  useEffect(() => {
    if (debouncedQuery.length < 2) {
      setResults([]);
      return;
    }

    const search = async () => {
      setLoading(true);
      try {
        // Reuse admin users endpoint but filter locally or use search endpoint if available.
        // Since we don't have a dedicated user search endpoint that returns JSON suitable here,
        // we'll use the admin list endpoint and filter. 
        // Ideally backend should support ?q=...
        // For now, let's fetch all and filter client side (not efficient but works for small user base)
        // OR use the /search endpoint if it returns users.
        // The /search endpoint returns SearchResult[], we need UserResult.
        // Let's try /admin/users first as it returns full user objects.
        const users = await api<UserResult[]>("/admin/users");
        const filtered = users.filter(u => 
          u.name.toLowerCase().includes(debouncedQuery.toLowerCase()) ||
          u.email.toLowerCase().includes(debouncedQuery.toLowerCase())
        ).slice(0, 10); // Limit to 10
        setResults(filtered);
      } catch (error) {
        console.error("Search failed", error);
        setResults([]);
      } finally {
        setLoading(false);
      }
    };

    search();
  }, [debouncedQuery]);

  const addMemberMutation = useMutation({
    mutationFn: (userId: string) => addMember(roomId, userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chat", "members", roomId] });
      // Also invalidate rooms list to update member counts if needed
      queryClient.invalidateQueries({ queryKey: ["chat", "rooms"] });
      onOpenChange(false);
      setQuery("");
    },
    onError: (error: any) => {
      console.error("Failed to add member", error);
      alert("Failed to add member");
    }
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="p-0 overflow-hidden max-w-md">
        <DialogHeader className="px-4 py-2 border-b">
          <DialogTitle>{t("chat.add_member", { defaultValue: "Add Member" })}</DialogTitle>
        </DialogHeader>
        <Command className="[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-muted-foreground [&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-group]]:px-2 [&_[cmdk-input-wrapper]_svg]:h-5 [&_[cmdk-input-wrapper]_svg]:w-5 [&_[cmdk-input]]:h-12 [&_[cmdk-item]]:px-2 [&_[cmdk-item]]:py-3 [&_[cmdk-item]_svg]:h-5 [&_[cmdk-item]_svg]:w-5">
          <div className="flex items-center border-b px-3" cmdk-input-wrapper="">
            <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
            <Command.Input
              value={query}
              onValueChange={setQuery}
              placeholder={t("chat.search_users", { defaultValue: "Search users by name or email..." })}
              className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>
          <Command.List className="max-h-[300px] overflow-y-auto overflow-x-hidden">
            <Command.Empty className="py-6 text-center text-sm">
              {loading ? (
                <div className="flex items-center justify-center gap-2">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  {t("common.searching", { defaultValue: "Searching..." })}
                </div>
              ) : (
                t("common.no_results", { defaultValue: "No results found." })
              )}
            </Command.Empty>
            
            {results.length > 0 && (
              <Command.Group>
                {results.map((user) => (
                  <Command.Item
                    key={user.id}
                    value={`${user.name} ${user.email}`}
                    onSelect={() => addMemberMutation.mutate(user.id)}
                    className="relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none aria-selected:bg-accent aria-selected:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                  >
                    <User className="mr-2 h-4 w-4" />
                    <div className="flex flex-col flex-1">
                      <span className="font-medium">{user.name}</span>
                      <span className="text-xs text-muted-foreground">{user.email}</span>
                    </div>
                    {addMemberMutation.isPending && (
                      <Loader2 className="h-4 w-4 animate-spin ml-2" />
                    )}
                  </Command.Item>
                ))}
              </Command.Group>
            )}
          </Command.List>
        </Command>
      </DialogContent>
    </Dialog>
  );
}
