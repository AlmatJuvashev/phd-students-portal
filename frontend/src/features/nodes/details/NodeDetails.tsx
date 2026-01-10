// components/node-details/NodeDetails.tsx
import { NodeVM, t } from "@/lib/playbook";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Pencil, X } from "lucide-react";
import { NodeDetailSwitch } from "./NodeDetailSwitch";
import { NodeAttachmentsSection } from "./NodeAttachmentsSection";
import { useRef, useState } from "react";
import { useSubmission } from "@/features/journey/hooks";
import { useTranslation } from "react-i18next";
import { useNodeDetailActions, useFocusOnOpen } from "./useNodeDetailActions";

export function NodeDetails({
  node,
  role = "student",
  onStateRefresh,
  onAdvance,
  closeOnComplete = false,
  isPreview = false,
}: {
  node: NodeVM;
  role?: "student" | "advisor" | "secretary" | "chair" | "admin";
  onStateRefresh?: () => void;
  onAdvance?: (nextNodeId: string | null, currentNodeId: string | null) => void;
  closeOnComplete?: boolean;
  isPreview?: boolean;
}) {
  const { t: T } = useTranslation("common");
  const [saving, setSaving] = useState(false);
  const [editing, setEditing] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const titleRef = useRef<HTMLDivElement | null>(null);
  
  // Conditionally use submission hook (or pass null/skip)
  const realSubmission = useSubmission(isPreview ? null : (node?.id || null));
  
  const submission = isPreview ? { state: 'not_started', slots: [] } : realSubmission.submission;
  const isLoading = isPreview ? false : realSubmission.isLoading;
  const save = isPreview ? { mutateAsync: async () => {} } : realSubmission.save;
  const refetch = isPreview ? async () => {} : realSubmission.refetch;

  useFocusOnOpen(titleRef, node?.id ?? null);

  // Stub onOpenChange since we are not controlling a sheet anymore, but actions might need it
  const onOpenChange = (open: boolean) => {
      // no-op or handle close?
  };

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

  const stateLabel = (n: NodeVM) => {
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

  const roleAllowed = !!node?.who_can_complete?.includes(role as any);

  return (
    <div className="space-y-4">
        {/* Header Actions for Editing */}
        <div className="flex justify-end gap-2 mb-2">
            {!import.meta.env.PROD && (
                <Badge variant="secondary" className="capitalize text-xs">
                    {node.type}
                </Badge>
            )}
            <Badge className="capitalize text-xs" variant={node.state === 'done' ? 'default' : 'outline'}>
                {stateLabel(node)}
            </Badge>
            
            {["form", "confirmTask"].includes(node.type) &&
                ["submitted", "done"].includes((submission as any)?.state as any) && (
                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 px-2 text-xs"
                        onClick={() => setEditing(!editing)}
                    >
                        {editing ? (
                            <>
                                <X className="h-3 w-3 mr-1" /> {T("common.cancel", "Cancel")}
                            </>
                        ) : (
                            <>
                                <Pencil className="h-3 w-3 mr-1" /> {T("common.edit", "Edit")}
                            </>
                        )}
                    </Button>
            )}
        </div>

      {errorMsg && (
        <div
          role="alert"
          aria-live="polite"
          className="rounded-lg border-2 border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive shadow-sm"
        >
          <div className="flex items-start gap-2">
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
            onAttachmentsRefresh={() => refetch()}
          />
          
          {node?.type !== "confirmTask" &&
            submission?.slots &&
            submission.slots.length > 0 && (
            <div className="pt-6 border-t border-border mt-6">
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
    </div>
  );
}
