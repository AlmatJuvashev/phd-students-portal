"use client"

import { useRouter } from "next/navigation"
import {
  ArrowLeft,
  Mail,
  Phone,
  Copy,
  CheckCircle2,
  Circle,
  Clock3,
  AlertTriangle,
  Lock,
  Send,
  Calendar,
  FileText,
  ExternalLink,
  Plus,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Progress } from "@/components/ui/progress"
import { Separator } from "@/components/ui/separator"
import { Textarea } from "@/components/ui/textarea"

type NodeState = "locked" | "active" | "submitted" | "waiting" | "needs_fixes" | "done"
type Language = "EN" | "RU" | "KZ"
type DocumentStatus = "pending_review" | "reviewed" | "needs_revision"

interface Document {
  id: string
  filename: string
  uploadDate: string
  status: DocumentStatus
  reviewedBy?: string
  reviewDate?: string
  comments?: string
}

interface Student {
  id: string
  name: string
  avatar?: string
  program: string
  department: string
  advisors: string[]
  cohort: string
  currentStage: string
  stageProgress: { current: number; total: number }
  overallProgress: number
  dueNext: string
  overdue: boolean
  lastUpdate: string
  rpRequired: boolean
  status: "normal" | "at-risk" | "overdue" | "completed"
  email: string
  phone: string
  documents?: Document[]
  stageHistory?: Array<{ stage: string; status: "completed" | "in-progress" | "at-risk" | "pending" }>
}

const stages = [
  { id: "W1", label: { EN: "I — Preparation", RU: "I — Подготовка", KZ: "I — Дайындық" } },
  {
    id: "W2",
    label: {
      EN: "II — Pre-examination (SC)",
      RU: "II — Предварительная экспертиза (НК)",
      KZ: "II — Алдын ала сараптама (ҒК)",
    },
  },
  {
    id: "W3",
    label: { EN: "III — RP (conditional)", RU: "III — RP (условно)", KZ: "III — RP (шартты)" },
    conditional: true,
  },
  {
    id: "W4",
    label: { EN: "IV — Submission to DC", RU: "IV — Подача в ДС", KZ: "IV — Диссертациялық кеңеске тапсыру" },
  },
  { id: "W5", label: { EN: "V — Restoration", RU: "V — Восстановление", KZ: "V — Дооформление" } },
  {
    id: "W6",
    label: { EN: "VI — After DC acceptance", RU: "VI — После принятия ДС", KZ: "VI — ДК қабылдағаннан кейін" },
  },
  {
    id: "W7",
    label: {
      EN: "VII — Defense & Post-defense",
      RU: "VII — Защита и После защиты",
      KZ: "VII — Қорғау және Қорғаудан кейін",
    },
  },
]

const nodeStates: Record<NodeState, { label: string; color: string; icon: any }> = {
  locked: { label: "Locked", color: "bg-muted text-muted-foreground", icon: Lock },
  active: { label: "Active", color: "bg-clinical-accent text-clinical-accent-foreground", icon: Circle },
  submitted: { label: "Submitted", color: "bg-clinical-info text-clinical-info-foreground", icon: Send },
  waiting: {
    label: "Waiting",
    color: "bg-clinical-warning/20 text-clinical-warning-foreground border border-clinical-warning/40",
    icon: Clock3,
  },
  needs_fixes: {
    label: "Needs Fixes",
    color: "bg-clinical-alert/20 text-clinical-alert-foreground border border-clinical-alert/40",
    icon: AlertTriangle,
  },
  done: { label: "Done", color: "bg-clinical-success text-clinical-success-foreground", icon: CheckCircle2 },
}

const documentStatusConfig: Record<
  DocumentStatus,
  { label: string; color: string; bgColor: string; borderColor: string }
> = {
  pending_review: {
    label: "Pending Review",
    color: "text-clinical-warning",
    bgColor: "bg-clinical-warning/10",
    borderColor: "border-clinical-warning/30",
  },
  reviewed: {
    label: "Reviewed",
    color: "text-clinical-success",
    bgColor: "bg-clinical-success/10",
    borderColor: "border-clinical-success/30",
  },
  needs_revision: {
    label: "Needs Revision",
    color: "text-clinical-alert",
    bgColor: "bg-clinical-alert/10",
    borderColor: "border-clinical-alert/30",
  },
}

