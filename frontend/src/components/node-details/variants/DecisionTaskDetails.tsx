import { Button } from "@/components/ui/button";
import { NodeVM, t } from "@/lib/playbook";
import { useTranslation } from "react-i18next";

export function DecisionTaskDetails({
  node,
  onSubmit,
  disabled = false,
}: {
  node: NodeVM;
  onSubmit?: () => void;
  disabled?: boolean;
}) {
  const { t: T } = useTranslation("common");
  return (
    <div className="space-y-4 text-sm">
      {node.description ? (
        <p className="whitespace-pre-wrap text-muted-foreground">
          {t(node.description as any, node.id)}
        </p>
      ) : (
        <p className="text-muted-foreground">{T("decision.instructions")}</p>
      )}
      <Button disabled={disabled} onClick={onSubmit}>
        {T("decision.submit")}
      </Button>
    </div>
  );
}
