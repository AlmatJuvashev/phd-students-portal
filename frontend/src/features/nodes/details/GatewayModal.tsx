import { NodeVM, t } from "@/lib/playbook";
import { Modal } from "@/components/ui/modal";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

export function GatewayModal({
  node,
  open,
  onClose,
  onUnlock,
}: {
  node: NodeVM | null;
  open: boolean;
  onClose: () => void;
  onUnlock: (node: NodeVM) => void;
}) {
  const { t: T } = useTranslation("common");
  if (!node) return null;
  return (
    <Modal open={open} onClose={onClose}>
      <div className="space-y-3">
        <h3 className="text-lg font-semibold">{t(node.title, node.id)}</h3>
        <p className="text-sm text-muted-foreground">
          {T("gateway.prompt", {
            defaultValue: "Proceed to the next chapter? You can unlock it now.",
          })}
        </p>
        <div className="flex justify-end gap-2 pt-2">
          <Button variant="secondary" onClick={onClose}>
            {T("common.cancel", { defaultValue: "Cancel" })}
          </Button>
          <Button onClick={() => onUnlock(node)}>
            {T("gateway.unlock", { defaultValue: "Unlock Chapter" })}
          </Button>
        </div>
      </div>
    </Modal>
  );
}

