
import { addDays, subDays } from 'date-fns';

export type ProgramType = 'university' | 'school' | 'prep_center';
export type ProgramStatus = 'draft' | 'active' | 'archived';
export type EnrollmentStatus = 'active' | 'paused' | 'completed' | 'dropped';

// --- INFRASTRUCTURE ---
export interface Building {
  id: string;
  name: string;
  code: string;
  location_lat?: number;
  location_lng?: number;
}

export interface Department {
  id: string;
  name: string;
  preferredBuildingId: string; // Proximity constraint
  color?: string; // UI Color
}

// 1. Global Course Definition
export interface GlobalCourse {
  id: string;
  code: string;
  title: string;
  credits: number;
  category: 'core' | 'elective' | 'research';
  defaultDurationDays: number;
  departmentId?: string; // Links to preferred building
  preferredRoomType?: RoomType;
  weeklyHours: number; // Derived from credits usually
  sequenceGroup?: string; // e.g., "CS_CORE_SEQUENCE"
  sequenceOrder?: number; // 1, 2, 3...
  owner?: string;
  lastUpdated?: Date;
}

// 2. Active Course Instance
export interface ActiveCourse {
  id: string;
  globalCourseId: string;
  globalCourseTitle: string; 
  sectionName: string; 
  instructor: string;
  instructorId?: string; // For conflict checking
  term: 'Fall 2024' | 'Spring 2025' | 'Summer 2025' | 'Fall 2025'; // NEW: Term context
  startDate: Date;
  endDate: Date;
  status: 'upcoming' | 'ongoing' | 'completed';
  enrolledCount: number;
  capacity: number;
}

export interface ProgramCourse {
  id: string;
  courseId: string;
  title: string;
  isRequired: boolean;
  passingGrade: number;
  deadlineDays: number;
  order: number;
}

export interface Program {
  id: string;
  title: string;
  type: ProgramType;
  status: ProgramStatus;
  cohortsCount: number;
  studentsCount: number;
  completionRate: number;
  lastUpdated: Date;
  description: string;
  courses: ProgramCourse[];
  owner?: string;
}

export interface Cohort {
  id: string;
  programId: string; 
  programName?: string; 
  name: string;
  startDate: Date;
  endDate: Date;
  location: string;
  studentsCount: number;
  status: 'planning' | 'active' | 'graduated';
}

export interface ScheduleEvent {
  id: string;
  targetId: string; // Cohort ID or Active Course ID
  targetType: 'cohort' | 'course';
  targetName: string;
  instructorId?: string; // Denormalized for fast conflict check
  instructorName?: string;
  title: string;
  type: 'lecture' | 'exam' | 'deadline' | 'holiday' | 'lab';
  date: Date;
  startTime?: string; // HH:mm
  durationMinutes?: number;
  location?: string;
  resourceId?: string; // Links to Room ID
  satisfactionScore?: number; // 0-100
  isLocked?: boolean; 
  isHistorical?: boolean; // NEW: Flag for copied past schedules
  warnings?: string[]; // Conflicts
}

export interface Student {
  id: string;
  name: string;
  email: string;
  avatar: string;
  status: 'active' | 'inactive';
}

export interface Enrollment {
  id: string;
  studentId: string;
  programId: string;
  cohortId: string;
  startDate: Date;
  status: EnrollmentStatus;
  progress: number;
  overdueTasks: number;
  lastActivity: Date;
  student: Student;
  program: Program;
}

// --- NEW: ROOM MANAGEMENT & MAINTENANCE ---
export type RoomType = 'lecture_hall' | 'classroom' | 'lab' | 'meeting_room' | 'auditorium';

export interface MaintenancePeriod {
  id: string;
  startDate: Date;
  endDate: Date;
  reason: string; // e.g., "Reconstruction", "Painting"
}

