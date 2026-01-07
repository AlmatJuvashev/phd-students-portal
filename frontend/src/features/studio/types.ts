export type ActivityType = 'text' | 'video' | 'quiz' | 'survey' | 'assignment' | 'resource' | 'live';
export type QuestionType = 'multiple_choice' | 'multi_select' | 'short_text' | 'ordering' | 'section_header' | 'page_break' | 'matrix';
export type FormFieldType = 'text' | 'number' | 'select' | 'date' | 'boolean' | 'upload' | 'collection' | 'note';

export interface FormField {
  id: string;
  key: string;
  type: FormFieldType;
  label: string;
  required?: boolean;
  placeholder?: string;
  help_text?: string;
  options?: { value: string; label: string }[]; 
  dictionaryKey?: string; 
  item_fields?: FormField[]; // For collections
}

export interface QuizQuestion {
  id: string;
  type: QuestionType;
  text: string;
  subtitle?: string;
  hint?: string;
  points: number;
  feedback_correct?: string;
  feedback_incorrect?: string;
  options?: { id: string; text: string; is_correct: boolean }[];
  correct_order?: string[];
  display_logic?: {
    depends_on_question_id: string;
    condition: 'equals' | 'not_equals' | 'contains';
    value: string;
  };
  matrixRows?: string[];
  matrixCols?: string[];
}

export interface Attachment {
  id: string;
  name: string;
  type: 'pdf' | 'word' | 'file' | 'image';
  url: string;
}

export interface Citation {
  id: string;
  text: string;
  url?: string;
}

export interface ChecklistItem {
  id: string;
  text: string;
  required: boolean;
  helpText?: string;
  is_completed?: boolean; // Runtime state
}

export interface ChecklistConfig {
  intro?: string;
  reviewer_role?: 'advisor' | 'secretary' | 'admin' | 'none';
  items: ChecklistItem[];
  templates?: { name: string; url: string; size?: string }[];
}

export interface AssignmentConfig {
  submission_types: string[]; // 'file_upload', 'text_entry', 'form', 'checklist'
  group_assignment: boolean;
  peer_review: boolean;
  rubric?: any[]; // Full rubric structure
  rubric_id?: string;
  points?: number;
  instruction_files?: { name: string; size: string; type: string }[];
  form_fields?: FormField[]; // Custom form structure for 'form' type
  checklist_config?: ChecklistConfig; // For 'checklist' type
  allowed_extensions?: string;
  peer_review_count?: number;
}

export interface Activity {
  id: string;
  title: string;
  type: ActivityType;
  points: number;
  is_optional: boolean;
  content: string;
  
  // Type-specific configs
  video_url?: string;
  video_urls?: string[];
  video_description?: string;
  
  quiz_config?: {
    time_limit_minutes: number;
    passing_score: number;
    shuffle_questions: boolean;
    show_results: boolean;
    questions: QuizQuestion[];
  };
  
  assignment_config?: AssignmentConfig;
  
  survey_config?: {
    title: string;
    anonymous: boolean;
    show_progress_bar: boolean;
    questions: any[]; // Using any for now to avoid circular dependency hell, or cleaner: SurveyQuestion[]
  };

  confirm_config?: {
    intro: string;
    button_text: string;
    reviewer: string;
    templates: { name: string; size: string }[];
    uploads: { id: string; label: string; required: boolean }[];
  };
  
  resource_config?: { url: string; file_name?: string };
  live_config?: { platform: 'zoom' | 'meet'; date: string; link: string };

  attachments: Attachment[];
  citations: Citation[];

  // Journey Map Visual Props
  position?: { x: number; y: number };
  worldId?: string; // Links to Module ID
}

export interface Lesson {
  id: string;
  title: string;
  order: number;
  activities: Activity[];
}

export interface Module {
  id: string;
  title: string;
  description?: string; // Added description
  order: number;
  isOpen?: boolean; // UI state
  color?: string; // For map visualization
  position?: { x: number; y: number }; // For map visualization
  condition?: {
    nodeId: string;
    fieldKey: string;
    operator: 'equals' | 'not_equals' | 'contains';
    value: string;
  } | null;
  lessons: Lesson[];
}

export interface FlowEdge {
  id: string;
  from: string;
  to: string;
  type: 'solid' | 'dashed';
}

export interface CourseContent {
  id: string;
  course_id: string;
  modules: Module[];
  edges?: FlowEdge[];
}

// --- Program Versions (Program Builder) ---

export type ProgramNodeType =
  | 'course'
  | 'form'
  | 'checklist'
  | 'milestone'
  | 'payment'
  | 'approval'
  | 'meeting'
  | 'survey'
  | 'sync_ops'
  | 'info'
  | 'confirmTask';

export interface ProgramVersionNode {
  id: string;
  program_version_id?: string;
  parent_node_id?: string;
  slug: string;
  type: ProgramNodeType;
  title: string;
  description?: string;
  module_key: string; // "I", "II" (Phase Identifier)
  coordinates: { x: number; y: number };
  config: any; // Dynamic config (e.g. { course_id: "..." })
  prerequisites?: string[];
  points?: number;
}

export interface ProgramPhase { // Frontend abstraction for grouping nodes by module_key
  id: string; // The module_key (e.g. "I")
  title: string;
  description?: string;
  color: string;
  order: number;
  position: { x: number; y: number };
}

export interface ProgramVersion {
  id: string;
  program_id: string;
  title: string;
  version: string;
  nodes: ProgramVersionNode[];
  edges: FlowEdge[]; // Edges are derived from prerequisites in DB, but managed as edges in UI
  phases: ProgramPhase[]; // Metadata often stored in config or inferred
}
