// components/node-details/NodeDetailsSheet.tsx

import { NodeVM, t } from "@/lib/playbook";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { NodeDetailSwitch } from "./NodeDetailSwitch";
import { useEffect, useState } from "react";
import { getNodeSubmission, NodeSubmissionDTO } from "@/api/journey";
import { useToast } from "@/components/toast";

export function NodeDetailsSheet({
  node,
  onOpenChange,
  role = "student",
  onEvent,
}: {
  node: NodeVM | null;
  onOpenChange: (open: boolean) => void;
  role?: "student" | "advisor" | "secretary" | "chair" | "admin";
  onEvent?: (evt: { type: string; payload?: any }) => void;
}) {
  const [submission, setSubmission] = useState<NodeSubmissionDTO | null>(null);
  const { push } = useToast();

  useEffect(() => {
    let mounted = true;
    if (node) {
      getNodeSubmission(node.id)
        .then((data) => {
          if (mounted) setSubmission(data);
        })
        .catch((err) => {
          push({
            title: "Failed to load",
            description: err.message || String(err),
          });
        });
    } else {
      setSubmission(null);
    }
    return () => {
      mounted = false;
    };
  }, [node]);

  return (
    <Sheet open={!!node} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full max-w-full sm:max-w-lg">
        {node && (
          <>
            <SheetHeader>
              <SheetTitle className="flex items-center gap-2">
                <span>{t(node.title, node.id)}</span>
                <Badge variant="secondary" className="capitalize">
                  {node.type}
                </Badge>
                <Badge className="capitalize">
                  {node.state?.replace("_", " ")}
                </Badge>
              </SheetTitle>
            </SheetHeader>

            <div className="mt-6">
              <NodeDetailSwitch
                node={node}
                role={role}
                submission={submission}
                onEvent={onEvent}
              />
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