export default function StudentDetailClient({ student }: { student: Student }) {
  const router = useRouter()
  const language: Language = "EN" // In production, this would come from context/state

  const getStageLabel = (stageId: string) => {
    const stage = stages.find((s) => s.id === stageId)
    return stage ? stage.label[language] : stageId
  }

  const handleDocumentReview = (documentId: string, newStatus: DocumentStatus) => {
    console.log(`[v0] Marking document ${documentId} for student ${student.id} as ${newStatus}`)
    // In production, this would update the backend
  }

  return (
    <div className="min-h-screen bg-clinical-background">
      {/* Header with Back Button */}
      <header className="sticky top-0 z-50 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80 border-b border-clinical-border">
        <div className="px-8 py-5">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => router.push("/")}
              className="gap-2 hover:bg-clinical-hover"
            >
              <ArrowLeft className="h-4 w-4" />
              Back to Students
            </Button>
            <Separator orientation="vertical" className="h-6" />
            <h1 className="text-xl font-semibold text-foreground">Student Details</h1>
            <Badge variant="outline" className="bg-clinical-primary/5 text-clinical-primary border-clinical-primary/20">
              {student.id}
            </Badge>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="px-8 py-8 max-w-6xl mx-auto">
        <div className="space-y-8">
          {/* Student Profile Header */}
          <Card className="border-clinical-border bg-clinical-card shadow-clinical">
            <CardContent className="p-8">
              <div className="flex items-start gap-6 mb-6">
                <Avatar className="h-24 w-24 border-2 border-clinical-border">
                  <AvatarImage src={student.avatar || "/placeholder.svg"} />
                  <AvatarFallback className="bg-clinical-primary/10 text-clinical-primary text-2xl">
                    {student.name
                      .split(" ")
                      .map((n) => n[0])
                      .join("")}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1">
                  <h2 className="text-3xl font-semibold text-foreground mb-2">{student.name}</h2>
                  <div className="flex flex-wrap gap-2 mb-4">
                    <Badge
                      variant="outline"
                      className="bg-clinical-primary/5 text-clinical-primary border-clinical-primary/20"
                    >
                      {student.program}
                    </Badge>
                    <Badge variant="outline" className="bg-clinical-muted/20 border-clinical-border">
                      {student.department}
                    </Badge>
                    <Badge variant="outline" className="bg-clinical-muted/20 border-clinical-border">
                      {student.cohort}
                    </Badge>
                    {student.rpRequired && (
                      <Badge
                        variant="outline"
                        className="bg-clinical-warning/10 text-clinical-warning border-clinical-warning/30"
                      >
                        RP Required
                      </Badge>
                    )}
                  </div>

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <Mail className="h-4 w-4" />
                      <a href={`mailto:${student.email}`} className="hover:text-clinical-primary">
                        {student.email}
                      </a>
                    </div>
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <Phone className="h-4 w-4" />
                      <a href={`tel:${student.phone}`} className="hover:text-clinical-primary">
                        {student.phone}
                      </a>
                    </div>
                  </div>
                </div>

                <div className="flex flex-col gap-2">
                  <Button size="sm" className="gap-2 bg-clinical-primary hover:bg-clinical-primary/90">
                    <Mail className="h-4 w-4" />
                    Send Email
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    className="gap-2 border-clinical-border hover:bg-clinical-hover bg-transparent"
                  >
                    <Phone className="h-4 w-4" />
                    Call
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    className="gap-2 border-clinical-border hover:bg-clinical-hover bg-transparent"
                  >
                    <Copy className="h-4 w-4" />
                    Copy Link
                  </Button>
                </div>
              </div>

              <Separator className="mb-6" />

              <div className="grid grid-cols-3 gap-6">
                <div>
                  <div className="text-sm text-muted-foreground mb-1">Advisors</div>
                  <div className="flex flex-wrap gap-1">
                    {student.advisors.map((advisor, idx) => (
                      <Badge
                        key={idx}
                        variant="outline"
                        className="text-xs bg-clinical-muted/20 border-clinical-border"
                      >
                        {advisor}
                      </Badge>
                    ))}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground mb-1">Overall Progress</div>
                  <div className="flex items-center gap-3">
                    <Progress value={student.overallProgress} className="h-2 flex-1" />
                    <span className="text-lg font-semibold text-foreground">{student.overallProgress}%</span>
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground mb-1">Status</div>
                  <Badge
                    variant="outline"
                    className={`text-xs ${
                      student.status === "overdue"
                        ? "bg-clinical-alert/10 text-clinical-alert border-clinical-alert/30"
                        : student.status === "at-risk"
                          ? "bg-clinical-warning/10 text-clinical-warning border-clinical-warning/30"
                          : student.status === "completed"
                            ? "bg-clinical-success/10 text-clinical-success border-clinical-success/30"
                            : "bg-clinical-muted/20 border-clinical-border"
                    }`}
                  >
                    {student.status.toUpperCase()}
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Journey Map */}
          <Card className="border-clinical-border bg-clinical-card shadow-clinical">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold text-foreground mb-6">Dissertation Journey Map</h3>
              <div className="flex items-center gap-2 overflow-x-auto pb-2">
                {stages
                  .filter((stage) => !stage.conditional || student.rpRequired)
                  .map((stage, idx, arr) => (
                    <div key={stage.id} className="flex items-center flex-shrink-0">
                      <div
                        className={`px-4 py-3 rounded-lg text-sm font-medium whitespace-nowrap transition-all ${
                          stage.id === student.currentStage
                            ? "bg-clinical-primary text-clinical-primary-foreground shadow-md scale-105"
                            : stages.findIndex((s) => s.id === student.currentStage) > idx
                              ? "bg-clinical-success/20 text-clinical-success"
                              : "bg-clinical-muted/30 text-muted-foreground"
                        }`}
                      >
                        {stage.label[language]}
                      </div>
                      {idx < arr.length - 1 && <div className="w-12 h-0.5 bg-clinical-border mx-2 flex-shrink-0" />}
                    </div>
                  ))}
              </div>
            </CardContent>
          </Card>

          {/* Stage Progress and Checklist */}
          <Card className="border-clinical-border bg-clinical-card shadow-clinical">
            <CardContent className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold text-foreground">
                  Current Stage: {getStageLabel(student.currentStage)}
                </h3>
                <Badge variant="outline" className="text-sm bg-clinical-muted/20 border-clinical-border">
                  {student.stageProgress.current}/{student.stageProgress.total} nodes completed
                </Badge>
              </div>

              <div className="grid gap-4">
                {student.currentStage === "W1" && (
                  <>
                    <NodeCard
                      id="S1_profile"
                      title={{ EN: "Student Profile", RU: "Профиль студента", KZ: "Студент профилі" }}
                      state="done"
                      language={language}
                      dueDate="2024-12-15"
                    />
                    <NodeCard
                      id="S1_text_ready"
                      title={{ EN: "Text Ready", RU: "Текст готов", KZ: "Мәтін дайын" }}
                      state="done"
                      language={language}
                      dueDate="2025-01-05"
                    />
                    <NodeCard
                      id="S1_antiplag"
                      title={{ EN: "Antiplagiarism ≥85%", RU: "Антиплагиат ≥85%", KZ: "Антиплагиат ≥85%" }}
                      state="done"
                      language={language}
                      dueDate="2025-01-10"
                      hasFiles
                    />
                    <NodeCard
                      id="S1_publications_list"
                      title={{ EN: "Publications List", RU: "Список публикаций", KZ: "Жарияланымдар тізімі" }}
                      state="active"
                      language={language}
                      dueDate="2025-02-01"
                    />
                  </>
                )}
                {student.currentStage === "W2" && (
                  <>
                    <NodeCard
                      id="E1_apply_omid"
                      title={{ EN: "Apply to OMID", RU: "Подать заявку в ОМИД", KZ: "ОМИД-ке өтініш беру" }}
                      state="done"
                      language={language}
                      dueDate="2025-01-05"
                    />
                    <NodeCard
                      id="NK_package"
                      title={{ EN: "SC Package", RU: "Пакет НК", KZ: "ҒК пакеті" }}
                      state="submitted"
                      language={language}
                      dueDate="2025-01-15"
                    />
                    <NodeCard
                      id="E3_hearing_nk"
                      title={{ EN: "SC Hearing", RU: "Слушание НК", KZ: "ҒК тыңдауы" }}
                      state="waiting"
                      language={language}
                      dueDate="2025-01-20"
                    />
                  </>
                )}
                {student.currentStage === "W4" && (
                  <>
                    <NodeCard
                      id="D1_normokontrol_ncste"
                      title={{ EN: "Normokontrol NCSTE", RU: "Нормоконтроль НЦНТИ", KZ: "ҒҚБЖ нормобақылау" }}
                      state="submitted"
                      language={language}
                      dueDate="2025-01-12"
                    />
                    <NodeCard
                      id="IV_rector_application"
                      title={{ EN: "Rector Application", RU: "Заявление ректору", KZ: "Ректорға өтініш" }}
                      state="done"
                      language={language}
                      dueDate="2025-01-08"
                    />
                    <NodeCard
                      id="IV3_publication_certificate_ncste"
                      title={{
                        EN: "Publication Certificate NCSTE",
                        RU: "Справка о публикациях НЦНТИ",
                        KZ: "ҒҚБЖ жариялау туралы анықтама",
                      }}
                      state="waiting"
                      language={language}
                      dueDate="2025-01-15"
                    />
                    <NodeCard
                      id="D2_apply_to_ds"
                      title={{ EN: "Apply to DC", RU: "Подача в ДС", KZ: "ДК-ға тапсыру" }}
                      state="locked"
                      language={language}
                      dueDate="2025-01-25"
                    />
                  </>
                )}
                {student.currentStage === "W6" && (
                  <>
                    <NodeCard
                      id="A1_post_acceptance_overview"
                      title={{
                        EN: "Post-Acceptance Overview",
                        RU: "Обзор после принятия",
                        KZ: "Қабылдаудан кейінгі шолу",
                      }}
                      state="active"
                      language={language}
                      dueDate="2025-01-18"
                    />
                  </>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Documents & Review */}
          <Card className="border-clinical-border bg-clinical-card shadow-clinical">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold text-foreground mb-6">Documents & Review</h3>
              {student.documents && student.documents.length > 0 ? (
                <div className="space-y-4">
                  {student.documents.map((doc) => {
                    const statusConfig = documentStatusConfig[doc.status]
                    return (
                      <div
                        key={doc.id}
                        className="p-5 rounded-lg border border-clinical-border bg-background space-y-4"
                      >
                        <div className="flex items-start justify-between">
                          <div className="flex items-start gap-3 flex-1">
                            <FileText className="h-5 w-5 text-clinical-primary mt-0.5" />
                            <div className="flex-1 min-w-0">
                              <div className="text-base font-medium text-foreground mb-1">{doc.filename}</div>
                              <div className="text-sm text-muted-foreground">
                                Uploaded {new Date(doc.uploadDate).toLocaleDateString()}
                              </div>
                              {doc.reviewedBy && (
                                <div className="text-sm text-muted-foreground mt-1">
                                  Reviewed by {doc.reviewedBy} on {new Date(doc.reviewDate!).toLocaleDateString()}
                                </div>
                              )}
                              {doc.comments && (
                                <div className="mt-3 p-3 rounded bg-clinical-muted/10 text-sm text-foreground">
                                  <strong>Review Comments:</strong> {doc.comments}
                                </div>
                              )}
                            </div>
                          </div>
                          <Badge
                            variant="outline"
                            className={`text-xs ${statusConfig.bgColor} ${statusConfig.color} border ${statusConfig.borderColor} ml-2`}
                          >
                            {statusConfig.label}
                          </Badge>
                        </div>

                        {doc.status === "pending_review" && (
                          <div className="flex gap-2 pt-3 border-t border-clinical-border">
                            <Button
                              size="sm"
                              className="flex-1 bg-clinical-success hover:bg-clinical-success/90 text-white"
                              onClick={() => handleDocumentReview(doc.id, "reviewed")}
                            >
                              <CheckCircle2 className="h-4 w-4 mr-2" />
                              Approve Document
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              className="flex-1 border-clinical-alert/30 text-clinical-alert hover:bg-clinical-alert/10 bg-transparent"
                              onClick={() => handleDocumentReview(doc.id, "needs_revision")}
                            >
                              <AlertTriangle className="h-4 w-4 mr-2" />
                              Request Changes
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              className="gap-2 hover:bg-clinical-hover"
                              onClick={() => window.open("#", "_blank")}
                            >
                              <ExternalLink className="h-4 w-4" />
                            </Button>
                          </div>
                        )}

                        {doc.status === "needs_revision" && (
                          <div className="flex gap-2 pt-3 border-t border-clinical-border">
                            <Button
                              size="sm"
                              className="flex-1 bg-clinical-success hover:bg-clinical-success/90 text-white"
                              onClick={() => handleDocumentReview(doc.id, "reviewed")}
                            >
                              <CheckCircle2 className="h-4 w-4 mr-2" />
                              Mark as Reviewed
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              className="gap-2 hover:bg-clinical-hover"
                              onClick={() => window.open("#", "_blank")}
                            >
                              <ExternalLink className="h-4 w-4 mr-2" />
                            </Button>
                          </div>
                        )}

                        {doc.status === "reviewed" && (
                          <div className="flex gap-2 pt-3 border-t border-clinical-border">
                            <Button
                              size="sm"
                              variant="outline"
                              className="flex-1 border-clinical-border hover:bg-clinical-hover bg-transparent"
                              onClick={() => window.open("#", "_blank")}
                            >
                              <ExternalLink className="h-4 w-4 mr-2" />
                              View Document
                            </Button>
                          </div>
                        )}
                      </div>
                    )
                  })}
                </div>
              ) : (
                <div className="text-sm text-muted-foreground p-8 text-center border border-dashed border-clinical-border rounded-lg">
                  No documents uploaded yet
                </div>
              )}
            </CardContent>
          </Card>

          {/* Comments & Notes */}
          <Card className="border-clinical-border bg-clinical-card shadow-clinical">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold text-foreground mb-6">Comments & Notes</h3>
              <div className="space-y-4 mb-6">
                <div className="p-4 rounded-lg border border-clinical-border bg-background">
                  <div className="flex items-center gap-3 mb-2">
                    <Avatar className="h-8 w-8">
                      <AvatarFallback className="bg-clinical-primary/10 text-clinical-primary text-xs">
                        AP
                      </AvatarFallback>
                    </Avatar>
                    <span className="text-sm font-medium text-foreground">Admin Petrov</span>
                    <span className="text-xs text-muted-foreground">Jan 10, 14:30</span>
                  </div>
                  <p className="text-sm text-foreground">
                    Please ensure antiplagiarism certificate is uploaded by Jan 12th.
                  </p>
                </div>
              </div>
              <div className="space-y-3">
                <Textarea
                  placeholder="Add a comment... Use @ to mention advisors"
                  className="min-h-[100px] bg-background border-clinical-border"
                />
                <div className="flex justify-end gap-2">
                  <Button size="sm" variant="outline" className="border-clinical-border bg-transparent">
                    Attach File
                  </Button>
                  <Button size="sm" className="bg-clinical-primary hover:bg-clinical-primary/90">
                    Add Comment
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Deadlines & Reminders */}
          <Card className="border-clinical-border bg-clinical-card shadow-clinical">
            <CardContent className="p-6">
              <h3 className="text-lg font-semibold text-foreground mb-6">Deadlines & Reminders</h3>
              <div className="space-y-3">
                <div className="flex items-center justify-between p-4 rounded-lg border border-clinical-border bg-background">
                  <div className="flex items-center gap-3">
                    <Calendar className="h-5 w-5 text-clinical-primary" />
                    <div>
                      <div className="text-sm font-medium text-foreground">Next Due: {student.dueNext}</div>
                      <div className="text-xs text-muted-foreground">Publications list submission</div>
                    </div>
                  </div>
                  {student.overdue && (
                    <Badge
                      variant="outline"
                      className="bg-clinical-alert/10 text-clinical-alert border-clinical-alert/30"
                    >
                      Overdue
                    </Badge>
                  )}
                </div>
                <Button
                  size="sm"
                  variant="outline"
                  className="w-full gap-2 border-clinical-border hover:bg-clinical-hover bg-transparent"
                >
                  <Plus className="h-4 w-4" />
                  Add New Reminder
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  )
}

