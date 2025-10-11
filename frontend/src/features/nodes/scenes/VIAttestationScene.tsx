import ChecklistDetails from "@/features/nodes/kinds/ChecklistDetails";

export default function VIAttestationScene({
  node,
  initial = {},
  disabled,
  onSubmit,
}: {
  node: any;
  initial?: Record<string, any>;
  disabled?: boolean;
  onSubmit?: (payload: any) => void;
}) {
  return (
    <ChecklistDetails
      node={node}
      initial={initial}
      disabled={disabled}
      onSubmit={onSubmit}
    />
  );
}
