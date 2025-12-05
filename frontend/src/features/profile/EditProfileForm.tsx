import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { UserProfile, updateProfile } from "@/api/user";
import { useToast } from "@/components/ui/use-toast";
import { useState } from "react";
import { useTranslation } from "react-i18next";

const profileSchema = z.object({
  email: z.string().email(),
  phone: z.string().optional(),
  bio: z.string().optional(),
  address: z.string().optional(),
  date_of_birth: z.string().optional(),
  current_password: z.string().min(1),
});

interface EditProfileFormProps {
  user: UserProfile;
  onSuccess: () => void;
}

export function EditProfileForm({ user, onSuccess }: EditProfileFormProps) {
  const { toast } = useToast();
  const { t } = useTranslation("common");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<z.infer<typeof profileSchema>>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      email: user.email || "",
      phone: user.phone || "",
      bio: user.bio || "",
      address: user.address || "",
      date_of_birth: user.date_of_birth
        ? new Date(user.date_of_birth).toISOString().split("T")[0]
        : "",
      current_password: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof profileSchema>) => {
    setIsSubmitting(true);
    try {
      const payload: any = { ...values };
      if (values.date_of_birth) {
        payload.date_of_birth = new Date(values.date_of_birth).toISOString();
      } else {
        payload.date_of_birth = null;
      }

      await updateProfile(payload);
      toast({ title: t("profile.profile_updated") });
      onSuccess();
    } catch (error) {
      console.error(error);
      toast({
        title: t("profile.update_failed"),
        description: t("profile.try_again"),
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <FormItem>
        <FormLabel>{t("profile.email")}</FormLabel>
        <Input type="email" {...register("email")} />
        {errors.email && <FormMessage>{errors.email.message}</FormMessage>}
      </FormItem>

      <FormItem>
        <FormLabel>{t("profile.phone_number")}</FormLabel>
        <Input placeholder={t("profile.phone_placeholder")} {...register("phone")} />
        {errors.phone && <FormMessage>{errors.phone.message}</FormMessage>}
      </FormItem>

      <FormItem>
        <FormLabel>{t("profile.date_of_birth")}</FormLabel>
        <Input type="date" {...register("date_of_birth")} />
        {errors.date_of_birth && (
          <FormMessage>{errors.date_of_birth.message}</FormMessage>
        )}
      </FormItem>

      <FormItem>
        <FormLabel>{t("profile.address")}</FormLabel>
        <Input placeholder={t("profile.address_placeholder")} {...register("address")} />
        {errors.address && <FormMessage>{errors.address.message}</FormMessage>}
      </FormItem>

      <FormItem>
        <FormLabel>{t("profile.bio")}</FormLabel>
        <Textarea
          placeholder={t("profile.bio_placeholder")}
          className="resize-none"
          {...register("bio")}
        />
        {errors.bio && <FormMessage>{errors.bio.message}</FormMessage>}
      </FormItem>

      <div className="border-t pt-4">
        <FormItem>
          <FormLabel>{t("profile.current_password")}</FormLabel>
          <Input type="password" {...register("current_password")} />
          {errors.current_password && (
            <FormMessage>{t("profile.password_required")}</FormMessage>
          )}
        </FormItem>
      </div>

      <div className="flex justify-end">
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? t("profile.saving") : t("profile.save_changes")}
        </Button>
      </div>
    </form>
  );
}