function NodeCard({
  id,
  title,
  state,
  language,
  dueDate,
  hasFiles,
}: {
  id: string
  title: { EN: string; RU: string; KZ: string }
  state: NodeState
  language: Language
  dueDate?: string
  hasFiles?: boolean
}) {
  const stateConfig = nodeStates[state]
  const StateIcon = stateConfig.icon

  return (
    <div className="p-5 rounded-lg border border-clinical-border bg-background hover:shadow-md transition-all">
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <code className="text-xs text-muted-foreground bg-clinical-muted/20 px-2 py-1 rounded">{id}</code>
            <Badge variant="outline" className={`text-xs ${stateConfig.color}`}>
              <StateIcon className="h-3 w-3 mr-1" />
              {stateConfig.label}
            </Badge>
          </div>
          <h4 className="text-base font-medium text-foreground">{title[language]}</h4>
        </div>
      </div>
      <div className="flex items-center gap-4 text-sm text-muted-foreground">
        {dueDate && (
          <div className="flex items-center gap-1">
            <Calendar className="h-4 w-4" />
            Due: {dueDate}
          </div>
        )}
        {hasFiles && (
          <div className="flex items-center gap-1">
            <FileText className="h-4 w-4" />2 files
          </div>
        )}
      </div>
    </div>
  )
}
