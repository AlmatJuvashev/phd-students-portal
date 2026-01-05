import { api } from '@/api/client';
import { Forum, Post, Topic, TopicWithPosts } from './types';

export const listCourseForums = (courseOfferingId: string) =>
  api.get<Forum[]>(`/courses/${courseOfferingId}/forums`);

export const listTopics = (forumId: string, limit = 20, offset = 0) =>
  api.get<Topic[]>(`/forums/${forumId}/topics?limit=${limit}&offset=${offset}`);

export const createTopic = (forumId: string, data: Pick<Topic, 'title' | 'content'>) =>
  api.post<Topic>(`/forums/${forumId}/topics`, data);

export const getTopic = (topicId: string) =>
  api.get<TopicWithPosts>(`/topics/${topicId}`);

export const createPost = (topicId: string, data: Pick<Post, 'content'> & { parent_id?: string | null }) =>
  api.post<Post>(`/topics/${topicId}/posts`, data);

