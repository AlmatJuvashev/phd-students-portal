import React from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import {
  fetchStudentDetails,
  fetchStudentJourney,
  fetchDeadlines,
  patchStudentNodeState,
  putDeadline,
  fetchStudentNodeFiles,
  reviewAttachment,
} from "../api";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Modal } from "@/components/ui/modal";
import {
  ArrowLeft,
  Mail,
  Phone,
  Copy,
  Calendar,
  CheckCircle2,
  Circle,
  Clock,
  AlertTriangle,
  Lock,
  Send,
  Plus,
  Download,
  Loader2,
} from "lucide-react";

const STAGES = [
  { id: "W1", label: "I — Preparation" },
  { id: "W2", label: "II — Pre-examination (SC)" },
  { id: "W3", label: "III — RP (conditional)", conditional: true },
  { id: "W4", label: "IV — Submission to DC" },
  { id: "W5", label: "V — Restoration" },
  { id: "W6", label: "VI — After DC acceptance" },
  { id: "W7", label: "VII — Defense & Post-defense" },
];

const attachmentStatuses: Record<
  string,
  { label: string; className: string }
> = {
  submitted: {
    label: "Pending",
    className: "bg-amber-50 text-amber-700 border border-amber-200",
  },
  approved: {
    label: "Approved",
    className: "bg-emerald-50 text-emerald-700 border border-emerald-200",
  },
  rejected: {
    label: "Needs fixes",
    className: "bg-rose-50 text-rose-700 border border-rose-200",
  },
};

const formatBytes = (bytes?: number) => {
  if (!bytes && bytes !== 0) return "";
  if (Math.abs(bytes) < 1024) return `${bytes} B`;
  const units = ["KB", "MB", "GB", "TB"];
  let value = bytes;
  let unitIndex = -1;
  do {
    value /= 1024;
    unitIndex += 1;
  } while (Math.abs(value) >= 1024 && unitIndex < units.length - 1);
  return `${value.toFixed(1)} ${units[unitIndex]}`;
};

