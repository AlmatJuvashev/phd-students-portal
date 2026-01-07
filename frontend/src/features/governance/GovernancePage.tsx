import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { Plus, Filter, FileText, CheckCircle2, XCircle, Clock, ArrowRight } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { listProposals } from './api';
import { Proposal } from './api';
import { useNavigate } from 'react-router-dom';
import { CreateProposalModal } from './CreateProposalModal';

const STATUS_COLORS: Record<string, string> = {
  pending: 'bg-amber-50 text-amber-700 border-amber-200',
  approved: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  rejected: 'bg-red-50 text-red-700 border-red-200',
  implemented: 'bg-slate-100 text-slate-700 border-slate-200',
};

const STATUS_ICONS: Record<string, React.ReactNode> = {
  pending: <Clock className="w-3 h-3" />,
  approved: <CheckCircle2 className="w-3 h-3" />,
  rejected: <XCircle className="w-3 h-3" />,
  implemented: <FileText className="w-3 h-3" />,
};

export const GovernancePage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const [filter, setFilter] = useState<string>('all');
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  
  const { data: proposals, isLoading } = useQuery({
    queryKey: ['governance', 'proposals', filter],
    queryFn: () => listProposals(filter === 'all' ? undefined : filter),
  });

  return (
    <div className="max-w-6xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-black text-slate-900 tracking-tight">Governance Proposals</h1>
          <p className="text-slate-500 font-medium mt-1">Review and manage change requests.</p>
        </div>
        <Button onClick={() => setIsCreateOpen(true)} className="bg-indigo-600 hover:bg-indigo-700">
          <Plus className="mr-2 h-4 w-4" /> New Proposal
        </Button>
      </div>

      <CreateProposalModal open={isCreateOpen} onOpenChange={setIsCreateOpen} />

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <div className="lg:col-span-3 space-y-4">
          <div className="flex items-center gap-2 mb-4 overflow-x-auto pb-2">
            {['all', 'pending', 'approved', 'rejected'].map((s) => (
              <button
                key={s}
                onClick={() => setFilter(s)}
                className={`px-4 py-2 rounded-full text-xs font-bold uppercase tracking-wider transition-colors ${
                  filter === s
                    ? 'bg-slate-900 text-white'
                    : 'bg-white border border-slate-200 text-slate-500 hover:bg-slate-50'
                }`}
              >
                {s}
              </button>
            ))}
          </div>

          {isLoading ? (
            <div className="p-12 text-center text-slate-400">Loading proposals...</div>
          ) : (
            <div className="grid gap-4">
              {proposals?.map((p) => (
                <div
                  key={p.id}
                  className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-all cursor-pointer group"
                  onClick={() => navigate(`/admin/governance/${p.id}`)}
                >
                  <div className="flex justify-between items-start gap-4">
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <Badge variant="outline" className={STATUS_COLORS[p.status]}>
                           <span className="flex items-center gap-1.5">
                             {STATUS_ICONS[p.status]} {p.status}
                           </span>
                        </Badge>
                        <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">
                          {p.type.replace('_', ' ')}
                        </span>
                      </div>
                      <h3 className="text-lg font-bold text-slate-900 group-hover:text-indigo-600 transition-colors">
                        {p.title}
                      </h3>
                      <p className="text-sm text-slate-500 line-clamp-2">{p.description}</p>
                    </div>
                    <div className="text-slate-300 group-hover:text-indigo-400 transition-colors">
                      <ArrowRight size={20} />
                    </div>
                  </div>
                  <div className="mt-4 pt-4 border-t border-slate-50 flex items-center justify-between text-xs text-slate-400 font-mono">
                    <span>ID: {p.id.substring(0, 8)}</span>
                    <span>{new Date(p.created_at).toLocaleDateString()}</span>
                  </div>
                </div>
              ))}
              {proposals?.length === 0 && (
                <div className="text-center p-12 bg-white rounded-2xl border border-dashed border-slate-300 text-slate-400">
                  No proposals found matching this filter.
                </div>
              )}
            </div>
          )}
        </div>

        <div className="space-y-6">
          <div className="bg-slate-900 text-white p-6 rounded-3xl shadow-lg">
            <h3 className="font-bold text-lg mb-2">Governance Stats</h3>
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-slate-400 text-sm">Pending Review</span>
                <span className="text-2xl font-black text-amber-400">
                  {proposals?.filter(p => p.status === 'pending').length || 0}
                </span>
              </div>
              <div className="w-full bg-white/10 h-px" />
              <div className="flex justify-between items-center">
                <span className="text-slate-400 text-sm">Total Approved</span>
                <span className="text-xl font-bold text-emerald-400">
                  {proposals?.filter(p => p.status === 'approved').length || 0}
                </span>
              </div>
            </div>
          </div>
          
          <div className="bg-white border border-slate-200 p-6 rounded-3xl shadow-sm">
             <div className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-4">Quick Actions</div>
             <div className="space-y-2">
               <Button variant="outline" className="w-full justify-start text-left h-auto py-3">
                 <FileText className="mr-2 h-4 w-4 text-slate-400" />
                 <div className="flex flex-col items-start">
                   <span className="font-bold text-slate-700">Download Report</span>
                   <span className="text-[10px] text-slate-400">PDF Summary of all changes</span>
                 </div>
               </Button>
             </div>
          </div>
        </div>
      </div>
    </div>
  );
};
