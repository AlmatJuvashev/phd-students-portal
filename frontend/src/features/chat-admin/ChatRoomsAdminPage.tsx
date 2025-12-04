import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import {
  Plus,
  Archive,
  RefreshCcw,
  MessageCircle,
  Check,
  Users,
  Search,
  UserPlus,
  X,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { ChatRoom } from "@/features/chat/api";
import {
  addRoomMember,
  createRoom,
  listAdminRooms,
  listRoomMembers,
  removeRoomMember,
  searchUsers,
  updateRoom,
  UserSearchResult,
  addRoomMembersBatch,
  removeRoomMembersBatch,
} from "./api";
import {
  listPrograms,
  listDepartments,
  listCohorts,
  listSpecialties,
  Program,
  Department,
  Cohort,
  Specialty,
} from "@/features/admin/dictionaries/api";
import { useToast } from "@/components/ui/use-toast";
import { Modal } from "@/components/ui/modal";

type RoomForm = {
  name: string;
  type: "cohort" | "advisory" | "other";
};

export function ChatRoomsAdminPage() {
  const { t } = useTranslation("common");
  const qc = useQueryClient();
  const { toast } = useToast();
  const [form, setForm] = useState<RoomForm>({ name: "", type: "cohort" });

  const {
    data: rooms = [],
    isLoading,
    isFetching,
  } = useQuery<ChatRoom[]>({
    queryKey: ["chat", "admin", "rooms"],
    queryFn: listAdminRooms,
  });

  const createMutation = useMutation({
    mutationFn: () => createRoom(form),
    onSuccess: () => {
      setForm({ name: "", type: "cohort" });
      qc.invalidateQueries({ queryKey: ["chat", "admin", "rooms"] });
    },
  });

  const archiveMutation = useMutation({
    mutationFn: (room: ChatRoom) =>
      updateRoom(room.id, { is_archived: !room.is_archived }),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ["chat", "admin", "rooms"] }),
  });

  const [manageRoomId, setManageRoomId] = useState<string | null>(null);
  const [memberSearch, setMemberSearch] = useState("");
  const [filters, setFilters] = useState({
    program: "",
    department: "",
    cohort: "",
    specialty: "",
  });

  const { data: programs = [] } = useQuery({
    queryKey: ["programs", "active"],
    queryFn: () => listPrograms(true),
    staleTime: 0,
  });
  const { data: departments = [] } = useQuery({
    queryKey: ["departments", "active"],
    queryFn: () => listDepartments(true),
    staleTime: 0,
  });
  const { data: cohorts = [] } = useQuery({
    queryKey: ["cohorts", "active"],
    queryFn: () => listCohorts(true),
    staleTime: 0,
  });
  const { data: specialties = [] } = useQuery({
    queryKey: ["specialties", "active"],
    queryFn: () => listSpecialties(true),
    staleTime: 0,
  });

  const membersQuery = useQuery({
    queryKey: ["chat", "admin", "members", manageRoomId],
    queryFn: () => listRoomMembers(manageRoomId || ""),
    enabled: !!manageRoomId,
  });
  const searchQuery = useQuery<UserSearchResult[]>({
    queryKey: ["chat", "admin", "user-search", memberSearch, filters],
    queryFn: () => searchUsers(memberSearch, filters),
    enabled:
      memberSearch.trim().length >= 2 || Object.values(filters).some(Boolean),
  });

  const addMemberMutation = useMutation({
    mutationFn: (input: { user_id: string }) =>
      addRoomMember(manageRoomId || "", { user_id: input.user_id }),
    onSuccess: () => {
      qc.invalidateQueries({
        queryKey: ["chat", "admin", "members", manageRoomId],
      });
      setMemberSearch("");
    },
  });

  const removeMemberMutation = useMutation({
    mutationFn: (userId: string) =>
      removeRoomMember(manageRoomId || "", userId),
    onSuccess: () =>
      qc.invalidateQueries({
        queryKey: ["chat", "admin", "members", manageRoomId],
      }),
  });

  const addBatchMutation = useMutation({
    mutationFn: () => addRoomMembersBatch(manageRoomId || "", filters),
    onSuccess: (data: any) => {
      qc.invalidateQueries({
        queryKey: ["chat", "admin", "members", manageRoomId],
      });
      toast({
        title: t("chat_admin.added_members", { defaultValue: "Added members" }),
        description: t("chat_admin.added_members_count", {
          defaultValue: "Added {{count}} members",
          count: data?.added_count ?? 0,
        }),
      });
    },
  });

  const removeBatchMutation = useMutation({
    mutationFn: () => removeRoomMembersBatch(manageRoomId || "", filters),
    onSuccess: (data: any) => {
      qc.invalidateQueries({
        queryKey: ["chat", "admin", "members", manageRoomId],
      });
      toast({
        title: t("chat_admin.removed_members", {
          defaultValue: "Removed members",
        }),
        description: t("chat_admin.removed_members_count", {
          defaultValue: "Removed {{count}} members",
          count: data?.removed_count ?? 0,
        }),
      });
    },
  });

  const hasRooms = rooms.length > 0;
  const currentRoom = useMemo(
    () => rooms.find((r) => r.id === manageRoomId),
    [rooms, manageRoomId]
  );

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-xl font-semibold flex items-center gap-2">
            <MessageCircle className="h-5 w-5 text-primary" />
            {t("chat_admin.title", { defaultValue: "Chat rooms" })}
          </h1>
          <p className="text-sm text-muted-foreground">
            {t("chat_admin.subtitle", {
              defaultValue:
                "Create cohort/advisory rooms and archive when done.",
            })}
          </p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() =>
            qc.invalidateQueries({ queryKey: ["chat", "admin", "rooms"] })
          }
          disabled={isFetching}
        >
          <RefreshCcw className="h-4 w-4 mr-2" />
          {t("common.refresh", { defaultValue: "Refresh" })}
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Plus className="h-4 w-4" />
            {t("chat_admin.create_heading", { defaultValue: "New room" })}
          </CardTitle>
        </CardHeader>
        <CardContent className="grid gap-3 sm:grid-cols-[1fr,200px,auto]">
          <div className="space-y-1.5">
            <Label>
              {t("chat_admin.room_name", { defaultValue: "Room name" })}
            </Label>
            <Input
              placeholder={t("chat_admin.room_placeholder", {
                defaultValue: "PhD 2025 • Public Health",
              })}
              value={form.name}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, name: e.target.value }))
              }
            />
          </div>
          <div className="space-y-1.5">
            <Label>{t("chat_admin.room_type", { defaultValue: "Type" })}</Label>
            <Select
              value={form.type}
              onValueChange={(v) =>
                setForm((prev) => ({ ...prev, type: v as RoomForm["type"] }))
              }
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="cohort">{t("chat.types.cohort")}</SelectItem>
                <SelectItem value="advisory">
                  {t("chat.types.advisory")}
                </SelectItem>
                <SelectItem value="other">{t("chat.types.admin")}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex items-end">
            <Button
              onClick={() => createMutation.mutate()}
              disabled={!form.name.trim() || createMutation.isPending}
            >
              <Plus className="h-4 w-4 mr-2" />
              {createMutation.isPending
                ? t("common.loading")
                : t("chat_admin.create")}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            {t("chat_admin.rooms_heading", { defaultValue: "Rooms" })}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          {isLoading ? (
            <div className="rounded-md border border-dashed p-3 text-sm text-muted-foreground">
              {t("chat.loading_rooms")}
            </div>
          ) : !hasRooms ? (
            <div className="rounded-md border border-dashed p-3 text-sm text-muted-foreground">
              {t("chat.no_rooms")}
            </div>
          ) : (
            rooms.map((room) => (
              <div
                key={room.id}
                className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 border rounded-md px-3 py-2"
              >
                <div className="space-y-0.5">
                  <div className="flex items-center gap-2">
                    <span className="font-semibold">{room.name}</span>
                    {room.is_archived && (
                      <Badge variant="outline" className="text-[10px]">
                        {t("chat.archived")}
                      </Badge>
                    )}
                  </div>
                  <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                    <Badge variant="outline" className="text-[10px]">
                      {t(`chat.types.${room.type}`)}
                    </Badge>
                    {room.created_by_role && (
                      <Badge variant="outline" className="text-[10px]">
                        {room.created_by_role}
                      </Badge>
                    )}
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setManageRoomId(room.id)}
                      className="flex items-center gap-2"
                    >
                      <Users className="h-4 w-4" />
                      {t("chat_admin.manage_members", {
                        defaultValue: "Members",
                      })}
                    </Button>
                    <Button
                      variant={room.is_archived ? "outline" : "ghost"}
                      size="sm"
                      onClick={() => archiveMutation.mutate(room)}
                      disabled={archiveMutation.isPending}
                      className="flex items-center gap-2"
                    >
                      {room.is_archived ? (
                        <Check className="h-4 w-4" />
                      ) : (
                        <Archive className="h-4 w-4" />
                      )}
                      {room.is_archived
                        ? t("chat_admin.unarchive", {
                            defaultValue: "Unarchive",
                          })
                        : t("chat_admin.archive", { defaultValue: "Archive" })}
                    </Button>
                  </div>
                </div>
              </div>
            ))
          )}
        </CardContent>
      </Card>

      <Modal
        open={!!manageRoomId}
        onOpenChange={(open) => !open && setManageRoomId(null)}
        size="lg"
      >
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <div>
              <div className="font-semibold">{currentRoom?.name}</div>
              <div className="text-sm text-muted-foreground">
                {t("chat_admin.members_title", { defaultValue: "Members" })}
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setManageRoomId(null)}
              aria-label={t("common.close", { defaultValue: "Close" })}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>

          <div className="space-y-2">
            <Label>
              {t("chat_admin.add_member", { defaultValue: "Add member" })}
            </Label>
            <div className="grid grid-cols-2 gap-2">
              <Select
                value={filters.program}
                onValueChange={(v) =>
                  setFilters((prev) => ({
                    ...prev,
                    program: v === "all" ? "" : v,
                  }))
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.forms.program", {
                      defaultValue: "Program",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">
                    {t("admin.forms.all_programs", {
                      defaultValue: "All programs",
                    })}
                  </SelectItem>
                  {programs.map((p: Program) => (
                    <SelectItem key={p.id} value={p.name}>
                      {p.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select
                value={filters.department}
                onValueChange={(v) =>
                  setFilters((prev) => ({
                    ...prev,
                    department: v === "all" ? "" : v,
                  }))
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.forms.department", {
                      defaultValue: "Department",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">
                    {t("admin.forms.all_departments", {
                      defaultValue: "All departments",
                    })}
                  </SelectItem>
                  {departments.map((d: Department) => (
                    <SelectItem key={d.id} value={d.name}>
                      {d.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select
                value={filters.cohort}
                onValueChange={(v) =>
                  setFilters((prev) => ({
                    ...prev,
                    cohort: v === "all" ? "" : v,
                  }))
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.forms.cohort", {
                      defaultValue: "Cohort",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">
                    {t("admin.forms.all_cohorts", {
                      defaultValue: "All cohorts",
                    })}
                  </SelectItem>
                  {cohorts.map((c: Cohort) => (
                    <SelectItem key={c.id} value={c.name}>
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select
                value={filters.specialty}
                onValueChange={(v) =>
                  setFilters((prev) => ({
                    ...prev,
                    specialty: v === "all" ? "" : v,
                  }))
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("admin.dictionaries.tabs.specialties", {
                      defaultValue: "Specialties",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">
                    {t("admin.dictionaries.specialties.all", {
                      defaultValue: "All specialties",
                    })}
                  </SelectItem>
                  {specialties.map((s: Specialty) => (
                    <SelectItem key={s.id} value={s.name}>
                      {s.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="flex gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder={t("chat_admin.search_user", {
                    defaultValue: "Search by name or email",
                  })}
                  value={memberSearch}
                  onChange={(e) => setMemberSearch(e.target.value)}
                  className="pl-9"
                />
              </div>
              <Button
                variant="secondary"
                onClick={() => addBatchMutation.mutate()}
                disabled={
                  addBatchMutation.isPending ||
                  !Object.values(filters).some(Boolean)
                }
              >
                {addBatchMutation.isPending
                  ? t("chat_admin.adding_all", { defaultValue: "Adding..." })
                  : t("chat_admin.add_all", { defaultValue: "Add all" })}
              </Button>
              <Button
                variant="destructive"
                onClick={() => removeBatchMutation.mutate()}
                disabled={
                  removeBatchMutation.isPending ||
                  !Object.values(filters).some(Boolean)
                }
              >
                {removeBatchMutation.isPending
                  ? t("chat_admin.removing_all", {
                      defaultValue: "Removing...",
                    })
                  : t("chat_admin.remove_all", { defaultValue: "Remove all" })}
              </Button>
            </div>
            {(memberSearch.trim().length >= 2 ||
              Object.values(filters).some(Boolean)) && (
              <div className="border rounded-md divide-y max-h-48 overflow-y-auto">
                {searchQuery.isFetching ? (
                  <div className="p-3 text-sm text-muted-foreground">
                    {t("common.loading", { defaultValue: "Loading…" })}
                  </div>
                ) : (searchQuery.data ?? []).length === 0 ? (
                  <div className="p-3 text-sm text-muted-foreground">
                    {t("chat_admin.no_users", {
                      defaultValue: "No users found",
                    })}
                  </div>
                ) : (
                  searchQuery.data?.map((u) => (
                    <div
                      key={u.id}
                      className="flex items-center justify-between px-3 py-2"
                    >
                      <div>
                        <div className="font-medium text-sm">
                          {u.name || u.email}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {u.email}
                        </div>
                      </div>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() =>
                          addMemberMutation.mutate({ user_id: u.id })
                        }
                        disabled={addMemberMutation.isPending}
                      >
                        <UserPlus className="h-4 w-4 mr-1" />
                        {t("chat_admin.add", { defaultValue: "Add" })}
                      </Button>
                    </div>
                  ))
                )}
              </div>
            )}
          </div>

          <div className="border rounded-md">
            <div className="px-3 py-2 text-sm font-semibold">
              {t("chat_admin.current_members", {
                defaultValue: "Current members",
              })}
            </div>
            <div className="divide-y max-h-64 overflow-y-auto">
              {membersQuery.isLoading ? (
                <div className="p-3 text-sm text-muted-foreground">
                  {t("chat.loading_rooms")}
                </div>
              ) : (membersQuery.data ?? []).length === 0 ? (
                <div className="p-3 text-sm text-muted-foreground">
                  {t("chat_admin.no_members", {
                    defaultValue: "No members yet",
                  })}
                </div>
              ) : (
                membersQuery.data?.map((m) => (
                  <div
                    key={m.user_id}
                    className="flex items-center justify-between px-3 py-2 text-sm"
                  >
                    <div className="space-y-0.5">
                      <div className="font-medium">
                        {[m.first_name, m.last_name]
                          .filter(Boolean)
                          .join(" ") ||
                          m.email ||
                          m.user_id}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {m.email}
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => removeMemberMutation.mutate(m.user_id)}
                      disabled={removeMemberMutation.isPending}
                    >
                      <X className="h-4 w-4 mr-1" />
                      {t("chat_admin.remove", { defaultValue: "Remove" })}
                    </Button>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      </Modal>
    </div>
  );
}
