export type QuestionType = 'multi_select' | 'short_answer' | 'osce' | 'true_false' | 'essay'; // Added standard types
export type QuestionStatus = 'draft' | 'review' | 'approved' | 'published' | 'retired' | 'archived';
export type Language = 'en' | 'ru' | 'kz';
export type DifficultyLevel = 'easy' | 'medium' | 'hard';

export interface BloomsTaxonomy {
  id?: string;
  level: string; // e.g. "Remember", "Understand"
}

export interface Tag {
  id: string;
  label: string;
  category: 'topic' | 'competency' | 'difficulty';
}

export interface Bank {
  id: string;
  tenant_id?: string;
  title: string;
  description?: string;
  subject?: string;
  item_count?: number; // Backend might need to compute this or return in list
  blooms_taxonomy?: BloomsTaxonomy;
  is_public?: boolean;
  created_by?: string;
  created_at?: string;
  updated_at?: string;
}

export interface QuestionOption {
  id?: string;
  text: string;
  is_correct: boolean;
  feedback?: string;
}

export interface Question {
  id: string;
  bank_id: string;
  type: QuestionType;
  stem: string; // The question text
  media_url?: string;
  points_default: number;
  difficulty_level?: DifficultyLevel;
  learning_outcome_id?: string;
  options?: QuestionOption[];
  // Legacy/v11 fields might need mapping or removal if not in backend
  status?: QuestionStatus; 
  tags?: Tag[];
  language?: Language;
  author?: string;
  created_at?: string;
  updated_at?: string;
}
