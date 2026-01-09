export interface Program {
  id: string;
  tenant_id: string;
  title: string | Record<string, string>;
  code: string;
  description: string | Record<string, string>;
  type: 'bachelor' | 'master' | 'doctoral' | 'certificate';
  total_credits: number;
  duration_semesters: number;
  status: 'draft' | 'active' | 'archived';
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Course {
  id: string;
  tenant_id: string;
  program_id?: string;
  code: string;
  title: string | Record<string, string>;
  description: string | Record<string, string>;
  credits: number;
  workload_hours?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
