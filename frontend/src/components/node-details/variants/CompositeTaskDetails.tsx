// components/node-details/variants/CompositeTaskDetails.tsx
import { NodeVM } from "@/lib/playbook";
import { OutcomeReviewDetails } from "./OutcomeReviewDetails";

export function CompositeTaskDetails({
  node,
  onFinalize,
}: {
  node: NodeVM;
  onFinalize?: (payload: {
    outcome: string;
    note?: string;
    files?: Record<string, File | null>;
  }) => void;
}) {
  return (
    <OutcomeReviewDetails
      node={node}
      canDecide
      canUpload
      onFinalize={onFinalize}
    />
  );
}
