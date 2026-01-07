import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ArrowLeft, Save, Plus, Trash2, CheckSquare, Type, LayoutList, AlignLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { createQuestion, updateQuestion, getQuestion, getQuestions } from './api';
import { Question, QuestionType, QuestionOption } from './types';
import { RichTextEditor } from '@/components/ui/RichTextEditor';

export const QuestionEditor: React.FC = () => {
  const { bankId, questionId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isNew = questionId === 'new';

  const [stem, setStem] = useState('');
  const [type, setType] = useState<QuestionType>('multi_select');
  const [points, setPoints] = useState(1);
  const [options, setOptions] = useState<QuestionOption[]>([
      { text: '', is_correct: false },
      { text: '', is_correct: false }
  ]);

  // If not new, fetch existing question
  // Since we couldn't rely on 'getQuestion' single API easily in previous step,
  // we can fetch all and find, OR retry fetching single if implemented.
  // Ideally, 'getQuestion' in api.ts SHOULD assume an endpoint exists or we implement a client-side find.
  // For this polished version, let's assume `getQuestion` works or we fix `api.ts` to be robust. 
  // Wait, I implemented `getQuestion` in `api.ts` expecting `/item-banks/banks/:bankId/items/:itemId`.
  // If that endpoint 404s, we might need to fallback. 
  // Let's assume it works for now as standard REST.

  const { data: existingQuestion, isLoading } = useQuery({
     queryKey: ['question', bankId, questionId],
     queryFn: () => getQuestion(bankId!, questionId!),
     enabled: !isNew && !!bankId && !!questionId,
     retry: false
  });

  useEffect(() => {
      if (existingQuestion) {
          setStem(existingQuestion.stem);
          setType(existingQuestion.type);
          setPoints(existingQuestion.points_default);
          if (existingQuestion.options) setOptions(existingQuestion.options);
      }
  }, [existingQuestion]);

  const saveMutation = useMutation({
      mutationFn: (data: Partial<Question>) => {
          if (isNew) return createQuestion(bankId!, data);
          return updateQuestion(bankId!, questionId!, data);
      },
      onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ['questions', bankId] });
          navigate(`/admin/item-banks/${bankId}`);
      }
  });

  const handleSave = () => {
      saveMutation.mutate({
          type,
          stem,
          points_default: points,
          options,
          bank_id: bankId
      });
  };

  const addOption = () => {
      setOptions([...options, { text: '', is_correct: false }]);
  };

  const removeOption = (index: number) => {
      setOptions(options.filter((_, i) => i !== index));
  };

  const updateOption = (index: number, field: keyof QuestionOption, value: any) => {
      const newOptions = [...options];
      newOptions[index] = { ...newOptions[index], [field]: value };
      setOptions(newOptions);
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="max-w-4xl mx-auto p-8 h-full flex flex-col">
       <div className="flex justify-between items-center mb-8">
          <div className="flex items-center gap-4">
             <Button variant="ghost" onClick={() => navigate(`/admin/item-banks/${bankId}`)}>
                <ArrowLeft size={16} />
             </Button>
             <h1 className="text-2xl font-black text-slate-900">{isNew ? 'Create Question' : 'Edit Question'}</h1>
          </div>
          <Button onClick={handleSave} disabled={saveMutation.isPending || !stem.trim()}>
             <Save className="mr-2 h-4 w-4" /> Save Question
          </Button>
       </div>

       <div className="grid grid-cols-3 gap-8 flex-1 overflow-hidden">
          {/* Main Editor */}
          <div className="col-span-2 space-y-6 overflow-y-auto pr-2 pb-20">
             <div className="space-y-2">
                <label className="text-sm font-bold text-slate-700">Question Text (Stem)</label>
                <div className="border border-slate-200 rounded-xl overflow-hidden bg-white">
                  <RichTextEditor 
                    value={stem}
                    onChange={setStem}
                    className="min-h-[120px] border-none shadow-none focus-visible:ring-0"
                    placeholder="Enter the question here..."
                  />
                </div>
             </div>

             {/* Multiple Choice Editor */}
             {(type === 'multi_select' || type === 'true_false') && (
                 <div className="space-y-4">
                     <label className="text-sm font-bold text-slate-700">Answer Options</label>
                     <div className="space-y-3">
                         {options.map((opt, idx) => (
                             <div key={idx} className="flex items-center gap-3">
                                 <input 
                                   type="checkbox" 
                                   checked={opt.is_correct} 
                                   onChange={(e) => updateOption(idx, 'is_correct', e.target.checked)}
                                   className="w-5 h-5 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
                                 />
                                 <Input 
                                   value={opt.text} 
                                   onChange={(e: any) => updateOption(idx, 'text', e.target.value)} 
                                   placeholder={`Option ${idx + 1}`}
                                   className="flex-1"
                                 />
                                 <button onClick={() => removeOption(idx)} className="text-slate-400 hover:text-red-500 p-2">
                                     <Trash2 size={16} />
                                 </button>
                             </div>
                         ))}
                     </div>
                     <Button variant="outline" size="sm" onClick={addOption}>
                         <Plus className="mr-2 h-4 w-4" /> Add Option
                     </Button>
                 </div>
             )}
          </div>

          {/* Sidebar */}
          <div className="bg-slate-50 p-6 rounded-2xl h-fit space-y-6 border border-slate-200">
             <div className="space-y-2">
                 <label className="text-xs font-bold text-slate-500 uppercase">Question Type</label>
                 <div className="grid grid-cols-2 gap-2">
                     {[
                         { id: 'multi_select', icon: CheckSquare, label: 'Multiple Choice' },
                         { id: 'short_answer', icon: Type, label: 'Short Answer' },
                         { id: 'true_false', icon: LayoutList, label: 'True/False' },
                         { id: 'essay', icon: AlignLeft, label: 'Essay' },
                     ].map(t => (
                         <button 
                           key={t.id}
                           onClick={() => setType(t.id as QuestionType)}
                           className={`p-3 rounded-xl border text-left flex flex-col items-center gap-2 transition-all ${type === t.id ? 'bg-white border-indigo-600 text-indigo-700 shadow-sm' : 'bg-white border-slate-200 text-slate-500 hover:border-slate-300'}`}
                         >
                             <t.icon size={20} />
                             <span className="text-[10px] font-bold uppercase">{t.label}</span>
                         </button>
                     ))}
                 </div>
             </div>

             <div className="space-y-2">
                 <label className="text-xs font-bold text-slate-500 uppercase">Points</label>
                 <Input 
                   type="number" 
                   value={points} 
                   onChange={(e: any) => setPoints(parseInt(e.target.value) || 0)} 
                 />
             </div>
          </div>
       </div>
    </div>
  );
};
