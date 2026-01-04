export type QuestionType = 'MCQ' | 'MRQ' | 'TRUE_FALSE' | 'TEXT' | 'LIKERT';
export type DifficultyLevel = 'EASY' | 'MEDIUM' | 'HARD';

export interface QuestionBank {
  id: string;
  tenant_id: string;
  title: string;
  description?: string | null;
  subject?: string | null;
  blooms_taxonomy?: string | null;
  is_public: boolean;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface QuestionOption {
  id?: string;
  question_id?: string;
  text: string;
  is_correct: boolean;
  sort_order?: number;
  feedback?: string | null;
}

export interface Question {
  id: string;
  bank_id: string;
  type: QuestionType;
  stem: string;
  media_url?: string | null;
  points_default: number;
  difficulty_level?: DifficultyLevel | null;
  learning_outcome_id?: string | null;
  created_at: string;
  updated_at: string;
  options?: QuestionOption[];
}

export interface CreateBankRequest {
  title: string;
  description?: string;
  subject?: string;
  blooms_taxonomy?: string;
  is_public?: boolean;
}

export interface CreateQuestionRequest {
  type: QuestionType;
  stem: string;
  points_default?: number;
  difficulty_level?: DifficultyLevel;
  options?: QuestionOption[];
}