export interface Room {
  id: string;
  name: string;
  buildingId: string; // For grouping
  buildingName: string; // Denormalized
  floor: number; // NEW: Floor level
  departmentId?: string; // NEW: Ownership
  capacity: number;
  type: RoomType;
  resources: string[]; 
  status: 'active' | 'maintenance' | 'closed';
  maintenanceLog?: MaintenancePeriod[]; // NEW: Track unavailability
}

// --- MOCK DATA GENERATION ---

let _buildings: Building[] = [
  { id: 'bld_1', name: 'Main Academic Block', code: 'MAB' },
  { id: 'bld_2', name: 'Science Center', code: 'SCI' },
  { id: 'bld_3', name: 'Medical School Wing', code: 'MED' },
];

let _departments: Department[] = [
  { id: 'dept_cs', name: 'Computer Science', preferredBuildingId: 'bld_2', color: '#6366f1' },
  { id: 'dept_med', name: 'General Medicine', preferredBuildingId: 'bld_3', color: '#ef4444' },
  { id: 'dept_hum', name: 'Humanities', preferredBuildingId: 'bld_1', color: '#f59e0b' },
];

let _rooms: Room[] = [
  { id: 'r1', name: 'Lecture Hall A', buildingId: 'bld_1', buildingName: 'Main Academic Block', floor: 1, departmentId: 'dept_hum', capacity: 150, type: 'lecture_hall', resources: ['Projector', 'Sound System', 'Recording'], status: 'active' },
  { 
    id: 'r2', 
    name: 'Computer Lab 302', 
    buildingId: 'bld_2', 
    buildingName: 'Science Center', 
    floor: 3,
    departmentId: 'dept_cs',
    capacity: 30, 
    type: 'lab', 
    resources: ['30 PCs', 'Whiteboard', 'Smart Screen'], 
    status: 'active',
    maintenanceLog: [
      { id: 'm1', startDate: addDays(new Date(), 2), endDate: addDays(new Date(), 5), reason: 'Hardware Upgrade' } // Unavailable in near future
    ]
  },
  { id: 'r3', name: 'Seminar Room 105', buildingId: 'bld_1', buildingName: 'Main Academic Block', floor: 1, departmentId: 'dept_hum', capacity: 25, type: 'classroom', resources: ['Whiteboard', 'Projector'], status: 'active' },
  { 
    id: 'r4', 
    name: 'Grand Auditorium', 
    buildingId: 'bld_3', 
    buildingName: 'Medical School Wing',
    floor: 2,
    departmentId: 'dept_med',
    capacity: 500, 
    type: 'auditorium', 
    resources: ['Stage', 'Lighting', 'Sound System'], 
    status: 'active',
    maintenanceLog: [
      { id: 'm2', startDate: subDays(new Date(), 30), endDate: subDays(new Date(), 10), reason: 'Roof Repair' } // Past maintenance
    ] 
  },
  { id: 'r5', name: 'Bio Lab 1', buildingId: 'bld_2', buildingName: 'Science Center', floor: 2, departmentId: 'dept_med', capacity: 20, type: 'lab', resources: ['Microscopes', 'Sink'], status: 'active' }
];

