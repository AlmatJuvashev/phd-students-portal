import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { CheckCircle2, Clock, Loader2, RotateCcw, Send, XCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { completeAttempt, getAttemptDetails, submitResponse } from './api';
import type { AttemptDetailsResponse, Question, QuestionType } from './types';

type ApiErrorPayload = {
  error?: string;
  code?: string;
  attempt?: { id?: string };
};

const parseApiError = (err: unknown): ApiErrorPayload | null => {
  if (!err || typeof err !== 'object') return null;
  const message = (err as { message?: unknown }).message;
  if (typeof message !== 'string') return null;
  try {
    return JSON.parse(message) as ApiErrorPayload;
  } catch {
    return null;
  }
};

type AnswerState = Record<string, { optionId?: string; text?: string }>;

const supported: Record<QuestionType, boolean> = {
  MCQ: true,
  TRUE_FALSE: true,
  TEXT: true,
  MRQ: false,
  LIKERT: false,
};

const formatSeconds = (totalSeconds: number) => {
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
};

export const AttemptPage: React.FC = () => {
  const { attemptId } = useParams();
  const navigate = useNavigate();
  const qc = useQueryClient();

  const detailsQuery = useQuery({
    queryKey: ['attempt', attemptId],
    queryFn: () => getAttemptDetails(attemptId!),
    enabled: Boolean(attemptId),
  });

  const details = detailsQuery.data as AttemptDetailsResponse | undefined;

  const initializedAttemptRef = useRef<string | null>(null);
  const [answers, setAnswers] = useState<AnswerState>({});

  useEffect(() => {
    if (!details?.attempt?.id) return;
    if (initializedAttemptRef.current === details.attempt.id) return;
    initializedAttemptRef.current = details.attempt.id;

    const next: AnswerState = {};
    for (const r of details.responses || []) {
      next[r.question_id] = {
        optionId: r.selected_option_id || undefined,
        text: r.text_response || undefined,
      };
    }
    setAnswers(next);
  }, [details?.attempt?.id, details?.responses]);

  const submitMutation = useMutation({
    mutationFn: (payload: { question_id: string; option_id?: string; text_response?: string }) =>
      submitResponse(attemptId!, payload),
    onError: (err) => {
      const payload = parseApiError(err);
      if (payload?.code === 'ATTEMPT_AUTO_SUBMITTED' && payload.attempt?.id) {
        qc.invalidateQueries({ queryKey: ['attempt', payload.attempt.id] });
      }
    },
  });

  const completeMutation = useMutation({
    mutationFn: () => completeAttempt(attemptId!),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['attempt', attemptId] }),
  });

  const responseByQuestion = useMemo(() => {
    const map = new Map<string, (AttemptDetailsResponse['responses'][number] | undefined)>();
    for (const r of details?.responses || []) {
      map.set(r.question_id, r);
    }
    return map;
  }, [details?.responses]);

  const totalQuestions = details?.questions?.filter((q) => supported[q.type]).length || 0;
  const answeredCount = useMemo(() => {
    if (!details?.questions) return 0;
    return details.questions.reduce((count, q) => {
      if (!supported[q.type]) return count;
      const a = answers[q.id];
      if (q.type === 'TEXT') return count + (a?.text?.trim() ? 1 : 0);
      return count + (a?.optionId ? 1 : 0);
    }, 0);
  }, [answers, details?.questions]);

  const timeLimitMinutes = details?.assessment?.time_limit_minutes || 0;
  const [nowMs, setNowMs] = useState(() => Date.now());
  useEffect(() => {
    if (!details?.attempt?.started_at) return;
    if (!timeLimitMinutes || timeLimitMinutes <= 0) return;
    if (details.attempt.status !== 'IN_PROGRESS') return;

    const id = window.setInterval(() => setNowMs(Date.now()), 1000);
    return () => window.clearInterval(id);
  }, [details?.attempt?.started_at, details?.attempt?.status, timeLimitMinutes]);

  const remainingSeconds = useMemo(() => {
    if (!details?.attempt?.started_at) return null;
    if (!timeLimitMinutes || timeLimitMinutes <= 0) return null;
    const startedAtMs = new Date(details.attempt.started_at).getTime();
    const expiresAtMs = startedAtMs + timeLimitMinutes * 60_000;
    return Math.max(0, Math.floor((expiresAtMs - nowMs) / 1000));
  }, [details?.attempt?.started_at, nowMs, timeLimitMinutes]);

  useEffect(() => {
    if (!details?.attempt?.id) return;
    if (details.attempt.status !== 'IN_PROGRESS') return;
    if (remainingSeconds === null) return;
    if (remainingSeconds > 0) return;
    if (completeMutation.isPending) return;
    completeMutation.mutate();
  }, [completeMutation, details?.attempt?.id, details?.attempt?.status, remainingSeconds]);

  const onSelectOption = (question: Question, optionId: string) => {
    setAnswers((prev) => ({ ...prev, [question.id]: { ...(prev[question.id] || {}), optionId } }));
    submitMutation.mutate({ question_id: question.id, option_id: optionId });
  };

  const onSaveText = (question: Question) => {
    const text = answers[question.id]?.text || '';
    submitMutation.mutate({ question_id: question.id, text_response: text });
  };

  const renderQuestion = (q: Question) => {
    const a = answers[q.id];
    const locked = details?.attempt?.status !== 'IN_PROGRESS';

    if (!supported[q.type]) {
      return (
        <div className="p-4 border border-slate-200 rounded-2xl bg-slate-50 text-slate-600 text-sm">
          <div className="font-bold">{q.stem}</div>
          <div className="mt-1 text-xs text-slate-500">
            Question type <span className="font-mono">{q.type}</span> is not supported yet.
          </div>
        </div>
      );
    }

    const response = responseByQuestion.get(q.id);
    const isCorrect = response?.is_correct;

    return (
      <div className="p-5 border border-slate-200 rounded-3xl bg-white shadow-sm">
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0">
            <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">Question</div>
            <div className="mt-1 font-bold text-slate-900">{q.stem}</div>
          </div>
          {details?.attempt?.status !== 'IN_PROGRESS' && (
            <Badge
              variant="secondary"
              className={cn(
                isCorrect ? 'bg-emerald-100 text-emerald-800' : 'bg-red-100 text-red-800'
              )}
            >
              {isCorrect ? 'Correct' : 'Incorrect'}
            </Badge>
          )}
        </div>

        {(q.type === 'MCQ' || q.type === 'TRUE_FALSE') && (
          <div className="mt-4 grid gap-2">
            {(q.options || []).map((opt) => {
              const selected = a?.optionId === opt.id;
              const showCorrectness = details?.attempt?.status !== 'IN_PROGRESS';
              const correct = opt.is_correct;
              return (
                <button
                  key={opt.id}
                  type="button"
                  disabled={locked}
                  onClick={() => onSelectOption(q, opt.id)}
                  className={cn(
                    'text-left p-3 rounded-2xl border transition-colors',
                    selected ? 'border-indigo-300 bg-indigo-50' : 'border-slate-200 hover:bg-slate-50',
                    locked ? 'cursor-not-allowed opacity-80' : ''
                  )}
                >
                  <div className="flex items-center justify-between gap-3">
                    <div className="text-sm text-slate-900">{opt.text}</div>
                    {showCorrectness && (
                      <div className="shrink-0">
                        {correct ? (
                          <CheckCircle2 className="h-5 w-5 text-emerald-600" />
                        ) : (
                          <XCircle className="h-5 w-5 text-slate-300" />
                        )}
                      </div>
                    )}
                  </div>
                  {details?.attempt?.status !== 'IN_PROGRESS' && selected && opt.feedback && (
                    <div className="mt-2 text-xs text-slate-500">{opt.feedback}</div>
                  )}
                </button>
              );
            })}
          </div>
        )}

        {q.type === 'TEXT' && (
          <div className="mt-4 space-y-2">
            <Textarea
              value={a?.text || ''}
              disabled={locked}
              onChange={(e) =>
                setAnswers((prev) => ({ ...prev, [q.id]: { ...(prev[q.id] || {}), text: e.target.value } }))
              }
              placeholder="Type your answer…"
              className="min-h-[120px]"
            />
            {details?.attempt?.status === 'IN_PROGRESS' && (
              <div className="flex justify-end">
                <Button variant="outline" onClick={() => onSaveText(q)} disabled={submitMutation.isPending}>
                  Save answer
                </Button>
              </div>
            )}
          </div>
        )}
      </div>
    );
  };

  if (!attemptId) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-white border border-slate-200 rounded-3xl p-10 shadow-sm text-slate-600">
          Missing attempt ID.
        </div>
      </div>
    );
  }

  if (detailsQuery.isLoading) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-white border border-slate-200 rounded-3xl p-10 shadow-sm flex items-center justify-center gap-3 text-slate-600">
          <Loader2 className="animate-spin" />
          Loading attempt…
        </div>
      </div>
    );
  }

  if (detailsQuery.isError || !details) {
    return (
      <div className="max-w-4xl mx-auto p-6 space-y-4">
        <div className="bg-red-50 border border-red-100 rounded-3xl p-6 text-red-900">
          <div className="font-black text-lg">Failed to load attempt</div>
          <div className="text-sm mt-1">{(detailsQuery.error as Error)?.message || 'Unknown error'}</div>
        </div>
        <Button variant="outline" onClick={() => navigate('/student/dashboard')}>
          Back to dashboard
        </Button>
      </div>
    );
  }

  const percent = Math.round(details.attempt.score * 10) / 10;
  const passed = percent >= (details.assessment.passing_score || 0);

  return (
    <div className="max-w-4xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div className="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm">
        <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
          <div className="min-w-0">
            <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">Assessment</div>
            <div className="mt-1 text-2xl font-black text-slate-900 tracking-tight truncate">
              {details.assessment.title}
            </div>
            {details.assessment.description && (
              <div className="mt-2 text-sm text-slate-600">{details.assessment.description}</div>
            )}
          </div>

          <div className="shrink-0 flex flex-col items-start gap-2">
            {details.attempt.status === 'IN_PROGRESS' && (
              <div className="inline-flex items-center gap-2 text-sm font-bold text-slate-700">
                <Clock className="h-4 w-4" />
                {remainingSeconds === null ? 'No time limit' : formatSeconds(remainingSeconds)}
              </div>
            )}
            {details.attempt.status !== 'IN_PROGRESS' && (
              <Badge
                variant="secondary"
                className={cn(passed ? 'bg-emerald-100 text-emerald-800' : 'bg-red-100 text-red-800')}
              >
                {passed ? 'Passed' : 'Not passed'} · {percent}%
              </Badge>
            )}
          </div>
        </div>

        {details.attempt.status === 'IN_PROGRESS' && (
          <div className="mt-6 space-y-3">
            <div className="flex items-center justify-between text-xs text-slate-500">
              <span>
                Answered <span className="font-bold text-slate-700">{answeredCount}</span> / {totalQuestions}
              </span>
              <span className="font-mono">
                {answeredCount === totalQuestions ? 'Ready to submit' : 'In progress'}
              </span>
            </div>
            <Progress value={totalQuestions > 0 ? (answeredCount / totalQuestions) * 100 : 0} />
            <div className="flex items-center justify-end gap-2">
              <Button
                onClick={() => completeMutation.mutate()}
                disabled={completeMutation.isPending}
                className="font-bold"
              >
                {completeMutation.isPending ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Submitting…
                  </>
                ) : (
                  <>
                    <Send className="mr-2 h-4 w-4" /> Submit attempt
                  </>
                )}
              </Button>
            </div>
          </div>
        )}

        {details.attempt.status !== 'IN_PROGRESS' && (
          <div className="mt-6 flex flex-wrap gap-2">
            <Button
              variant="outline"
              onClick={() => navigate(`/student/assessments/${details.assessment.id}`)}
              className="font-bold"
            >
              <RotateCcw className="mr-2 h-4 w-4" /> Try again
            </Button>
            <Button variant="outline" onClick={() => navigate('/student/grades')}>
              View grades
            </Button>
          </div>
        )}
      </div>

      <div className="grid gap-4">
        {details.questions.map((q) => (
          <div key={q.id}>{renderQuestion(q)}</div>
        ))}
      </div>
    </div>
  );
};

