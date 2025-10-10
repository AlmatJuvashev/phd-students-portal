import { FormTaskDetails } from "@/components/node-details/variants/FormTaskDetails";
import type { NodeVM } from "@/lib/playbook";

// Simple wrapper; current FormTaskDetails already supports checklist UX
export default function ChecklistDetails(props: {
  node: NodeVM;
  initial?: Record<string, any>;
  disabled?: boolean;
  onSubmit?: (payload: any) => void;
}) {
  const { node, initial, disabled, onSubmit } = props;
  return (
    <FormTaskDetails
      node={node}
      initial={initial}
      disabled={disabled}
      canEdit={!disabled}
      onSubmit={onSubmit}
    />
  );
}

