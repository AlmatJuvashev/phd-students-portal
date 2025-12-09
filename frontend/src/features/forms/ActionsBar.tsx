import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

export function ActionsBar({
  onSubmit,
  onDraft,
  disabled,
}: {
  onSubmit: () => void;
  onDraft: () => void;
  disabled?: boolean;
}) {
  const { t: T } = useTranslation("common");
  return (
    <div className="flex flex-col-reverse sm:flex-row gap-2">
      <Button onClick={onSubmit} disabled={disabled} aria-busy={disabled} className="w-full sm:w-auto" data-testid="node-submit-button">
        {T("forms.save_submit")}
      </Button>
      <Button variant="secondary" onClick={onDraft} disabled={disabled} aria-busy={disabled} className="w-full sm:w-auto">
        {T("forms.save_draft")}
      </Button>
    </div>
  );
}

