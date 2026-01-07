export interface Assessment {
  id: string;
  tenant_id: string;
  course_offering_id: string;
  title: string;
  description?: string;
  time_limit_minutes?: number;
  available_from?: string;
  available_until?: string;
  shuffle_questions: boolean;
  grading_policy: 'AUTOMATIC' | 'MANUAL_REVIEW';
  passing_score: number;
  created_by: string;
  created_at: string;
  updated_at: string;
  
  // Joins (optional depending on API response)
  questions_count?: number; 
  course_offering?: {
    id: string;
    section: string;
    course?: {
      title: string;
      code: string;
    }
  }
}

export interface AssessmentListFilters {
  course_id?: string;
  status?: string;
  search?: string;
}
