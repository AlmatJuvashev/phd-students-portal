import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

export function ModuleGuardModal({
  open,
  title,
  onConfirm,
  onClose,
}: {
  open: boolean;
  title: string;
  onConfirm: () => void;
  onClose: () => void;
}) {
  const { t: T } = useTranslation("common");
  return (
    <Modal open={open} onClose={onClose}>
      <div className="space-y-3">
        <h3 className="text-lg font-semibold">{title}</h3>
        <p className="text-sm text-muted-foreground">
          {T("module.guard_prompt", {
            defaultValue: "Ready to start the next module? You can unlock it now.",
          })}
        </p>
        <div className="flex justify-end gap-2 pt-2">
          <Button variant="secondary" onClick={onClose}>
            {T("common.cancel", { defaultValue: "Cancel" })}
          </Button>
          <Button onClick={onConfirm}>{T("module.unlock", { defaultValue: "Unlock Module" })}</Button>
        </div>
      </div>
    </Modal>
  );
}

export default ModuleGuardModal;

