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

// Lazy heavy variants
const UploadTaskDetails = React.lazy(() => import("./variants/UploadTaskDetails").then(m => ({ default: m.UploadTaskDetails })));
const CompositeTaskDetails = React.lazy(() => import("./variants/CompositeTaskDetails").then(m => ({ default: m.CompositeTaskDetails })));
const ConfirmUploadTaskDetails = React.lazy(() => import("./variants/ConfirmUploadTaskDetails").then(m => ({ default: m.default })));
const ExternalProcessDetails = React.lazy(() => import("./variants/ExternalProcessDetails").then(m => ({ default: m.ExternalProcessDetails })));
const WaitLockDetails = React.lazy(() => import("./variants/WaitLockDetails").then(m => ({ default: m.WaitLockDetails })));

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
        return (
          <OutcomeReviewDetails
            node={node}
            canDecide={canDecide}
            canUpload={canUpload}
            onFinalize={(payload) =>
              onEvent?.({ type: "finalize-outcome", payload })
            }
          />
        );
      case "wait":
        return (
          <Suspense fallback={<div className="p-2 text-sm">Loading…</div>}>
            <WaitLockDetails
              node={node}
              onSubscribe={() => onEvent?.({ type: "subscribe-timer" })}
            />
          </Suspense>
        );
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
  if (kinds.includes("composite")) {
        return (
          <Suspense fallback={<div className="p-2 text-sm">Loading…</div>}>
            <CompositeTaskDetails
              node={node}
              onFinalize={(payload) => onEvent?.({ type: "finalize-composite", payload })}
            />
          </Suspense>
        );
  }

  // Fallback: render in priority order (no recursion)
  const order = [
    "outcome",
    "upload",
    "form",
    "external",
    "wait",
    "gateway",
  ] as const;
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
        <OutcomeReviewDetails
          node={node}
          canDecide={canDecide}
          canUpload={canUpload}
          onFinalize={(payload) =>
            onEvent?.({ type: "finalize-outcome", payload })
          }
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
    case "wait":
      return (
        <WaitLockDetails
          node={node}
          onSubscribe={() => onEvent?.({ type: "subscribe-timer" })}
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
