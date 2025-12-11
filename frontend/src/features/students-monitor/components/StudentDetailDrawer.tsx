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
  fetchStudentNodeFiles,
  type NodeFileRow,
  reviewAttachment,
} from "../api";
import { Button } from "@/components/ui/button";
import {
  Calendar,
  FileText,
  Mail,
  Copy,
  CheckCircle2,
  Circle,
  Clock,
  AlertTriangle,
  Lock,
  Send,
  FileSearch,
  Paperclip,
  Download,
  ThumbsUp,
  ThumbsDown,
  Loader2,
} from "lucide-react";
import { useTranslation } from "react-i18next";
import { useToast } from "@/components/ui/use-toast";
import PlaybookData from "../../../playbooks/playbook.json";
import { t as tPlaybook } from "@/lib/playbook";

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
  { id: "W7 (PhD)", label: "Defense & Post-defense" },
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
  const { t } = useTranslation("common");
  const { toast } = useToast();
  const [nodes, setNodes] = React.useState<JourneyNode[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [deadlines, setDeadlines] = React.useState<Record<string, string>>({});
  const [comment, setComment] = React.useState("");
  
  // Selection & Docs
  const [selectedNodeId, setSelectedNodeId] = React.useState<string | null>(null);
  const [nodeFiles, setNodeFiles] = React.useState<NodeFileRow[]>([]);
  const [loadingFiles, setLoadingFiles] = React.useState(false);

  // When drawer opens or student changes, fetch journey
  React.useEffect(() => {
    if (!open || !student) {
      setNodes([]);
      setDeadlines({});
      setSelectedNodeId(null);
      setNodeFiles([]);
      return;
    }
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

  // When selected node changes, fetch its files
  React.useEffect(() => {
    if (!student || !selectedNodeId) {
      setNodeFiles([]);
      return;
    }
    setLoadingFiles(true);
    fetchStudentNodeFiles(student.id, selectedNodeId)
      .then(setNodeFiles)
      .catch((err) => {
        console.error("Failed to fetch node files", err);
        toast({ title: "Error", description: "Failed to load files", variant: "destructive" });
      })
      .finally(() => setLoadingFiles(false));
  }, [student?.id, selectedNodeId]);

  async function confirm(nodeId: string) {
    if (!student) return;
    await patchStudentNodeState(student.id, nodeId, "done");
    const d = await fetchStudentJourney(student.id);
    setNodes(d.nodes);
    toast({ title: "Success", description: "Node marked as done" });
  }

  async function setDue(nodeId: string, due: string) {
    if (!student) return;
    await putDeadline(student.id, nodeId, due);
    setDeadlines((prev) => ({ ...prev, [nodeId]: due }));
    toast({ title: "Success", description: "Deadline set" });
  }

  async function handleReview(attachmentId: string, status: "approved" | "rejected") {
    try {
      if (!student || !selectedNodeId) return;
      await reviewAttachment(attachmentId, { status });
      toast({ title: "Success", description: status === "approved" ? "Document approved" : "Changes requested" });
      // refresh files
      const files = await fetchStudentNodeFiles(student.id, selectedNodeId);
      setNodeFiles(files);
      // refresh journey to update state if needed
      const d = await fetchStudentJourney(student.id);
      setNodes(d.nodes);
    } catch (e) {
      console.error(e);
      toast({ title: "Error", description: "Action failed", variant: "destructive" });
    }
  }

  if (!student) return null;

  const selectedNode = nodes.find(n => n.node_id === selectedNodeId);

  const getNodeTitle = React.useCallback((nodeId: string) => {
    // Find node in playbook
    for (const world of PlaybookData.worlds) {
      const n = world.nodes.find((nod) => nod.id === nodeId);
      if (n) return tPlaybook(n.title as any);
    }
    return nodeId;
  }, []);

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-[700px] w-full overflow-y-auto p-0 gap-0 flex flex-col h-full bg-background/95 backdrop-blur-xl">
        <SheetHeader className="px-6 py-6 border-b bg-muted/40 sticky top-0 z-10">
          {/* ... Header content ... */}
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              <Avatar className="h-16 w-16 border-2 ring-2 ring-background shadow-sm">
                <AvatarFallback className="bg-primary/10 text-primary text-xl font-medium">
                  {student.name
                    ?.split(" ")
                    .map((n) => n[0])
                    .join("")
                    .slice(0, 2) || "??"}
                </AvatarFallback>
              </Avatar>
              <div>
                <SheetTitle className="text-2xl font-bold tracking-tight">{student.name}</SheetTitle>
                <div className="flex items-center gap-2 text-muted-foreground mt-1">
                  <Badge variant="secondary" className="font-normal text-xs px-2 py-0.5 h-auto">
                    {student.id}
                  </Badge>
                  {student.email && (
                    <span className="text-sm flex items-center gap-1">
                      <Mail className="h-3 w-3" />
                      {student.email}
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap gap-2 mt-4">
            {student.program && (
              <Badge variant="outline" className="bg-background/50 backdrop-blur-sm">
                {student.program}
              </Badge>
            )}
            {student.department && (
              <Badge variant="outline" className="bg-background/50 backdrop-blur-sm truncate max-w-[200px]" title={student.department}>
                {student.department}
              </Badge>
            )}
            {student.cohort && (
              <Badge variant="outline" className="bg-background/50 backdrop-blur-sm">
                {student.cohort}
              </Badge>
            )}
            {student.rp_required && (
              <Badge
                variant="outline"
                className="bg-amber-50 text-amber-700 border-amber-200"
              >
                {t("monitor.filters.rp_only", "RP Required")}
              </Badge>
            )}
          </div>
        </SheetHeader>

        <div className="flex-1 overflow-y-auto px-6 py-6 space-y-8">
          {/* Journey Map */}
          <div>
            <h3 className="text-sm font-semibold mb-4 text-muted-foreground uppercase tracking-wider">{t("monitor.journey_map", "Journey Map")}</h3>
            <div className="flex items-center gap-2 overflow-x-auto pb-4 px-1 -mx-1 snap-x">
              {STAGES.filter(
                (stage) => stage.id !== "W3" || student.rp_required
              ).map((stage, idx, arr) => {
                const currentStageIdx = STAGES.findIndex(
                  (s) => s.id === (student.current_stage || "W1")
                );
                const thisStageIdx = STAGES.findIndex((s) => s.id === stage.id);
                const isCurrent = student.current_stage === stage.id;
                const isCompleted = currentStageIdx > thisStageIdx;

                // We should translate stage labels too if possible, but they are hardcoded in STAGES array above.
                // Assuming we leave them or map them later. For now, keep as is or map if simple.
                // STAGES array is hardcoded at top of file. I should probably move it inside component or translate it.
                // Let's defer stage translation to keep complexity low, or replace STAGES constant usage with a translated version?
                // The prompt asked for Node IDs/Names.
                
                return (
                  <div
                    key={stage.id}
                    className="flex items-center flex-shrink-0 snap-start"
                  >
                    <div
                      className={`px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border ${
                        isCurrent
                          ? "bg-primary text-primary-foreground border-primary shadow-md ring-2 ring-primary/20"
                          : isCompleted
                          ? "bg-green-100 text-green-800 border-green-200"
                          : "bg-muted/50 text-muted-foreground border-transparent"
                      }`}
                    >
                      {stage.label}
                    </div>
                    {idx < arr.length - 1 && (
                      <div className={`w-8 h-[2px] mx-1 flex-shrink-0 rounded-full ${isCompleted ? 'bg-green-200' : 'bg-border'}`} />
                    )}
                  </div>
                );
              })}
            </div>
          </div>

          <Separator />

          {/* Stage Nodes */}
          <div>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
                {t("monitor.current_stage_nodes", "Current Stage Nodes")}
              </h3>
              <Badge variant="secondary" className="text-xs font-normal">
                {t("monitor.nodes_count", { count: nodes.length, defaultValue: `${nodes.length} nodes` })}
              </Badge>
            </div>

            <div className="grid gap-3 sm:grid-cols-2">
              {loading ? (
                <div className="col-span-2 flex justify-center py-12 text-muted-foreground">
                  <Loader2 className="h-6 w-6 animate-spin mr-2" />
                  {t("monitor.loading_nodes", "Loading nodes...")}
                </div>
              ) : nodes.length === 0 ? (
                <div className="col-span-2 text-center py-12 text-sm text-muted-foreground border-2 border-dashed rounded-xl">
                  {t("monitor.no_nodes", "No nodes available for this stage")}
                </div>
              ) : (
                nodes.map((node) => (
                  <NodeCard
                    key={node.node_id}
                    id={node.node_id}
                    title={getNodeTitle(node.node_id)}
                    state={node.state as NodeState}
                    dueDate={deadlines[node.node_id]}
                    attachments={node.attachments || 0}
                    selected={selectedNodeId === node.node_id}
                    onSelect={() => setSelectedNodeId(node.node_id)}
                    onSetDue={(due) => setDue(node.node_id, due)}
                    onConfirm={() => confirm(node.node_id)}
                    t={t}
                  />
                ))
              )}
            </div>
          </div>

          <Separator />

          {/* Documents & Review */}
          <div id="documents-section" className="scroll-mt-6">
             <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider mb-1">Documents & Reviewing</h3>
                <p className="text-xs text-muted-foreground">
                  {selectedNodeId ? `Files for node ${selectedNodeId}` : "Select a node above to view files"}
                </p>
              </div>
            </div>

            <div className="bg-muted/30 border rounded-xl p-4 min-h-[150px]">
              {!selectedNodeId ? (
                <div className="h-full flex flex-col items-center justify-center text-muted-foreground/60 py-8">
                  <FileSearch className="h-10 w-10 mb-2 opacity-50" />
                  <p className="text-sm">Select a node from the list above to view uploaded documents</p>
                </div>
              ) : loadingFiles ? (
                <div className="h-full flex items-center justify-center text-muted-foreground py-8">
                  <Loader2 className="h-5 w-5 animate-spin mr-2" />
                  Loading files...
                </div>
              ) : nodeFiles.length === 0 ? (
                <div className="h-full flex flex-col items-center justify-center text-muted-foreground py-8">
                  <FileText className="h-10 w-10 mb-2 opacity-50" />
                  <p className="text-sm">No files uploaded for {selectedNodeId}</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {nodeFiles.map((file) => (
                    <div key={file.attachment_id} className="bg-background border rounded-lg p-3 shadow-sm transition-all hover:shadow-md">
                      <div className="flex items-start justify-between gap-3">
                         <div className="flex items-center gap-3 overflow-hidden">
                           <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0 text-primary">
                             <FileText className="h-5 w-5" />
                           </div>
                           <div className="min-w-0">
                             <div className="text-sm font-medium truncate" title={file.filename}>
                               {file.filename}
                             </div>
                             <div className="flex items-center gap-2 text-xs text-muted-foreground">
                               <span>{(file.size_bytes / 1024 / 1024).toFixed(2)} MB</span>
                               <span>•</span>
                               <span>{new Date(file.attached_at || "").toLocaleDateString()}</span>
                               {file.uploaded_by && (
                                 <>
                                  <span>•</span>
                                  <span>{file.uploaded_by}</span>
                                 </>
                               )}
                             </div>
                           </div>
                         </div>
                         
                         <div className="flex-shrink-0">
                            {file.status === "approved" ? (
                              <Badge className="bg-green-100 text-green-800 hover:bg-green-100 border-green-200">
                                <CheckCircle2 className="h-3 w-3 mr-1" />
                                Approved
                              </Badge>
                            ) : file.status === "rejected" ? (
                              <Badge variant="destructive" className="bg-red-100 text-red-800 hover:bg-red-100 border-red-200">
                                <AlertTriangle className="h-3 w-3 mr-1" />
                                Changes Requested
                              </Badge>
                            ) : (
                              <Badge variant="secondary" className="bg-amber-100 text-amber-800 hover:bg-amber-100 border-amber-200">
                                <Clock className="h-3 w-3 mr-1" />
                                Review Pending
                              </Badge>
                            )}
                         </div>
                      </div>

                      <div className="mt-3 flex items-center justify-between pt-3 border-t">
                        <a 
                          href={file.download_url} 
                          target="_blank" 
                          rel="noreferrer"
                          className="text-xs flex items-center gap-1 text-primary hover:underline font-medium"
                        >
                          <Download className="h-3 w-3" />
                          Download
                        </a>
                        
                        {(file.status === "submitted" || file.status === "rejected") && (
                          <div className="flex items-center gap-2">
                            <Button 
                              size="sm" 
                              variant="ghost" 
                              className="h-7 text-xs text-red-600 hover:text-red-700 hover:bg-red-50"
                              onClick={() => handleReview(file.attachment_id, "rejected")}
                            >
                              <ThumbsDown className="h-3 w-3 mr-1" />
                              Request Changes
                            </Button>
                            <Button 
                              size="sm" 
                              className="h-7 text-xs bg-green-600 hover:bg-green-700 text-white"
                              onClick={() => handleReview(file.attachment_id, "approved")}
                            >
                              <ThumbsUp className="h-3 w-3 mr-1" />
                              Approve
                            </Button>
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          <Separator />

          {/* Comments */}
          <div>
            <h3 className="text-sm font-semibold mb-3 text-muted-foreground uppercase tracking-wider">Comments & Notes</h3>
            <div className="space-y-2">
              <Textarea
                placeholder="Add a comment... Use @ to mention advisors"
                className="min-h-[80px] resize-none"
                value={comment}
                onChange={(e) => setComment(e.target.value)}
              />
              <div className="flex justify-end gap-2">
                <Button size="sm" variant="outline" className="gap-2">
                  <Paperclip className="h-3.5 w-3.5" />
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
  title,
  state,
  dueDate,
  attachments = 0,
  selected = false,
  onSelect,
  onSetDue,
  onConfirm,
  t,
}: {
  id: string;
  title?: string;
  state: NodeState;
  dueDate?: string;
  attachments?: number;
  selected?: boolean;
  onSelect: () => void;
  onSetDue: (due: string) => void;
  onConfirm: () => void;
  t: any;
}) {
  const stateConfig = nodeStates[state] || nodeStates.active;
  const StateIcon = stateConfig?.icon || Circle;

  // Visual hint: orange dot if pending review (submitted/under_review)
  const isPendingReview = state === "submitted" || state === "under_review";
  
  // Translate state label
  const stateLabel = t ? t(`admin.monitor.detail.node_states.${state}`, stateConfig.label) : stateConfig.label;

  return (
    <div 
      className={`
        p-4 rounded-xl border transition-all cursor-pointer relative group
        ${selected 
          ? "bg-primary/5 border-primary ring-1 ring-primary shadow-sm" 
          : "bg-card hover:border-primary/50 hover:shadow-md"
        }
      `}
      onClick={onSelect}
    >
      {/* Pending status dot */}
      {isPendingReview && (
        <span className="absolute top-3 right-3 h-2.5 w-2.5 rounded-full bg-amber-500 ring-2 ring-white shadow-sm animate-pulse" title="Needs Review" />
      )}

      <div className="flex flex-col h-full justify-between gap-3">
        <div>
          <div className="flex items-center gap-2 mb-2">
            <code className="text-[10px] font-mono font-semibold text-muted-foreground bg-muted px-1.5 py-0.5 rounded border">
              {id}
            </code>
            {title && (
              <span className="text-xs font-medium truncate flex-1" title={title}>
                {title}
              </span>
            )}
            <Badge variant="secondary" className={`text-[10px] font-medium px-1.5 py-0 border ${stateConfig.color}`}>
              <StateIcon className="h-3 w-3 mr-1" />
              {stateLabel}
            </Badge>
          </div>
          
          <div className="flex items-center gap-3 text-xs text-muted-foreground">
             {attachments > 0 && (
               <div className="flex items-center gap-1 font-medium text-foreground">
                 <Paperclip className="h-3 w-3" />
                 {attachments} files
               </div>
             )}
             {dueDate && (
              <div className="flex items-center gap-1">
                <Calendar className="h-3 w-3 text-muted-foreground" />
                {new Date(dueDate).toLocaleDateString()}
              </div>
             )}
             {!attachments && !dueDate && (
               <span className="text-muted-foreground/50 italic">No details</span>
             )}
          </div>
        </div>

        <div className="flex gap-2 pt-2 mt-auto border-t border-border/50" onClick={e => e.stopPropagation()}>
          <div className="relative flex-1">
             <input
              type="date"
              aria-label="Set due date"
              className="w-full text-xs border rounded-md px-2 py-1.5 bg-background hover:bg-muted/50 transition-colors focus:outline-none focus:ring-2 focus:ring-primary/20"
              value={dueDate ? dueDate.slice(0, 10) : ""}
              onChange={(e) => onSetDue(e.target.value)}
            />
          </div>
          <Button 
            size="sm" 
            variant="ghost" 
            className="h-8 px-2 text-xs hover:bg-green-50 hover:text-green-700"
            onClick={onConfirm}
            title="Mark as Done"
          >
            <CheckCircle2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
