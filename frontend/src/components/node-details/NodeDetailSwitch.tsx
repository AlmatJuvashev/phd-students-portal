// components/node-details/NodeDetailSwitch.tsx
import { NodeVM, detectActionKinds } from "@/lib/playbook";
import { FormTaskDetails } from "./variants/FormTaskDetails";
import { UploadTaskDetails } from "./variants/UploadTaskDetails";
import { OutcomeReviewDetails } from "./variants/OutcomeReviewDetails";
import { WaitLockDetails } from "./variants/WaitLockDetails";
import { ExternalProcessDetails } from "./variants/ExternalProcessDetails";
import { GatewayInfoDetails } from "./variants/GatewayInfoDetails";
import { CompositeTaskDetails } from "./variants/CompositeTaskDetails";
import type { NodeSubmissionDTO } from "@/api/journey";
import { DecisionTaskDetails } from "./variants/DecisionTaskDetails";
import ConfirmTaskDetails from "./variants/ConfirmTaskDetails";
import ConfirmUploadTaskDetails from "./variants/ConfirmUploadTaskDetails";
import InfoDetails from "./variants/InfoDetails";

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

  // Prefer rendering a form when node is of type 'form' and includes fields,
  // even if outcomes also exist (e.g., checklist with completion rule).
  if (kinds.includes("form") && node.type === "form") {
    const initialForm = submission?.form?.data ?? {};
    return (
      <FormTaskDetails
        node={node}
        canEdit={!saving}
        initial={initialForm}
        disabled={saving}
        onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
      />
    );
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
        // If this is a simple confirm-style upload task (no real uploads), render ConfirmUploadTaskDetails
        if (node.type === "uploadTask") {
          return <ConfirmUploadTaskDetails node={node} />;
        }
        return (
          <UploadTaskDetails
            node={node}
            canEdit={!saving}
            existing={attachmentsBySlot}
            onSubmit={(payload) =>
              onEvent?.({ type: "submit-upload", payload })
            }
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
      case "wait":
        return (
          <WaitLockDetails
            node={node}
            onSubscribe={() => onEvent?.({ type: "subscribe-timer" })}
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
        // Support explicit lightweight node types without fields/uploads/outcomes
        if (node.type === "info")
          return (
            <InfoDetails
              node={node}
              onContinue={() => onEvent?.({ type: "continue" })}
            />
          );
        if (node.type === "confirmTask") return <ConfirmTaskDetails node={node} />;
        if (node.type === "uploadTask") return <ConfirmUploadTaskDetails node={node} />;
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
      <CompositeTaskDetails
        node={node}
        onFinalize={(payload) =>
          onEvent?.({ type: "finalize-composite", payload })
        }
      />
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
