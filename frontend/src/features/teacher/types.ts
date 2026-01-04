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

export type { GradebookEntry, GradeSubmissionRequest } from '@/features/grading/types';
