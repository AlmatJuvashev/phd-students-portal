
export type UserRole = 'admin' | 'teacher' | 'student';

export interface SectionScore {
  id: string;
  score: number;
  max: number;
  feedback: string;
}

export interface AnswerRecord {
  user_id: number;
  name_rus: string; // Used as category/subject if topic missing
  topic: string;
  question: string;
  answer: string;
  comment: string;
  score: number; // Teacher score
  ai_score?: number; // Model score for audit comparison
  examiner_id: number;
  examiner_name: string;
  attempt_id: string;
  created_at: string;
  // New fields for detailed feedback
  section_scores?: SectionScore[];
  strengths?: string[];
  weaknesses?: string[];
  suggestions?: string[];
}

export interface PendingAnswer {
  id: string;
  student_id: number; // Anonymized in UI
  category: string;
  question: string;
  answer: string;
  submitted_at: string;
}

export interface TeacherStat {
  examiner_id: number;
  examiner_name: string;
  answers_count: number;
  average_score: number;
  std_dev: number;
  deviation_from_global: number;
  status: 'lenient' | 'strict' | 'neutral';
}

export interface Protocol {
  id: string;
  name: string;
  type: 'PDF' | 'Text';
  language: 'RU' | 'KZ';
  status: 'indexed' | 'not_indexed';
  last_updated: string;
  active: boolean;
}

export interface RubricItem {
  id: string;
  criteria: string;
  description: string;
  max_score: number;
}

export interface Question {
  id: string;
  text: string;
  difficulty: 'Beginner' | 'Intermediate' | 'Advanced';
  estimated_time_mins: number;
  rubric?: RubricItem[];
}

export interface Topic {
  id: string;
  title: string;
  description: string;
  question_count: number;
  objectives: string[];
  questions: Question[];
}

export interface ModelConfig {
  id: string;
  name: string;
  tags: string[];
}

export interface Feedback {
  score: number;
  strengths: string[];
  weaknesses: string[];
  suggestions: string[];
  rag_snippets: { title: string; text: string }[];
}
