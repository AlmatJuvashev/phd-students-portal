export type Proposal = {
  id: string;
  tenant_id: string;
  requester_id: string;
  type: string;
  target_id: string;
  title: string;
  description: string;
  status: 'pending' | 'approved' | 'rejected' | 'implemented';
  data: any;
  current_step: number;
  created_at: string;
  updated_at: string;
};

export type ProposalReview = {
  id: string;
  proposal_id: string;
  reviewer_id: string;
  status: 'approved' | 'rejected';
  comment: string;
  created_at: string;
};

import { api } from '@/api/client';

export const listProposals = (status?: string) => 
  api.get<Proposal[]>(`/governance/proposals${status ? `?status=${status}` : ''}`);

export const getProposal = (id: string) => 
  api.get<Proposal>(`/governance/proposals/${id}`);

export const submitProposal = (data: Partial<Proposal>) => 
  api.post<Proposal>('/governance/proposals', data);

export const reviewProposal = (id: string, data: { status: string; comment: string }) => 
  api.post<{ message: string; status: string }>(`/governance/proposals/${id}/review`, data);

export const listReviews = (id: string) => 
  api.get<ProposalReview[]>(`/governance/proposals/${id}/reviews`);
