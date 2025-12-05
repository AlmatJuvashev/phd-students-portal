import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAuth } from "@/contexts/AuthContext";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { API_URL } from "@/api/client";
import { Video, MapPin } from "lucide-react";
import { cn } from "@/lib/utils";

type EventFormValues = {
  title: string;
  description?: string;
  start_time: string;
  end_time: string;
  event_type: "meeting" | "deadline" | "academic";
  location?: string;
  meeting_type: "online" | "offline";
  meeting_url?: string;
  physical_address?: string;
  color?: string;
};

interface EventModalProps {
  isOpen: boolean;
  onClose: () => void;
  event: any | null;
  onSuccess: () => void;
  defaultStart?: Date;
  defaultEnd?: Date;
}

// Color options for events
const EVENT_COLORS = [
  { key: "blue", label: "Blue (Online)", hex: "#3b82f6" },
  { key: "purple", label: "Purple (Offline)", hex: "#8b5cf6" },
  { key: "red", label: "Red (Deadline)", hex: "#ef4444" },
  { key: "green", label: "Green (Academic)", hex: "#22c55e" },
  { key: "orange", label: "Orange", hex: "#f97316" },
  { key: "pink", label: "Pink", hex: "#ec4899" },
];

export const EventModal: React.FC<EventModalProps> = ({
  isOpen,
  onClose,
  event,
  onSuccess,
  defaultStart,
  defaultEnd,
}) => {
  const { token } = useAuth();
  const { t } = useTranslation("common");
  const [meetingType, setMeetingType] = useState<"online" | "offline">("offline");
  const [selectedColor, setSelectedColor] = useState("blue");

  const eventSchema = React.useMemo(
    () =>
      z.object({
        title: z
          .string()
          .min(
            1,
            t("calendar.validation.title", {
              defaultValue: "Title is required",
            })
          ),
        description: z.string().optional(),
        start_time: z
          .string()
          .min(
            1,
            t("calendar.validation.start", {
              defaultValue: "Start time is required",
            })
          ),
        end_time: z
          .string()
          .min(
            1,
            t("calendar.validation.end", {
              defaultValue: "End time is required",
            })
          ),
        event_type: z.enum(["meeting", "deadline", "academic"]),
        location: z.string().optional(),
        meeting_type: z.enum(["online", "offline"]),
        meeting_url: z.string().optional(),
        physical_address: z.string().optional(),
        color: z.string().optional(),
      }),
    [t]
  );

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<EventFormValues>({
    resolver: zodResolver(eventSchema),
    defaultValues: {
      event_type: "meeting",
      meeting_type: "offline",
      color: "blue",
    },
  });

  const watchedMeetingType = watch("meeting_type");

  // Helper to format date as local datetime for input[type="datetime-local"]
  const formatLocalDatetime = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  };

  useEffect(() => {
    if (event) {
      setValue("title", event.title);
      setValue("description", event.description || "");
      // Use local time formatting for existing events
      setValue("start_time", formatLocalDatetime(new Date(event.start_time)));
      setValue("end_time", formatLocalDatetime(new Date(event.end_time)));
      setValue("event_type", event.event_type);
      setValue("location", event.location || "");
      setValue("meeting_type", event.meeting_type || "offline");
      setValue("meeting_url", event.meeting_url || "");
      setValue("physical_address", event.physical_address || "");
      setValue("color", event.color || "blue");
      setMeetingType(event.meeting_type || "offline");
      setSelectedColor(event.color || "blue");
    } else {
      const startTime = defaultStart || new Date();
      const endTime = defaultEnd || new Date(startTime.getTime() + 3600000);
      reset({
        event_type: "meeting",
        meeting_type: "offline",
        color: "blue",
        // Use local time formatting for new events
        start_time: formatLocalDatetime(startTime),
        end_time: formatLocalDatetime(endTime),
      });
      setMeetingType("offline");
      setSelectedColor("blue");
    }
  }, [event, isOpen, reset, setValue, defaultStart, defaultEnd]);

  const mutation = useMutation({
    mutationFn: async (data: EventFormValues) => {
      const url = event ? `${API_URL}/events/${event.id}` : `${API_URL}/events`;
      const method = event ? "PUT" : "POST";

      const res = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...data,
          start_time: new Date(data.start_time).toISOString(),
          end_time: new Date(data.end_time).toISOString(),
        }),
      });

      if (!res.ok)
        throw new Error(
          t("calendar.errors.save", { defaultValue: "Failed to save event" })
        );
      return res.json();
    },
    onSuccess: () => {
      onSuccess();
      onClose();
    },
  });

  const onSubmit = (data: EventFormValues) => {
    mutation.mutate(data);
  };

  const handleDelete = async () => {
    if (
      !event ||
      !confirm(
        t("calendar.confirm_delete", {
          defaultValue: "Are you sure you want to delete this event?",
        })
      )
    )
      return;

    try {
      const res = await fetch(`${API_URL}/events/${event.id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok)
        throw new Error(
          t("calendar.errors.delete", {
            defaultValue: "Failed to delete event",
          })
        );
      onSuccess();
      onClose();
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {event
              ? t("calendar.edit_event", { defaultValue: "Edit Event" })
              : t("calendar.create_event", { defaultValue: "Create Event" })}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">
              {t("calendar.fields.title", { defaultValue: "Title" })}
            </Label>
            <Input id="title" {...register("title")} />
            {errors.title && (
              <span className="text-red-500 text-sm">
                {errors.title.message}
              </span>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="start_time">
                {t("calendar.fields.start", { defaultValue: "Start" })}
              </Label>
              <Input
                id="start_time"
                type="datetime-local"
                {...register("start_time")}
              />
              {errors.start_time && (
                <span className="text-red-500 text-sm">
                  {errors.start_time.message}
                </span>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="end_time">
                {t("calendar.fields.end", { defaultValue: "End" })}
              </Label>
              <Input
                id="end_time"
                type="datetime-local"
                {...register("end_time")}
              />
              {errors.end_time && (
                <span className="text-red-500 text-sm">
                  {errors.end_time.message}
                </span>
              )}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="event_type">
                {t("calendar.fields.type", { defaultValue: "Type" })}
              </Label>
              <Select
                onValueChange={(val) => setValue("event_type", val as any)}
                defaultValue={event?.event_type || "meeting"}
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={t("calendar.fields.type_placeholder", {
                      defaultValue: "Select type",
                    })}
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="meeting">
                    {t("calendar.types.meeting", { defaultValue: "Meeting" })}
                  </SelectItem>
                  <SelectItem value="deadline">
                    {t("calendar.types.deadline", { defaultValue: "Deadline" })}
                  </SelectItem>
                  <SelectItem value="academic">
                    {t("calendar.types.academic", { defaultValue: "Academic" })}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>
                {t("calendar.fields.meeting_type", { defaultValue: "Format" })}
              </Label>
              <div className="flex gap-2">
                <Button
                  type="button"
                  variant={watchedMeetingType === "online" ? "default" : "outline"}
                  size="sm"
                  className="flex-1"
                  onClick={() => {
                    setValue("meeting_type", "online");
                    setMeetingType("online");
                  }}
                >
                  <Video className="h-4 w-4 mr-1" />
                  Online
                </Button>
                <Button
                  type="button"
                  variant={watchedMeetingType === "offline" ? "default" : "outline"}
                  size="sm"
                  className="flex-1"
                  onClick={() => {
                    setValue("meeting_type", "offline");
                    setMeetingType("offline");
                  }}
                >
                  <MapPin className="h-4 w-4 mr-1" />
                  Offline
                </Button>
              </div>
            </div>
          </div>

          {/* Conditional fields based on meeting type */}
          {watchedMeetingType === "online" ? (
            <div className="space-y-2">
              <Label htmlFor="meeting_url">
                {t("calendar.fields.meeting_url", { defaultValue: "Meeting URL (Zoom/Google Meet)" })}
              </Label>
              <Input
                id="meeting_url"
                type="url"
                placeholder="https://zoom.us/j/..."
                {...register("meeting_url")}
              />
            </div>
          ) : (
            <div className="space-y-2">
              <Label htmlFor="physical_address">
                {t("calendar.fields.physical_address", { defaultValue: "Physical Address" })}
              </Label>
              <Input
                id="physical_address"
                placeholder={t("calendar.fields.address_placeholder", { defaultValue: "Enter location address" })}
                {...register("physical_address")}
              />
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="location">
              {t("calendar.fields.location", { defaultValue: "Location/Room" })}
            </Label>
            <Input id="location" {...register("location")} />
          </div>

          {/* Color picker */}
          <div className="space-y-2">
            <Label>
              {t("calendar.fields.color", { defaultValue: "Event Color" })}
            </Label>
            <div className="flex gap-2 flex-wrap">
              {EVENT_COLORS.map((color) => (
                <button
                  key={color.key}
                  type="button"
                  onClick={() => {
                    setValue("color", color.key);
                    setSelectedColor(color.key);
                  }}
                  className={cn(
                    "w-8 h-8 rounded-full transition-all border-2",
                    selectedColor === color.key
                      ? "ring-2 ring-offset-2 ring-primary scale-110"
                      : "border-transparent hover:scale-105"
                  )}
                  style={{ backgroundColor: color.hex }}
                  title={color.label}
                />
              ))}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">
              {t("calendar.fields.description", {
                defaultValue: "Description",
              })}
            </Label>
            <Textarea id="description" {...register("description")} />
          </div>

          <DialogFooter className="flex justify-between sm:justify-between">
            {event && (
              <Button
                type="button"
                variant="destructive"
                onClick={handleDelete}
              >
                {t("calendar.actions.delete", { defaultValue: "Delete" })}
              </Button>
            )}
            <div className="flex gap-2">
              <Button type="button" variant="outline" onClick={onClose}>
                {t("common.cancel", { defaultValue: "Cancel" })}
              </Button>
              <Button type="submit" disabled={mutation.isPending}>
                {mutation.isPending
                  ? t("calendar.actions.saving", { defaultValue: "Saving..." })
                  : t("calendar.actions.save", { defaultValue: "Save" })}
              </Button>
            </div>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};
