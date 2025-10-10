import { FormTaskDetails } from "@/components/node-details/variants/FormTaskDetails";
import type { NodeVM } from "@/lib/playbook";

export default function FormEntryDetails(props: {
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

