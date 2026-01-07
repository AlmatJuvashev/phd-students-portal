export interface ClassSession {
  id: string;
  course_offering_id: string;
  title: string;
  date: string;
  start_time: string;
  end_time: string;
  room_id?: string | null;
  instructor_id?: string | null;
  type: string;
  session_format?: string | null;
  meeting_url?: string | null;
  is_cancelled: boolean;
  created_at: string;
  updated_at: string;
}

export interface TeacherDashboardStats {
  next_class?: ClassSession;
  active_courses: number;
  pending_grading: number;
  today_classes_count: number;
  at_risk_count?: number;
}

export interface CourseOffering {
  id: string;
  course_id: string;
  term_id: string;
  tenant_id: string;
  section: string;
  delivery_format: string;
  max_capacity: number;
  virtual_capacity?: number | null;
  current_enrolled: number;
  meeting_url?: string | null;
  target_cohorts: string[];
  is_active: boolean;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CourseEnrollment {
  id: string;
  course_offering_id: string;
  student_id: string;
  status: string;
  method: string;
  enrolled_at: string;
  updated_at: string;
  student_name?: string;
  student_email?: string;
}

export interface ActivitySubmission {
  id: string;
  activity_id: string;
  student_id: string;
  course_offering_id: string;
  content: any;
  submitted_at: string;
  status: string;
  activity_title?: string;
  student_name?: string;
  student_email?: string;
}

export type RiskLevel = 'low' | 'medium' | 'high' | 'critical';

export interface StudentRiskProfile {
  student_id: string;
  student_name: string;
  overall_progress: number;
  assignments_completed: number;
  assignments_total: number;
  assignments_overdue: number;
  last_activity: string;
  days_inactive: number;
  average_grade: number;
  risk_level: RiskLevel;
  risk_factors: string[];
  suggested_actions: string[];
}

export interface TeacherStudentActivityEvent {
  kind: 'submission' | 'grade';
  occurred_at: string;
  title: string;
  status?: string | null;
  activity_id?: string | null;
  submission_id?: string | null;
  score?: number | null;
  max_score?: number | null;
  grade?: string | null;
}

export type { GradebookEntry, GradeSubmissionRequest } from '@/features/grading/types';
