import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Dice5, Check, RefreshCw } from "lucide-react";
import { updateAvatar } from "@/api/user";
import { useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/components/ui/use-toast";
import { useTranslation } from "react-i18next";

interface AvatarPickerModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const AvatarPickerModal = ({ isOpen, onClose }: AvatarPickerModalProps) => {
  const { t } = useTranslation("common");
  const [seed, setSeed] = useState(() => Math.random().toString(36).substring(7));
  const [loading, setLoading] = useState(false);
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const avatarUrl = `https://api.dicebear.com/7.x/avataaars/svg?seed=${seed}`;

  const handleRandomize = () => {
    setSeed(Math.random().toString(36).substring(7));
  };

  const handleSave = async () => {
    setLoading(true);
    try {
      await updateAvatar(avatarUrl);
      console.log("Avatar updated, invalidating cache...");
      await queryClient.invalidateQueries({ queryKey: ["me"] });
      await queryClient.refetchQueries({ queryKey: ["me"] });
      console.log("Cache invalidated and refetched.");
      toast({ title: t("profile.avatar_updated", "Avatar updated successfully") });
      onClose();
    } catch (error) {
      console.error(error);
      toast({ 
        title: t("profile.avatar_update_failed", "Failed to update avatar"), 
        variant: "destructive" 
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("profile.choose_avatar", "Choose Avatar")}</DialogTitle>
        </DialogHeader>
        
        <div className="flex flex-col items-center justify-center p-6 gap-6">
          <div className="relative">
            <div className="w-40 h-40 rounded-full border-4 border-slate-100 dark:border-slate-800 overflow-hidden shadow-xl">
               <img 
                 src={avatarUrl} 
                 alt="Avatar Preview" 
                 className="w-full h-full object-cover bg-slate-50"
               />
            </div>
            <Button
              size="icon"
              variant="outline"
              className="absolute bottom-0 right-0 rounded-full shadow-lg h-10 w-10 bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-700"
              onClick={handleRandomize}
            >
              <RefreshCw className="h-4 w-4" />
            </Button>
          </div>
          
          <div className="text-center text-sm text-muted-foreground">
            {t("profile.avatar_hint", "Click the shuffle button to generate a new unique look.")}
          </div>
        </div>

        <DialogFooter className="flex gap-2 sm:justify-center">
            <Button variant="outline" onClick={handleRandomize} className="gap-2">
                <Dice5 className="w-4 h-4" />
                {t("profile.randomize", "Randomize")}
            </Button>
            <Button onClick={handleSave} disabled={loading} className="gap-2">
                {loading ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Check className="w-4 h-4" />}
                {t("profile.save_avatar", "Use this Avatar")}
            </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