let _globalCourses: GlobalCourse[] = [
  { id: 'gc1', code: 'RES-101', title: 'Research Methodology', credits: 5, category: 'core', defaultDurationDays: 30, owner: 'Dr. Smith', lastUpdated: subDays(new Date(), 1), departmentId: 'dept_cs', weeklyHours: 3, preferredRoomType: 'classroom', sequenceGroup: 'RES_SEQ', sequenceOrder: 1 },
  { id: 'gc2', code: 'WRT-202', title: 'Academic Writing', credits: 3, category: 'core', defaultDurationDays: 45, owner: 'Jane Smith', lastUpdated: subDays(new Date(), 3), departmentId: 'dept_hum', weeklyHours: 2, preferredRoomType: 'classroom' },
  { id: 'gc3', code: 'STAT-300', title: 'Advanced Statistics', credits: 5, category: 'research', defaultDurationDays: 60, owner: 'Dr. Smith', lastUpdated: subDays(new Date(), 7), departmentId: 'dept_cs', weeklyHours: 3, preferredRoomType: 'lab', sequenceGroup: 'RES_SEQ', sequenceOrder: 2 },
  { id: 'gc4', code: 'ETH-100', title: 'Research Ethics', credits: 2, category: 'core', defaultDurationDays: 14, owner: 'Dept. Head', lastUpdated: subDays(new Date(), 30), departmentId: 'dept_med', weeklyHours: 1, preferredRoomType: 'lecture_hall' },
  { id: 'gc5', code: 'AI-500', title: 'Intro to AI Systems', credits: 4, category: 'elective', defaultDurationDays: 45, owner: 'Alex Rivera', lastUpdated: new Date(), departmentId: 'dept_cs', weeklyHours: 4, preferredRoomType: 'lab', sequenceGroup: 'AI_SEQ', sequenceOrder: 1 },
];

let _activeCourses: ActiveCourse[] = [
  { 
    id: 'ac1', 
    globalCourseId: 'gc1', 
    globalCourseTitle: 'Research Methodology',
    sectionName: 'Fall 2025 - Group A', 
    instructor: 'Prof. Alimov', 
    instructorId: 'inst_1',
    term: 'Fall 2025',
    startDate: subDays(new Date(), 10), 
    endDate: addDays(new Date(), 20), 
    status: 'ongoing', 
    enrolledCount: 24, 
    capacity: 30 
  },
  { 
    id: 'ac2', 
    globalCourseId: 'gc5', 
    globalCourseTitle: 'Intro to AI Systems',
    sectionName: 'Intensive Weekend', 
    instructor: 'Dr. Watson', 
    instructorId: 'inst_2',
    term: 'Fall 2025',
    startDate: addDays(new Date(), 5), 
    endDate: addDays(new Date(), 10), 
    status: 'upcoming', 
    enrolledCount: 12, 
    capacity: 20 
  },
  {
    id: 'ac3',
    globalCourseId: 'gc3',
    globalCourseTitle: 'Advanced Statistics',
    sectionName: 'Spring 2026 - Grp 1',
    instructor: 'Dr. Smith',
    instructorId: 'inst_1',
    term: 'Spring 2025', 
    startDate: addDays(new Date(), 120),
    endDate: addDays(new Date(), 160),
    status: 'upcoming',
    enrolledCount: 0,
    capacity: 25
  },
  {
    id: 'ac4',
    globalCourseId: 'gc2',
    globalCourseTitle: 'Academic Writing',
    sectionName: 'Summer Elective',
    instructor: 'Jane Smith',
    instructorId: 'inst_4',
    term: 'Summer 2025',
    startDate: addDays(new Date(), 60),
    endDate: addDays(new Date(), 90),
    status: 'upcoming',
    enrolledCount: 5,
    capacity: 20
  }
];

let _programs: Program[] = [
  {
    id: 'prog_1',
    title: 'PhD in Computer Science',
    type: 'university',
    status: 'active',
    cohortsCount: 3,
    studentsCount: 142,
    completionRate: 68,
    lastUpdated: new Date(),
    description: 'Comprehensive doctoral pathway focusing on AI and Data Systems.',
    owner: 'Alex Rivera',
    courses: [
        { id: 'pc1', courseId: 'gc1', title: 'Research Methodology', isRequired: true, passingGrade: 80, deadlineDays: 30, order: 0 },
        { id: 'pc2', courseId: 'gc2', title: 'Academic Writing', isRequired: true, passingGrade: 75, deadlineDays: 60, order: 1 }
    ]
  },
  {
    id: 'prog_2',
    title: 'Pre-Med Entrance Prep',
    type: 'prep_center',
    status: 'active',
    cohortsCount: 5,
    studentsCount: 320,
    completionRate: 45,
    lastUpdated: subDays(new Date(), 2),
    description: 'Intensive preparation for medical school entrance exams.',
    owner: 'Dr. House',
    courses: []
  },
  {
    id: 'prog_3',
    title: 'Data Science Masterclass',
    type: 'school',
    status: 'draft',
    cohortsCount: 0,
    studentsCount: 0,
    completionRate: 0,
    lastUpdated: subDays(new Date(), 10),
    description: 'Draft curriculum for 2026.',
    owner: 'Alex Rivera',
    courses: []
  }
];

