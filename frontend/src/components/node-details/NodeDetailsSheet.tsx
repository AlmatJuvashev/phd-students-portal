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
import { useEffect, useState, useCallback } from "react";
import {
  getNodeSubmission,
  NodeSubmissionDTO,
  saveNodeSubmission,
} from "@/api/journey";
import { useToast } from "@/components/toast";
import { useTranslation } from "react-i18next";

export function NodeDetailsSheet({
  node,
  onOpenChange,
  role = "student",
  onStateRefresh,
}: {
  node: NodeVM | null;
  onOpenChange: (open: boolean) => void;
  role?: "student" | "advisor" | "secretary" | "chair" | "admin";
  onStateRefresh?: () => void;
}) {
  const { t: T } = useTranslation("common");
  const [submission, setSubmission] = useState<NodeSubmissionDTO | null>(null);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const { push } = useToast();

  const loadSubmission = useCallback(
    async (id: string) => {
      setLoading(true);
      try {
        const data = await getNodeSubmission(id);
        setSubmission(data);
      } catch (err: any) {
        push({
          title: T("common.error", { defaultValue: "Error" }),
          description: err?.message ?? String(err),
        });
      } finally {
        setLoading(false);
      }
    },
    [T, push],
  );

  useEffect(() => {
    if (node) {
      loadSubmission(node.id);
    } else {
      setSubmission(null);
      setLoading(false);
    }
  }, [node?.id, loadSubmission]);

  const handleEvent = useCallback(
    async (evt: { type: string; payload?: any }) => {
      if (!node) return;
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
            push({
              title: isDraft ? T("forms.save_draft") : T("forms.save_submit"),
              description: T("common.success", { defaultValue: "Saved." }),
            });
            if (!isDraft) {
              onStateRefresh?.();
            }
          } catch (err: any) {
            push({
              title: T("common.error", { defaultValue: "Error" }),
              description: err?.message ?? String(err),
            });
          } finally {
            setSaving(false);
          }
          break;
        }
        case "submit-upload": {
          push({
            title: T("common.info", { defaultValue: "Info" }),
            description: T("upload.not_supported", {
              defaultValue: "File uploads will be available soon.",
            }),
          });
          break;
        }
        case "finalize-composite":
        case "finalize-outcome":
          push({
            title: T("common.info", { defaultValue: "Info" }),
            description: T("common.not_implemented", {
              defaultValue: "Action not yet available.",
            }),
          });
          break;
        default:
          break;
      }
    },
    [node, onStateRefresh, push, T],
  );

  return (
    <Sheet open={!!node} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full max-w-full sm:max-w-lg">
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

            <div className="mt-6 min-h-[120px]">
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
