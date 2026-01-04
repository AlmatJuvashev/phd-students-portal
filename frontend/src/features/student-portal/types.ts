export interface ProgramProgress {
  title: string;
  progress_percent: number;
  completed_nodes: number;
  total_nodes: number;
  overdue_count: number;
}

export interface StudentDeadline {
  id: string;
  title: string;
  due_at?: string | null;
  source: string;
  status?: string;
  severity?: string;
  link?: string | null;
}

export interface StudentAnnouncement {
  id: string;
  title: string;
  body: string;
  created_at: string;
  link?: string | null;
}

export interface StudentGradeEntry {
  id: string;
  course_offering_id: string;
  course_id?: string;
  course_code?: string | null;
  course_title?: string | null;
  activity_id: string;
  student_id: string;
  score: number;
  max_score: number;
  grade: string;
  feedback: string;
  graded_by_id: string;
  graded_at: string;
}

export interface StudentCourseNextSession {
  id: string;
  date: string;
  start_time: string;
  end_time: string;
  room_id?: string | null;
  meeting_url?: string | null;
  type: string;
}

export interface StudentCourse {
  enrollment_id: string;
  course_offering_id: string;
  status: string;
  course_id: string;
  code: string;
  title: string;
  section: string;
  term_id: string;
  delivery_format: string;
  instructor_name?: string | null;
  progress_percent: number;
  next_session?: StudentCourseNextSession | null;
}

export interface StudentAssignment {
  id: string;
  title: string;
  source: string;
  status: string;
  due_at?: string | null;
  link?: string | null;
  severity?: string;
}

export interface StudentDashboard {
  program: ProgramProgress;
  upcoming_deadlines: StudentDeadline[];
  recent_grades: StudentGradeEntry[];
  announcements: StudentAnnouncement[];
}

