import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { useFormContext } from "./FormProvider";
import { useState } from "react";
import { ConfirmModal } from "./ConfirmModal";
import { t } from "@/lib/playbook";

export interface FormActionsProps {
  showConfirm?: boolean;
  submitLabel?: string;
  draftLabel?: string;
  hideSubmit?: boolean;
  hideDraft?: boolean;
  className?: string;
}

/**
 * Standard form action buttons with unified state management
 * Uses FormProvider context for submit/draft logic
 */
export function FormActions({
  showConfirm = false,
  submitLabel,
  draftLabel,
  hideSubmit = false,
  hideDraft = false,
  className = "",
}: FormActionsProps) {
  const { t: T } = useTranslation("common");
  const { node, canEdit, disabled, isValid, submit, saveDraft } = useFormContext();
  const [confirmOpen, setConfirmOpen] = useState(false);

  if (!canEdit) return null;

  const handleSubmitClick = () => {
    if (showConfirm) {
      setConfirmOpen(true);
    } else {
      submit();
    }
  };

  const handleConfirm = () => {
    setConfirmOpen(false);
    submit();
  };

  return (
    <>
      <div className={`flex gap-2 pt-4 border-t bg-background/80 backdrop-blur-sm sticky bottom-0 z-10 ${className}`}>
        {!hideSubmit && (
          <Button
            onClick={handleSubmitClick}
            disabled={!isValid || disabled}
            aria-busy={disabled}
            className="touch-manipulation min-h-[44px]"
          >
            {submitLabel || T("forms.proceed_next", { defaultValue: "Submit" })}
          </Button>
        )}
        {!hideDraft && (
          <Button
            variant="secondary"
            onClick={() => saveDraft()}
            disabled={disabled}
            aria-busy={disabled}
            className="touch-manipulation min-h-[44px]"
          >
            {draftLabel || T("forms.save_draft", { defaultValue: "Save Draft" })}
          </Button>
        )}
      </div>

      {showConfirm && (
        <ConfirmModal
          open={confirmOpen}
          onOpenChange={setConfirmOpen}
          message={t((node as any).description || node.title, "")}
          confirmLabel={submitLabel || T("forms.proceed_next", { defaultValue: "Submit" })}
          cancelLabel={T("common.cancel", { defaultValue: "Cancel" })}
          onConfirm={handleConfirm}
        />
      )}
    </>
  );
}
