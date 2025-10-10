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
import { useEffect, useRef, useState } from "react";
import { NodeSubmissionDTO } from "@/api/journey";
import { useSubmission } from "@/features/journey/hooks";
import { useTranslation } from "react-i18next";
import { patchJourneyState } from "@/features/journey/session";

export function NodeDetailsSheet({
  node,
  onOpenChange,
  role = "student",
  onStateRefresh,
  onAdvance,
}: {
  node: NodeVM | null;
  onOpenChange: (open: boolean) => void;
  role?: "student" | "advisor" | "secretary" | "chair" | "admin";
  onStateRefresh?: () => void;
  onAdvance?: (nextNodeId: string | null) => void;
}) {
  const { t: T } = useTranslation("common");
  const [saving, setSaving] = useState(false);
  const [editing, setEditing] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const titleRef = useRef<HTMLDivElement | null>(null);
  const { submission, isLoading, save } = useSubmission(node?.id || null);

  useEffect(() => {
    // focus the title on open
    if (node && titleRef.current) titleRef.current.focus();
  }, [node?.id]);

  useEffect(() => {
    if (node && titleRef.current) {
      // focus the title for screen readers when sheet opens
      titleRef.current.focus();
    }
  }, [node?.id]);

  const handleEvent = async (evt: { type: string; payload?: any }) => {
    if (!node || saving) return;
    switch (evt.type) {
      case "submit-form": {
        const payload = { ...(evt.payload ?? {}) };
        const isDraft = !!payload.__draft;
        delete payload.__draft;
        const nextOverride: string | undefined =
          (evt.payload && evt.payload.__nextOverride) || undefined;
        setSaving(true);
        try {
          const res = await save.mutateAsync({
            form_data: payload,
            state: isDraft ? "active" : "submitted",
          });
          // toast removed -> optionally log success
          console.info(
            isDraft ? T("forms.save_draft") : T("forms.save_submit"),
            T("common.success", { defaultValue: "Saved." })
          );
          setErrorMsg(null);
          if (!isDraft) {
            // persist session progress for this node
            patchJourneyState({ [node.id]: "submitted" });
            onStateRefresh?.();
            const nextId =
              nextOverride ||
              (Array.isArray(node.next) ? node.next[0] : undefined);
            onOpenChange(false);
            if (nextId) {
              onAdvance?.(nextId);
            } else {
              onAdvance?.(null);
            }
          }
        } catch (err: any) {
          // toast removed -> log error
          console.error(
            T("common.error", { defaultValue: "Error" }),
            err?.message ?? String(err)
          );
          setErrorMsg(err?.message ?? String(err));
        } finally {
          setSaving(false);
        }
        break;
      }
      case "submit-decision": {
        setSaving(true);
        try {
          const res = await save.mutateAsync({
            form_data: evt.payload ?? {},
            state: "submitted",
          });
          // toast removed -> optionally log success
          console.info(
            T("decision.submit"),
            T("common.success", { defaultValue: "Saved." })
          );
          setErrorMsg(null);
          onStateRefresh?.();
          const nextId = Array.isArray(node.next) ? node.next[0] : undefined;
          onOpenChange(false);
          if (nextId) {
            onAdvance?.(nextId);
          } else {
            onAdvance?.(null);
          }
        } catch (err: any) {
          console.error("submit decision failed", err);
          setErrorMsg(err?.message ?? String(err));
        } finally {
          setSaving(false);
        }
        break;
      }
      case "submit-upload": {
        // deferred; no toast
        break;
      }
      case "finalize-composite":
      case "finalize-outcome":
        // not implemented; no toast
        break;
      default:
        break;
    }
  };

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

            <div className="flex-1 overflow-y-auto px-6 py-5 space-y-4">
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
                  role={role}
                  submission={submission as any}
                  onEvent={handleEvent}
                  saving={saving}
                  canEdit={
                    editing || (submission as any)?.state !== "submitted"
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
