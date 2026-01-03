export interface Department {
  id: string;
  name: string;
  color?: string;
  preferred_building_id?: string;
}

export type RoomType = 'lecture_hall' | 'classroom' | 'lab' | 'meeting_room' | 'auditorium' | 'seminar_room';

export interface RoomAttribute {
  room_id: string;
  key: string;
  value: string;
}

export interface Room {
  id: string;
  building_id: string;
  building_name?: string; // May need to join this or fetch separately? API response usually has ID.
  name: string;
  capacity: number;
  floor: number;
  department_id?: string;
  type: RoomType;
  attributes?: RoomAttribute[];
  is_active: boolean;
}

export interface Building {
  id: string;
  name: string;
  address?: string;
  is_active: boolean;
}

export interface AcademicTerm {
  id: string;
  name: string; // "Fall 2025"
  code: string; // "2025-FA"
  start_date: string;
  end_date: string;
  is_active: boolean;
}

export interface CourseOffering {
  id: string;
  course_id: string; 
  term_id: string;
  section: string;
  delivery_format: string;
  current_enrolled: number;
  max_capacity: number;
  // Included via join often
  course_title?: string;
  instructor_names?: string[]; 
}

export interface ClassSession {
  id: string;
  course_offering_id: string;
  title: string;
  date: string; // ISO Date
  start_time: string; // "14:00"
  end_time: string;   // "15:30"
  room_id?: string;
  room_name?: string; // Ideally joined
  instructor_id?: string;
  instructor_name?: string; // Ideally joined
  type: string; // LECTURE, LAB, etc.
  is_cancelled: boolean;
  warnings?: string[]; // Frontend specific for solver
  
  // Frontend helpers
  resourceId?: string; // Alias for room_id for UI logic
  durationMinutes?: number;
}

export interface SolverConfig {
  max_iterations: number;
  utilization_weight: number;
  satisfaction_weight: number;
  allow_overtime: boolean;
  prioritize_buildings: boolean;
  enable_department_constraints: boolean;
}
