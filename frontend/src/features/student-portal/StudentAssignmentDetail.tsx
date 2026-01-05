import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { ArrowLeft, FileText, Loader2, Save, Send } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { getStudentAssignmentDetail, submitAssignment } from './api';

export const StudentAssignmentDetail: React.FC = () => {
  const { assignmentId } = useParams<{ assignmentId: string }>();
  const [searchParams] = useSearchParams();
  const courseOfferingId = searchParams.get('course_offering_id') || undefined;
  const navigate = useNavigate();
  const qc = useQueryClient();

  const detailQuery = useQuery({
    queryKey: ['student', 'assignment', assignmentId, courseOfferingId],
    queryFn: () => getStudentAssignmentDetail(assignmentId!, courseOfferingId),
    enabled: Boolean(assignmentId),
  });

  const detail = detailQuery.data;
  const activity = detail?.activity;
  const submission = detail?.submission || null;

  const initialText = useMemo(() => {
    const text = submission?.content?.text;
    return typeof text === 'string' ? text : '';
  }, [submission?.content]);

  const initialFileUrl = useMemo(() => {
    const url = submission?.content?.file_url || submission?.content?.fileUrl;
    return typeof url === 'string' ? url : '';
  }, [submission?.content]);

  const [text, setText] = useState('');
  const [fileUrl, setFileUrl] = useState('');

  useEffect(() => {
    setText(initialText);
    setFileUrl(initialFileUrl);
  }, [initialFileUrl, initialText]);

  const submitMutation = useMutation({
    mutationFn: async (status: 'DRAFT' | 'SUBMITTED') => {
      if (!assignmentId) throw new Error('Missing assignment id');
      return submitAssignment(assignmentId, {
        course_offering_id: detail?.course_offering_id || courseOfferingId,
        status,
        content: { text, file_url: fileUrl || undefined },
      });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['student', 'assignment', assignmentId] });
      await qc.invalidateQueries({ queryKey: ['student', 'assignments'] });
    },
  });

  if (detailQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (detailQuery.isError || !activity) {
    return (
      <div className="max-w-4xl mx-auto p-6 space-y-4">
        <div className="bg-red-50 border border-red-100 rounded-3xl p-6 text-red-900">
          <div className="font-black text-lg">Failed to load assignment.</div>
          <div className="text-sm mt-1">{(detailQuery.error as Error)?.message || 'Unknown error'}</div>
        </div>
        <Button variant="outline" onClick={() => navigate('/student/assignments')}>
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to assignments
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm">
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
                <ArrowLeft className="h-5 w-5" />
              </Button>
              <div className="w-9 h-9 rounded-xl bg-slate-50 flex items-center justify-center text-slate-500">
                <FileText className="h-5 w-5" />
              </div>
              <div className="min-w-0">
                <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">Assignment</div>
                <div className="font-black text-slate-900 truncate">{activity.title}</div>
              </div>
            </div>

            <div className="mt-3 flex flex-wrap items-center gap-2">
              <Badge variant="outline" className="text-[10px] uppercase">
                {activity.type}
              </Badge>
              {activity.points > 0 && (
                <Badge variant="secondary" className="text-[10px] font-mono">
                  {activity.points} pts
                </Badge>
              )}
              <Badge variant="outline" className="text-[10px] font-mono">
                Offering {detail.course_offering_id}
              </Badge>
              {submission && (
                <Badge
                  variant="secondary"
                  className={cn(
                    submission.status === 'GRADED'
                      ? 'bg-emerald-100 text-emerald-800'
                      : submission.status === 'SUBMITTED'
                        ? 'bg-indigo-100 text-indigo-800'
                        : 'bg-slate-100 text-slate-700'
                  )}
                >
                  {submission.status}
                </Badge>
              )}
            </div>
          </div>

          <div className="shrink-0 flex flex-col items-end gap-2">
            <Button
              variant="outline"
              onClick={() => submitMutation.mutate('DRAFT')}
              disabled={submitMutation.isPending}
            >
              {submitMutation.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving…
                </>
              ) : (
                <>
                  <Save className="mr-2 h-4 w-4" /> Save draft
                </>
              )}
            </Button>
            <Button onClick={() => submitMutation.mutate('SUBMITTED')} disabled={submitMutation.isPending}>
              {submitMutation.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Submitting…
                </>
              ) : (
                <>
                  <Send className="mr-2 h-4 w-4" /> Submit
                </>
              )}
            </Button>
          </div>
        </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm space-y-5">
        <div>
          <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">Instructions</div>
          <div className="mt-2 text-sm text-slate-700 whitespace-pre-wrap">
            {(() => {
              try {
                const parsed = activity.content ? JSON.parse(activity.content) : null;
                const instructions = parsed?.instructions || parsed?.prompt || parsed?.content;
                return typeof instructions === 'string' && instructions.trim() ? instructions : 'No instructions provided.';
              } catch {
                return activity.content?.trim() ? activity.content : 'No instructions provided.';
              }
            })()}
          </div>
        </div>

        <div className="grid gap-4">
          <div>
            <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">Text Submission</div>
            <Textarea
              value={text}
              onChange={(e) => setText(e.target.value)}
              placeholder="Write your answer…"
              className="min-h-[160px] mt-2"
            />
          </div>

          <div>
            <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">File URL (optional)</div>
            <Input
              value={fileUrl}
              onChange={(e) => setFileUrl(e.target.value)}
              placeholder="https://..."
              className="mt-2"
            />
            <div className="text-xs text-slate-500 mt-1">
              Use S3 upload / attachments if configured; this field stores a link in your submission JSON.
            </div>
          </div>
        </div>

        {submission && (
          <div className="pt-4 border-t border-slate-100">
            <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">Last submission</div>
            <div className="mt-2 text-sm text-slate-600">
              Submitted at: <span className="font-mono">{new Date(submission.submitted_at).toLocaleString()}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

