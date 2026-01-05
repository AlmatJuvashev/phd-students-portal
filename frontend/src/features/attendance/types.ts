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

export interface ClassAttendance {
  id: string;
  class_session_id: string;
  student_id: string;
  status: string;
  notes: string;
  recorded_by_id: string;
  created_at: string;
  updated_at: string;
}

export interface AttendanceUpdate {
  student_id: string;
  status: string;
  notes?: string;
}