const formatDateLabel = (value?: string) => {
  if (!value) return "—";
  const d = new Date(value);
  return d.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

type NodeState =
  | "locked"
  | "active"
  | "submitted"
  | "waiting"
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


export function StudentDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { i18n, t } = useTranslation('common');
  const [comment, setComment] = React.useState("");
  const [stageNodeIds, setStageNodeIds] = React.useState<string[] | null>(null);
  const [nodeTitles, setNodeTitles] = React.useState<Record<string, string>>(
    {}
  );
  const [allNodeTitles, setAllNodeTitles] = React.useState<Record<string, string>>({});
  const [selectedNodeId, setSelectedNodeId] = React.useState<string | null>(
    null,
  );
  const [reviewDialog, setReviewDialog] = React.useState<
    { attachmentId: string; filename: string } | null
  >(null);
  const [reviewNote, setReviewNote] = React.useState("");
  const [reviewMessage, setReviewMessage] = React.useState<
    { text: string; tone: "success" | "error" } | null
  >(null);
  const [pendingAttachment, setPendingAttachment] = React.useState<string | null>(
    null,
  );

  const {
    data: student,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["student", id],
    queryFn: () => fetchStudentDetails(id!),
    enabled: !!id,
  });

  const { data: journeyData, refetch: refetchJourney } = useQuery({
    queryKey: ["student-journey", id],
    queryFn: () => fetchStudentJourney(id!),
    enabled: !!id,
  });

  const { data: deadlinesList } = useQuery({
    queryKey: ["student-deadlines", id],
    queryFn: () => fetchDeadlines(id!),
    enabled: !!id,
  });

  const {
    data: nodeFiles = [],
    isLoading: nodeFilesLoading,
    refetch: refetchNodeFiles,
  } = useQuery({
    queryKey: ["student-node-files", id, selectedNodeId],
    queryFn: () => fetchStudentNodeFiles(id!, selectedNodeId!),
    enabled: !!id && !!selectedNodeId,
  });

  const deadlines = React.useMemo(() => {
    const map: Record<string, string> = {};
    deadlinesList?.forEach((d) => {
      map[d.node_id] = d.due_at;
    });
    return map;
  }, [deadlinesList]);

  // Load playbook and extract nodes for current stage
  React.useEffect(() => {
    let mounted = true;
    const stage = student?.current_stage;
    if (!stage) {
      setStageNodeIds(null);
      return;
    }

    import("@/playbooks/playbook.json")
      .then((mod: any) => {
        if (!mounted) return;
        const pb = (mod && (mod.default || mod)) as any;
        const worlds = (pb.worlds || pb.Worlds || []) as any[];
        const world = worlds.find((w: any) => w.id === stage || w.ID === stage);
        if (world) {
          const nodesArr = (world.nodes || world.Nodes || []) as any[];
          const ids = nodesArr.map((n: any) => n.id || n.ID);
          setStageNodeIds(ids);
          // Build titles map (prefer EN -> RU -> KZ if available)
          const titleFor = (n: any) => {
            const t = n.title || n.Title || {};
            return t.en || t.EN || t.En || t.ru || t.RU || t.kz || t.KZ || "";
          };
          const map: Record<string, string> = {};
          nodesArr.forEach((n: any) => {
            const id = n.id || n.ID;
            map[id] = titleFor(n);
          });
          setNodeTitles(map);
        } else {
          setStageNodeIds(null);
          setNodeTitles({});
        }
      })
      .catch(() => {
        if (mounted) {
          setStageNodeIds(null);
          setNodeTitles({});
        }
      });

    return () => {
      mounted = false;
    };
  }, [student?.current_stage]);

  // Build global map of node id -> human title for Upcoming Deadlines
  React.useEffect(() => {
    let mounted = true;
    import("@/playbooks/playbook.json")
      .then((mod: any) => {
        if (!mounted) return;
        const pb = (mod && (mod.default || mod)) as any;
        const worlds = (pb.worlds || pb.Worlds || []) as any[];
        const lang = (i18n?.language || 'en').toLowerCase();
        const pick = (obj: any, key: string) => obj?.[key] || obj?.[key?.toUpperCase?.()] || (key ? obj?.[key.charAt(0).toUpperCase()+key.slice(1)] : undefined);
        const titles: Record<string, string> = {};
        worlds.forEach((w: any) => {
          const nodesArr = (w.nodes || w.Nodes || []) as any[];
          nodesArr.forEach((n: any) => {
            const id = n.id || n.ID;
            const t = n.title || n.Title || {};
            titles[id] = pick(t, lang) || pick(t, 'en') || pick(t, 'ru') || pick(t, 'kz') || id;
          });
        });
        setAllNodeTitles(titles);
      })
      .catch(() => setAllNodeTitles({}));
    return () => {
      mounted = false;
    };
  }, [i18n?.language]);

  // Compute stage nodes using playbook data
  const nodes = journeyData?.nodes || [];
  const stageNodes = React.useMemo(() => {
    // Filter nodes by current stage using playbook node IDs
    if (!stageNodeIds || stageNodeIds.length === 0) return [];

    const set = new Set(stageNodeIds);
    return nodes.filter((n: any) => set.has(n.node_id));
  }, [nodes, stageNodeIds]);

  React.useEffect(() => {
    if (!selectedNodeId && stageNodes.length > 0) {
      setSelectedNodeId(stageNodes[0].node_id);
    }
  }, [selectedNodeId, stageNodes]);

  const handleMarkDone = async (nodeId: string) => {
    if (!id) return;
    await patchStudentNodeState(id, nodeId, "done");
    refetchJourney();
  };

  const handleSetDeadline = async (nodeId: string, due: string) => {
    if (!id) return;
    await putDeadline(id, nodeId, due);
  };
  const reviewMutation = useMutation({
    mutationFn: ({
      attachmentId,
      status,
      note,
    }: {
      attachmentId: string;
      status: "approved" | "rejected";
      note?: string;
    }) => reviewAttachment(attachmentId, { status, note }),
    onSuccess: (_data, variables) => {
    setReviewMessage({
      tone: "success",
      text:
        variables.status === "approved"
          ? t("admin.review.approved_toast", {
              defaultValue: "Document approved",
            })
          : t("admin.review.requested_toast", {
              defaultValue: "Requested changes sent",
            }),
    });
      refetchNodeFiles();
      refetchJourney();
    },
    onError: (err: any) => {
      setReviewMessage({
        tone: "error",
        text:
          err?.message ||
          t("admin.review.error_toast", {
            defaultValue: "Unable to update status",
          }),
      });
    },
  });

  const handleApprove = (attachmentId: string) => {
    setPendingAttachment(attachmentId);
    reviewMutation.mutate(
      { attachmentId, status: "approved" },
      {
        onSettled: () => setPendingAttachment(null),
      },
    );
  };

  const submitRequestChanges = () => {
    if (!reviewDialog) return;
    setPendingAttachment(reviewDialog.attachmentId);
    reviewMutation.mutate(
      {
        attachmentId: reviewDialog.attachmentId,
        status: "rejected",
        note: reviewNote || undefined,
      },
      {
        onSettled: () => {
          setPendingAttachment(null);
          setReviewDialog(null);
          setReviewNote("");
        },
      },
    );
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-muted-foreground">Loading student details...</div>
      </div>
    );
  }

  if (error || !student) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-destructive">Failed to load student details.</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header with Back Button */}
      <header className="sticky top-0 z-50 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80 border-b">
        <div className="px-8 py-5">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate("/admin/students-monitor")}
              className="gap-2"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Students
            </Button>
            <Separator orientation="vertical" className="h-6" />
            <h1 className="text-xl font-semibold">Student Details</h1>
            <Badge
              variant="outline"
              className="bg-primary/5 text-primary border-primary/20"
            >
              {student.id}
            </Badge>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="px-8 py-8 max-w-6xl mx-auto">
        <div className="space-y-8">
          {/* Student Profile Header */}
          <Card className="border shadow-sm">
            <CardContent className="p-8">
              <div className="flex items-start gap-6 mb-6">
                <Avatar className="h-24 w-24 border-2">
                  <AvatarFallback className="bg-primary/10 text-primary text-2xl">
                    {student.name
                      .split(" ")
                      .map((n: string) => n[0])
                      .join("")
                      .slice(0, 2)}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1">
                  <h2 className="text-3xl font-semibold mb-2">
                    {student.name}
                  </h2>
                  <div className="flex flex-wrap gap-2 mb-4">
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

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    {student.email && (
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <Mail className="h-4 w-4" />
                        <a
                          href={`mailto:${student.email}`}
                          className="hover:text-primary"
                        >
                          {student.email}
                        </a>
                      </div>
                    )}
                    {student.phone && (
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <Phone className="h-4 w-4" />
                        {student.phone}
                      </div>
                    )}
                  </div>
                </div>

                <div className="flex flex-col gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    className="gap-2"
                    disabled
                  >
                    <Copy className="h-4 w-4" />
                    Copy Link
                  </Button>
                </div>
              </div>

              <Separator className="mb-6" />

              <div className="grid grid-cols-3 gap-6">
                <div>
                  <div className="text-sm text-muted-foreground mb-1">
                    Advisors
                  </div>
                  <div className="flex flex-wrap gap-1">
                    {(student.advisors || []).map((advisor: any) => (
                      <Badge
                        key={advisor.id}
                        variant="outline"
                        className="text-xs"
                      >
                        {advisor.name}
                      </Badge>
                    ))}
                    {(!student.advisors || student.advisors.length === 0) && (
                      <span className="text-sm text-muted-foreground">
                        No advisor assigned
                      </span>
                    )}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground mb-1">
                    Overall Progress
                  </div>
                  <div className="flex items-center gap-3">
                    <Progress
                      value={student.overall_progress_pct || 0}
                      className="h-2 flex-1"
                    />
                    <span className="text-lg font-semibold">
                      {Math.round(student.overall_progress_pct || 0)}%
                    </span>
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground mb-1">
                    Status
                  </div>
                  <Badge
                    variant="outline"
                    className={`text-xs ${
                      student.overdue
                        ? "bg-red-50 text-red-700 border-red-200"
                        : "bg-muted/20"
                    }`}
                  >
                    {student.overdue ? "OVERDUE" : "NORMAL"}
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Journey Map */}
          <Card className="border shadow-sm">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold mb-6">
                Dissertation Journey Map
              </h3>
              <div className="flex items-center gap-2 overflow-x-auto pb-2">
                {STAGES.filter(
                  (stage) => !stage.conditional || student.rp_required
                ).map((stage, idx, arr) => {
                  const currentStageIdx = STAGES.findIndex(
                    (s) => s.id === student.current_stage
                  );
                  const thisStageIdx = STAGES.findIndex(
                    (s) => s.id === stage.id
                  );
                  const isCurrent = stage.id === student.current_stage;
                  const isCompleted = currentStageIdx > thisStageIdx;

                  return (
                    <div
                      key={stage.id}
                      className="flex items-center flex-shrink-0"
                    >
                      <div
                        className={`px-4 py-3 rounded-lg text-sm font-medium whitespace-nowrap transition-all ${
                          isCurrent
                            ? "bg-primary text-primary-foreground shadow-md scale-105"
                            : isCompleted
                            ? "bg-green-100 text-green-800"
                            : "bg-muted/30 text-muted-foreground"
                        }`}
                      >
                        {stage.label}
                      </div>
                      {idx < arr.length - 1 && (
                        <div className="w-12 h-0.5 bg-border mx-2 flex-shrink-0" />
                      )}
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          {/* Stage Progress and Checklist */}
          <Card className="border shadow-sm">
            <CardContent className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold">
                  Current Stage:{" "}
                  {STAGES.find((s) => s.id === student.current_stage)?.label ||
                    "—"}
                </h3>
                <Badge variant="outline" className="text-sm">
                  {stageNodes.length} nodes
                </Badge>
              </div>

              <div className="grid gap-4">
                {stageNodes.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No nodes available for this stage
                  </div>
                ) : (
                  stageNodes.map((node: any) => (
                    <NodeCard
                      key={node.node_id}
                      id={node.node_id}
                      title={nodeTitles[node.node_id]}
                      state={node.state as NodeState}
                      dueDate={deadlines[node.node_id]}
                      onSetDue={(due) => handleSetDeadline(node.node_id, due)}
                    />
                  ))
                )}
              </div>
            </CardContent>
          </Card>

          {/* Comments */}
          <Card className="border shadow-sm">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold mb-4">Comments & Notes</h3>
              <div className="space-y-3">
                <Textarea
                  placeholder="Add a comment... Use @ to mention advisors"
                  className="min-h-[100px]"
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
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
            </CardContent>
          </Card>

          {/* Documents & Review */}
          <Card className="border shadow-sm">
            <CardContent className="p-6 space-y-4">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <h3 className="text-lg font-semibold">
                  {t("admin.review.title", { defaultValue: "Documents & Review" })}
                </h3>
                {stageNodes.length > 0 && (
                  <Select
                    value={selectedNodeId ?? ""}
                    onValueChange={setSelectedNodeId}
                  >
                    <SelectTrigger className="w-full sm:w-72">
                      <SelectValue
                        placeholder={t("admin.review.select_node", {
                          defaultValue: "Select node",
                        })}
                      />
                    </SelectTrigger>
                    <SelectContent>
                      {stageNodes.map((node: any) => (
                        <SelectItem key={node.node_id} value={node.node_id}>
                          {nodeTitles[node.node_id] || node.node_id}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              </div>
              {reviewMessage && (
                <div
                  className={`rounded-md border p-3 text-sm ${
                    reviewMessage.tone === "error"
                      ? "border-destructive/50 bg-destructive/5 text-destructive"
                      : "border-emerald-200 bg-emerald-50 text-emerald-700"
                  }`}
                >
                  {reviewMessage.text}
                </div>
              )}
              {!selectedNodeId ? (
                <p className="text-sm text-muted-foreground">
                  {t("admin.review.hint", {
                    defaultValue: "Select a node above to inspect uploaded files.",
                  })}
                </p>
              ) : nodeFilesLoading ? (
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  {t("admin.review.loading", { defaultValue: "Loading files..." })}
                </div>
              ) : nodeFiles.length === 0 ? (
                <div className="rounded-md border border-dashed p-6 text-center text-sm text-muted-foreground">
                  {t("admin.review.empty", {
                    defaultValue: "No files uploaded for this node yet.",
                  })}
                </div>
              ) : (
                <div className="space-y-3">
                  {nodeFiles.map((file) => {
                    const statusBadge = attachmentStatuses[file.status] || {
                      label: file.status,
                      className:
                        "bg-muted text-muted-foreground border border-border",
                    };
                    const working =
                      pendingAttachment === file.attachment_id &&
                      reviewMutation.isPending;
                    const statusLabel = t(`uploads.status.${file.status}`, {
                      defaultValue: statusBadge.label,
                    });
                    const uploadedByText = file.uploaded_by
                      ? t("admin.review.uploaded_by", {
                          defaultValue: ` · ${file.uploaded_by}`,
                          name: file.uploaded_by,
                        })
                      : "";
                    const uploadMeta = t("admin.review.uploaded_meta", {
                      defaultValue: `Uploaded ${formatDateLabel(file.attached_at)}${uploadedByText}`,
                      date: formatDateLabel(file.attached_at),
                      by: uploadedByText,
                    });
                    return (
                      <div
                        key={file.attachment_id}
                        className="rounded-lg border border-dashed px-4 py-3 space-y-3"
                      >
                        <div className="flex flex-wrap items-center justify-between gap-3">
                          <div>
                            <p className="text-sm font-medium text-foreground break-all">
                              {file.filename}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              {formatBytes(file.size_bytes)} · {uploadMeta}
                            </p>
                            {file.review_note && (
                              <p className="text-xs text-amber-700 mt-1">
                                {t("uploads.note", {
                                  defaultValue: "Reviewer note:",
                                })}{" "}
                                {file.review_note}
                              </p>
                            )}
                          </div>
                          <div className="flex items-center gap-2">
                            <Badge variant="outline" className={statusBadge.className}>
                              {statusLabel}
                            </Badge>
                            <Button variant="ghost" size="icon" asChild>
                              <a
                                href={file.download_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                aria-label={t("uploads.download", {
                                  defaultValue: "Download",
                                })}
                              >
                                <Download className="h-4 w-4" />
                              </a>
                            </Button>
                          </div>
                        </div>
                        <div className="flex flex-wrap gap-2">
                          {file.status !== "approved" && (
                            <Button
                              size="sm"
                              className="bg-green-600 hover:bg-green-700 text-white"
                              disabled={working}
                              onClick={() => handleApprove(file.attachment_id)}
                            >
                              {working ? (
                                <Loader2 className="h-4 w-4 animate-spin mr-2" />
                              ) : (
                                <CheckCircle2 className="h-4 w-4 mr-2" />
                              )}
                              {t("admin.review.approve", { defaultValue: "Approve" })}
                            </Button>
                          )}
                          {file.status !== "rejected" && (
                            <Button
                              size="sm"
                              variant="outline"
                              className="border-red-200 text-red-700 hover:bg-red-50"
                              disabled={working}
                              onClick={() =>
                                setReviewDialog({
                                  attachmentId: file.attachment_id,
                                  filename: file.filename,
                                })
                              }
                            >
                              <AlertTriangle className="h-4 w-4 mr-2" />
                              {t("admin.review.request_changes", {
                                defaultValue: "Request changes",
                              })}
                            </Button>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Deadlines & Reminders */}
          <Card className="border shadow-sm">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold mb-6">
                Deadlines & Reminders
              </h3>
              <div className="space-y-3">
                {student.due_next && (
                  <div className="flex items-center justify-between p-4 rounded-lg border bg-background">
                    <div className="flex items-center gap-3">
                      <Calendar className="h-5 w-5 text-primary" />
                      <div>
                        <div className="text-sm font-medium text-foreground">
                          Next Due:{" "}
                          {new Date(student.due_next).toLocaleDateString()}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {student.stage_done}/{student.stage_total} nodes
                          completed in current stage
                        </div>
                      </div>
                    </div>
                    {student.overdue && (
                      <Badge
                        variant="outline"
                        className="bg-red-50 text-red-700 border-red-200"
                      >
                        Overdue
                      </Badge>
                    )}
                  </div>
                )}

                {deadlinesList && deadlinesList.length > 0 && (
                  <div className="space-y-2 mt-4">
                    <div className="text-sm font-medium text-foreground mb-2">
                      Upcoming Deadlines
                    </div>
                    {deadlinesList.slice(0, 5).map((deadline) => (
                      <div
                        key={deadline.node_id}
                        className="flex items-center justify-between p-3 rounded-lg border bg-muted/10"
                      >
                        <div className="flex items-center gap-2">
                          <Clock className="h-4 w-4 text-muted-foreground" />
                          <div className="flex flex-col">
                            <div className="text-sm font-medium text-foreground">
                              {allNodeTitles[deadline.node_id] || deadline.node_id}
                            </div>
                            <code className="text-[10px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded font-mono">
                              {deadline.node_id}
                            </code>
                          </div>
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {new Date(deadline.due_at).toLocaleString()}
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                <Button
                  size="sm"
                  variant="outline"
                  className="w-full gap-2 mt-4"
                >
                  <Plus className="h-4 w-4" />
                  Add New Reminder
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
      <Modal
        open={!!reviewDialog}
        onClose={() => {
          setReviewDialog(null);
          setReviewNote("");
        }}
      >
        <div className="space-y-4">
          <div>
            <h4 className="text-lg font-semibold">
              {t("admin.review.modal_title", { defaultValue: "Request changes" })}
            </h4>
            <p className="text-sm text-muted-foreground">
              {t("admin.review.modal_hint", {
                defaultValue: reviewDialog?.filename || "",
                filename: reviewDialog?.filename ?? "",
              })}
            </p>
          </div>
          <Textarea
            placeholder={t("admin.review.modal_placeholder", {
              defaultValue: "Let the student know what needs to be fixed",
            })}
            rows={4}
            value={reviewNote}
            onChange={(event) => setReviewNote(event.target.value)}
          />
          <div className="flex justify-end gap-2">
            <Button
              variant="ghost"
              onClick={() => {
                setReviewDialog(null);
                setReviewNote("");
              }}
            >
              {t("common.cancel", { defaultValue: "Cancel" })}
            </Button>
            <Button
              onClick={submitRequestChanges}
              disabled={
                !reviewNote.trim() ||
                (pendingAttachment === reviewDialog?.attachmentId &&
                  reviewMutation.isPending)
              }
            >
              {pendingAttachment === reviewDialog?.attachmentId &&
              reviewMutation.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
              ) : null}
              {t("admin.review.send", { defaultValue: "Send request" })}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

function NodeCard({
  id,
  title,
  state,
  dueDate,
  onSetDue,
}: {
  id: string;
  title?: string;
  state: NodeState;
  dueDate?: string;
  onSetDue: (due: string) => void;
}) {
  const stateConfig = nodeStates[state];
  const StateIcon = stateConfig.icon;

  return (
    <div className="p-4 rounded-lg border bg-card hover:bg-muted/10 transition-colors">
      <div className="flex items-start justify-between">
        <div className="min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <code className="text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded font-mono">
              {id}
            </code>
            <Badge variant="outline" className={`text-xs ${stateConfig.color}`}>
              <StateIcon className="h-3 w-3 mr-1" />
              {stateConfig.label}
            </Badge>
          </div>
          <div className="text-sm font-medium text-foreground truncate">
            {title || "—"}
          </div>
          {dueDate && (
            <div className="mt-1 flex items-center gap-1 text-xs text-muted-foreground">
              <Calendar className="h-3.5 w-3.5" />
              Due: {new Date(dueDate).toLocaleString()}
            </div>
          )}
        </div>
        <div className="w-48 ml-4">
          <label className="block text-xs text-muted-foreground mb-1">
            Set deadline
          </label>
          <input
            type="datetime-local"
            aria-label="Set due date"
            className="w-full border rounded-md px-2 py-1 text-xs"
            value={dueDate ? dueDate.slice(0, 16) : ""}
            onChange={(e) => onSetDue(e.target.value)}
          />
        </div>
      </div>
    </div>
  );
}

export default StudentDetailPage;
