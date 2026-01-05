import { api } from '@/api/client';
import type {
  ActivitySubmission,
  CourseActivity,
  CourseModule,
  StudentAssignment,
  StudentAssignmentDetail,
  StudentCourse,
  StudentCourseDetail,
  StudentDashboard,
  StudentGradeEntry,
  StudentAnnouncement,
} from './types';

export const getStudentDashboard = () => api.get<StudentDashboard>('/student/dashboard');
export const getStudentCourses = () => api.get<StudentCourse[]>('/student/courses');
export const getStudentAssignments = () => api.get<StudentAssignment[]>('/student/assignments');
export const getStudentGrades = () => api.get<StudentGradeEntry[]>('/student/grades');

export const getStudentCourseDetail = (courseOfferingId: string) =>
  api.get<StudentCourseDetail>(`/student/courses/${courseOfferingId}`);

export const getStudentCourseModules = (courseOfferingId: string) =>
  api.get<CourseModule[]>(`/student/courses/${courseOfferingId}/modules`);

export const getStudentCourseAnnouncements = (courseOfferingId: string) =>
  api.get<StudentAnnouncement[]>(`/student/courses/${courseOfferingId}/announcements`);

export const getStudentCourseResources = (courseOfferingId: string) =>
  api.get<CourseActivity[]>(`/student/courses/${courseOfferingId}/resources`);

export const getStudentAssignmentDetail = (activityId: string, courseOfferingId?: string) =>
  api.get<StudentAssignmentDetail>(
    `/student/assignments/${activityId}${courseOfferingId ? `?course_offering_id=${courseOfferingId}` : ''}`
  );

export const getMyAssignmentSubmission = (activityId: string, courseOfferingId?: string) =>
  api.get<{ submission: ActivitySubmission | null; course_offering_id: string }>(
    `/student/assignments/${activityId}/submission${courseOfferingId ? `?course_offering_id=${courseOfferingId}` : ''}`
  );

export const submitAssignment = (
  activityId: string,
  payload: { course_offering_id?: string; content: any; status?: string }
) => api.post<ActivitySubmission>(`/student/assignments/${activityId}/submit`, payload);
