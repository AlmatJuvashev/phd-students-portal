// components/node-details/NodeDetailSwitch.tsx
import { NodeVM } from "@/lib/playbook";
import type { NodeSubmissionDTO } from "@/api/journey";
import ConfirmTaskDetails from "./variants/ConfirmTaskDetails";
import InfoDetails from "./variants/InfoDetails";
import { deriveNodeKind } from "@/features/nodes/deriveNodeKind";
import FormEntryDetails from "@/features/nodes/kinds/FormEntryDetails";
import ChecklistDetails from "@/features/nodes/kinds/ChecklistDetails";
import CardsDetails from "@/features/nodes/kinds/CardsDetails";
import React from "react";
import D2ApplyToDsScene from "@/features/nodes/scenes/D2ApplyToDsScene";
import V1ReinstatementScene from "@/features/nodes/scenes/V1ReinstatementScene";
import RP2HearingPrepScene from "@/features/nodes/scenes/RP2HearingPrepScene";
import VIAttestationScene from "@/features/nodes/scenes/VIAttestationScene";
import { FormTaskDetails } from "./variants/FormTaskDetails";
import useGuideForNode from "@/features/guides/useGuideForNode";

type Props = {
  node: NodeVM;
  onEvent?: (evt: { type: string; payload?: any }) => void; // bubble up submit/finalize/etc.
  submission?: NodeSubmissionDTO | null;
  saving?: boolean;
  canEdit?: boolean;
};

export function NodeDetailSwitch({
  node,
  onEvent,
  submission,
  saving = false,
  canEdit,
}: Props) {
  const uiKind = deriveNodeKind(node);
  const renderGuide = useGuideForNode(node) || undefined;
  const initialForm = submission?.form?.data ?? {};

  // Explicit scenes (specialized UX) by node id
  if (node.id === "D2_apply_to_ds") {
    return (
      <D2ApplyToDsScene
        node={node}
        initial={initialForm}
        disabled={saving}
        onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
      />
    );
  }
  if (node.id === "V1_reinstatement_package") {
    return (
      <V1ReinstatementScene
        node={node}
        initial={initialForm}
        disabled={saving}
        onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
      />
    );
  }
  if (node.id === "RP2_sc_hearing_prep") {
    return (
      <RP2HearingPrepScene
        node={node}
        initial={initialForm}
        disabled={saving}
        onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
      />
    );
  }
  if (node.id === "VI_attestation_file") {
    return (
      <VIAttestationScene
        node={node}
        initial={initialForm}
        disabled={saving}
        onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
      />
    );
  }

  // Force E3_hearing_nk to use the specialized flow in FormTaskDetails (yes/no cards, back navigation)
  if (node.id === "E3_hearing_nk") {
    return (
      <FormTaskDetails
        node={node}
        canEdit={canEdit ?? !saving}
        initial={initialForm}
        disabled={saving}
        onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
      />
    );
  }

  // permissions (rough defaults, adjust as you wire real RBAC)
  // Prefer UI-specific kinds first for form-like nodes
  if (node.type === "form") {
    const initialForm = submission?.form?.data ?? {};
    if (uiKind === "formEntry") {
      return (
        <FormEntryDetails
          node={node}
          initial={initialForm}
          disabled={saving || canEdit === false}
          onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
        />
      );
    }
    if (uiKind === "checklist") {
      return (
        <ChecklistDetails
          node={node}
          initial={initialForm}
          disabled={saving}
          onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
        />
      );
    }
    if (uiKind === "cards") {
      return (
        <CardsDetails
          node={node}
          disabled={saving}
          onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
        />
      );
    }
  }

  if (node.type === "confirmTask") {
    const nextOverride = (node as any)?.states?.completed?.next_node;
    return (
      <ConfirmTaskDetails
        node={node}
        onComplete={() =>
          onEvent?.({
            type: "submit-decision",
            payload: nextOverride ? { __nextOverride: nextOverride } : {},
          })
        }
        onReset={() => onEvent?.({ type: "reset-node" })}
      />
    );
  }

  if (node.type === "info") {
    return (
      <InfoDetails
        node={node}
        renderGuide={renderGuide}
        onContinue={() => onEvent?.({ type: "continue" })}
      />
    );
  }

  // Default to a simple form view when type is unrecognised.
  return (
    <FormTaskDetails
      node={node}
      canEdit={canEdit ?? !saving}
      initial={initialForm}
      disabled={saving}
      onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
    />
  );
}
