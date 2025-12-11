import React, { useState } from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useQuery } from "@tanstack/react-query";
import { getRoomMembers } from "../api";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { UserPlus, Archive, LogOut, Clock, Shield } from "lucide-react";
import { format } from "date-fns";
import { AddMemberDialog } from "../AddMemberDialog";
import { Separator } from "@/components/ui/separator";
import { useTranslation } from 'react-i18next';

interface ChatRoomDetailsProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    roomId: string;
    roomName: string;
    currentUser: { id: string; role?: string };
    onArchive?: () => void;
}

export const ChatRoomDetails: React.FC<ChatRoomDetailsProps> = ({ 
    open, 
    onOpenChange, 
    roomId, 
    roomName,
    currentUser,
    onArchive
}) => {
    const { t } = useTranslation("common");
    const [addMemberOpen, setAddMemberOpen] = useState(false);

    const { data: members = [], isLoading } = useQuery({
        queryKey: ["chat", "members", roomId],
        queryFn: () => getRoomMembers(roomId),
        enabled: open, // Only fetch when open
    });

    const isAdmin = currentUser.role === 'admin' || currentUser.role === 'superadmin';

    // Simple initials
    const getInitials = (name: string) => name.substring(0, 2).toUpperCase();

    return (
        <>
            <Sheet open={open} onOpenChange={onOpenChange}>
                <SheetContent className="w-full sm:max-w-md overflow-y-auto">
                    <SheetHeader className="mb-6">
                        <SheetTitle>{t("chat.room_details", "Group Info")}</SheetTitle>
                    </SheetHeader>

                    <div className="flex flex-col items-center mb-6">
                        <div className="w-24 h-24 rounded-full bg-slate-200 dark:bg-slate-800 flex items-center justify-center text-3xl font-bold text-slate-500 mb-3">
                            {getInitials(roomName)}
                        </div>
                        <h2 className="text-xl font-bold">{roomName}</h2>
                        <p className="text-sm text-slate-500">
                            {t("chat.group", "Group")} · {members.length} {t("chat.participants", "participants")}
                        </p>
                    </div>

                    <div className="space-y-6">
                        {/* Actions */}
                        <div className="space-y-2">
                            {isAdmin && (
                                <Button variant="outline" className="w-full justify-start text-primary" onClick={() => setAddMemberOpen(true)}>
                                    <UserPlus className="mr-2 h-4 w-4" /> {t("chat.add_member", "Add Participants")}
                                </Button>
                            )}
                            {isAdmin && onArchive && (
                                <Button variant="outline" className="w-full justify-start text-red-500 hover:text-red-600" onClick={onArchive}>
                                    <Archive className="mr-2 h-4 w-4" /> {t("chat.archive_group", "Archive Group")}
                                </Button>
                            )}
                        </div>

                        <Separator />

                        {/* Members List */}
                        <div>
                            <h3 className="text-sm font-semibold mb-3 text-slate-500 uppercase tracking-wider">{t("chat.participants", "Participants")}</h3>
                            {isLoading ? (
                                <div className="text-center py-4 text-slate-400">Loading...</div>
                            ) : (
                                <div className="space-y-3">
                                    {members.map((m: any) => (
                                        <div key={m.user_id} className="flex items-center justify-between">
                                            <div className="flex items-center gap-3">
                                                <Avatar className="h-9 w-9">
                                                    <AvatarFallback className="bg-slate-100 text-slate-600 font-medium">
                                                        {getInitials(m.first_name || m.email || "?")}
                                                    </AvatarFallback>
                                                </Avatar>
                                                <div>
                                                    <div className="text-sm font-medium">
                                                        {m.first_name ? `${m.first_name} ${m.last_name || ''}` : m.username || m.email}
                                                        {m.user_id === currentUser.id && <span className="text-slate-400 font-normal ml-1">({t("chat.you", "You")})</span>}
                                                    </div>
                                                    <div className="text-xs text-slate-500 capitalize">{m.role_in_room}</div>
                                                </div>
                                            </div>
                                            <div className="flex items-center gap-2">
                                                {m.role_in_room === 'admin' && <Shield className="h-3 w-3 text-primary" />}
                                                {/* Logic to show last seen could go here if available */}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>
                </SheetContent>
            </Sheet>

            <AddMemberDialog 
                open={addMemberOpen} 
                onOpenChange={setAddMemberOpen} 
                roomId={roomId}
            />
        </>
    );
};
