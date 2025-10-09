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
import { useEffect, useState } from "react";
import {
  getNodeSubmission,
  NodeSubmissionDTO,
  saveNodeSubmission,
} from "@/api/journey";
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
  const [submission, setSubmission] = useState<NodeSubmissionDTO | null>(null);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!node) {
      setSubmission(null);
      setLoading(false);
      return;
    }
    let cancelled = false;
    const run = async () => {
      setLoading(true);
      try {
        const data = await getNodeSubmission(node.id);
        if (!cancelled) {
          setSubmission(data);
        }
      } catch (err: any) {
        if (!cancelled) {
          // toast removed -> log error instead
          console.error(
            T("common.error", { defaultValue: "Error" }),
            err?.message ?? String(err)
          );
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };
    run();
    return () => {
      cancelled = true;
    };
  }, [node?.id, T]);

  const handleEvent = async (evt: { type: string; payload?: any }) => {
    if (!node || saving) return;
    switch (evt.type) {
      case "submit-form": {
        const payload = { ...(evt.payload ?? {}) };
        const isDraft = !!payload.__draft;
        delete payload.__draft;
        setSaving(true);
        try {
          const res = await saveNodeSubmission(node.id, {
            form_data: payload,
            state: isDraft ? "active" : "submitted",
          });
          setSubmission(res);
          // toast removed -> optionally log success
          console.info(
            isDraft ? T("forms.save_draft") : T("forms.save_submit"),
            T("common.success", { defaultValue: "Saved." })
          );
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
        } finally {
          setSaving(false);
        }
        break;
      }
      case "submit-decision": {
        setSaving(true);
        try {
          const res = await saveNodeSubmission(node.id, {
            form_data: evt.payload ?? {},
            state: "submitted",
          });
          setSubmission(res);
          // toast removed -> optionally log success
          console.info(
            T("decision.submit"),
            T("common.success", { defaultValue: "Saved." })
          );
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
              <SheetTitle className="flex items-center gap-2">
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
              {loading ? (
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
