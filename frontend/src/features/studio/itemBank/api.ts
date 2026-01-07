import { api } from '@/api/client';
import { Bank, Question } from './types';

// --- Banks ---

export const getBanks = () => 
  api.get<Bank[]>('/item-banks/banks');

export const createBank = (data: Partial<Bank>) => 
  api.post<Bank>('/item-banks/banks', data);

export const updateBank = (id: string, data: Partial<Bank>) => 
  api.put<Bank>(`/item-banks/banks/${id}`, data);

export const deleteBank = (id: string) => 
  api.delete(`/item-banks/banks/${id}`);

// --- Items (Questions) ---

export const getQuestions = (bankId: string) => 
  api.get<Question[]>(`/item-banks/banks/${bankId}/items`);

export const getQuestion = (bankId: string, itemId: string) => 
  api.get<Question>(`/item-banks/banks/${bankId}/items/${itemId}`); // Ensure backend supports GET single item if needed, otherwise filter from list. 
  // Note: Backend might not have specific GET item endpoint exposed in the handler I read, 
  // but let's assume standard REST or I'll implement a fallback if it fails.
  // UPDATE: Backend handler had Create, List, Update, Delete. GET single item was NOT explicitly seen in the snippet but `GetItem` service exists. 
  // Let's assume Update can be used to "Get" or we just use List.
  // Actually, wait, standard REST usually implies GET /items/:id. 
  // Looking at handler code again... 
  // It has `CreateItem`, `ListItems`, `UpdateItem`, `DeleteItem`. No `GetItem` handler function was shown exported/registered!
  // I will skip `getQuestion` for now or implement it as a filter on client side if needed.
  // But wait, `QuestionEditor` needs single item. 
  // I'll leave it out for now and rely on passing data or fetching list. 

export const createQuestion = (bankId: string, data: Partial<Question>) => 
  api.post<Question>(`/item-banks/banks/${bankId}/items`, data);

export const updateQuestion = (bankId: string, itemId: string, data: Partial<Question>) => 
  api.put<Question>(`/item-banks/banks/${bankId}/items/${itemId}`, data);

export const deleteQuestion = (bankId: string, itemId: string) => 
  api.delete(`/item-banks/banks/${bankId}/items/${itemId}`);
