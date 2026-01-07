import React, { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { Loader2 } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useToast } from '@/components/ui/use-toast';
import { submitProposal } from './api';

interface CreateProposalModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

type FormValues = {
  title: string;
  type: string;
  description: string;
  target_id?: string;
};

export const CreateProposalModal: React.FC<CreateProposalModalProps> = ({ open, onOpenChange }) => {
  const { toast } = useToast();
  const qc = useQueryClient();
  const { register, handleSubmit, reset, setValue } = useForm<FormValues>();
  
  const createMutation = useMutation({
    mutationFn: (data: FormValues) => submitProposal({ ...data, status: 'pending', current_step: 1 }),
    onSuccess: () => {
      toast({ title: 'Proposal Created', description: 'Your proposal has been submitted for review.' });
      qc.invalidateQueries({ queryKey: ['governance', 'proposals'] });
      reset();
      onOpenChange(false);
    },
    onError: (err: any) => {
      toast({ title: 'Error', description: err.response?.data?.error || 'Failed to create proposal', variant: 'destructive' });
    }
  });

  const onSubmit = (data: FormValues) => {
    createMutation.mutate(data);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Create New Proposal</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>
            <Input id="title" {...register('title', { required: true })} placeholder="e.g. Update Research Methodology Course" />
          </div>

          <div className="space-y-2">
            <Label htmlFor="type">Type</Label>
            <Select onValueChange={(v) => setValue('type', v)}>
               <SelectTrigger>
                  <SelectValue placeholder="Select type..." />
               </SelectTrigger>
               <SelectContent>
                  <SelectItem value="curriculum_change">Curriculum Change</SelectItem>
                  <SelectItem value="grade_change">Grade Change</SelectItem>
                  <SelectItem value="policy_update">Policy Update</SelectItem>
                  <SelectItem value="other">Other</SelectItem>
               </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
             <Label htmlFor="description">Description & Justification</Label>
             <Textarea 
                id="description" 
                {...register('description', { required: true })} 
                placeholder="Explain what needs to change and why..."
                className="min-h-[100px]"
             />
          </div>

          <div className="space-y-2">
             <Label htmlFor="target_id">Target ID (Optional)</Label>
             <Input id="target_id" {...register('target_id')} placeholder="e.g. Course ID or Student ID" />
          </div>

          <div className="pt-4 flex justify-end gap-3">
             <Button type="button" variant="ghost" onClick={() => onOpenChange(false)}>Cancel</Button>
             <Button type="submit" disabled={createMutation.isPending} className="bg-indigo-600 hover:bg-indigo-700">
                {createMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Submit Proposal
             </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
