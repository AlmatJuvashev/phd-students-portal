import React, { useMemo } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, CalendarClock, ExternalLink, FileText, Loader2, Megaphone, Package, PlayCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Accordion, AccordionItem } from '@/components/ui/accordion';
import { cn } from '@/lib/utils';
import {
  getStudentCourseAnnouncements,
  getStudentCourseDetail,
  getStudentCourseModules,
  getStudentCourseResources,
} from './api';
import type { CourseActivity, CourseModule, StudentAnnouncement } from './types';

const tryParse = (raw: string): any => {
  if (!raw) return null;
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
};

export const StudentCourseDetail: React.FC = () => {
  const { courseOfferingId } = useParams<{ courseOfferingId: string }>();
  const navigate = useNavigate();

  const courseQuery = useQuery({
    queryKey: ['student', 'courses', courseOfferingId],
    queryFn: () => getStudentCourseDetail(courseOfferingId!),
    enabled: Boolean(courseOfferingId),
  });

  const modulesQuery = useQuery({
    queryKey: ['student', 'courses', courseOfferingId, 'modules'],
    queryFn: () => getStudentCourseModules(courseOfferingId!),
    enabled: Boolean(courseOfferingId),
  });

  const announcementsQuery = useQuery({
    queryKey: ['student', 'courses', courseOfferingId, 'announcements'],
    queryFn: () => getStudentCourseAnnouncements(courseOfferingId!),
    enabled: Boolean(courseOfferingId),
  });

  const resourcesQuery = useQuery({
    queryKey: ['student', 'courses', courseOfferingId, 'resources'],
    queryFn: () => getStudentCourseResources(courseOfferingId!),
    enabled: Boolean(courseOfferingId),
  });

  const modules = (modulesQuery.data || []) as CourseModule[];
  const announcements = (announcementsQuery.data || []) as StudentAnnouncement[];
  const resources = (resourcesQuery.data || []) as CourseActivity[];

  const openActivity = (activity: CourseActivity) => {
    const parsed = tryParse(activity.content);
    if (activity.type === 'assignment') {
      navigate(`/student/assignments/${activity.id}?course_offering_id=${courseOfferingId}`);
      return;
    }
    if (activity.type === 'quiz') {
      const assessmentId =
        parsed?.assessment_id || parsed?.assessmentId || parsed?.assessmentID || parsed?.assessment;
      if (assessmentId) {
        navigate(`/student/assessments/${assessmentId}`);
      }
      return;
    }
    if (activity.type === 'resource') {
      const url = parsed?.url || parsed?.resource_url;
      if (url) window.open(url, '_blank', 'noopener,noreferrer');
    }
  };

  const hasLoading =
    courseQuery.isLoading || modulesQuery.isLoading || announcementsQuery.isLoading || resourcesQuery.isLoading;

  const hasError =
    courseQuery.isError || modulesQuery.isError || announcementsQuery.isError || resourcesQuery.isError;

  const course = courseQuery.data?.course;
  const sessions = courseQuery.data?.sessions || [];

  const schedulePreview = useMemo(() => sessions.slice(0, 6), [sessions]);

  if (hasLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (hasError || !course) {
    return (
      <div className="max-w-5xl mx-auto p-6 space-y-4">
        <div className="bg-red-50 border border-red-100 rounded-3xl p-6 text-red-900">
          <div className="font-black text-lg">Failed to load course.</div>
          <div className="text-sm mt-1">
            {(courseQuery.error as Error)?.message || (modulesQuery.error as Error)?.message || 'Unknown error'}
          </div>
        </div>
        <Button variant="outline" onClick={() => navigate('/student/courses')}>
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to courses
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm">
        <div className="flex flex-col gap-4">
          <div className="flex items-start justify-between gap-4">
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <Button variant="ghost" size="icon" onClick={() => navigate('/student/courses')}>
                  <ArrowLeft className="h-5 w-5" />
                </Button>
                <Badge variant="secondary" className="font-mono">
                  {course.code}
                </Badge>
                <Badge variant="outline" className="text-[10px]">
                  Section {course.section}
                </Badge>
              </div>
              <div className="mt-2 text-2xl font-black text-slate-900 tracking-tight truncate">{course.title}</div>
              <div className="mt-1 text-sm text-slate-500 truncate">{course.instructor_name || '—'}</div>
            </div>

            <div className="shrink-0 flex flex-col items-end gap-2">
              <Badge variant="outline" className="text-[10px]">
                {course.delivery_format}
              </Badge>
              {course.next_session && (
                <div className="flex items-start gap-2 text-xs text-slate-600">
                  <CalendarClock className="h-4 w-4 text-slate-400 mt-0.5" />
                  <div>
                    <div className="font-bold">{course.next_session.date}</div>
                    <div className="text-[11px] text-slate-500">
                      {course.next_session.start_time}–{course.next_session.end_time} · {course.next_session.type}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>

          {schedulePreview.length > 0 && (
            <div className="mt-2 grid grid-cols-1 md:grid-cols-2 gap-3">
              {schedulePreview.map((s) => (
                <div key={s.id} className="p-4 rounded-2xl border border-slate-200 bg-slate-50">
                  <div className="text-xs text-slate-500 flex items-center justify-between">
                    <span className="font-mono">{s.date}</span>
                    <span className="text-[10px] uppercase font-bold">{s.type}</span>
                  </div>
                  <div className="mt-1 font-bold text-slate-900">
                    {s.start_time}–{s.end_time}
                  </div>
                  {s.meeting_url && (
                    <button
                      type="button"
                      onClick={() => window.open(s.meeting_url!, '_blank', 'noopener,noreferrer')}
                      className="mt-2 inline-flex items-center text-xs font-bold text-indigo-700 hover:underline"
                    >
                      <ExternalLink className="mr-1 h-3 w-3" /> Join meeting
                    </button>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <Tabs defaultValue="content" className="bg-white border border-slate-200 rounded-3xl shadow-sm">
        <TabsList className="w-full justify-start rounded-none rounded-t-3xl border-b border-slate-100 bg-white px-4 py-3">
          <TabsTrigger value="content" className="font-bold">
            Content
          </TabsTrigger>
          <TabsTrigger value="announcements" className="font-bold">
            Announcements
          </TabsTrigger>
          <TabsTrigger value="forums" className="font-bold">
            Forums
          </TabsTrigger>
          <TabsTrigger value="resources" className="font-bold">
            Resources
          </TabsTrigger>
        </TabsList>

        <TabsContent value="content" className="p-6">
          {modules.length === 0 ? (
            <div className="p-10 text-center text-slate-400 italic border-2 border-dashed border-slate-200 rounded-3xl">
              No modules yet.
            </div>
          ) : (
            <div className="space-y-3">
              <Accordion>
                {modules.map((m) => (
                  <AccordionItem
                    key={m.id}
                    header={
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-xl bg-slate-50 flex items-center justify-center text-slate-500">
                          <Package className="h-4 w-4" />
                        </div>
                        <div className="min-w-0">
                          <div className="truncate">{m.title}</div>
                          <div className="text-xs text-slate-500 font-normal">Module</div>
                        </div>
                      </div>
                    }
                  >
                    <div className="space-y-4">
                      {m.lessons.map((l) => (
                        <div key={l.id} className="p-4 border border-slate-100 rounded-2xl bg-slate-50">
                          <div className="font-bold text-slate-900">{l.title}</div>
                          <div className="mt-3 grid gap-2">
                            {l.activities.map((a) => {
                              const clickable = ['assignment', 'quiz', 'resource'].includes(a.type);
                              return (
                                <button
                                  key={a.id}
                                  type="button"
                                  onClick={() => (clickable ? openActivity(a) : undefined)}
                                  className={cn(
                                    'w-full text-left p-3 rounded-2xl border flex items-center justify-between gap-3',
                                    clickable
                                      ? 'border-slate-200 bg-white hover:bg-slate-50 transition-colors'
                                      : 'border-slate-200 bg-white opacity-70 cursor-default'
                                  )}
                                >
                                  <div className="min-w-0">
                                    <div className="flex items-center gap-2">
                                      {a.type === 'quiz' ? (
                                        <PlayCircle className="h-4 w-4 text-indigo-600" />
                                      ) : a.type === 'assignment' ? (
                                        <FileText className="h-4 w-4 text-emerald-600" />
                                      ) : a.type === 'resource' ? (
                                        <Package className="h-4 w-4 text-slate-600" />
                                      ) : (
                                        <FileText className="h-4 w-4 text-slate-400" />
                                      )}
                                      <span className="font-bold text-slate-900 truncate">{a.title}</span>
                                      <Badge variant="outline" className="text-[10px] uppercase">
                                        {a.type}
                                      </Badge>
                                    </div>
                                    {a.points > 0 && <div className="text-xs text-slate-500 mt-1">{a.points} pts</div>}
                                  </div>
                                  {clickable && <ExternalLink className="h-4 w-4 text-slate-400 shrink-0" />}
                                </button>
                              );
                            })}
                            {l.activities.length === 0 && <div className="text-xs text-slate-500 italic">No activities.</div>}
                          </div>
                        </div>
                      ))}
                      {m.lessons.length === 0 && <div className="text-xs text-slate-500 italic">No lessons.</div>}
                    </div>
                  </AccordionItem>
                ))}
              </Accordion>
            </div>
          )}
        </TabsContent>

        <TabsContent value="announcements" className="p-6">
          {announcements.length === 0 ? (
            <div className="p-10 text-center text-slate-400 italic border-2 border-dashed border-slate-200 rounded-3xl">
              No announcements.
            </div>
          ) : (
            <div className="space-y-3">
              {announcements.map((a) => (
                <div key={a.id} className="p-5 border border-slate-200 rounded-3xl bg-white">
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <div className="flex items-center gap-2">
                        <Megaphone className="h-4 w-4 text-indigo-600" />
                        <div className="font-black text-slate-900 truncate">{a.title}</div>
                      </div>
                      <div className="mt-1 text-xs text-slate-500">{new Date(a.created_at).toLocaleString()}</div>
                    </div>
                  </div>
                  <div className="mt-3 text-sm text-slate-700 whitespace-pre-wrap">{a.body}</div>
                </div>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="forums" className="p-6">
          <div className="bg-slate-50 border border-slate-200 rounded-3xl p-6 flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
            <div>
              <div className="font-black text-slate-900">Course forums</div>
              <div className="text-sm text-slate-600 mt-1">Ask questions, discuss topics, and read updates.</div>
            </div>
            <Button onClick={() => navigate(`/forums/course/${courseOfferingId}`)}>Open Forums</Button>
          </div>
        </TabsContent>

        <TabsContent value="resources" className="p-6">
          {resources.length === 0 ? (
            <div className="p-10 text-center text-slate-400 italic border-2 border-dashed border-slate-200 rounded-3xl">
              No resources.
            </div>
          ) : (
            <div className="grid gap-3">
              {resources.map((r) => (
                <button
                  key={r.id}
                  type="button"
                  onClick={() => openActivity(r)}
                  className="p-4 border border-slate-200 rounded-2xl bg-white hover:bg-slate-50 transition-colors flex items-center justify-between gap-3"
                >
                  <div className="min-w-0">
                    <div className="font-bold text-slate-900 truncate">{r.title}</div>
                    <div className="text-xs text-slate-500 mt-1">Resource</div>
                  </div>
                  <ExternalLink className="h-4 w-4 text-slate-400 shrink-0" />
                </button>
              ))}
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};
