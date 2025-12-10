"use client"
import StudentDetailClient from "./student-detail-client"
import { AlertTriangle } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import Link from "next/link"

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

// Mock data - in production this would come from an API
const mockStudents: Record<string, Student> = {
  S001: {
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
  S002: {
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
  S003: {
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
  S004: {
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
}

export default async function StudentDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  const student = mockStudents[id]

  if (!student) {
    return (
      <div className="min-h-screen bg-clinical-background flex items-center justify-center">
        <Card className="border-clinical-border bg-clinical-card shadow-clinical p-8">
          <CardContent className="text-center space-y-4">
            <AlertTriangle className="h-12 w-12 text-clinical-alert mx-auto" />
            <h2 className="text-xl font-semibold text-foreground">Student Not Found</h2>
            <p className="text-muted-foreground">Student with ID {id} does not exist.</p>
            <Button asChild className="mt-4">
              <Link href="/">Back to Dashboard</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return <StudentDetailClient student={student} />
}
