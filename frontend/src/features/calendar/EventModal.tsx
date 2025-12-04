import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useAuth } from "@/contexts/AuthContext";
import { useMutation } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

type EventFormValues = {
  title: string;
  description?: string;
  start_time: string;
  end_time: string;
  event_type: "meeting" | "deadline" | "academic";
  location?: string;
};

interface EventModalProps {
  isOpen: boolean;
  onClose: () => void;
  event: any | null;
  onSuccess: () => void;
}

export const EventModal: React.FC<EventModalProps> = ({ isOpen, onClose, event, onSuccess }) => {
  const { token } = useAuth();
  const { t } = useTranslation("common");
  const eventSchema = React.useMemo(
    () =>
      z.object({
        title: z.string().min(1, t("calendar.validation.title", { defaultValue: "Title is required" })),
        description: z.string().optional(),
        start_time: z
          .string()
          .min(1, t("calendar.validation.start", { defaultValue: "Start time is required" })),
        end_time: z
          .string()
          .min(1, t("calendar.validation.end", { defaultValue: "End time is required" })),
        event_type: z.enum(["meeting", "deadline", "academic"]),
        location: z.string().optional(),
      }),
    [t]
  );
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    formState: { errors },
  } = useForm<EventFormValues>({
    resolver: zodResolver(eventSchema),
    defaultValues: {
      event_type: "meeting",
    },
  });

  useEffect(() => {
    if (event) {
      setValue('title', event.title);
      setValue('description', event.description || '');
      setValue('start_time', new Date(event.start_time).toISOString().slice(0, 16));
      setValue('end_time', new Date(event.end_time).toISOString().slice(0, 16));
      setValue('event_type', event.event_type);
      setValue('location', event.location || '');
    } else {
      reset({
        event_type: 'meeting',
        start_time: new Date().toISOString().slice(0, 16),
        end_time: new Date(new Date().getTime() + 3600000).toISOString().slice(0, 16),
      });
    }
  }, [event, isOpen, reset, setValue]);

  const mutation = useMutation({
    mutationFn: async (data: EventFormValues) => {
      const url = event
        ? `${API_URL}/events/${event.id}`
        : `${API_URL}/events`;

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

      if (!res.ok) throw new Error(t("calendar.errors.save", { defaultValue: "Failed to save event" }));
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
    if (!event || !confirm(t("calendar.confirm_delete", { defaultValue: "Are you sure you want to delete this event?" })))
      return;
    
    try {
      const res = await fetch(`${API_URL}/events/${event.id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error(t("calendar.errors.delete", { defaultValue: "Failed to delete event" }));
      onSuccess();
      onClose();
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
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
            <Input id="title" {...register('title')} />
            {errors.title && <span className="text-red-500 text-sm">{errors.title.message}</span>}
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="start_time">
                {t("calendar.fields.start", { defaultValue: "Start" })}
              </Label>
              <Input id="start_time" type="datetime-local" {...register('start_time')} />
              {errors.start_time && <span className="text-red-500 text-sm">{errors.start_time.message}</span>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="end_time">
                {t("calendar.fields.end", { defaultValue: "End" })}
              </Label>
              <Input id="end_time" type="datetime-local" {...register('end_time')} />
              {errors.end_time && <span className="text-red-500 text-sm">{errors.end_time.message}</span>}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="event_type">
              {t("calendar.fields.type", { defaultValue: "Type" })}
            </Label>
            <Select onValueChange={(val) => setValue('event_type', val as any)} defaultValue={event?.event_type || 'meeting'}>
              <SelectTrigger>
                <SelectValue placeholder={t("calendar.fields.type_placeholder", { defaultValue: "Select type" })} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="meeting">{t("calendar.types.meeting", { defaultValue: "Meeting" })}</SelectItem>
                <SelectItem value="deadline">{t("calendar.types.deadline", { defaultValue: "Deadline" })}</SelectItem>
                <SelectItem value="academic">{t("calendar.types.academic", { defaultValue: "Academic" })}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="location">
              {t("calendar.fields.location", { defaultValue: "Location" })}
            </Label>
            <Input id="location" {...register('location')} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">
              {t("calendar.fields.description", { defaultValue: "Description" })}
            </Label>
            <Textarea id="description" {...register('description')} />
          </div>

          <DialogFooter className="flex justify-between sm:justify-between">
            {event && (
              <Button type="button" variant="destructive" onClick={handleDelete}>
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
