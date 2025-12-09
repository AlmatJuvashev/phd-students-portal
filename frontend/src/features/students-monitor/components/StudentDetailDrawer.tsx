import React from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import {
  fetchStudentJourney,
  type JourneyNode,
  patchStudentNodeState,
  fetchDeadlines,
  putDeadline,
} from "../api";
import { Button } from "@/components/ui/button";
import {
  Calendar,
  FileText,
  Mail,
  Copy,
  Plus,
  ExternalLink,
  CheckCircle2,
  Circle,
  Clock,
  AlertTriangle,
  Lock,
  Send,
  FileSearch,
} from "lucide-react";

type NodeState =
  | "locked"
  | "active"
  | "submitted"
  | "waiting"
  | "under_review"
  | "needs_fixes"
  | "done";

const nodeStates: Record<
  NodeState,
  { label: string; color: string; icon: typeof CheckCircle2 }
> = {
  locked: {
    label: "Locked",
    color: "bg-muted text-muted-foreground",
    icon: Lock,
  },
  active: { label: "Active", color: "bg-blue-100 text-blue-800", icon: Circle },
  submitted: {
    label: "Submitted",
    color: "bg-purple-100 text-purple-800",
    icon: Send,
  },
  waiting: {
    label: "Waiting",
    color: "bg-amber-50 text-amber-700 border border-amber-200",
    icon: Clock,
  },
  under_review: {
    label: "Under Review",
    color: "bg-purple-50 text-purple-700 border border-purple-200",
    icon: FileSearch,
  },
  needs_fixes: {
    label: "Needs Fixes",
    color: "bg-red-50 text-red-700 border border-red-200",
    icon: AlertTriangle,
  },
  done: {
    label: "Done",
    color: "bg-green-100 text-green-800",
    icon: CheckCircle2,
  },
};

const STAGES = [
  { id: "W1", label: "Preparation" },
  { id: "W2", label: "Pre-examination" },
  { id: "W3", label: "RP" },
  { id: "W4", label: "Submission to DC" },
  { id: "W5", label: "Restoration" },
  { id: "W6", label: "After DC acceptance" },
  { id: "W7", label: "Defense & Post-defense" },
];

