import { api } from '@/api/client';
import { AttendanceUpdate, ClassAttendance, ClassSession } from './types';

export const listSessions = (offeringId: string, start?: Date, end?: Date) => {
  const params = new URLSearchParams();
  params.set('offering_id', offeringId);
  if (start) params.set('start', start.toISOString());
  if (end) params.set('end', end.toISOString());
  return api.get<ClassSession[]>(`/scheduling/sessions?${params.toString()}`);
};

export const getSessionAttendance = (sessionId: string) =>
  api.get<ClassAttendance[]>(`/teacher/sessions/${sessionId}/attendance`);

export const batchRecordAttendance = (sessionId: string, updates: AttendanceUpdate[]) =>
  api.post(`/teacher/sessions/${sessionId}/attendance`, { updates });

