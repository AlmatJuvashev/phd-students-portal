// components/node-details/NodeDetailsSheet.tsx

import { NodeVM, t } from "@/lib/playbook";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { NodeDetailSwitch } from "./NodeDetailSwitch";
import { useRef, useState } from "react";
import { useSubmission } from "@/features/journey/hooks";
import { useTranslation } from "react-i18next";
import { useNodeDetailActions, useFocusOnOpen } from "./useNodeDetailActions";

export function NodeDetailsSheet({
  node,
  onOpenChange,
  role = "student",
  onStateRefresh,
  onAdvance,
  closeOnComplete = false,
}: {
  node: NodeVM | null;
  onOpenChange: (open: boolean) => void;
  role?: "student" | "advisor" | "secretary" | "chair" | "admin";
  onStateRefresh?: () => void;
  onAdvance?: (nextNodeId: string | null, currentNodeId: string | null) => void;
  closeOnComplete?: boolean;
}) {
  const { t: T } = useTranslation("common");
  const [saving, setSaving] = useState(false);
  const [editing, setEditing] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const titleRef = useRef<HTMLDivElement | null>(null);
  const { submission, isLoading, save } = useSubmission(node?.id || null);

  useFocusOnOpen(titleRef, node?.id ?? null);

  const { handleEvent } = useNodeDetailActions({
    node,
    saving,
    setSaving,
    save,
    onStateRefresh,
    onOpenChange,
    onAdvance,
    setErrorMsg,
    closeOnComplete,
  });

  return (
    <Sheet open={!!node} onOpenChange={onOpenChange}>
      <SheetContent
        side="right"
        className="w-full max-w-full sm:max-w-6xl p-0 flex flex-col overflow-hidden bg-gradient-to-br from-background via-background to-muted/10 border-l-2 border-primary/20 shadow-2xl"
      >
        {node && (
          <>
            <SheetHeader className="px-6 py-5 border-b border-border/50 bg-card/80 backdrop-blur-md sticky top-0 z-10">
              <div className="flex flex-col sm:flex-row sm:items-start gap-3">
                <div className="flex-1 min-w-0">
                  <SheetTitle
                    ref={titleRef as any}
                    tabIndex={-1}
                    className="text-xl sm:text-2xl font-bold outline-none bg-gradient-to-r from-primary via-primary/90 to-primary/70 bg-clip-text text-transparent leading-tight pr-2"
                  >
                    {t(node.title, node.id)}
                  </SheetTitle>
                  {(node as any).description && (
                    <p className="text-sm text-muted-foreground mt-2 line-clamp-2">
                      {t((node as any).description, "")}
                    </p>
                  )}
                </div>
                <div className="flex items-center gap-2 flex-shrink-0">
                  <Badge
                    variant="secondary"
                    className="capitalize shadow-sm hover:shadow transition-shadow"
                  >
                    {node.type}
                  </Badge>
                  <Badge className="capitalize shadow-sm hover:shadow transition-shadow">
                    {node.state?.replace("_", " ")}
                  </Badge>
                  {node.type === "form" &&
                    (submission as any)?.state === "submitted" &&
                    (!editing ? (
                      <button
                        className="ml-1 text-xs font-medium text-primary hover:text-primary/80 underline underline-offset-2 transition-colors"
                        onClick={() => setEditing(true)}
                      >
                        {T("common.edit", { defaultValue: "Edit" })}
                      </button>
                    ) : (
                      <button
                        className="ml-1 text-xs font-medium text-muted-foreground hover:text-foreground underline underline-offset-2 transition-colors"
                        onClick={() => setEditing(false)}
                      >
                        {T("common.cancel_edit", { defaultValue: "Cancel" })}
                      </button>
                    ))}
                </div>
              </div>
            </SheetHeader>

            <div className="flex-1 min-h-0 overflow-y-auto px-6 py-5 space-y-4">
              {errorMsg && (
                <div
                  role="alert"
                  aria-live="polite"
                  className="rounded-lg border-2 border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive shadow-sm animate-in fade-in slide-in-from-top-2 duration-300"
                >
                  <div className="flex items-start gap-2">
                    <svg
                      className="h-5 w-5 shrink-0 mt-0.5"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                        clipRule="evenodd"
                      />
                    </svg>
                    <span className="font-medium">{errorMsg}</span>
                  </div>
                </div>
              )}
              {isLoading ? (
                <div className="flex flex-col items-center justify-center py-12 space-y-3">
                  <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary"></div>
                  <p className="text-sm text-muted-foreground animate-pulse">
                    {T("common.loading")}
                  </p>
                </div>
              ) : (
                <NodeDetailSwitch
                  node={node}
                  submission={submission as any}
                  onEvent={handleEvent}
                  saving={saving}
                  canEdit={
                    editing ||
                    !["submitted", "done"].includes(
                      (submission as any)?.state as any
                    )
                  }
                />
              )}
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
