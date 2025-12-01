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
import {
  presignReviewedDocumentUpload,
  attachReviewedDocument,
} from "@/api/admin";
import { stageLabel } from "../utils";
import { API_URL } from "@/api/client";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

// Extract base URL without /api suffix
const BASE_URL = API_URL.replace(/\/api$/, "");
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
  FileText,
  File,
  FileCheck,
  AlertCircle,
  Eye,
  FileSearch,
} from "lucide-react";

const STAGES: Array<{ id: string; conditional?: boolean }> = [
  { id: "W1" },
  { id: "W2" },
  { id: "W3", conditional: true },
  { id: "W4" },
  { id: "W5" },
  { id: "W6" },
  { id: "W7" },
];

const attachmentStatuses: Record<string, { label: string; className: string }> =
  {
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

export function StudentDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { i18n, t } = useTranslation("common");
  const language = i18n.language || "en";
  const [comment, setComment] = React.useState("");
  const [stageNodeIds, setStageNodeIds] = React.useState<string[] | null>(null);
  const [nodeTitles, setNodeTitles] = React.useState<Record<string, string>>(
    {}
  );
  const [allNodeTitles, setAllNodeTitles] = React.useState<
    Record<string, string>
  >({});
  const [selectedNodeId, setSelectedNodeId] = React.useState<string | null>(
    null
  );
  const [slotLabels, setSlotLabels] = React.useState<
    Record<string, { label: string; required?: boolean }>
  >({});
  const [reviewDialog, setReviewDialog] = React.useState<{
    attachmentId: string;
    filename: string;
  } | null>(null);
  const [reviewNote, setReviewNote] = React.useState("");
  const [reviewedFile, setReviewedFile] = React.useState<File | null>(null);
  const [isUploadingReviewed, setIsUploadingReviewed] = React.useState(false);
  const [reviewMessage, setReviewMessage] = React.useState<{
    text: string;
    tone: "success" | "error";
  } | null>(null);
  const [pendingAttachment, setPendingAttachment] = React.useState<
    string | null
  >(null);

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
        const lang = (i18n?.language || "en").toLowerCase();
        const pick = (obj: any, key: string) =>
          obj?.[key] ||
          obj?.[key?.toUpperCase?.()] ||
          (key ? obj?.[key.charAt(0).toUpperCase() + key.slice(1)] : undefined);
        const titles: Record<string, string> = {};
        worlds.forEach((w: any) => {
          const nodesArr = (w.nodes || w.Nodes || []) as any[];
          nodesArr.forEach((n: any) => {
            const id = n.id || n.ID;
            const t = n.title || n.Title || {};
            titles[id] =
              pick(t, lang) ||
              pick(t, "en") ||
              pick(t, "ru") ||
              pick(t, "kz") ||
              id;
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

  // Load slot labels for the selected node from the playbook (uploads spec)
  React.useEffect(() => {
    let mounted = true;
    if (!selectedNodeId) {
      setSlotLabels({});
      return;
    }
    import("@/playbooks/playbook.json")
      .then((mod: any) => {
        if (!mounted) return;
        const pb = (mod && (mod.default || mod)) as any;
        const worlds = (pb.worlds || pb.Worlds || []) as any[];
        const lang = (i18n?.language || "en").toLowerCase();
        const findNode = () => {
          for (const w of worlds) {
            const nodesArr = (w.nodes || w.Nodes || []) as any[];
            for (const n of nodesArr) {
              const id = n.id || n.ID;
              if (id === selectedNodeId) return n;
            }
          }
          return null;
        };
        const node = findNode();
        const uploads =
          node?.requirements?.uploads || node?.Requirements?.Uploads || [];
        const map: Record<string, { label: string; required?: boolean }> = {};
        for (const up of uploads) {
          const key = up.key || up.Key;
          const lbl = up.label || up.Label || {};
          const label =
            lbl[lang] || lbl[lang?.toUpperCase?.()] || lbl.en || lbl.EN || key;
          map[key] = { label, required: !!(up.required ?? up.Required) };
        }
        setSlotLabels(map);
      })
      .catch(() => setSlotLabels({}));
    return () => {
      mounted = false;
    };
  }, [selectedNodeId, i18n?.language]);

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
      }
    );
  };

  const submitRequestChanges = async () => {
    if (!reviewDialog) return;

    setPendingAttachment(reviewDialog.attachmentId);
    setIsUploadingReviewed(true);

    try {
      // Step 1: Upload reviewed document if file is selected
      if (reviewedFile) {
        // Get presigned URL
        const presignData = await presignReviewedDocumentUpload(
          reviewDialog.attachmentId,
          {
            filename: reviewedFile.name,
            content_type: reviewedFile.type,
            size_bytes: reviewedFile.size,
          }
        );

        // Upload file to S3
        const uploadResponse = await fetch(presignData.upload_url, {
          method: "PUT",
          body: reviewedFile,
          headers: {
            "Content-Type": reviewedFile.type,
          },
        });

        if (!uploadResponse.ok) {
          throw new Error("Failed to upload file to S3");
        }

        const etag =
          uploadResponse.headers.get("ETag")?.replace(/"/g, "") || "";

        // Attach the uploaded file
        await attachReviewedDocument(reviewDialog.attachmentId, {
          object_key: presignData.object_key,
          filename: reviewedFile.name,
          content_type: reviewedFile.type,
          size_bytes: reviewedFile.size,
          etag,
        });
      }

      // Step 2: Submit the review with rejection status
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
            setReviewedFile(null);
            setIsUploadingReviewed(false);
          },
        }
      );
    } catch (error) {
      console.error("Error uploading reviewed document:", error);
      setReviewMessage({
        text: t("admin.review.upload_error", {
          defaultValue: "Failed to upload reviewed document",
        }),
        tone: "error",
      });
      setPendingAttachment(null);
      setIsUploadingReviewed(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-muted-foreground">
          {t("admin.monitor.detail.loading", {
            defaultValue: "Loading student details...",
          })}
        </div>
      </div>
    );
  }

  if (error || !student) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-destructive">
          {t("admin.monitor.detail.error", {
            defaultValue: "Failed to load student details.",
          })}
        </div>
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
              {t("admin.monitor.detail.back", {
                defaultValue: "Back to Students",
              })}
            </Button>
            <Separator orientation="vertical" className="h-6" />
            <h1 className="text-xl font-semibold">
              {t("admin.monitor.detail.title", {
                defaultValue: "Student Details",
              })}
            </h1>
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
                        {t("admin.monitor.kanban.rp_required", {
                          defaultValue: "RP Required",
                        })}
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
                    {t("admin.monitor.detail.copy_link", {
                      defaultValue: "Copy Link",
                    })}
                  </Button>
                </div>
              </div>

              <Separator className="mb-6" />

              <div className="grid grid-cols-3 gap-6">
                <div>
                  <div className="text-sm text-muted-foreground mb-1">
                    {t("admin.monitor.detail.advisors", {
                      defaultValue: "Advisors",
                    })}
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
                        {t("admin.monitor.detail.no_advisor", {
                          defaultValue: "No advisor assigned",
                        })}
                      </span>
                    )}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground mb-1">
                    {t("admin.monitor.detail.overall_progress", {
                      defaultValue: "Overall Progress",
                    })}
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
                    {t("admin.monitor.detail.status", {
                      defaultValue: "Status",
                    })}
                  </div>
                  <Badge
                    variant="outline"
                    className="bg-blue-100 text-blue-800 font-medium"
                  >
                    {t("admin.monitor.detail.status_in_progress", {
                      defaultValue: "In Progress",
                    })}
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Journey Map */}
          <Card className="border shadow-sm">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold mb-6">
                {t("admin.monitor.detail.journey_card", {
                  defaultValue: "Dissertation Journey Map",
                })}
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
                        {stageLabel(stage.id, language)}
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
                  {t("admin.monitor.detail.current_stage", {
                    defaultValue: "Current Stage: {{stage}}",
                    stage: stageLabel(student.current_stage, language) || "—",
                  })}
                </h3>
                <Badge variant="outline" className="text-sm">
                  {t("admin.monitor.detail.stage_nodes_badge", {
                    defaultValue: "{{count}} nodes",
                    count: stageNodes.length,
                  })}
                </Badge>
              </div>

              <div className="grid gap-4">
                {stageNodes.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    {t("admin.monitor.detail.stage_empty", {
                      defaultValue: "No nodes available for this stage",
                    })}
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
              <h3 className="text-lg font-semibold mb-4">
                {t("admin.monitor.detail.comments.title", {
                  defaultValue: "Comments & Notes",
                })}
              </h3>
              <div className="space-y-3">
                <Textarea
                  placeholder={t("admin.monitor.detail.comments.placeholder", {
                    defaultValue: "Add a comment... Use @ to mention advisors",
                  })}
                  className="min-h-[100px]"
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
                />
                <div className="flex justify-end gap-2">
                  <Button size="sm" variant="outline">
                    {t("admin.monitor.detail.comments.attach", {
                      defaultValue: "Attach File",
                    })}
                  </Button>
                  <Button size="sm" disabled={!comment.trim()}>
                    {t("admin.monitor.detail.comments.add", {
                      defaultValue: "Add Comment",
                    })}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Documents & Review */}
          <Card className="border shadow-sm">
            <CardContent className="p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold">
                  {t("admin.review.title", {
                    defaultValue: "Documents & Review",
                  })}
                </h3>
                <p className="text-sm text-muted-foreground mt-1">
                  {t("admin.review.subtitle", {
                    defaultValue: "Review and approve student submissions",
                  })}
                </p>
              </div>

              {/* Node Selection Cards */}
              {stageNodes.length > 0 && (
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
                  {stageNodes.map((node: any) => {
                    const nodeId = node.node_id;
                    const isSelected = selectedNodeId === nodeId;
                    const nodeData = nodes.find(
                      (n: any) => n.node_id === nodeId
                    );
                    const filesList = nodeData?.files || [];

                    // Calculate status indicators
                    const allApproved =
                      filesList.length > 0 &&
                      filesList.every((f: any) => f.status === "approved");
                    const hasRejected = filesList.some(
                      (f: any) => f.status === "rejected"
                    );
                    const hasPending = filesList.some(
                      (f: any) => f.status === "pending"
                    );
                    const isEmpty = filesList.length === 0;

                    return (
                      <button
                        key={nodeId}
                        onClick={() => setSelectedNodeId(nodeId)}
                        className={`p-4 rounded-lg border-2 text-left transition-all ${
                          isSelected
                            ? "border-primary bg-primary/5 shadow-sm"
                            : "border-border hover:border-primary/50 hover:shadow-sm"
                        }`}
                      >
                        <div className="flex items-start gap-3">
                          <div
                            className={`p-2 rounded-lg flex-shrink-0 ${
                              allApproved
                                ? "bg-emerald-100 text-emerald-700"
                                : hasRejected
                                ? "bg-rose-100 text-rose-700"
                                : hasPending
                                ? "bg-amber-100 text-amber-700"
                                : "bg-muted text-muted-foreground"
                            }`}
                          >
                            {allApproved ? (
                              <CheckCircle2 className="h-4 w-4" />
                            ) : hasRejected ? (
                              <AlertCircle className="h-4 w-4" />
                            ) : hasPending ? (
                              <FileText className="h-4 w-4" />
                            ) : (
                              <File className="h-4 w-4" />
                            )}
                          </div>
                          <div className="flex-1 min-w-0">
                            <code className="text-xs bg-muted px-1.5 py-0.5 rounded font-mono block mb-1">
                              {nodeId}
                            </code>
                            <h4 className="font-medium text-sm truncate">
                              {nodeTitles[nodeId] || nodeId}
                            </h4>
                            <p className="text-xs text-muted-foreground mt-1">
                              {isEmpty
                                ? t("admin.review.no_files", {
                                    defaultValue: "No files",
                                  })
                                : t("admin.review.files_count", {
                                    defaultValue: "{{count}} files",
                                    count: filesList.length,
                                  })}
                            </p>
                            {allApproved && (
                              <div className="mt-1">
                                <Badge className="text-xs bg-emerald-100 text-emerald-700 hover:bg-emerald-100 border-emerald-200">
                                  ✓{" "}
                                  {t("admin.review.approved", {
                                    defaultValue: "Approved",
                                  })}
                                </Badge>
                              </div>
                            )}
                            {hasRejected && (
                              <div className="mt-1">
                                <Badge
                                  variant="destructive"
                                  className="text-xs"
                                >
                                  !{" "}
                                  {t("admin.review.rejected", {
                                    defaultValue: "Rejected",
                                  })}
                                </Badge>
                              </div>
                            )}
                            {!allApproved && !hasRejected && hasPending && (
                              <div className="mt-1">
                                <Badge
                                  variant="outline"
                                  className="text-xs bg-amber-50 text-amber-700 border-amber-200"
                                >
                                  ⏱{" "}
                                  {t("admin.review.pending", {
                                    defaultValue: "Pending",
                                  })}
                                </Badge>
                              </div>
                            )}
                          </div>
                        </div>
                      </button>
                    );
                  })}
                </div>
              )}

              {reviewMessage && (
                <div
                  className={`rounded-lg border p-4 text-sm flex items-start gap-3 ${
                    reviewMessage.tone === "error"
                      ? "border-destructive/50 bg-destructive/5 text-destructive"
                      : "border-emerald-200 bg-emerald-50 text-emerald-700"
                  }`}
                >
                  {reviewMessage.tone === "error" ? (
                    <AlertCircle className="h-5 w-5 flex-shrink-0 mt-0.5" />
                  ) : (
                    <CheckCircle2 className="h-5 w-5 flex-shrink-0 mt-0.5" />
                  )}
                  <span>{reviewMessage.text}</span>
                </div>
              )}

              {!selectedNodeId ? (
                <div className="rounded-xl border-2 border-dashed bg-muted/20 p-12 text-center">
                  <FileText className="h-12 w-12 text-muted-foreground/50 mx-auto mb-4" />
                  <p className="text-sm font-medium text-muted-foreground">
                    {t("admin.review.hint", {
                      defaultValue:
                        "Select a node above to inspect uploaded files.",
                    })}
                  </p>
                </div>
              ) : nodeFilesLoading ? (
                <div className="flex items-center justify-center gap-3 py-12 text-muted-foreground">
                  <Loader2 className="h-6 w-6 animate-spin" />
                  <span className="text-sm font-medium">
                    {t("admin.review.loading", {
                      defaultValue: "Loading files...",
                    })}
                  </span>
                </div>
              ) : nodeFiles.length === 0 ? (
                <div className="rounded-xl border-2 border-dashed bg-muted/20 p-12 text-center">
                  <File className="h-12 w-12 text-muted-foreground/50 mx-auto mb-4" />
                  <h4 className="font-semibold text-foreground mb-2">
                    {t("admin.review.empty_title", {
                      defaultValue: "No submissions yet",
                    })}
                  </h4>
                  <p className="text-sm text-muted-foreground">
                    {t("admin.review.empty", {
                      defaultValue:
                        "Student hasn't uploaded any files for this node.",
                    })}
                  </p>
                </div>
              ) : (
                <div className="space-y-6">
                  {Object.entries(
                    nodeFiles.reduce(
                      (acc: Record<string, typeof nodeFiles>, f) => {
                        const k = (f as any).slot_key || "default";
                        (acc[k] ||= []).push(f);
                        return acc;
                      },
                      {}
                    )
                  ).map(([slot, files]) => {
                    const meta = slotLabels[slot];
                    const allApproved = files.every(
                      (f: any) => f.status === "approved"
                    );
                    const hasRejected = files.some(
                      (f: any) => f.status === "rejected"
                    );

                    return (
                      <div key={slot} className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div
                              className={`p-2 rounded-lg ${
                                allApproved
                                  ? "bg-emerald-100 text-emerald-700"
                                  : hasRejected
                                  ? "bg-rose-100 text-rose-700"
                                  : "bg-amber-100 text-amber-700"
                              }`}
                            >
                              {allApproved ? (
                                <FileCheck className="h-5 w-5" />
                              ) : hasRejected ? (
                                <AlertCircle className="h-5 w-5" />
                              ) : (
                                <FileText className="h-5 w-5" />
                              )}
                            </div>
                            <div>
                              <h4 className="text-base font-semibold">
                                {meta?.label || slot}
                              </h4>
                              <p className="text-xs text-muted-foreground">
                                {t("admin.monitor.detail.files_count", {
                                  defaultValue: "{{count}} files",
                                  count: files.length,
                                })}
                                {meta?.required && (
                                  <span>
                                    {" · "}
                                    {t("common.required", {
                                      defaultValue: "Required",
                                    })}
                                  </span>
                                )}
                              </p>
                            </div>
                          </div>
                          {meta?.required && (
                            <Badge
                              variant="outline"
                              className="text-xs bg-primary/5 text-primary border-primary/20"
                            >
                              {t("common.required", {
                                defaultValue: "Required",
                              })}
                            </Badge>
                          )}
                        </div>

                        <div className="space-y-3">
                          {files.map((file) => {
                            const statusBadge = attachmentStatuses[
                              file.status
                            ] || {
                              label: file.status,
                              className:
                                "bg-muted text-muted-foreground border border-border",
                            };
                            const working =
                              pendingAttachment === file.attachment_id &&
                              reviewMutation.isPending;
                            const statusLabel = t(
                              `uploads.status.${file.status}`,
                              {
                                defaultValue: statusBadge.label,
                              }
                            );
                            const uploadedByText = file.uploaded_by
                              ? t("admin.review.uploaded_by", {
                                  defaultValue: ` · ${file.uploaded_by}`,
                                  name: file.uploaded_by,
                                })
                              : "";
                            const uploadMeta = t("admin.review.uploaded_meta", {
                              defaultValue: `Uploaded ${formatDateLabel(
                                file.attached_at
                              )}${uploadedByText}`,
                              date: formatDateLabel(file.attached_at),
                              by: uploadedByText,
                            });

                            return (
                              <div
                                key={file.attachment_id}
                                className="group rounded-lg border bg-card hover:border-primary/40 hover:shadow-md transition-all duration-200"
                              >
                                <div className="p-4 space-y-4">
                                  <div className="flex items-start gap-4">
                                    <div
                                      className={`p-3 rounded-lg flex-shrink-0 ${
                                        file.status === "approved"
                                          ? "bg-emerald-50"
                                          : file.status === "rejected"
                                          ? "bg-rose-50"
                                          : "bg-amber-50"
                                      }`}
                                    >
                                      {file.mime_type?.includes("pdf") ? (
                                        <FileText className="h-6 w-6 text-red-600" />
                                      ) : (
                                        <File className="h-6 w-6 text-blue-600" />
                                      )}
                                    </div>

                                    <div className="flex-1 min-w-0">
                                      <div className="flex items-start justify-between gap-3 mb-2">
                                        <div className="flex-1 min-w-0">
                                          <h5 className="font-medium text-foreground truncate group-hover:text-primary transition-colors">
                                            {file.filename}
                                          </h5>
                                          <div className="flex flex-wrap items-center gap-2 mt-1 text-xs text-muted-foreground">
                                            <span>
                                              {formatBytes(file.size_bytes)}
                                            </span>
                                            <span>•</span>
                                            <span>{uploadMeta}</span>
                                          </div>
                                        </div>

                                        <div className="flex items-center gap-2 flex-shrink-0">
                                          <Badge
                                            variant="outline"
                                            className={`${statusBadge.className} font-medium`}
                                          >
                                            {file.status === "approved" && (
                                              <CheckCircle2 className="h-3 w-3 mr-1" />
                                            )}
                                            {file.status === "rejected" && (
                                              <AlertCircle className="h-3 w-3 mr-1" />
                                            )}
                                            {file.status === "submitted" && (
                                              <Clock className="h-3 w-3 mr-1" />
                                            )}
                                            {statusLabel}
                                          </Badge>

                                          <Button
                                            variant="ghost"
                                            size="icon"
                                            className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
                                            onClick={async () => {
                                              try {
                                                // Fetch the download URL which will redirect to presigned S3 URL
                                                const response = await fetch(
                                                  `${BASE_URL}${file.download_url}`,
                                                  {
                                                    headers: {
                                                      Authorization: `Bearer ${localStorage.getItem(
                                                        "token"
                                                      )}`,
                                                    },
                                                  }
                                                );

                                                // Download as blob to avoid redirect issues
                                                const blob =
                                                  await response.blob();
                                                const url =
                                                  URL.createObjectURL(blob);
                                                const a =
                                                  document.createElement("a");
                                                a.href = url;
                                                a.download = file.filename;
                                                document.body.appendChild(a);
                                                a.click();
                                                document.body.removeChild(a);
                                                URL.revokeObjectURL(url);
                                              } catch (err) {
                                                console.error(
                                                  "Download failed:",
                                                  err
                                                );
                                              }
                                            }}
                                            aria-label={t("uploads.download", {
                                              defaultValue: "Download",
                                            })}
                                          >
                                            <Download className="h-4 w-4" />
                                          </Button>
                                        </div>
                                      </div>

                                      {file.review_note && (
                                        <div className="rounded-md bg-amber-50 border border-amber-200 p-3 mt-3">
                                          <div className="flex items-start gap-2">
                                            <AlertTriangle className="h-4 w-4 text-amber-700 flex-shrink-0 mt-0.5" />
                                            <div className="flex-1 min-w-0">
                                              <p className="text-xs font-medium text-amber-900 mb-1">
                                                {t("uploads.reviewer_note", {
                                                  defaultValue:
                                                    "Reviewer feedback",
                                                })}
                                              </p>
                                              <p className="text-xs text-amber-800">
                                                {file.review_note}
                                              </p>
                                            </div>
                                          </div>
                                        </div>
                                      )}

                                      {file.approved_at && file.approved_by && (
                                        <div className="flex items-center gap-2 mt-2 text-xs text-muted-foreground">
                                          <CheckCircle2 className="h-3.5 w-3.5 text-emerald-600" />
                                          <span>
                                            Approved by {file.approved_by} on{" "}
                                            {formatDateLabel(file.approved_at)}
                                          </span>
                                        </div>
                                      )}

                                      <div className="flex flex-wrap gap-2 mt-4 pt-4 border-t">
                                        {file.status !== "approved" && (
                                          <Button
                                            size="sm"
                                            className="bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm"
                                            disabled={working}
                                            onClick={() =>
                                              handleApprove(file.attachment_id)
                                            }
                                          >
                                            {working ? (
                                              <Loader2 className="h-4 w-4 animate-spin mr-2" />
                                            ) : (
                                              <CheckCircle2 className="h-4 w-4 mr-2" />
                                            )}
                                            {t("admin.review.approve", {
                                              defaultValue: "Approve",
                                            })}
                                          </Button>
                                        )}
                                        {file.status !== "rejected" && (
                                          <Button
                                            size="sm"
                                            variant="outline"
                                            className="border-red-200 text-red-700 hover:bg-red-50 hover:text-red-800"
                                            disabled={working}
                                            onClick={() =>
                                              setReviewDialog({
                                                attachmentId:
                                                  file.attachment_id,
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
                                        <Button
                                          size="sm"
                                          variant="ghost"
                                          className="gap-2"
                                          onClick={async () => {
                                            try {
                                              const response = await fetch(
                                                `${BASE_URL}${file.download_url}`,
                                                {
                                                  headers: {
                                                    Authorization: `Bearer ${localStorage.getItem(
                                                      "token"
                                                    )}`,
                                                  },
                                                }
                                              );

                                              // The response.url will be the presigned URL after redirect
                                              if (
                                                response.url &&
                                                response.url !==
                                                  `${BASE_URL}${file.download_url}`
                                              ) {
                                                // We got redirected to presigned URL
                                                // For PDF files, open in browser. For others, trigger download
                                                const isPdf =
                                                  file.mime_type?.includes(
                                                    "pdf"
                                                  ) ||
                                                  file.filename
                                                    .toLowerCase()
                                                    .endsWith(".pdf");

                                                if (isPdf) {
                                                  // Open PDF in new tab for inline viewing
                                                  window.open(
                                                    response.url,
                                                    "_blank"
                                                  );
                                                } else {
                                                  // For Word docs, download the file
                                                  const link =
                                                    document.createElement("a");
                                                  link.href = response.url;
                                                  link.download = file.filename;
                                                  link.target = "_blank";
                                                  document.body.appendChild(
                                                    link
                                                  );
                                                  link.click();
                                                  document.body.removeChild(
                                                    link
                                                  );
                                                }
                                              } else if (response.ok) {
                                                // Direct file response, create blob
                                                const blob =
                                                  await response.blob();
                                                const url =
                                                  URL.createObjectURL(blob);
                                                const isPdf =
                                                  file.mime_type?.includes(
                                                    "pdf"
                                                  ) ||
                                                  file.filename
                                                    .toLowerCase()
                                                    .endsWith(".pdf");

                                                if (isPdf) {
                                                  window.open(url, "_blank");
                                                } else {
                                                  const link =
                                                    document.createElement("a");
                                                  link.href = url;
                                                  link.download = file.filename;
                                                  document.body.appendChild(
                                                    link
                                                  );
                                                  link.click();
                                                  document.body.removeChild(
                                                    link
                                                  );
                                                }
                                                setTimeout(
                                                  () =>
                                                    URL.revokeObjectURL(url),
                                                  10000
                                                );
                                              }
                                            } catch (err) {
                                              console.error(
                                                "Preview failed:",
                                                err
                                              );
                                            }
                                          }}
                                        >
                                          <Eye className="h-4 w-4" />
                                          {t("common.preview", {
                                            defaultValue: "Preview",
                                          })}
                                        </Button>
                                      </div>
                                    </div>
                                  </div>
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </main>
      <Modal
        open={!!reviewDialog}
        onClose={() => {
          setReviewDialog(null);
          setReviewNote("");
          setReviewedFile(null);
        }}
      >
        <div className="space-y-4">
          <div>
            <h4 className="text-lg font-semibold">
              {t("admin.review.modal_title", {
                defaultValue: "Request changes",
              })}
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

          {/* Optional file upload section */}
          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">
              {t("admin.review.upload_edited", {
                defaultValue: "Upload edited document (optional)",
              })}
            </label>
            <p className="text-xs text-muted-foreground">
              {t("admin.review.upload_hint", {
                defaultValue: "Upload a corrected version with your comments",
              })}
            </p>
            {reviewedFile ? (
              <div className="flex items-center gap-2 p-3 rounded-lg border bg-muted/50">
                <FileText className="h-4 w-4 text-muted-foreground flex-shrink-0" />
                <span className="text-sm flex-1 truncate">
                  {reviewedFile.name}
                </span>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0"
                  onClick={() => setReviewedFile(null)}
                >
                  ✕
                </Button>
              </div>
            ) : (
              <Button
                variant="outline"
                size="sm"
                className="w-full"
                onClick={() => {
                  const input = document.createElement("input");
                  input.type = "file";
                  input.accept = ".pdf,.doc,.docx";
                  input.onchange = (e) => {
                    const file = (e.target as HTMLInputElement).files?.[0];
                    if (file) setReviewedFile(file);
                  };
                  input.click();
                }}
              >
                <Download className="h-4 w-4 mr-2" />
                {t("admin.review.choose_file", {
                  defaultValue: "Choose file",
                })}
              </Button>
            )}
          </div>

          <div className="flex justify-end gap-2">
            <Button
              variant="ghost"
              onClick={() => {
                setReviewDialog(null);
                setReviewNote("");
                setReviewedFile(null);
              }}
            >
              {t("common.cancel", { defaultValue: "Cancel" })}
            </Button>
            <Button
              onClick={submitRequestChanges}
              disabled={
                !reviewNote.trim() ||
                isUploadingReviewed ||
                (pendingAttachment === reviewDialog?.attachmentId &&
                  reviewMutation.isPending)
              }
            >
              {isUploadingReviewed ||
              (pendingAttachment === reviewDialog?.attachmentId &&
                reviewMutation.isPending) ? (
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
  const { t, i18n } = useTranslation("common");
  const locale = i18n.language || "en";
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
              {t(`admin.monitor.detail.node_states.${state}`, {
                defaultValue: stateConfig.label,
              })}
            </Badge>
          </div>
          <div className="text-sm font-medium text-foreground truncate">
            {title || "—"}
          </div>
          {dueDate && (
            <div className="mt-1 flex items-center gap-1 text-xs text-muted-foreground">
              <Calendar className="h-3.5 w-3.5" />
              {t("admin.monitor.detail.node_card.due", {
                defaultValue: "Due:",
              })}{" "}
              {new Date(dueDate).toLocaleString(locale)}
            </div>
          )}
        </div>
        <div className="w-48 ml-4">
          <label className="block text-xs text-muted-foreground mb-1">
            {t("admin.monitor.detail.node_card.set_label", {
              defaultValue: "Set deadline",
            })}
          </label>
          <input
            type="datetime-local"
            aria-label={t("admin.monitor.detail.node_card.set_aria", {
              defaultValue: "Set due date",
            })}
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
