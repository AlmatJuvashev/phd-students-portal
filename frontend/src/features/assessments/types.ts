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

export type QuestionType = 'MCQ' | 'TRUE_FALSE' | 'TEXT' | 'MRQ' | 'LIKERT';

export interface QuestionOption {
  id: string;
  text: string;
  is_correct: boolean;
  feedback?: string;
}

export interface Question {
  id: string;
  type: QuestionType;
  stem: string;
  options?: QuestionOption[];
  sort_order: number;
}

export interface Attempt {
  id: string;
  user_id: string;
  assessment_id: string;
  status: 'IN_PROGRESS' | 'SUBMITTED' | 'GRADED';
  started_at: string;
  completed_at?: string;
  score: number;
  grade?: string;
}

export interface AttemptDetailsResponse {
  assessment: Assessment;
  questions: Question[];
  attempt: Attempt;
  responses: Array<{
    id: string;
    question_id: string;
    selected_option_id?: string;
    text_response?: string;
    is_correct?: boolean;
    score?: number;
  }>;
}
