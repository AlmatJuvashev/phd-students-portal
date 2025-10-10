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
            onStateRefresh?.();
            const nextId = Array.isArray(node.next) ? node.next[0] : undefined;
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
      <SheetContent side="right" className="w-full max-w-full sm:max-w-6xl">
        {node && (
          <>
            <SheetHeader>
              <SheetTitle
                ref={titleRef as any}
                tabIndex={-1}
                className="flex items-center gap-2 outline-none"
              >
                <span>{t(node.title, node.id)}</span>
                <Badge variant="secondary" className="capitalize">
                  {node.type}
                </Badge>
                <Badge className="capitalize">
                  {node.state?.replace("_", " ")}
                </Badge>
              </SheetTitle>
            </SheetHeader>

            <div className="mt-6 h-[calc(100vh-8rem)] overflow-hidden">
              {errorMsg && (
                <div
                  role="alert"
                  aria-live="polite"
                  className="mb-3 rounded-md border border-red-300 bg-red-50 p-3 text-sm text-red-800"
                >
                  {errorMsg}
                </div>
              )}
              {isLoading ? (
                <div className="text-sm text-muted-foreground">
                  {T("common.loading")}
                </div>
              ) : (
                <NodeDetailSwitch
                  node={node}
                  role={role}
                  submission={submission}
                  onEvent={handleEvent}
                  saving={saving}
                />
              )}
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