let _cohorts: Cohort[] = [
  { id: 'coh_1', programId: 'prog_1', programName: 'PhD in CS', name: 'Fall 2024 Intake', startDate: subDays(new Date(), 45), endDate: addDays(new Date(), 120), location: 'Main Campus', studentsCount: 45, status: 'active' },
  { id: 'coh_2', programId: 'prog_1', programName: 'PhD in CS', name: 'Spring 2025 Intake', startDate: addDays(new Date(), 15), endDate: addDays(new Date(), 180), location: 'Remote', studentsCount: 97, status: 'planning' },
  { id: 'coh_3', programId: 'prog_2', programName: 'Pre-Med Prep', name: 'Summer Camp', startDate: addDays(new Date(), 60), endDate: addDays(new Date(), 90), location: 'Lab B', studentsCount: 0, status: 'planning' },
];

let _scheduleEvents: ScheduleEvent[] = [
  { id: 'evt_1', targetId: 'coh_1', targetType: 'cohort', targetName: 'Fall 2024 Intake', title: 'Thesis Proposal Deadline', type: 'deadline', date: addDays(new Date(), 5), startTime: '09:00', durationMinutes: 60, resourceId: undefined },
  { id: 'evt_2', targetId: 'ac1', targetType: 'course', targetName: 'Res. Method - Grp A', instructorId: 'inst_1', instructorName: 'Prof. Alimov', title: 'Guest Lecture: Dr. Chen', type: 'lecture', date: new Date(), startTime: '10:00', durationMinutes: 90, location: 'Lecture Hall A', resourceId: 'r1', satisfactionScore: 95 },
  { id: 'evt_3', targetId: 'coh_2', targetType: 'cohort', targetName: 'Spring 2025 Intake', instructorId: 'inst_3', instructorName: 'Dept. Head', title: 'Orientation Day', type: 'lecture', date: addDays(new Date(), 1), startTime: '14:00', durationMinutes: 120, location: 'Auditorium', resourceId: 'r4', satisfactionScore: 80 },
];

export const MOCK_ALL_STUDENTS: Student[] = [
  { id: 's1', name: 'Alikhan Ivanov', email: 'alikhan@example.com', avatar: 'AI', status: 'active' },
  { id: 's2', name: 'Sarah Connor', email: 'sarah@example.com', avatar: 'SC', status: 'active' },
  { id: 's3', name: 'John Doe', email: 'john@example.com', avatar: 'JD', status: 'active' },
  { id: 's4', name: 'Elena Gilbert', email: 'elena@example.com', avatar: 'EG', status: 'active' },
  { id: 's5', name: 'Michael Ross', email: 'mike@example.com', avatar: 'MR', status: 'inactive' },
  { id: 's6', name: 'Rachel Zane', email: 'rachel@example.com', avatar: 'RZ', status: 'active' },
];

