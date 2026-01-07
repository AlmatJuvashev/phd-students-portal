import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { MoreVertical, FolderOpen, Plus, FileText, Clock, BarChart3, CheckCircle2, AlertTriangle, Layers, Edit2, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { getBanks, createBank, updateBank, deleteBank } from './api';
import { Bank } from './types';
import { BankDialog } from './components/BankDialog';

const StatCard = ({ label, value, icon: Icon, color }: any) => (
  <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex items-center gap-4">
    <div className={`p-3 rounded-lg ${color}`}>
      <Icon size={24} />
    </div>
    <div>
      <p className="text-xs font-bold text-slate-400 uppercase tracking-wider">{label}</p>
      <p className="text-2xl font-black text-slate-900">{value}</p>
    </div>
  </div>
);

export const ItemBanksPage: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: banks = [], isLoading } = useQuery({ queryKey: ['item-banks'], queryFn: getBanks });
  
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingBank, setEditingBank] = useState<Bank | undefined>(undefined);
  const [activeMenuId, setActiveMenuId] = useState<string | null>(null);

  const createMutation = useMutation({
    mutationFn: createBank,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-banks'] });
      setIsDialogOpen(false);
    }
  });

  const updateMutation = useMutation({
    mutationFn: (data: Partial<Bank>) => updateBank(data.id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-banks'] });
      setIsDialogOpen(false);
    }
  });

  const deleteMutation = useMutation({
    mutationFn: deleteBank,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-banks'] });
    }
  });

  const handleCreate = () => {
    setEditingBank(undefined);
    setIsDialogOpen(true);
  };

  const handleEdit = (bank: Bank, e: React.MouseEvent) => {
    e.stopPropagation();
    setEditingBank(bank);
    setIsDialogOpen(true);
    setActiveMenuId(null);
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (confirm('Are you sure you want to delete this bank? All items inside will be archived.')) {
      deleteMutation.mutate(id);
    }
    setActiveMenuId(null);
  };

  const handleSave = (bank: Partial<Bank>) => {
    if (bank.id) {
      updateMutation.mutate(bank);
    } else {
      createMutation.mutate(bank);
    }
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="p-8 max-w-7xl mx-auto space-y-8 overflow-y-auto h-full" onClick={() => setActiveMenuId(null)}>
      
      {/* Dashboard Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <StatCard label="Total Items" value={(banks || []).reduce((acc, b) => acc + (b.item_count || 0), 0).toLocaleString()} icon={Layers} color="bg-indigo-50 text-indigo-600" />
        <StatCard label="Active Banks" value={banks.length} icon={FolderOpen} color="bg-blue-50 text-blue-600" />
        <StatCard label="In Review" value="—" icon={AlertTriangle} color="bg-amber-50 text-amber-600" />
        <StatCard label="Published" value="—" icon={CheckCircle2} color="bg-emerald-50 text-emerald-600" />
      </div>

      <div className="flex justify-between items-end border-b border-slate-200 pb-4">
        <div>
          <h1 className="text-2xl font-black text-slate-900">Item Banks</h1>
          <p className="text-slate-500 mt-1">Manage collections of questions and assessment content.</p>
        </div>
        <Button onClick={handleCreate}>
            <Plus className="mr-2 h-4 w-4" /> Create Bank
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {banks.map(bank => (
          <div 
            key={bank.id} 
            className="group bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-lg hover:border-indigo-200 transition-all cursor-pointer flex flex-col relative"
            onClick={() => navigate('/admin/item-banks/' + bank.id)}
          >
            <div className="absolute top-0 left-0 w-full h-1.5 bg-gradient-to-r from-indigo-500 to-purple-500 opacity-0 group-hover:opacity-100 transition-opacity rounded-t-2xl" />
            
            <div className="flex justify-between items-start mb-4 relative">
               <div className="w-12 h-12 bg-slate-50 group-hover:bg-indigo-50 text-slate-400 group-hover:text-indigo-600 rounded-xl flex items-center justify-center transition-colors">
                 <FolderOpen size={24} />
               </div>
               <div className="relative">
                 <button 
                   onClick={(e) => { e.stopPropagation(); setActiveMenuId(activeMenuId === bank.id ? null : bank.id); }}
                   className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
                 >
                   <MoreVertical size={18} />
                 </button>
                 
                 {/* Context Menu */}
                 {activeMenuId === bank.id && (
                   <div className="absolute right-0 top-full mt-1 w-32 bg-white rounded-lg shadow-xl border border-slate-200 z-20 py-1 animate-in fade-in zoom-in-95 duration-100">
                     <button 
                       onClick={(e) => handleEdit(bank, e)}
                       className="w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 hover:text-indigo-600 flex items-center gap-2"
                     >
                       <Edit2 size={14} /> Edit
                     </button>
                     <button 
                       onClick={(e) => handleDelete(bank.id, e)}
                       className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center gap-2"
                     >
                       <Trash2 size={14} /> Delete
                     </button>
                   </div>
                 )}
               </div>
            </div>
            
            <h3 className="text-lg font-bold text-slate-900 mb-2 group-hover:text-indigo-600 transition-colors">{bank.title}</h3>
            <p className="text-sm text-slate-500 mb-6 flex-1 line-clamp-2">{bank.description}</p>
            
            <div className="pt-4 border-t border-slate-100 flex items-center justify-between text-xs text-slate-400 font-medium">
               <span className="flex items-center gap-1.5"><FileText size={14} /> {(bank.item_count || 0).toLocaleString()} items</span>
               <span className="flex items-center gap-1.5"><Clock size={14} /> {new Date(bank.updated_at || Date.now()).toLocaleDateString()}</span>
            </div>
          </div>
        ))}
        
        {/* Import CTA */}
        <button 
          onClick={() => navigate('/admin/item-banks/imports')}
          className="border-2 border-dashed border-slate-300 rounded-2xl p-6 flex flex-col items-center justify-center gap-3 text-slate-400 hover:text-indigo-600 hover:border-indigo-300 hover:bg-indigo-50/50 transition-all min-h-[240px]"
        >
           <div className="w-12 h-12 bg-slate-100 rounded-full flex items-center justify-center shadow-inner">
             <Plus size={24} />
           </div>
           <span className="font-bold">Import from QTI / CSV</span>
        </button>
      </div>

      <BankDialog 
        isOpen={isDialogOpen} 
        onClose={() => setIsDialogOpen(false)} 
        onSave={handleSave} 
        initialBank={editingBank}
      />
    </div>
  );
};
