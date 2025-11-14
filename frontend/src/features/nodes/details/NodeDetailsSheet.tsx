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
import { NodeAttachmentsSection } from "./NodeAttachmentsSection";
import { useEffect, useRef, useState } from "react";
import { useSubmission } from "@/features/journey/hooks";
import { useTranslation } from "react-i18next";
import { useNodeDetailActions, useFocusOnOpen } from "./useNodeDetailActions";
import { useSwipeToClose } from "@/hooks/useSwipeToClose";
import { AnimatePresence, motion } from "framer-motion";

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
  const { submission, isLoading, save, refetch } = useSubmission(
    node?.id || null,
  );

  // Detect mobile for bottom sheet pattern
  const [isMobile, setIsMobile] = useState(false);
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 640); // sm breakpoint
    };
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  // Enable swipe-to-close on mobile bottom sheet
  useSwipeToClose({
    onClose: () => onOpenChange(false),
    enabled: isMobile && !!node,
    threshold: 80,
  });

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

  const stateLabel = (n: NodeVM | null) => {
    if (!n) return "";
    const state = n.state || "";
    if (n.type === "confirmTask" && state === "done") {
      return t(
        { ru: "Шаг подтверждён", kz: "Қадам расталды", en: "Step confirmed" },
        "Шаг подтверждён"
      );
    }
    if (state === "active") {
      return t({ ru: "Активно", kz: "Белсенді", en: "Active" }, "Активно");
    }
    if (state === "submitted") {
      return t(
        { ru: "Отправлено", kz: "Жіберілді", en: "Submitted" },
        "Отправлено"
      );
    }
    if (state === "done") {
      return t({ ru: "Готово", kz: "Дайын", en: "Done" }, "Готово");
    }
    return (n.state || "").replace("_", " ");
  };

  const roleAllowed = !!node?.who_can_complete?.includes(role as any)

  return (
    <Sheet open={!!node} onOpenChange={onOpenChange}>
      <SheetContent
        side={isMobile ? "bottom" : "right"}
        className={`p-0 flex flex-col overflow-hidden bg-gradient-to-br from-background via-background to-muted/10 shadow-2xl
          ${
            isMobile
              ? "w-full h-[90vh] max-h-[90vh] rounded-t-3xl border-t-2 border-primary/20"
              : "w-full max-w-full sm:max-w-6xl border-l-2 border-primary/20"
          }`}
      >
        {node && (
          <>
            {/* Drag handle for mobile bottom sheet */}
            {isMobile && (
              <div
                data-drag-handle
                className="flex justify-center pt-3 pb-1 cursor-grab active:cursor-grabbing"
                aria-label="Drag to close"
              >
                <div className="w-12 h-1.5 bg-muted-foreground/30 rounded-full" />
              </div>
            )}
            <SheetHeader
              data-sheet-header
              className={`px-4 sm:px-6 border-b border-border/50 bg-card/80 backdrop-blur-md sticky top-0 z-10 ${
                isMobile ? "py-3" : "py-5"
              }`}
            >
              <div className="flex flex-col gap-3">
                {/* Title section */}
                <div className="flex-1 min-w-0 pr-10">
                  <SheetTitle
                    ref={titleRef as any}
                    tabIndex={-1}
                    className="text-lg sm:text-xl md:text-2xl font-bold outline-none bg-gradient-to-r from-primary via-primary/90 to-primary/70 bg-clip-text text-transparent leading-tight"
                  >
                    {t(node.title, node.id)}
                  </SheetTitle>
                  {(node as any).description && (
                    <p className="text-xs sm:text-sm text-muted-foreground mt-1.5 sm:mt-2 line-clamp-2">
                      {t((node as any).description, "")}
                    </p>
                  )}
                </div>

                {/* Badges and actions - separate row for better spacing */}
                <div className="flex items-center gap-2 flex-wrap pr-10">
                  <Badge
                    variant="secondary"
                    className="capitalize shadow-sm hover:shadow transition-shadow text-xs"
                  >
                    {node.type}
                  </Badge>
                  <Badge className="capitalize shadow-sm hover:shadow transition-shadow text-xs">
                    {stateLabel(node)}
                  </Badge>
                  {node.type === "form" &&
                    ["submitted", "done"].includes(
                      (submission as any)?.state as any
                    ) &&
                    (!editing ? (
                      <button
                        className="text-xs font-medium text-primary hover:text-primary/80 underline underline-offset-2 transition-colors touch-manipulation min-h-[32px] px-2"
                        onClick={() => setEditing(true)}
                      >
                        {T("common.edit", { defaultValue: "Edit" })}
                      </button>
                    ) : (
                      <button
                        className="text-xs font-medium text-muted-foreground hover:text-foreground underline underline-offset-2 transition-colors touch-manipulation min-h-[32px] px-2"
                        onClick={() => setEditing(false)}
                      >
                        {T("common.cancel_edit", { defaultValue: "Cancel" })}
                      </button>
                    ))}
                </div>
              </div>
            </SheetHeader>

            <div
              data-sheet-content
              className="flex-1 min-h-0 overflow-y-auto px-4 sm:px-6 py-4 sm:py-5 overscroll-contain"
            >
              <AnimatePresence mode="wait" initial={false}>
                <motion.div
                  key={node.id}
                  initial={{ opacity: 0, x: 24 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -24 }}
                  transition={{ type: "tween", duration: 0.25, ease: "easeInOut" }}
                  className="space-y-4"
                >
                  {errorMsg && (
                    <div
                      role="alert"
                      aria-live="polite"
                      className="rounded-lg border-2 border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive shadow-sm"
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
                    <>
                      {!roleAllowed && (
                        <div className="rounded-md border border-amber-300 bg-amber-50 text-amber-900 p-3 text-sm mb-3">
                          <div className="font-medium mb-1">Доступ ограничен</div>
                          Только {node?.who_can_complete?.join(', ')} могут заполнить эту форму
                        </div>
                      )}
                      <NodeDetailSwitch
                        node={node}
                        submission={submission as any}
                        onEvent={handleEvent}
                        saving={saving}
                        canEdit={
                          roleAllowed &&
                          (editing ||
                            !["submitted", "done"].includes(
                              (submission as any)?.state as any
                            ))
                        }
                      />
                      {submission?.slots && submission.slots.length > 0 && (
                        <div className="pt-6">
                          <NodeAttachmentsSection
                            nodeId={node.id}
                            slots={submission.slots}
                            canEdit={roleAllowed}
                            onRefresh={() => refetch()}
                          />
                        </div>
                      )}
                    </>
                  )}
                </motion.div>
              </AnimatePresence>
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
