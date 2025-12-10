"use client"

import { useState } from "react"
import { Search, Calendar, Clock, AlertCircle, Users, CheckCircle2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Progress } from "@/components/ui/progress"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Separator } from "@/components/ui/separator"
import { Checkbox } from "@/components/ui/checkbox"
import { useRouter } from "next/navigation" // Import useRouter

// Mock data types
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

const mockStudents: Student[] = [
  {
    id: "S001",
    name: "Айгерим Нұрғалиева",
    avatar: "/professional-woman.png",
    program: "PhD Clinical Medicine",
    department: "Cardiology",
    advisors: ["Dr. Petrov A.V.", "Prof. Kim S.Y."],
    cohort: "2022-2025",
    currentStage: "W2",
    stageProgress: { current: 2, total: 3 },
    overallProgress: 45,
    dueNext: "2025-01-20",
    overdue: false,
    lastUpdate: "2025-01-10 14:32",
    rpRequired: false,
    status: "normal",
    email: "a.nurgaliyeva@med.edu.kz",
    phone: "+7 701 234 5678",
    documents: [
      {
        id: "doc1",
        filename: "Antiplagiarism_Certificate.pdf",
        uploadDate: "2025-01-10",
        status: "reviewed",
        reviewedBy: "Dr. Petrov A.V.",
        reviewDate: "2025-01-11",
      },
      {
        id: "doc2",
        filename: "Publications_List.pdf",
        uploadDate: "2025-01-08",
        status: "pending_review",
      },
    ],
    stageHistory: [
      { stage: "W1", status: "completed" },
      { stage: "W2", status: "in-progress" },
      { stage: "W3", status: "pending" },
      { stage: "W4", status: "pending" },
      { stage: "W5", status: "pending" },
      { stage: "W6", status: "pending" },
      { stage: "W7", status: "pending" },
    ],
  },
  {
    id: "S002",
    name: "Dmitry Sokolov",
    avatar: "/man-professional.jpg",
    program: "DBA Healthcare Management",
    department: "Public Health",
    advisors: ["Prof. Akhmetova Z.K."],
    cohort: "2020-2024",
    currentStage: "W4",
    stageProgress: { current: 3, total: 4 },
    overallProgress: 72,
    dueNext: "2025-01-15",
    overdue: true,
    lastUpdate: "2025-01-08 09:15",
    rpRequired: true,
    status: "overdue",
    email: "d.sokolov@med.edu.kz",
    phone: "+7 702 345 6789",
    documents: [
      {
        id: "doc3",
        filename: "Normokontrol_Report.pdf",
        uploadDate: "2025-01-12",
        status: "needs_revision",
        reviewedBy: "Prof. Akhmetova Z.K.",
        reviewDate: "2025-01-13",
        comments: "Please fix formatting in section 3.2",
      },
    ],
    stageHistory: [
      { stage: "W1", status: "completed" },
      { stage: "W2", status: "completed" },
      { stage: "W3", status: "completed" },
      { stage: "W4", status: "at-risk" },
      { stage: "W5", status: "pending" },
      { stage: "W6", status: "pending" },
      { stage: "W7", status: "pending" },
    ],
  },
  {
    id: "S003",
    name: "Sara Omarova",
    avatar: "/woman-doctor.jpg",
    program: "PhD Neurology",
    department: "Neuroscience",
    advisors: ["Dr. Lee M.J.", "Prof. Ivanov P.S."],
    cohort: "2023-2026",
    currentStage: "W1",
    stageProgress: { current: 4, total: 4 },
    overallProgress: 28,
    dueNext: "2025-02-01",
    overdue: false,
    lastUpdate: "2025-01-11 16:45",
    rpRequired: false,
    status: "normal",
    email: "s.omarova@med.edu.kz",
    phone: "+7 705 456 7890",
    documents: [],
    stageHistory: [
      { stage: "W1", status: "in-progress" },
      { stage: "W2", status: "pending" },
      { stage: "W4", status: "pending" },
      { stage: "W5", status: "pending" },
      { stage: "W6", status: "pending" },
      { stage: "W7", status: "pending" },
    ],
  },
  {
    id: "S004",
    name: "Алексей Волков",
    avatar: "/man-scientist.jpg",
    program: "PhD Pharmacology",
    department: "Clinical Pharmacology",
    advisors: ["Prof. Chen W."],
    cohort: "2021-2024",
    currentStage: "W6",
    stageProgress: { current: 1, total: 1 },
    overallProgress: 88,
    dueNext: "2025-01-18",
    overdue: false,
    lastUpdate: "2025-01-12 11:20",
    rpRequired: true,
    status: "at-risk",
    email: "a.volkov@med.edu.kz",
    phone: "+7 707 567 8901",
    documents: [
      {
        id: "doc4",
        filename: "Defense_Presentation.pdf",
        uploadDate: "2025-01-09",
        status: "reviewed",
        reviewedBy: "Prof. Chen W.",
        reviewDate: "2025-01-10",
      },
    ],
    stageHistory: [
      { stage: "W1", status: "completed" },
      { stage: "W2", status: "completed" },
      { stage: "W3", status: "completed" },
      { stage: "W4", status: "completed" },
      { stage: "W5", status: "completed" },
      { stage: "W6", status: "in-progress" },
      { stage: "W7", status: "pending" },
    ],
  },
]

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
  locked: { label: "Locked", color: "bg-muted text-muted-foreground", icon: null },
  active: { label: "Active", color: "bg-clinical-accent text-clinical-accent-foreground", icon: null },
  submitted: { label: "Submitted", color: "bg-clinical-info text-clinical-info-foreground", icon: null },
  waiting: {
    label: "Waiting",
    color: "bg-clinical-warning/20 text-clinical-warning-foreground border border-clinical-warning/40",
    icon: null,
  },
  needs_fixes: {
    label: "Needs Fixes",
    color: "bg-clinical-alert/20 text-clinical-alert-foreground border border-clinical-alert/40",
    icon: null,
  },
  done: { label: "Done", color: "bg-clinical-success text-clinical-success-foreground", icon: null },
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

