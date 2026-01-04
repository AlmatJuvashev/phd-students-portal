import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { startAttempt } from './api';

type ApiErrorPayload = {
  error?: string;
  code?: string;
  attempt?: { id?: string };
  retry_after_seconds?: number;
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

export const StartAssessmentPage: React.FC = () => {
  const { assessmentId } = useParams();
  const navigate = useNavigate();
  const startedRef = useRef(false);
  const [errorText, setErrorText] = useState<string | null>(null);

  const startMutation = useMutation({
    mutationFn: async () => {
      if (!assessmentId) throw new Error('missing assessmentId');
      return startAttempt(assessmentId);
    },
    onSuccess: (attempt) => {
      navigate(`/student/attempts/${attempt.id}`, { replace: true });
    },
    onError: (err) => {
      const payload = parseApiError(err);
      if (payload?.code === 'ATTEMPT_IN_PROGRESS' && payload.attempt?.id) {
        navigate(`/student/attempts/${payload.attempt.id}`, { replace: true });
        return;
      }
      if (payload?.code === 'COOLDOWN_ACTIVE') {
        const seconds = payload.retry_after_seconds ?? 0;
        setErrorText(`Cooldown active. Try again in ${Math.ceil(seconds / 60)} minutes.`);
        return;
      }
      if (payload?.code === 'MAX_ATTEMPTS_REACHED') {
        setErrorText('Max attempts reached for this assessment.');
        return;
      }
      const message = (err as Error)?.message || 'Failed to start assessment.';
      setErrorText(payload?.error || message);
    },
  });

  const canAutoStart = useMemo(() => Boolean(assessmentId) && !startedRef.current, [assessmentId]);

  useEffect(() => {
    if (!canAutoStart) return;
    startedRef.current = true;
    startMutation.mutate();
  }, [canAutoStart, startMutation]);

  if (startMutation.isPending) {
    return (
      <div className="max-w-3xl mx-auto p-6">
        <div className="bg-white border border-slate-200 rounded-3xl p-10 shadow-sm flex items-center justify-center gap-3 text-slate-600">
          <Loader2 className="animate-spin" />
          Starting assessment…
        </div>
      </div>
    );
  }

  if (errorText) {
    return (
      <div className="max-w-3xl mx-auto p-6 space-y-4">
        <div className="bg-red-50 border border-red-100 rounded-3xl p-6 text-red-900">
          <div className="font-black text-lg">Unable to start assessment</div>
          <div className="text-sm mt-1">{errorText}</div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate('/student/assignments')}>
            Back to assignments
          </Button>
          <Button onClick={() => startMutation.mutate()}>Try again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto p-6">
      <div className="bg-white border border-slate-200 rounded-3xl p-10 shadow-sm text-slate-600">
        Unable to start assessment.
      </div>
    </div>
  );
};

