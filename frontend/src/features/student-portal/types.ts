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

export interface StudentCourseDetail {
  course: StudentCourse;
  sessions: StudentCourseNextSession[];
}

export interface CourseActivity {
  id: string;
  lesson_id: string;
  type: string;
  title: string;
  order: number;
  points: number;
  is_optional: boolean;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface CourseLesson {
  id: string;
  module_id: string;
  title: string;
  order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  activities: CourseActivity[];
}

export interface CourseModule {
  id: string;
  course_id: string;
  title: string;
  order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  lessons: CourseLesson[];
}

export interface ActivitySubmission {
  id: string;
  activity_id: string;
  student_id: string;
  course_offering_id: string;
  content: any;
  submitted_at: string;
  status: string;
}

export interface StudentAssignmentDetail {
  activity: CourseActivity;
  submission?: ActivitySubmission | null;
  course_offering_id: string;
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
