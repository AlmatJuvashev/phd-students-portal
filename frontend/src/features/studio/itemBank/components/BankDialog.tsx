import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Save, FolderOpen } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Bank } from '../types';

interface BankDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (bank: Partial<Bank>) => void;
  initialBank?: Bank;
}

export const BankDialog: React.FC<BankDialogProps> = ({ isOpen, onClose, onSave, initialBank }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  useEffect(() => {
    if (isOpen) {
      setTitle(initialBank?.title || '');
      setDescription(initialBank?.description || '');
    }
  }, [isOpen, initialBank]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({ 
      ...(initialBank?.id ? { id: initialBank.id } : {}),
      title, 
      description 
    });
    onClose();
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="bg-white w-full max-w-md rounded-2xl shadow-2xl overflow-hidden"
          >
            <div className="flex justify-between items-center p-4 border-b border-slate-100">
              <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
                <FolderOpen className="text-indigo-600" size={20} />
                {initialBank ? 'Edit Bank' : 'Create Question Bank'}
              </h2>
              <button onClick={onClose} className="p-2 hover:bg-slate-100 rounded-full text-slate-400">
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-500 uppercase">Bank Title</label>
                <Input 
                  value={title} 
                  onChange={(e: any) => setTitle(e.target.value)} 
                  placeholder="e.g. General Surgery 2025" 
                  autoFocus
                />
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-500 uppercase">Description</label>
                <textarea 
                  className="w-full p-3 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-100 outline-none resize-none"
                  rows={3}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Optional description of contents..."
                />
              </div>

              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" variant="ghost" onClick={onClose}>Cancel</Button>
                <Button type="submit" disabled={!title.trim()}>
                    <Save className="mr-2 h-4 w-4" /> Save Bank
                </Button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
};
