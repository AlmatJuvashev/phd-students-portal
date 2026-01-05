import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { ArrowLeft, Loader2, Send } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { createPost, getTopic } from './api';

export const TopicDetailPage: React.FC = () => {
  const navigate = useNavigate();
  const qc = useQueryClient();
  const { courseOfferingId, topicId } = useParams<{ courseOfferingId: string; topicId: string }>();
  const [reply, setReply] = useState('');

  const topicQuery = useQuery({
    queryKey: ['forums', 'topic', topicId],
    queryFn: () => getTopic(topicId!),
    enabled: !!topicId,
  });

  const postMutation = useMutation({
    mutationFn: () => createPost(topicId!, { content: reply.trim() }),
    onSuccess: async () => {
      setReply('');
      await qc.invalidateQueries({ queryKey: ['forums', 'topic', topicId] });
    },
  });

  if (!courseOfferingId || !topicId) {
    return <div className="p-8 text-center text-slate-500">Topic not found.</div>;
  }

  if (topicQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  const data = topicQuery.data;
  if (!data?.topic) {
    return <div className="p-8 text-center text-slate-500">Topic not found.</div>;
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="space-y-4">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} /> Back
        </button>

        <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm">
          <div className="flex items-start justify-between gap-4">
            <div>
              <div className="flex items-center gap-2">
                {data.topic.is_pinned ? (
                  <Badge variant="secondary" className="uppercase text-[10px]">
                    Pinned
                  </Badge>
                ) : null}
                <h1 className="text-2xl font-black text-slate-900">{data.topic.title}</h1>
              </div>
              <div className="text-xs text-slate-500 mt-1">
                {data.topic.author_name || data.topic.author_id} • {new Date(data.topic.created_at).toLocaleString()}
              </div>
            </div>
            <Badge variant="outline" className="uppercase text-[10px]">
              {data.topic.views_count} views
            </Badge>
          </div>

          <div className="mt-4 text-sm text-slate-700 whitespace-pre-wrap">{data.topic.content}</div>
        </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm space-y-4">
        <div className="font-black text-slate-900">Reply</div>
        <Textarea value={reply} onChange={(e) => setReply(e.target.value)} placeholder="Write a reply..." />
        <div className="flex justify-end">
          <Button disabled={!reply.trim() || postMutation.isPending} onClick={() => postMutation.mutate()}>
            {postMutation.isPending ? <Loader2 className="animate-spin mr-2 h-4 w-4" /> : <Send className="mr-2 h-4 w-4" />}
            Send
          </Button>
        </div>
      </div>

      <div className="space-y-3">
        {(data.posts || []).map((p) => (
          <div key={p.id} className="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm">
            <div className="flex items-center justify-between gap-4">
              <div className="font-bold text-slate-900">{p.author_name || p.author_id}</div>
              <div className="text-xs text-slate-500">{new Date(p.created_at).toLocaleString()}</div>
            </div>
            <div className="mt-3 text-sm text-slate-700 whitespace-pre-wrap">{p.content}</div>
          </div>
        ))}
        {(data.posts || []).length === 0 && (
          <div className="text-center py-12 text-slate-400 italic">No replies yet.</div>
        )}
      </div>
    </div>
  );
};

