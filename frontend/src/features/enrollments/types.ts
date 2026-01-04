import { Program } from '../curriculum/types';

export interface Student {
  id: string;
  name: string;
  email: string;
  avatar: string; // Initial/Letter
  program?: string;
  department?: string;
  cohort?: string;
}

export interface Enrollment {
  id: string;
  student_id: string;
  program_id: string;
  cohort_id: string;
  status: 'active' | 'paused' | 'completed' | 'dropped';
  progress: number;
  overdue_tasks: number;
  last_activity: string;
  created_at: string;
  
  // Joined fields for UI
  student: Student;
  program: Program;
}

export interface EnrollmentCreateRequest {
  program_id: string;
  cohort_id: string;
  student_ids: string[];
  start_date: string;
}
