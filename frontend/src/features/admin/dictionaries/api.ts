import { api } from "@/api/client";

export type Program = {
  id: string;
  name: string;
  code: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type Specialty = {
  id: string;
  name: string;
  code: string;
  program_ids: string[]; // Multiple programs
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export const listPrograms = async (activeOnly = false): Promise<Program[]> => {
  return api.get(`/admin/dictionaries/programs?active=${activeOnly}`);
};
export type Cohort = {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};
export const createProgram = async (data: { name: string; code?: string }) => {
  return api.post("/admin/dictionaries/programs", data);
};

export const updateProgram = async (id: string, data: { name?: string; code?: string; is_active?: boolean }) => {
  return api.put(`/admin/dictionaries/programs/${id}`, data);
};

export const deleteProgram = async (id: string) => {
  return api.delete(`/admin/dictionaries/programs/${id}`);
};

export const listSpecialties = async (activeOnly = false, programId?: string): Promise<Specialty[]> => {
  let url = `/admin/dictionaries/specialties?active=${activeOnly}`;
  if (programId) url += `&program_id=${programId}`;
  return api.get(url);
};

export const createSpecialty = async (data: { name: string; code?: string; program_ids?: string[] }) => {
  return api.post("/admin/dictionaries/specialties", data);
};

export const updateSpecialty = async (id: string, data: { name?: string; code?: string; program_ids?: string[]; is_active?: boolean }) => {
  return api.put(`/admin/dictionaries/specialties/${id}`, data);
};

export const deleteSpecialty = async (id: string) => {
  return api.delete(`/admin/dictionaries/specialties/${id}`);
};

// --- Cohorts ---

export const listCohorts = async (activeOnly: boolean = true) => {
  return api.get(`/admin/dictionaries/cohorts?active=${activeOnly}`);
};

export const createCohort = async (data: { name: string; start_date?: string; end_date?: string }) => {
  return api.post("/admin/dictionaries/cohorts", data);
};

export const updateCohort = async (id: string, data: { name?: string; start_date?: string; end_date?: string; is_active?: boolean }) => {
  return api.put(`/admin/dictionaries/cohorts/${id}`, data);
};

export const deleteCohort = async (id: string) => {
  return api.delete(`/admin/dictionaries/cohorts/${id}`);
};

// --- Departments ---

export type Department = {
  id: string;
  name: string;
  code: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export const listDepartments = async (activeOnly: boolean = true) => {
  return api.get(`/admin/dictionaries/departments?active=${activeOnly}`);
};

export const createDepartment = async (data: { name: string; code?: string }) => {
  return api.post("/admin/dictionaries/departments", data);
};

export const updateDepartment = async (id: string, data: { name?: string; code?: string; is_active?: boolean }) => {
  return api.put(`/admin/dictionaries/departments/${id}`, data);
};

export const deleteDepartment = async (id: string) => {
  return api.delete(`/admin/dictionaries/departments/${id}`);
};
