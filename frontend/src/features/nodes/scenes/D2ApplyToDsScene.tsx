import ChecklistDetails from "@/features/nodes/kinds/ChecklistDetails";

export default function D2ApplyToDsScene({
  node,
  initial = {},
  disabled,
  canEdit,
  onSubmit,
}: {
  node: any;
  initial?: Record<string, any>;
  disabled?: boolean;
  canEdit?: boolean;
  onSubmit?: (payload: any) => void;
}) {
  return (
    <ChecklistDetails
      node={node}
      initial={initial}
      disabled={disabled}
      canEdit={canEdit}
      onSubmit={onSubmit}
    />
  );
}
