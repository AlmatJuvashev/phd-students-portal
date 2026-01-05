import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { ArrowLeft, Loader2, Plus } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { createTopic, listTopics } from './api';

export const ForumTopicsPage: React.FC = () => {
  const navigate = useNavigate();
  const qc = useQueryClient();
  const { courseOfferingId, forumId } = useParams<{ courseOfferingId: string; forumId: string }>();

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');

  const topicsQuery = useQuery({
    queryKey: ['forums', forumId, 'topics'],
    queryFn: () => listTopics(forumId!),
    enabled: !!forumId,
  });

  const createMutation = useMutation({
    mutationFn: () => createTopic(forumId!, { title: title.trim(), content: content.trim() }),
    onSuccess: async () => {
      setTitle('');
      setContent('');
      await qc.invalidateQueries({ queryKey: ['forums', forumId, 'topics'] });
    },
  });

  if (!courseOfferingId || !forumId) {
    return <div className="p-8 text-center text-slate-500">Forum not found.</div>;
  }

  if (topicsQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  const topics = topicsQuery.data || [];

  return (
    <div className="max-w-5xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="space-y-4">
        <button
          onClick={() => navigate(`/forums/course/${courseOfferingId}`)}
          className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} /> Back
        </button>

        <div>
          <h1 className="text-3xl font-black text-slate-900 tracking-tight">Topics</h1>
          <p className="text-slate-500 font-medium mt-1">Forum: {forumId}</p>
        </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm space-y-4">
        <div className="font-black text-slate-900">Start a new topic</div>
        <Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Title" />
        <Textarea value={content} onChange={(e) => setContent(e.target.value)} placeholder="Write your post..." />
        <div className="flex justify-end">
          <Button
            disabled={!title.trim() || !content.trim() || createMutation.isPending}
            onClick={() => createMutation.mutate()}
          >
            {createMutation.isPending ? <Loader2 className="animate-spin mr-2 h-4 w-4" /> : <Plus className="mr-2 h-4 w-4" />}
            Post
          </Button>
        </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm">
        <table className="w-full text-sm text-left">
          <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
            <tr>
              <th className="px-6 py-4">Topic</th>
              <th className="px-6 py-4">Author</th>
              <th className="px-6 py-4">Replies</th>
              <th className="px-6 py-4">Views</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {topics.map((t) => (
              <tr
                key={t.id}
                className="hover:bg-slate-50 transition-colors cursor-pointer"
                onClick={() => navigate(`/forums/course/${courseOfferingId}/topics/${t.id}`)}
              >
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2">
                    {t.is_pinned ? (
                      <Badge variant="secondary" className="uppercase text-[10px]">
                        Pinned
                      </Badge>
                    ) : null}
                    <span className="font-bold text-slate-900">{t.title}</span>
                  </div>
                </td>
                <td className="px-6 py-4 text-slate-600">{t.author_name || t.author_id}</td>
                <td className="px-6 py-4">{t.reply_count ?? 0}</td>
                <td className="px-6 py-4">{t.views_count ?? 0}</td>
              </tr>
            ))}
            {topics.length === 0 && (
              <tr>
                <td colSpan={4} className="p-12 text-center text-slate-400 italic">
                  No topics yet.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

