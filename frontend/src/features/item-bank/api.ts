import { api } from '@/api/client';
import { CreateBankRequest, CreateQuestionRequest, Question, QuestionBank } from './types';

export const listBanks = () => api.get<QuestionBank[]>('/item-banks/banks');
export const createBank = (data: CreateBankRequest) => api.post<QuestionBank>('/item-banks/banks', data);

export const listQuestions = (bankId: string) => api.get<Question[]>(`/item-banks/banks/${bankId}/items`);
export const createQuestion = (bankId: string, data: CreateQuestionRequest) =>
  api.post<Question>(`/item-banks/banks/${bankId}/items`, data);

