export type QuestionType = 'MCQ' | 'MRQ' | 'TRUE_FALSE' | 'TEXT' | 'LIKERT';
export type AttemptStatus = 'IN_PROGRESS' | 'SUBMITTED' | 'GRADED';
export type GradingPolicy = 'AUTOMATIC' | 'MANUAL_REVIEW';

export interface Assessment {
  id: string;
  tenant_id?: string;
  course_offering_id: string;
  title: string;
  description?: string | null;
  time_limit_minutes?: number | null;
  available_from?: string | null;
  available_until?: string | null;
  shuffle_questions: boolean;
  grading_policy: GradingPolicy;
  security_settings?: unknown;
  passing_score: number;
  created_by?: string;
  created_at?: string;
  updated_at?: string;
}

export interface QuestionOption {
  id: string;
  question_id?: string;
  text: string;
  is_correct: boolean;
  sort_order?: number;
  feedback?: string | null;
}

export interface Question {
  id: string;
  bank_id?: string;
  type: QuestionType;
  stem: string;
  media_url?: string | null;
  points_default: number;
  options?: QuestionOption[];
}

export interface AssessmentAttempt {
  id: string;
  assessment_id: string;
  student_id: string;
  started_at: string;
  finished_at?: string | null;
  score: number;
  status: AttemptStatus;
}

export interface ItemResponse {
  id?: string;
  attempt_id: string;
  question_id: string;
  selected_option_id?: string | null;
  text_response?: string | null;
  score: number;
  is_correct: boolean;
  graded_at?: string | null;
}

export interface AssessmentForTakingResponse {
  assessment: Assessment;
  questions: Question[];
}

export interface AttemptDetailsResponse {
  attempt: AssessmentAttempt;
  assessment: Assessment;
  questions: Question[];
  responses: ItemResponse[];
}

