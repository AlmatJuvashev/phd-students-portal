import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Check, X, User, MessageSquare, Clock, FileText, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
import { getProposal, listReviews, reviewProposal } from './api';
import { useToast } from '@/components/ui/use-toast';
import { useAuth } from '@/contexts/AuthContext';

export const ProposalDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation('common');
  const { toast } = useToast();
  const qc = useQueryClient();
  const { user } = useAuth();
  
  const [comment, setComment] = useState('');

  const { data: proposal, isLoading: pLoading } = useQuery({
    queryKey: ['governance', 'proposal', id],
    queryFn: () => getProposal(id!),
    enabled: !!id,
  });

  const { data: reviews, isLoading: rLoading } = useQuery({
    queryKey: ['governance', 'reviews', id],
    queryFn: () => listReviews(id!),
    enabled: !!id,
  });

  const reviewMutation = useMutation({
    mutationFn: (status: string) => reviewProposal(id!, { status, comment }),
    onSuccess: () => {
      toast({ title: 'Review Submitted', description: 'Your decision has been recorded.' });
      qc.invalidateQueries({ queryKey: ['governance', 'proposal', id] });
      qc.invalidateQueries({ queryKey: ['governance', 'reviews', id] });
      setComment('');
    },
     onError: (err: any) => {
        toast({ title: 'Error', description: err.response?.data?.error || 'Failed to submit review', variant: 'destructive' });
     }
  });

  if (pLoading || rLoading || !proposal) {
     return <div className="p-12 flex justify-center"><Loader2 className="animate-spin text-slate-400" /></div>;
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in duration-500">
      <button
        onClick={() => navigate('/admin/governance')}
        className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors"
      >
        <ArrowLeft size={16} /> Back to Proposals
      </button>

      <div className="bg-white border border-slate-200 rounded-3xl overflow-hidden shadow-sm">
        <div className="p-8 border-b border-slate-100">
           <div className="flex justify-between items-start gap-4">
              <div className="space-y-2">
                 <div className="flex items-center gap-3">
                    <Badge variant="secondary" className="bg-slate-100 text-slate-600">
                       {proposal.type.toUpperCase().replace('_', ' ')}
                    </Badge>
                    <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">
                       {new Date(proposal.created_at).toLocaleString()}
                    </span>
                 </div>
                 <h1 className="text-3xl font-black text-slate-900 tracking-tight">{proposal.title}</h1>
              </div>
              <Badge 
                 className={`text-sm px-3 py-1 ${
                    proposal.status === 'approved' ? 'bg-emerald-100 text-emerald-700' :
                    proposal.status === 'rejected' ? 'bg-red-100 text-red-700' :
                    proposal.status === 'pending' ? 'bg-amber-100 text-amber-700' : 
                    'bg-slate-100 text-slate-700'
                 }`}
              >
                 {proposal.status.toUpperCase()}
              </Badge>
           </div>
        </div>
        
        <div className="p-8 grid grid-cols-1 lg:grid-cols-3 gap-8">
           <div className="lg:col-span-2 space-y-8">
              <div>
                 <h3 className="font-bold text-slate-900 mb-2 flex items-center gap-2">
                    <FileText className="w-4 h-4 text-slate-400" /> Description
                 </h3>
                 <p className="text-slate-600 leading-relaxed whitespace-pre-wrap">{proposal.description}</p>
              </div>

              {proposal.data && (
                 <div>
                    <h3 className="font-bold text-slate-900 mb-2">Technical Data</h3>
                    <pre className="bg-slate-50 p-4 rounded-xl text-xs font-mono text-slate-600 overflow-x-auto border border-slate-200">
                       {JSON.stringify(proposal.data, null, 2)}
                    </pre>
                 </div>
              )}

              <div className="border-t border-slate-100 pt-8">
                 <h3 className="font-bold text-slate-900 mb-6 flex items-center gap-2">
                    <MessageSquare className="w-4 h-4 text-slate-400" /> Reviews & Comments
                 </h3>
                 <div className="space-y-4">
                    {reviews?.map(r => (
                       <div key={r.id} className="bg-slate-50 p-4 rounded-xl border border-slate-100">
                          <div className="flex justify-between items-center mb-2">
                             <div className="flex items-center gap-2">
                                <div className="w-6 h-6 rounded-full bg-indigo-100 flex items-center justify-center text-xs font-bold text-indigo-600">
                                   <User className="w-3 h-3" />
                                </div>
                                <span className="font-bold text-sm text-slate-700">Reviewer ({r.reviewer_id.substring(0,6)})</span>
                             </div>
                             <span className={`text-xs font-bold uppercase ${
                                r.status === 'approved' ? 'text-emerald-600' : 'text-red-600'
                             }`}>
                                {r.status}
                             </span>
                          </div>
                          <p className="text-sm text-slate-600">{r.comment}</p>
                          <div className="mt-2 text-[10px] text-slate-400 font-bold uppercase tracking-wider">
                             {new Date(r.created_at).toLocaleString()}
                          </div>
                       </div>
                    ))}
                    {reviews?.length === 0 && (
                       <div className="text-center text-slate-400 text-sm italic">No reviews yet.</div>
                    )}
                 </div>
              </div>
           </div>

           <div className="space-y-6">
              {proposal.status === 'pending' && (
                 <div className="bg-slate-50 border border-slate-200 rounded-2xl p-6 sticky top-6">
                    <h3 className="font-bold text-slate-900 mb-4">Submit Review</h3>
                    <div className="space-y-4">
                       <Textarea 
                          placeholder="Add a comment explaining your decision..."
                          className="min-h-[100px] bg-white"
                          value={comment}
                          onChange={e => setComment(e.target.value)}
                       />
                       <div className="grid grid-cols-2 gap-3">
                          <Button 
                             onClick={() => reviewMutation.mutate('approved')}
                             disabled={reviewMutation.isPending}
                             className="bg-emerald-600 hover:bg-emerald-700 text-white"
                          >
                             <Check className="mr-2 w-4 h-4" /> Approve
                          </Button>
                          <Button 
                             onClick={() => reviewMutation.mutate('rejected')}
                             disabled={reviewMutation.isPending}
                             variant="destructive"
                          >
                             <X className="mr-2 w-4 h-4" /> Reject
                          </Button>
                       </div>
                    </div>
                 </div>
              )}
              
              <div className="bg-white border border-slate-200 rounded-2xl p-6">
                 <h3 className="font-bold text-slate-900 mb-4 text-sm uppercase tracking-wider">Meta</h3>
                 <div className="space-y-3 text-sm">
                    <div className="flex justify-between">
                       <span className="text-slate-500">Requester</span>
                       <span className="font-mono text-slate-700">{proposal.requester_id.substring(0,8)}</span>
                    </div>
                    <div className="flex justify-between">
                       <span className="text-slate-500">Current Step</span>
                       <span className="font-bold text-slate-700">{proposal.current_step}</span>
                    </div>
                 </div>
              </div>
           </div>
        </div>
      </div>
    </div>
  );
};
