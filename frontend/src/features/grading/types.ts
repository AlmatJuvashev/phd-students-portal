export interface GradingSchema {
  id: string;
  tenant_id: string;
  name: string;
  scale: any;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface GradebookEntry {
  id: string;
  course_offering_id: string;
  activity_id: string;
  student_id: string;
  score: number;
  max_score: number;
  grade: string;
  feedback: string;
  graded_by_id: string;
  graded_at: string;
  created_at: string;
  updated_at: string;
}

export interface GradeSubmissionRequest {
  course_offering_id: string;
  activity_id: string;
  student_id: string;
  score: number;
  max_score: number;
  feedback?: string;
}

