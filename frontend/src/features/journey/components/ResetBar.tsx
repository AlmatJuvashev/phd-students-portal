import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useJourneyState } from "../hooks";
import { useTranslation } from "react-i18next";
import { AlertTriangle } from "lucide-react";

export function ResetBar() {
  const { t: T } = useTranslation("common");
  const { reset } = useJourneyState();
  const [isOpen, setIsOpen] = useState(false);
  const [confirmText, setConfirmText] = useState("");

  // The phrase to type for confirmation (internationalized)
  const confirmPhrase = T("journey.reset_confirm_phrase", { defaultValue: "reset map" });
  const isConfirmValid = confirmText.toLowerCase().trim() === confirmPhrase.toLowerCase().trim();

  const handleReset = () => {
    if (isConfirmValid) {
      reset();
      setIsOpen(false);
      setConfirmText("");
    }
  };

  const handleClose = () => {
    setIsOpen(false);
    setConfirmText("");
  };

  return (
    <>
      <div className="fixed inset-x-0 bottom-0 z-40">
        <div className="mx-auto max-w-6xl px-4 pb-4">
          <div className="rounded-xl bg-muted/70 backdrop-blur border p-3 flex items-center justify-between">
            <div className="text-xs text-muted-foreground">
              {T("journey.reset_hint", { defaultValue: "You can reset your journey to start over." })}
            </div>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => setIsOpen(true)}
            >
              {T("journey.reset", { defaultValue: "Reset Journey" })}
            </Button>
          </div>
        </div>
      </div>

      <Dialog open={isOpen} onOpenChange={handleClose}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-destructive">
              <AlertTriangle className="h-5 w-5" />
              {T("journey.reset_title", { defaultValue: "Reset Journey" })}
            </DialogTitle>
            <DialogDescription>
              {T("journey.reset_warning", { 
                defaultValue: "This action cannot be undone. All your progress will be lost and you will start from the beginning." 
              })}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
              <p className="text-sm text-destructive">
                {T("journey.reset_type_to_confirm", { 
                  defaultValue: "To confirm, type \"{{phrase}}\" below:",
                  phrase: confirmPhrase
                })}
              </p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirm-input" className="sr-only">
                {T("journey.reset_confirm_label", { defaultValue: "Confirmation" })}
              </Label>
              <Input
                id="confirm-input"
                value={confirmText}
                onChange={(e) => setConfirmText(e.target.value)}
                placeholder={confirmPhrase}
                className="font-mono"
                autoComplete="off"
              />
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-0">
            <Button variant="outline" onClick={handleClose}>
              {T("common.cancel", { defaultValue: "Cancel" })}
            </Button>
            <Button
              variant="destructive"
              onClick={handleReset}
              disabled={!isConfirmValid}
            >
              {T("journey.reset_confirm_button", { defaultValue: "Reset My Journey" })}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