export default function StudentsProgress() {
  const [language, setLanguage] = useState<Language>("EN")
  // REMOVED: const [selectedStudent, setSelectedStudent] = useState<Student | null>(null)
  const [activeTab, setActiveTab] = useState("table")
  const [searchQuery, setSearchQuery] = useState("")
  const router = useRouter() // ADDED: router for navigation

  const getStageLabel = (stageId: string) => {
    const stage = stages.find((s) => s.id === stageId)
    return stage ? stage.label[language] : stageId
  }

  const filteredStudents = mockStudents.filter(
    (student) =>
      student.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      student.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
      student.email.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  // REMOVED: const handleDocumentReview = (studentId: string, documentId: string, newStatus: DocumentStatus) => {
  //   console.log(`[v0] Marking document ${documentId} for student ${studentId} as ${newStatus}`)
  //   // In production, this would update the backend
  // }

  return (
    <div className="min-h-screen bg-clinical-background">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80 border-b border-clinical-border">
        <div className="px-8 py-5">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <h1 className="text-2xl font-semibold text-foreground tracking-tight">Students Progress</h1>
              <Badge
                variant="outline"
                className="bg-clinical-primary/5 text-clinical-primary border-clinical-primary/20"
              >
                {filteredStudents.length} students
              </Badge>
            </div>

            <div className="flex items-center gap-3">
              {/* Language Toggle */}
              <Select value={language} onValueChange={(val) => setLanguage(val as Language)}>
                <SelectTrigger className="w-[100px] h-9 bg-background border-clinical-border">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="EN">English</SelectItem>
                  <SelectItem value="RU">Русский</SelectItem>
                  <SelectItem value="KZ">Қазақша</SelectItem>
                </SelectContent>
              </Select>

              <Separator orientation="vertical" className="h-6" />

              <Button
                variant="outline"
                size="sm"
                className="gap-2 bg-background border-clinical-border hover:bg-clinical-hover"
              >
                Export CSV
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="gap-2 bg-background border-clinical-border hover:bg-clinical-hover"
              >
                Bulk Message
              </Button>
              <Button
                size="sm"
                className="gap-2 bg-clinical-primary hover:bg-clinical-primary/90 text-clinical-primary-foreground"
              >
                New Reminder
              </Button>
            </div>
          </div>

          {/* Filters Bar */}
          <div className="mt-5 flex items-center gap-3 flex-wrap">
            <div className="relative flex-1 min-w-[320px]">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by name, ID, phone, or email..."
                className="pl-10 bg-background border-clinical-border h-10"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>

            <Select defaultValue="all-programs">
              <SelectTrigger className="w-[180px] bg-background border-clinical-border h-10">
                <SelectValue placeholder="Program" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all-programs">All Programs</SelectItem>
                <SelectItem value="phd">PhD Programs</SelectItem>
                <SelectItem value="dba">DBA Programs</SelectItem>
              </SelectContent>
            </Select>

            <Select defaultValue="all-stages">
              <SelectTrigger className="w-[180px] bg-background border-clinical-border h-10">
                <SelectValue placeholder="Stage" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all-stages">All Stages</SelectItem>
                {stages.map((stage) => (
                  <SelectItem key={stage.id} value={stage.id}>
                    {stage.label[language]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select defaultValue="all-cohorts">
              <SelectTrigger className="w-[150px] bg-background border-clinical-border h-10">
                <SelectValue placeholder="Cohort" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all-cohorts">All Cohorts</SelectItem>
                <SelectItem value="2023-2026">2023-2026</SelectItem>
                <SelectItem value="2022-2025">2022-2025</SelectItem>
                <SelectItem value="2021-2024">2021-2024</SelectItem>
                <SelectItem value="2020-2024">2020-2024</SelectItem>
              </SelectContent>
            </Select>

            <Button
              variant="outline"
              size="sm"
              className="gap-2 bg-background border-clinical-border hover:bg-clinical-hover h-10"
            >
              More Filters
            </Button>

            <div className="flex items-center gap-2 ml-2">
              <Checkbox id="rp-only" />
              <label htmlFor="rp-only" className="text-sm text-muted-foreground cursor-pointer">
                RP required only
              </label>
            </div>

            <div className="flex items-center gap-2">
              <Checkbox id="overdue-only" />
              <label htmlFor="overdue-only" className="text-sm text-muted-foreground cursor-pointer">
                Overdue only
              </label>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="px-8 py-6">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="bg-clinical-card border border-clinical-border mb-6">
            <TabsTrigger
              value="table"
              className="data-[state=active]:bg-clinical-primary data-[state=active]:text-clinical-primary-foreground"
            >
              Table View
            </TabsTrigger>
            <TabsTrigger
              value="kanban"
              className="data-[state=active]:bg-clinical-primary data-[state=active]:text-clinical-primary-foreground"
            >
              Kanban View
            </TabsTrigger>
            <TabsTrigger
              value="analytics"
              className="data-[state=active]:bg-clinical-primary data-[state=active]:text-clinical-primary-foreground"
            >
              Cohort Analytics
            </TabsTrigger>
          </TabsList>

          <TabsContent value="table" className="mt-0">
            <div className="bg-background rounded-lg border border-clinical-border overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full border-collapse">
                  <thead>
                    <tr className="bg-clinical-muted/10 border-b border-clinical-border">
                      <th className="text-left py-3 px-4 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                        Student
                      </th>
                      <th className="text-left py-3 px-4 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                        Program
                      </th>
                      <th className="text-left py-3 px-4 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                        Stage
                      </th>
                      <th className="text-left py-3 px-4 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                        Progress
                      </th>
                      <th className="text-center py-3 px-4 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                        Cohort
                      </th>
                      <th className="text-center py-3 px-4 text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                        Due
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {filteredStudents.map((student, index) => (
                      <tr
                        key={student.id}
                        className={`border-b border-clinical-border/50 hover:bg-clinical-hover/50 cursor-pointer transition-colors ${
                          index % 2 === 0 ? "bg-white" : "bg-clinical-muted/5"
                        }`}
                        onClick={() => router.push(`/students/${student.id}`)} // CHANGED: Updated onClick to navigate to student detail page
                      >
                        <td className="py-2.5 px-4">
                          <div className="flex items-center gap-3">
                            <Avatar className="h-8 w-8 border border-clinical-border">
                              <AvatarImage src={student.avatar || "/placeholder.svg"} />
                              <AvatarFallback className="bg-clinical-primary/10 text-clinical-primary text-xs">
                                {student.name
                                  .split(" ")
                                  .map((n) => n[0])
                                  .join("")}
                              </AvatarFallback>
                            </Avatar>
                            <div>
                              <div className="font-medium text-sm text-foreground">{student.name}</div>
                              <div className="text-xs text-muted-foreground">{student.id}</div>
                            </div>
                          </div>
                        </td>
                        <td className="py-2.5 px-4">
                          <div className="text-sm text-foreground">{student.program.split(" ")[0]}</div>
                          <div className="text-xs text-muted-foreground">{student.department}</div>
                        </td>
                        <td className="py-2.5 px-4">
                          <div className="text-sm font-medium text-clinical-primary">{student.currentStage}</div>
                          <div className="text-xs text-muted-foreground">
                            {student.stageProgress.current}/{student.stageProgress.total} nodes
                          </div>
                        </td>
                        <td className="py-2.5 px-4">
                          <div className="flex items-center gap-1.5">
                            {student.stageHistory?.map((stage, idx) => (
                              <div
                                key={idx}
                                className={`h-5 w-7 rounded-sm flex items-center justify-center text-xs font-medium ${
                                  stage.status === "completed"
                                    ? "bg-clinical-success text-white"
                                    : stage.status === "in-progress"
                                      ? "bg-clinical-primary text-white"
                                      : stage.status === "at-risk"
                                        ? "bg-clinical-alert text-white"
                                        : "bg-clinical-muted/30 text-muted-foreground"
                                }`}
                                title={`${stages.find((s) => s.id === stage.stage)?.label[language] || stage.stage}: ${stage.status}`}
                              >
                                {stage.stage.replace("W", "")}
                              </div>
                            ))}
                            <span className="ml-2 text-sm font-semibold text-foreground">
                              {student.overallProgress}%
                            </span>
                          </div>
                        </td>
                        <td className="py-2.5 px-4 text-center">
                          <div className="text-sm text-foreground">{student.cohort}</div>
                        </td>
                        <td className="py-2.5 px-4 text-center">
                          <div
                            className={`text-sm font-medium ${student.overdue ? "text-clinical-alert" : "text-foreground"}`}
                          >
                            {student.dueNext}
                          </div>
                          {student.overdue && <div className="text-xs text-clinical-alert">Overdue</div>}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </TabsContent>

          {/* Kanban View */}
          <TabsContent value="kanban" className="mt-0">
            <div className="flex gap-4 overflow-x-auto pb-4">
              {stages
                .filter((stage) => !stage.conditional || filteredStudents.some((s) => s.rpRequired))
                .map((stage) => (
                  <div key={stage.id} className="flex-shrink-0 w-[320px]">
                    <Card className="border-clinical-border bg-clinical-card shadow-clinical h-full">
                      <CardHeader className="pb-3 bg-clinical-muted/20 border-b border-clinical-border">
                        <CardTitle className="text-base font-medium">{stage.label[language]}</CardTitle>
                        <CardDescription className="text-sm">
                          {filteredStudents.filter((s) => s.currentStage === stage.id).length} students
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="p-3 space-y-3">
                        {filteredStudents
                          .filter((s) => s.currentStage === stage.id)
                          .map((student) => (
                            <Card
                              key={student.id}
                              className="border-clinical-border bg-background hover:shadow-md transition-all cursor-pointer"
                              onClick={() => router.push(`/students/${student.id}`)} // CHANGED: Updated onClick to navigate to student detail page
                            >
                              <CardContent className="p-4">
                                <div className="flex items-start gap-3 mb-3">
                                  <Avatar className="h-9 w-9 border border-clinical-border">
                                    <AvatarImage src={student.avatar || "/placeholder.svg"} />
                                    <AvatarFallback className="bg-clinical-primary/10 text-clinical-primary text-xs">
                                      {student.name
                                        .split(" ")
                                        .map((n) => n[0])
                                        .join("")}
                                    </AvatarFallback>
                                  </Avatar>
                                  <div className="flex-1 min-w-0">
                                    <div className="font-medium text-sm text-foreground truncate">{student.name}</div>
                                    <div className="text-xs text-muted-foreground truncate">{student.program}</div>
                                  </div>
                                </div>

                                <div className="flex flex-wrap gap-1 mb-3">
                                  {student.advisors.slice(0, 2).map((advisor, idx) => (
                                    <Badge
                                      key={idx}
                                      variant="outline"
                                      className="text-xs bg-clinical-muted/20 border-clinical-border"
                                    >
                                      {advisor.split(" ")[1]}
                                    </Badge>
                                  ))}
                                </div>

                                <div className="space-y-2">
                                  <div className="flex items-center justify-between text-xs">
                                    <span className="text-muted-foreground">Progress</span>
                                    <span className="font-medium text-foreground">{student.overallProgress}%</span>
                                  </div>
                                  <Progress value={student.overallProgress} className="h-1.5" />
                                </div>

                                <div className="flex items-center justify-between mt-3 pt-3 border-t border-clinical-border">
                                  <div
                                    className={`flex items-center gap-1 text-xs ${student.overdue ? "text-clinical-alert" : "text-muted-foreground"}`}
                                  >
                                    <Clock className="h-3 w-3" />
                                    {student.dueNext}
                                  </div>
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
                                    {student.status}
                                  </Badge>
                                </div>
                              </CardContent>
                            </Card>
                          ))}
                      </CardContent>
                    </Card>
                  </div>
                ))}
            </div>
          </TabsContent>

          {/* Cohort Analytics */}
          <TabsContent value="analytics" className="mt-0">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <Card className="border-clinical-border bg-clinical-card shadow-clinical">
                <CardHeader>
                  <CardTitle className="text-base font-medium flex items-center gap-2">
                    <CheckCircle2 className="h-5 w-5 text-clinical-success" />
                    Antiplagiarism Compliance
                  </CardTitle>
                  <CardDescription>Students with ≥85% confirmed</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-semibold text-clinical-primary">87%</div>
                    <div className="text-sm text-muted-foreground">
                      ({(filteredStudents.length * 0.87) | 0}/{filteredStudents.length})
                    </div>
                  </div>
                  <Progress value={87} className="mt-4 h-2" />
                </CardContent>
              </Card>

              <Card className="border-clinical-border bg-clinical-card shadow-clinical">
                <CardHeader>
                  <CardTitle className="text-base font-medium flex items-center gap-2">
                    <Clock className="h-5 w-5 text-clinical-info" />
                    Median Days in W2
                  </CardTitle>
                  <CardDescription>Pre-examination stage duration</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-semibold text-clinical-primary">42</div>
                    <div className="text-sm text-muted-foreground">days</div>
                  </div>
                  <div className="mt-4 text-sm text-muted-foreground">
                    <span className="text-clinical-success">↓ 8 days</span> vs last cohort
                  </div>
                </CardContent>
              </Card>

              <Card className="border-clinical-border bg-clinical-card shadow-clinical">
                <CardHeader>
                  <CardTitle className="text-base font-medium flex items-center gap-2">
                    <AlertCircle className="h-5 w-5 text-clinical-warning" />
                    Bottleneck Node
                  </CardTitle>
                  <CardDescription>Most common waiting/needs_fixes</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="text-lg font-medium text-foreground mb-1">D1_normokontrol_ncste</div>
                  <div className="text-sm text-muted-foreground mb-3">Normokontrol submission</div>
                  <div className="flex items-center gap-2">
                    <Badge
                      variant="outline"
                      className="bg-clinical-warning/10 text-clinical-warning border-clinical-warning/30"
                    >
                      12 students
                    </Badge>
                    <span className="text-xs text-muted-foreground">awaiting review</span>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-clinical-border bg-clinical-card shadow-clinical">
                <CardHeader>
                  <CardTitle className="text-base font-medium flex items-center gap-2">
                    <Users className="h-5 w-5 text-clinical-info" />
                    RP Required
                  </CardTitle>
                  <CardDescription>Students with conditional stage</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-semibold text-clinical-primary">
                      {filteredStudents.filter((s) => s.rpRequired).length}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      students (
                      {((filteredStudents.filter((s) => s.rpRequired).length / filteredStudents.length) * 100).toFixed(
                        0,
                      )}
                      %)
                    </div>
                  </div>
                  <div className="mt-4 text-sm text-muted-foreground">Years since graduation {">"} 3</div>
                </CardContent>
              </Card>

              <Card className="border-clinical-border bg-clinical-card shadow-clinical">
                <CardHeader>
                  <CardTitle className="text-base font-medium flex items-center gap-2">
                    <AlertCircle className="h-5 w-5 text-clinical-alert" />
                    Overdue Items
                  </CardTitle>
                  <CardDescription>Students with overdue tasks</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-baseline gap-2">
                    <div className="text-4xl font-semibold text-clinical-alert">
                      {filteredStudents.filter((s) => s.overdue).length}
                    </div>
                    <div className="text-sm text-muted-foreground">students</div>
                  </div>
                  <div className="mt-4 flex gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      className="flex-1 border-clinical-alert/30 text-clinical-alert hover:bg-clinical-alert/10 bg-transparent"
                    >
                      Send Reminders
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>
      </main>

      {/* REMOVED: Student Detail Drawer */}
    </div>
  )
}

// Helper component for node cards
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
    <div className="p-4 rounded-lg border border-clinical-border bg-background hover:shadow-md transition-all">
      <div className="flex items-start justify-between mb-2">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <code className="text-xs text-muted-foreground bg-clinical-muted/20 px-2 py-0.5 rounded">{id}</code>
            <Badge variant="outline" className={`text-xs ${stateConfig.color}`}>
              {stateConfig.label}
            </Badge>
          </div>
          <h4 className="text-sm font-medium text-foreground">{title[language]}</h4>
        </div>
      </div>
      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        {dueDate && (
          <div className="flex items-center gap-1">
            <Calendar className="h-3 w-3" />
            {dueDate}
          </div>
        )}
        {hasFiles && (
          <div className="flex items-center gap-1">
            <CheckCircle2 className="h-3 w-3" />2 files
          </div>
        )}
      </div>
    </div>
  )
}
