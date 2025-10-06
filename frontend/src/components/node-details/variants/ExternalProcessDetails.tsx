// components/node-details/variants/ExternalProcessDetails.tsx
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { NodeVM } from "@/lib/playbook";
import { UploadTaskDetails } from "./UploadTaskDetails";
import { AssetsDownloads } from "../AssetsDownloads";
import { useTranslation } from "react-i18next";

export function ExternalProcessDetails({
  node,
  onComplete,
}: {
  node: NodeVM;
  onComplete?: (payload: any) => void;
}) {
  const { t: T } = useTranslation("common");
  const hasUploads = !!node.requirements?.uploads?.length;
  return (
    <Card className="p-4 space-y-4">
      {node.requirements?.notes && (
        <p className="text-sm text-muted-foreground">
          {node.requirements.notes}
        </p>
      )}

      {/* Templates / Downloads (if any) */}
      <AssetsDownloads node={node} />

      {!!node.requirements?.checklist?.length && (
        <div>
          <div className="mb-2 font-medium">{T("external.checklist_title")}</div>
          <ul className="list-inside list-disc text-sm">
            {node.requirements.checklist.map((s, i) => (
              <li key={i}>{s}</li>
            ))}
          </ul>
        </div>
      )}

      {hasUploads && (
        <>
          <Separator />
          <UploadTaskDetails
            node={node}
            onSubmit={(payload) => onComplete?.(payload)}
          />
        </>
      )}

      {!hasUploads && (
        <Button onClick={() => onComplete?.({ completed: true })}>
          {T("external.mark_done")}
        </Button>
      )}
    </Card>
  );
}
