export interface Program {
  id: string;
  tenant_id: string;
  title: string;
  code: string;
  description: string;
  type: 'bachelor' | 'master' | 'doctoral' | 'certificate';
  total_credits: number;
  duration_semesters: number;
  status: 'draft' | 'active' | 'archived';
  created_at: string;
  updated_at: string;
}

export interface Course {
  id: string;
  tenant_id: string;
  program_id?: string;
  code: string;
  title: string;
  description: string;
  credits: number;
  category: 'core' | 'elective' | 'research';
  prerequisites?: string[];
  status: 'draft' | 'active' | 'archived';
  created_at: string;
  updated_at: string;
}
