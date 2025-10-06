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

type Props = {
  node: NodeVM;
  role?: "student" | "advisor" | "secretary" | "chair" | "admin";
  onEvent?: (evt: { type: string; payload?: any }) => void; // bubble up submit/finalize/etc.
  submission?: NodeSubmissionDTO | null;
};

export function NodeDetailSwitch({
  node,
  role = "student",
  onEvent,
  submission,
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

  // Single dominant kind
  if (kinds.length === 1) {
    switch (kinds[0]) {
      case "form":
        return (
          <FormTaskDetails
            node={node}
            canEdit
            initial={initialForm}
            onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
          />
        );
      case "upload":
        return (
          <UploadTaskDetails
            node={node}
            canEdit
            existing={attachmentsBySlot}
            onSubmit={(payload) =>
              onEvent?.({ type: "submit-upload", payload })
            }
          />
        );
      case "outcome":
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
          canEdit
          onSubmit={(payload) => onEvent?.({ type: "submit-form", payload })}
        />
      );
    case "upload":
      return (
        <UploadTaskDetails
          node={node}
          canEdit
          onSubmit={(payload) => onEvent?.({ type: "submit-upload", payload })}
        />
      );
    case "outcome":
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
      return (
        <GatewayInfoDetails
          node={node}
          onContinue={() => onEvent?.({ type: "continue" })}
        />
      );
  }
}
