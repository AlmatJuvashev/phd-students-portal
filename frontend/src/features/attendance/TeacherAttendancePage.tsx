import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Loader2, Save } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ScrollArea } from '@/components/ui/scroll-area';
import { cn } from '@/lib/utils';
import { getCourses } from '@/features/curriculum/api';
import { getCourseRoster, getTeacherCourses } from '@/features/teacher/api';
import { batchRecordAttendance, getSessionAttendance, listSessions } from './api';
import { AttendanceUpdate } from './types';

const STATUS_OPTIONS = ['PRESENT', 'ABSENT', 'LATE', 'EXCUSED'] as const;

export const TeacherAttendancePage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { courseId: offeringId } = useParams<{ courseId: string }>();
  const qc = useQueryClient();

  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null);
  const [draft, setDraft] = useState<Record<string, { status: string; notes: string }>>({});

  const offeringsQuery = useQuery({ queryKey: ['teacher', 'courses'], queryFn: getTeacherCourses });
  const catalogQuery = useQuery({
    queryKey: ['curriculum', 'courses'],
    queryFn: () => getCourses(),
    staleTime: 5 * 60 * 1000,
    retry: false,
  });

  const offering = useMemo(
    () => (offeringsQuery.data || []).find((o) => o.id === offeringId),
    [offeringsQuery.data, offeringId]
  );
  const course = useMemo(() => {
    if (!offering) return undefined;
    return (catalogQuery.data || []).find((c) => c.id === offering.course_id);
  }, [catalogQuery.data, offering]);

  const timeWindow = useMemo(() => {
    const start = new Date();
    start.setDate(start.getDate() - 30);
    const end = new Date();
    end.setDate(end.getDate() + 30);
    return { start, end };
  }, []);

  const sessionsQuery = useQuery({
    queryKey: ['attendance', 'sessions', offeringId, timeWindow.start.toISOString(), timeWindow.end.toISOString()],
    queryFn: () => listSessions(offeringId!, timeWindow.start, timeWindow.end),
    enabled: !!offeringId,
  });

  const rosterQuery = useQuery({
    queryKey: ['teacher', 'courses', offeringId, 'roster'],
    queryFn: () => getCourseRoster(offeringId!),
    enabled: !!offeringId,
  });

  const sessions = useMemo(() => {
    const list = sessionsQuery.data || [];
    return [...list].sort((a, b) => (a.date > b.date ? 1 : -1));
  }, [sessionsQuery.data]);

  useEffect(() => {
    if (!selectedSessionId && sessions.length > 0) {
      setSelectedSessionId(sessions[0].id);
    }
  }, [selectedSessionId, sessions]);

  const attendanceQuery = useQuery({
    queryKey: ['attendance', 'session', selectedSessionId],
    queryFn: () => getSessionAttendance(selectedSessionId!),
    enabled: !!selectedSessionId,
  });

  useEffect(() => {
    const roster = rosterQuery.data || [];
    const records = attendanceQuery.data || [];
    const byStudent = new Map(records.map((r) => [r.student_id, r]));

    const next: Record<string, { status: string; notes: string }> = {};
    for (const e of roster) {
      const r = byStudent.get(e.student_id);
      next[e.student_id] = {
        status: r?.status || 'PRESENT',
        notes: r?.notes || '',
      };
    }
    setDraft(next);
  }, [selectedSessionId, rosterQuery.data, attendanceQuery.data]);

  const saveMutation = useMutation({
    mutationFn: (updates: AttendanceUpdate[]) => batchRecordAttendance(selectedSessionId!, updates),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['attendance', 'session', selectedSessionId] });
    },
  });

  if (offeringsQuery.isLoading || sessionsQuery.isLoading || rosterQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (!offeringId || !offering) {
    return <div className="p-8 text-center text-slate-500">{t('teacher.detail.not_found')}</div>;
  }

  const selectedSession = sessions.find((s) => s.id === selectedSessionId) || null;
  const roster = rosterQuery.data || [];

  return (
    <div className="max-w-6xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="space-y-4">
        <button
          onClick={() => navigate(`/admin/teacher/courses/${offeringId}`)}
          className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} /> {t('teacher.tracker.back')}
        </button>

        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
          <div>
            <div className="flex items-center gap-3 mb-1">
              <Badge variant="outline" className="bg-white border-slate-200 text-slate-500">
                {course?.code || 'COURSE'}
              </Badge>
              <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">{offering.section}</span>
              <Badge variant="secondary" className="bg-slate-100 text-slate-700">
                {t('teacher.courses.detail.attendance')}
              </Badge>
            </div>
            <h1 className="text-3xl font-black text-slate-900 tracking-tight">{t('teacher.courses.detail.attendance')}</h1>
            <p className="text-slate-500 font-medium mt-1">{course?.title || offering.course_id}</p>
          </div>

          <Button
            disabled={!selectedSessionId || saveMutation.isPending}
            onClick={() => {
              const updates: AttendanceUpdate[] = Object.entries(draft).map(([studentId, v]) => ({
                student_id: studentId,
                status: v.status,
                notes: v.notes,
              }));
              saveMutation.mutate(updates);
            }}
          >
            {saveMutation.isPending ? <Loader2 className="animate-spin mr-2 h-4 w-4" /> : <Save className="mr-2 h-4 w-4" />}
            Save Attendance
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="bg-white border border-slate-200 rounded-2xl shadow-sm overflow-hidden">
          <div className="px-5 py-4 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
            Sessions
          </div>
          <ScrollArea className="h-[520px]">
            <div className="p-3 space-y-2">
              {sessions.map((s) => {
                const active = s.id === selectedSessionId;
                return (
                  <button
                    key={s.id}
                    onClick={() => setSelectedSessionId(s.id)}
                    className={cn(
                      'w-full text-left rounded-xl border p-3 transition-colors',
                      active ? 'border-indigo-200 bg-indigo-50' : 'border-slate-200 hover:bg-slate-50'
                    )}
                  >
                    <div className="font-bold text-slate-900">{s.title}</div>
                    <div className="text-xs text-slate-500 mt-1">
                      {new Date(s.date).toLocaleDateString()} • {s.start_time}–{s.end_time}
                    </div>
                  </button>
                );
              })}
              {sessions.length === 0 && <div className="p-6 text-center text-slate-400 italic">No sessions found.</div>}
            </div>
          </ScrollArea>
        </div>

        <div className="lg:col-span-2 bg-white border border-slate-200 rounded-2xl shadow-sm overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-200 flex items-center justify-between">
            <div>
              <div className="text-xs font-bold text-slate-500 uppercase">Roster</div>
              <div className="text-sm font-black text-slate-900 mt-1">{selectedSession?.title || 'Select a session'}</div>
            </div>
            {attendanceQuery.isFetching ? (
              <div className="text-xs text-slate-500 flex items-center gap-2">
                <Loader2 className="animate-spin h-3 w-3" /> Loading
              </div>
            ) : null}
          </div>

          <table className="w-full text-sm text-left">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
              <tr>
                <th className="px-6 py-3">Student</th>
                <th className="px-6 py-3">Status</th>
                <th className="px-6 py-3">Notes</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {roster.map((e) => (
                <tr key={e.id} className="hover:bg-slate-50 transition-colors">
                  <td className="px-6 py-4 font-bold text-slate-900">{e.student_name || e.student_id}</td>
                  <td className="px-6 py-4">
                    <select
                      className="h-9 px-3 rounded-md border border-slate-200 bg-white text-sm"
                      value={draft[e.student_id]?.status || 'PRESENT'}
                      onChange={(ev) =>
                        setDraft((prev) => ({
                          ...prev,
                          [e.student_id]: { status: ev.target.value, notes: prev[e.student_id]?.notes || '' },
                        }))
                      }
                    >
                      {STATUS_OPTIONS.map((s) => (
                        <option key={s} value={s}>
                          {s}
                        </option>
                      ))}
                    </select>
                  </td>
                  <td className="px-6 py-4">
                    <Input
                      value={draft[e.student_id]?.notes || ''}
                      onChange={(ev) =>
                        setDraft((prev) => ({
                          ...prev,
                          [e.student_id]: { status: prev[e.student_id]?.status || 'PRESENT', notes: ev.target.value },
                        }))
                      }
                      placeholder="Optional note"
                    />
                  </td>
                </tr>
              ))}
              {roster.length === 0 && (
                <tr>
                  <td colSpan={3} className="p-12 text-center text-slate-400 italic">
                    {t('teacher.detail.no_students')}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
