// components/node-details/NodeDetailSwitch.tsx
import { NodeVM, detectActionKinds } from "@/lib/playbook";
import { GatewayInfoDetails } from "./variants/GatewayInfoDetails";
import type { NodeSubmissionDTO } from "@/api/journey";
import { DecisionTaskDetails } from "./variants/DecisionTaskDetails";
import ConfirmTaskDetails from "./variants/ConfirmTaskDetails";
import InfoDetails from "./variants/InfoDetails";
import { deriveNodeKind } from "@/features/nodes/deriveNodeKind";
import FormEntryDetails from "@/features/nodes/kinds/FormEntryDetails";
import ChecklistDetails from "@/features/nodes/kinds/ChecklistDetails";
import CardsDetails from "@/features/nodes/kinds/CardsDetails";
import React, { Suspense } from "react";
import D2ApplyToDsScene from "@/features/nodes/scenes/D2ApplyToDsScene";
import V1ReinstatementScene from "@/features/nodes/scenes/V1ReinstatementScene";
import RP2HearingPrepScene from "@/features/nodes/scenes/RP2HearingPrepScene";
import VIAttestationScene from "@/features/nodes/scenes/VIAttestationScene";
import { FormTaskDetails } from "./variants/FormTaskDetails";

// Lazy heavy variants
const UploadTaskDetails = React.lazy(() => import("./variants/UploadTaskDetails").then(m => ({ default: m.UploadTaskDetails })));
const ConfirmUploadTaskDetails = React.lazy(() => import("./variants/ConfirmUploadTaskDetails").then(m => ({ default: m.default })));
const ExternalProcessDetails = React.lazy(() => import("./variants/ExternalProcessDetails").then(m => ({ default: m.ExternalProcessDetails })));

type Props = {
  node: NodeVM;
  role?: "student" | "advisor" | "secretary" | "chair" | "admin";
  onEvent?: (evt: { type: string; payload?: any }) => void; // bubble up submit/finalize/etc.
  submission?: NodeSubmissionDTO | null;
  saving?: boolean;
};

export function NodeDetailSwitch({
  node,
  role = "student",
  onEvent,
  submission,
  saving = false,
}: Props) {
  const kinds = detectActionKinds(node);
  const uiKind = deriveNodeKind(node);
  const initialForm = submission?.form?.data ?? {};
  const attachmentsBySlot = new Map<
    string,
    NodeSubmissionDTO["slots"][number]["attachments"]
  >();
  submission?.slots.forEach((slot) => {
    attachmentsBySlot.set(slot.key, slot.attachments);
  });

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

  // permissions (rough defaults, adjust as you wire real RBAC)
  const canDecide = role === "secretary" || role === "chair";
  const canUpload = role !== "admin"; // example
  const canComplete = node.who_can_complete?.includes(role);

  // Prefer UI-specific kinds first for form-like nodes
  if (node.type === "form") {
    const initialForm = submission?.form?.data ?? {};
    if (uiKind === "formEntry") {
      return (
        <FormEntryDetails
          node={node}
          initial={initialForm}
          disabled={saving}
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

  // Single dominant kind
  if (kinds.length === 1) {
    switch (kinds[0]) {
      case "form":
        return (
          <FormTaskDetails
            node={node}
            canEdit={!saving}
            initial={initialForm}
            disabled={saving}
            onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
          />
        );
      case "upload":
        // Confirm-style upload task (instructions + confirmation only)
        if (node.type === "uploadTask") {
          return (
            <Suspense fallback={<div className="p-2 text-sm">Loading…</div>}>
              <ConfirmUploadTaskDetails node={node} />
            </Suspense>
          );
        }
        return (
          <Suspense fallback={<div className="p-2 text-sm">Loading…</div>}>
            <UploadTaskDetails
              node={node}
              canEdit={!saving}
              existing={attachmentsBySlot}
              onSubmit={(payload) => onEvent?.({ type: "submit-upload", payload })}
            />
          </Suspense>
        );
      case "outcome":
        if (node.type === "decision" && canComplete) {
          return (
            <DecisionTaskDetails
              node={node}
              disabled={saving}
              onSubmit={() =>
                onEvent?.({
                  type: "submit-decision",
                  payload: { acknowledged: true },
                })
              }
            />
          );
        }
        // no generic outcome renderer needed now
        break;
      case "external":
        return (
          <Suspense fallback={<div className="p-2 text-sm">Loading…</div>}>
            <ExternalProcessDetails
              node={node}
              onComplete={(payload) => onEvent?.({ type: "complete-external", payload })}
            />
          </Suspense>
        );
      case "gateway":
        // Support explicit lightweight node types without fields/uploads/outcomes
        if (node.type === "info")
          return (
            <InfoDetails
              node={node}
              onContinue={() => onEvent?.({ type: "continue" })}
            />
          );
        if (node.type === "confirmTask")
          return <ConfirmTaskDetails node={node} />;
        if (node.type === "uploadTask")
          return <ConfirmUploadTaskDetails node={node} />;
      default:
        return (
          <GatewayInfoDetails
            node={node}
            onContinue={() => onEvent?.({ type: "continue" })}
          />
        );
    }
  }

  // Composite preference (outcome + upload)
  // no composite renderer in current playbook

  // Fallback: render in priority order (no recursion)
  const order = ["upload", "form", "external", "gateway"] as const;
  const first = order.find((k) => kinds.includes(k as any)) ?? "gateway";
  switch (first) {
    case "form":
      return (
        <FormTaskDetails
          node={node}
          canEdit={!saving}
          disabled={saving}
          initial={initialForm}
          onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
        />
      );
    case "upload":
      if (node.type === "uploadTask") {
        return <ConfirmUploadTaskDetails node={node} />;
      }
      return (
        <UploadTaskDetails
          node={node}
          canEdit={!saving}
          existing={attachmentsBySlot}
          onSubmit={(payload) => onEvent?.({ type: "submit-upload", payload })}
        />
      );
    case "outcome":
      if (node.type === "decision" && canComplete) {
        return (
          <DecisionTaskDetails
            node={node}
            disabled={saving}
            onSubmit={() =>
              onEvent?.({
                type: "submit-decision",
                payload: { acknowledged: true },
              })
            }
          />
        );
      }
      return (
        <GatewayInfoDetails
          node={node}
          onContinue={() => onEvent?.({ type: "continue" })}
        />
      );
    case "external":
      return (
        <ExternalProcessDetails
          node={node}
          onComplete={(payload) =>
            onEvent?.({ type: "complete-external", payload })
          }
        />
      );
    
    case "gateway":
    default:
      if (node.type === "info") {
        return (
          <InfoDetails
            node={node}
            onContinue={() => onEvent?.({ type: "continue" })}
          />
        );
      }
      if (node.type === "confirmTask") {
        return <ConfirmTaskDetails node={node} />;
      }
      return (
        <GatewayInfoDetails
          node={node}
          onContinue={() => onEvent?.({ type: "continue" })}
        />
      );
  }
}