export function StudentDetailDrawer({
  open,
  onOpenChange,
  student,
}: {
  open: boolean;
  onOpenChange: (b: boolean) => void;
  student: {
    id: string;
    name: string;
    email?: string;
    phone?: string;
    program?: string;
    department?: string;
    cohort?: string;
    current_stage?: string;
    rp_required?: boolean;
    advisors?: { id: string; name: string }[];
  } | null;
}) {
  const [nodes, setNodes] = React.useState<JourneyNode[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [deadlines, setDeadlines] = React.useState<Record<string, string>>({});
  const [comment, setComment] = React.useState("");

  React.useEffect(() => {
    if (!open || !student) return;
    setLoading(true);
    Promise.all([
      fetchStudentJourney(student.id).then((d) => setNodes(d.nodes)),
      fetchDeadlines(student.id)
        .then((list) => {
          const m: Record<string, string> = {};
          list.forEach((it) => {
            m[it.node_id] = it.due_at;
          });
          return m;
        })
        .then((m) => setDeadlines(m)),
    ]).finally(() => setLoading(false));
  }, [open, student?.id]);

  async function confirm(nodeId: string) {
    if (!student) return;
    await patchStudentNodeState(student.id, nodeId, "done");
    const d = await fetchStudentJourney(student.id);
    setNodes(d.nodes);
  }

  async function setDue(nodeId: string, due: string) {
    if (!student) return;
    await putDeadline(student.id, nodeId, due);
    setDeadlines((prev) => ({ ...prev, [nodeId]: due }));
  }

  if (!student) return null;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-[600px] w-full overflow-y-auto">
        <SheetHeader className="space-y-4">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              <Avatar className="h-16 w-16 border-2">
                <AvatarFallback className="bg-primary/10 text-primary text-lg">
                  {student.name
                    .split(" ")
                    .map((n) => n[0])
                    .join("")
                    .slice(0, 2)}
                </AvatarFallback>
              </Avatar>
              <div>
                <SheetTitle className="text-2xl">{student.name}</SheetTitle>
                <SheetDescription className="text-base">
                  {student.id}
                </SheetDescription>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap gap-2">
            {student.program && (
              <Badge
                variant="outline"
                className="bg-primary/5 text-primary border-primary/20"
              >
                {student.program}
              </Badge>
            )}
            {student.department && (
              <Badge variant="outline">{student.department}</Badge>
            )}
            {student.cohort && (
              <Badge variant="outline">{student.cohort}</Badge>
            )}
            {student.rp_required && (
              <Badge
                variant="outline"
                className="bg-amber-50 text-amber-700 border-amber-200"
              >
                RP Required
              </Badge>
            )}
          </div>

          <div className="flex gap-2">
            <Button
              size="sm"
              variant="outline"
              className="flex-1 gap-2"
              disabled
            >
              <Copy className="h-4 w-4" />
              Copy Link
            </Button>
          </div>

          <div className="text-sm space-y-1">
            {student.email && (
              <div className="flex items-center gap-2 text-muted-foreground">
                <Mail className="h-4 w-4" />
                {student.email}
              </div>
            )}
          </div>

          <Separator />
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {/* Journey Map */}
          <div>
            <h3 className="text-sm font-medium mb-4">Journey Map</h3>
            <div className="flex items-center gap-2 overflow-x-auto pb-2">
              {STAGES.filter(
                (stage) => stage.id !== "W3" || student.rp_required
              ).map((stage, idx, arr) => {
                const currentStageIdx = STAGES.findIndex(
                  (s) => s.id === student.current_stage
                );
                const thisStageIdx = STAGES.findIndex((s) => s.id === stage.id);
                const isCurrent = stage.id === student.current_stage;
                const isCompleted = currentStageIdx > thisStageIdx;

                return (
                  <div
                    key={stage.id}
                    className="flex items-center flex-shrink-0"
                  >
                    <div
                      className={`px-3 py-2 rounded-lg text-xs font-medium whitespace-nowrap transition-all ${
                        isCurrent
                          ? "bg-primary text-primary-foreground shadow-md scale-105"
                          : isCompleted
                          ? "bg-green-100 text-green-800"
                          : "bg-muted text-muted-foreground"
                      }`}
                    >
                      {stage.label}
                    </div>
                    {idx < arr.length - 1 && (
                      <div className="w-8 h-0.5 bg-border mx-1 flex-shrink-0" />
                    )}
                  </div>
                );
              })}
            </div>
          </div>

          <Separator />

          {/* Stage Checklist */}
          <div>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium">
                Current Stage:{" "}
                {STAGES.find((s) => s.id === student.current_stage)?.label ||
                  "—"}
              </h3>
              <Badge variant="outline" className="text-xs">
                {nodes.length} nodes
              </Badge>
            </div>

            <div className="space-y-3">
              {loading ? (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  Loading nodes...
                </div>
              ) : nodes.length === 0 ? (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  No nodes available
                </div>
              ) : (
                nodes.map((node) => (
                  <NodeCard
                    key={node.node_id}
                    id={node.node_id}
                    state={node.state as NodeState}
                    dueDate={deadlines[node.node_id]}
                    onSetDue={(due) => setDue(node.node_id, due)}
                    onConfirm={() => confirm(node.node_id)}
                  />
                ))
              )}
            </div>
          </div>

          <Separator />

          {/* Documents */}
          <div>
            <h3 className="text-sm font-medium mb-3">Documents & Templates</h3>
            <div className="text-center py-8 text-sm text-muted-foreground">
              No documents uploaded yet
            </div>
          </div>

          <Separator />

          {/* Comments */}
          <div>
            <h3 className="text-sm font-medium mb-3">Comments & Notes</h3>
            <div className="space-y-2">
              <Textarea
                placeholder="Add a comment... Use @ to mention advisors"
                className="min-h-[80px]"
                value={comment}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                  setComment(e.target.value)
                }
              />
              <div className="flex justify-end gap-2">
                <Button size="sm" variant="outline">
                  Attach File
                </Button>
                <Button size="sm" disabled={!comment.trim()}>
                  Add Comment
                </Button>
              </div>
            </div>
          </div>


        </div>
      </SheetContent>
    </Sheet>
  );
}

function NodeCard({
  id,
  state,
  dueDate,
  onSetDue,
  onConfirm,
}: {
  id: string;
  state: NodeState;
  dueDate?: string;
  onSetDue: (due: string) => void;
  onConfirm: () => void;
}) {
  const stateConfig = nodeStates[state];
  const StateIcon = stateConfig.icon;

  return (
    <div className="p-4 rounded-lg border hover:shadow-md transition-all">
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <code className="text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded">
              {id}
            </code>
            <Badge variant="outline" className={`text-xs ${stateConfig.color}`}>
              <StateIcon className="h-3 w-3 mr-1" />
              {stateConfig.label}
            </Badge>
          </div>
        </div>
      </div>
      <div className="flex items-center gap-4 text-xs text-muted-foreground mb-3">
        {dueDate && (
          <div className="flex items-center gap-1">
            <Calendar className="h-3 w-3" />
            {new Date(dueDate).toLocaleString()}
          </div>
        )}
      </div>
      <div className="flex gap-2">
        <input
          type="datetime-local"
          aria-label="Set due date"
          className="flex-1 border rounded-md px-3 py-2 text-sm"
          value={dueDate ? dueDate.slice(0, 16) : ""}
          onChange={(e) => onSetDue(e.target.value)}
        />
        <Button size="sm" variant="outline" onClick={onConfirm}>
          Mark Done
        </Button>
      </div>
    </div>
  );
}
