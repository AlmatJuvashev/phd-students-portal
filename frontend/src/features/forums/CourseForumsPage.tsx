import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Loader2, MessageSquare } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { listCourseForums } from './api';

export const CourseForumsPage: React.FC = () => {
  const navigate = useNavigate();
  const { courseOfferingId } = useParams<{ courseOfferingId: string }>();

  const forumsQuery = useQuery({
    queryKey: ['forums', 'course', courseOfferingId],
    queryFn: () => listCourseForums(courseOfferingId!),
    enabled: !!courseOfferingId,
  });

  if (!courseOfferingId) {
    return <div className="p-8 text-center text-slate-500">Course not found.</div>;
  }

  if (forumsQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  const forums = forumsQuery.data || [];

  return (
    <div className="max-w-5xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="space-y-4">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} /> Back
        </button>

        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-3xl font-black text-slate-900 tracking-tight">Forums</h1>
            <p className="text-slate-500 font-medium mt-1">Course offering: {courseOfferingId}</p>
          </div>
          <Button variant="secondary" onClick={() => navigate('/chat')}>
            <MessageSquare className="mr-2 h-4 w-4" />
            Open Chat
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {forums.map((f) => (
          <button
            key={f.id}
            onClick={() => navigate(`/forums/course/${courseOfferingId}/forums/${f.id}`)}
            className="group bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-lg hover:border-indigo-200 transition-all p-6 text-left"
          >
            <div className="flex items-center justify-between gap-3">
              <div className="font-black text-slate-900">{f.title}</div>
              <Badge variant="secondary" className="uppercase text-[10px]">
                {f.type}
              </Badge>
            </div>
            <div className="text-sm text-slate-600 mt-2">{f.description || '—'}</div>
            <div className="mt-4 text-xs text-indigo-600 font-bold group-hover:translate-x-1 transition-transform">
              Open →
            </div>
          </button>
        ))}

        {forums.length === 0 && (
          <div className="col-span-full text-center py-12 text-slate-400 border-2 border-dashed border-slate-200 rounded-2xl">
            No forums yet.
          </div>
        )}
      </div>
    </div>
  );
};