export const MOCK_ENROLLMENTS: Enrollment[] = [
  { 
    id: 'enr_1', 
    studentId: 's1', 
    programId: 'prog_1', 
    cohortId: 'coh_1', 
    startDate: subDays(new Date(), 45), 
    status: 'active', 
    progress: 35, 
    overdueTasks: 1, 
    lastActivity: subDays(new Date(), 1),
    student: MOCK_ALL_STUDENTS[0],
    program: _programs[0]
  },
  { 
    id: 'enr_2', 
    studentId: 's2', 
    programId: 'prog_1', 
    cohortId: 'coh_1', 
    startDate: subDays(new Date(), 45), 
    status: 'paused', 
    progress: 12, 
    overdueTasks: 0, 
    lastActivity: subDays(new Date(), 14),
    student: MOCK_ALL_STUDENTS[1],
    program: _programs[0]
  },
  { 
    id: 'enr_3', 
    studentId: 's3', 
    programId: 'prog_2', 
    cohortId: 'coh_x', 
    startDate: subDays(new Date(), 10), 
    status: 'active', 
    progress: 88, 
    overdueTasks: 0, 
    lastActivity: new Date(),
    student: MOCK_ALL_STUDENTS[2],
    program: _programs[1]
  }
];

// --- API SIMULATION HELPERS ---

export const getPrograms = () => [..._programs];
export const getGlobalCourses = () => [..._globalCourses];
export const getActiveCourses = () => [..._activeCourses];
export const getCohorts = () => [..._cohorts];
export const getScheduleEvents = () => [..._scheduleEvents];
export const getRooms = () => [..._rooms];
export const getBuildings = () => [..._buildings];
export const getDepartments = () => [..._departments];

// Programs
export const addProgram = (program: Program) => {
    _programs = [program, ..._programs];
    return program;
};
export const updateProgram = (id: string, updates: Partial<Program>) => {
    _programs = _programs.map(p => p.id === id ? { ...p, ...updates } : p);
};

// Global Courses
export const addGlobalCourse = (course: GlobalCourse) => {
    _globalCourses = [course, ..._globalCourses];
    return course;
};
export const updateGlobalCourse = (id: string, updates: Partial<GlobalCourse>) => {
    _globalCourses = _globalCourses.map(c => c.id === id ? { ...c, ...updates } : c);
};
export const deleteGlobalCourse = (id: string) => {
    _globalCourses = _globalCourses.filter(c => c.id !== id);
}

// Active Courses (CRUD)
export const addActiveCourse = (course: ActiveCourse) => {
    _activeCourses = [course, ..._activeCourses];
};
export const updateActiveCourse = (id: string, updates: Partial<ActiveCourse>) => {
    _activeCourses = _activeCourses.map(c => c.id === id ? { ...c, ...updates } : c);
};
export const deleteActiveCourse = (id: string) => {
    _activeCourses = _activeCourses.filter(c => c.id !== id);
};

// Cohorts (CRUD)
export const addCohort = (cohort: Cohort) => {
    _cohorts = [cohort, ..._cohorts];
};
export const updateCohort = (id: string, updates: Partial<Cohort>) => {
    _cohorts = _cohorts.map(c => c.id === id ? { ...c, ...updates } : c);
};
export const deleteCohort = (id: string) => {
    _cohorts = _cohorts.filter(c => c.id !== id);
};

// Schedules (CRUD)
export const addScheduleEvent = (evt: ScheduleEvent) => {
    _scheduleEvents = [evt, ..._scheduleEvents];
};
export const updateScheduleEvent = (id: string, updates: Partial<ScheduleEvent>) => {
    _scheduleEvents = _scheduleEvents.map(e => e.id === id ? { ...e, ...updates } : e);
};
export const deleteScheduleEvent = (id: string) => {
    _scheduleEvents = _scheduleEvents.filter(e => e.id !== id);
};

// Rooms (CRUD)
export const addRoom = (room: Room) => {
    _rooms = [room, ..._rooms];
};
export const updateRoom = (id: string, updates: Partial<Room>) => {
    _rooms = _rooms.map(r => r.id === id ? { ...r, ...updates } : r);
};
export const deleteRoom = (id: string) => {
    _rooms = _rooms.filter(r => r.id !== id);
};

// Exports for backward compatibility
export const MOCK_PROGRAMS = _programs; 
export const MOCK_GLOBAL_COURSES = _globalCourses;
export const MOCK_COHORTS = _cohorts;
