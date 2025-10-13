import { useCallback, useEffect, type RefObject } from "react";
import { NodeVM } from "@/lib/playbook";
import { patchJourneyState } from "@/features/journey/session";
import { api } from "@/api/client";
import { useConditions } from "@/features/journey/useConditions";

export type SaveMutation = {
  mutateAsync: (payload: {
    form_data: Record<string, any>;
    state: string;
  }) => Promise<unknown>;
};

export type NodeDetailActionArgs = {
  node: NodeVM | null;
  saving: boolean;
  setSaving: (value: boolean) => void;
  save: SaveMutation;
  onStateRefresh?: () => void;
  onOpenChange: (open: boolean) => void;
  onAdvance?: (nextId: string | null, currentNodeId: string | null) => void;
  setErrorMsg: (msg: string | null) => void;
  closeOnComplete?: boolean;
};

const resolveNextNode = (
  node: NodeVM | null,
  rpRequired: boolean,
  override?: string
): string | null => {
  if (!node) return null;
  if (override) return override;
  if (!node.next || !Array.isArray(node.next) || node.next.length === 0) {
    return null;
  }
  if (node.condition === "rp_required" && node.next.length >= 2) {
    return rpRequired ? node.next[0] : node.next[1];
  }
  return node.next[0] ?? null;
};

export function useNodeDetailActions({
  node,
  saving,
  setSaving,
  save,
  onStateRefresh,
  onOpenChange,
  onAdvance,
  setErrorMsg,
  closeOnComplete = false,
}: NodeDetailActionArgs) {
  const { rp_required } = useConditions();

  const handleEvent = useCallback(
    async (evt: { type: string; payload?: any }) => {
      if (!node || saving) return;

      const onComplete = async (nextOverride?: string) => {
        patchJourneyState({ [node.id]: "done" });
        try {
          await api("/journey/state", {
            method: "PUT",
            body: JSON.stringify({ node_id: node.id, state: "done" }),
          });
        } catch (error) {
          console.warn("state upsert failed", error);
        }
        onStateRefresh?.();
        if (closeOnComplete) {
          onOpenChange(false);
        }
        onAdvance?.(
          resolveNextNode(node, !!rp_required, nextOverride),
          node?.id ?? null
        );
      };

      switch (evt.type) {
        case "reset-node": {
          setSaving(true);
          try {
            patchJourneyState({ [node.id]: "active" });
            await api("/journey/state", {
              method: "PUT",
              body: JSON.stringify({ node_id: node.id, state: "active" }),
            });
            onStateRefresh?.();
          } catch (error: any) {
            console.error("reset node failed", error);
            setErrorMsg(error?.message ?? String(error));
          } finally {
            setSaving(false);
          }
          break;
        }
        case "continue": {
          setSaving(true);
          try {
            await save.mutateAsync({ form_data: {}, state: "done" });
            setErrorMsg(null);
            await onComplete(evt.payload?.__nextOverride);
          } catch (error: any) {
            console.error("continue failed", error);
            setErrorMsg(error?.message ?? String(error));
          } finally {
            setSaving(false);
          }
          break;
        }
        case "submit-form": {
          const payload = { ...(evt.payload ?? {}) };
          const isDraft = !!payload.__draft;
          const nextOverride = payload.__nextOverride;
          // Remove metadata keys before saving to backend
          if ("__draft" in payload) delete payload.__draft;
          if ("__nextOverride" in payload) delete payload.__nextOverride;
          if ("__submittedAt" in payload) delete payload.__submittedAt;
          setSaving(true);
          try {
            await save.mutateAsync({
              form_data: payload,
              state: isDraft ? "active" : "done",
            });
            setErrorMsg(null);
            if (isDraft) {
              setSaving(false);
              return;
            }
            await onComplete(nextOverride);
          } catch (error: any) {
            console.error("submit form failed", error);
            setErrorMsg(error?.message ?? String(error));
          } finally {
            setSaving(false);
          }
          break;
        }
        case "submit-decision": {
          setSaving(true);
          try {
            await save.mutateAsync({
              form_data: evt.payload ?? {},
              state: "done",
            });
            setErrorMsg(null);
            await onComplete(evt.payload?.__nextOverride);
          } catch (error: any) {
            console.error("decision submit failed", error);
            setErrorMsg(error?.message ?? String(error));
          } finally {
            setSaving(false);
          }
          break;
        }
        default:
          break;
      }
    },
    [
      node,
      saving,
      setSaving,
      save,
      onStateRefresh,
      onOpenChange,
      onAdvance,
      setErrorMsg,
      rp_required,
    ]
  );

  return { handleEvent };
}

export function useFocusOnOpen<T extends HTMLElement>(
  ref: RefObject<T>,
  dependencyKey: string | null
) {
  useEffect(() => {
    if (dependencyKey && ref.current) {
      ref.current.focus();
    }
  }, [dependencyKey, ref]);
}
