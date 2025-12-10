import { api } from "./client";

export interface Program {
  id: string;
  name: string;
  code: string;
  is_active: boolean;
}

export interface Specialty {
  id: string;
  name: string;
  code: string;
  is_active: boolean;
  program_ids: string[];
}

export interface Department {
  id: string;
  name: string;
  code: string;
  is_active: boolean;
}

export interface Cohort {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  is_active: boolean;
}

export async function getPrograms() {
  return api<Program[]>("/dictionaries/programs?active=true");
}

export async function getSpecialties() {
  return api<Specialty[]>("/dictionaries/specialties?active=true");
}

export async function getDepartments() {
  return api<Department[]>("/dictionaries/departments?active=true");
}

export async function getCohorts() {
  return api<Cohort[]>("/dictionaries/cohorts?active=true");
}
