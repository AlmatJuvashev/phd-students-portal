import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Plus, Search, Filter, MoreVertical, Edit2, Trash2, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { getQuestions, deleteQuestion, createQuestion } from './api';
import { Question } from './types';

export const QuestionsPage: React.FC = () => {
  const { bankId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [search, setSearch] = useState('');
  const [filterType, setFilterType] = useState<string>('all');

  const { data: questions = [], isLoading } = useQuery({
    queryKey: ['questions', bankId],
    queryFn: () => getQuestions(bankId!),
    enabled: !!bankId
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteQuestion(bankId!, id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['questions', bankId] });
    }
  });

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this question?')) {
      deleteMutation.mutate(id);
    }
  };

  const handleCreate = () => {
      // Navigate to editor with 'new' ID
      navigate(`/admin/item-banks/${bankId}/questions/new`);
  };

  const filtered = questions.filter(q => {
    const matchesSearch = q.stem.toLowerCase().includes(search.toLowerCase());
    const matchesType = filterType === 'all' || q.type === filterType;
    return matchesSearch && matchesType;
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="h-full flex flex-col space-y-6 animate-in fade-in duration-500 p-8">
       <div className="flex justify-between items-end flex-shrink-0">
          <div>
            <Button variant="ghost" size="sm" onClick={() => navigate('/admin/item-banks')} className="mb-2 -ml-2 text-slate-500">
                <ArrowLeft size={16} className="mr-1" /> Back to Banks
            </Button>
             <h1 className="text-2xl font-black text-slate-900 tracking-tight">Questions</h1>
             <p className="text-slate-500 text-sm mt-1">Manage items in this bank.</p>
          </div>
          <Button onClick={handleCreate}>
            <Plus className="mr-2 h-4 w-4" /> Add Question
          </Button>
       </div>

       {/* Toolbar */}
       <div className="bg-white p-2 rounded-2xl border border-slate-200 shadow-sm flex flex-col sm:flex-row gap-4 flex-shrink-0">
          <div className="relative flex-1">
             <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
             <input 
               value={search}
               onChange={(e) => setSearch(e.target.value)}
               placeholder="Search questions..." 
               className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus:ring-0 text-sm outline-none"
             />
          </div>
          <div className="h-10 w-px bg-slate-200 hidden sm:block" />
          <div className="flex gap-1 bg-slate-100 p-1 rounded-xl">
             {['all', 'multi_select', 'short_answer', 'true_false'].map(f => (
               <button
                 key={f}
                 onClick={() => setFilterType(f)}
                 className={`px-4 py-1.5 rounded-lg text-xs font-bold capitalize transition-all ${filterType === f ? "bg-white text-slate-900 shadow-sm" : "text-slate-500 hover:text-slate-700"}`}
               >
                 {f.replace('_', ' ')}
               </button>
             ))}
          </div>
       </div>

       {/* List */}
       <div className="flex-1 bg-white border border-slate-200 rounded-3xl overflow-hidden shadow-sm flex flex-col">
          <div className="overflow-y-auto flex-1">
             <table className="w-full text-sm text-left">
                <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase sticky top-0 z-10">
                   <tr>
                      <th className="px-6 py-4">Question Stem</th>
                      <th className="px-6 py-4">Type</th>
                      <th className="px-6 py-4">Points</th>
                      <th className="px-6 py-4 text-right">Actions</th>
                   </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                   {filtered.map(q => (
                      <tr key={q.id} className="hover:bg-slate-50 transition-colors group">
                         <td className="px-6 py-4 font-medium text-slate-900 line-clamp-2 max-w-md">
                            {q.stem}
                         </td>
                         <td className="px-6 py-4">
                            <Badge variant="secondary" className="uppercase text-[10px]">{q.type.replace('_', ' ')}</Badge>
                         </td>
                         <td className="px-6 py-4 text-slate-500">
                            {q.points_default}
                         </td>
                         <td className="px-6 py-4 text-right">
                            <div className="flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                               <button 
                                 onClick={() => navigate(`/admin/item-banks/${bankId}/questions/${q.id}`)}
                                 className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                               >
                                  <Edit2 size={16} />
                               </button>
                               <button 
                                 onClick={() => handleDelete(q.id)}
                                 className="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                               >
                                  <Trash2 size={16} />
                               </button>
                            </div>
                         </td>
                      </tr>
                   ))}
                   {filtered.length === 0 && (
                      <tr>
                         <td colSpan={4} className="px-6 py-12 text-center text-slate-400 italic">No questions found.</td>
                      </tr>
                   )}
                </tbody>
             </table>
          </div>
       </div>
    </div>
  );
};
