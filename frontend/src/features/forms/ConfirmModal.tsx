import * as Dialog from "@radix-ui/react-dialog";
import { Button } from "@/components/ui/button";

export function ConfirmModal({
  open,
  onOpenChange,
  message,
  onConfirm,
  confirmLabel,
  cancelLabel,
  busy,
}: {
  open: boolean;
  onOpenChange: (o: boolean) => void;
  message: string;
  onConfirm: () => void;
  confirmLabel: string;
  cancelLabel: string;
  busy?: boolean;
}) {
  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50 z-[70]" />
        <Dialog.Content className="fixed z-[70] left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg p-6 w-full max-w-md shadow-lg outline-none">
          <div className="mb-4 text-sm text-muted-foreground whitespace-pre-line">
            {message}
          </div>
          <div className="flex gap-2 justify-end">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              {cancelLabel}
            </Button>
            <Button onClick={onConfirm} aria-busy={!!busy}>
              {confirmLabel}
            </Button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}

