export type ForumType = 'ANNOUNCEMENT' | 'QNA' | 'DISCUSSION';

export interface Forum {
  id: string;
  course_offering_id: string;
  title: string;
  description: string;
  type: ForumType;
  is_locked: boolean;
  created_at: string;
  updated_at: string;
}

export interface Topic {
  id: string;
  forum_id: string;
  author_id: string;
  title: string;
  content: string;
  is_pinned: boolean;
  is_locked: boolean;
  views_count: number;
  created_at: string;
  updated_at: string;
  author_name?: string;
  reply_count?: number;
  last_post_at?: string | null;
}

export interface Post {
  id: string;
  topic_id: string;
  author_id: string;
  parent_id?: string | null;
  content: string;
  created_at: string;
  updated_at: string;
  author_name?: string;
  author_role?: string;
}

export interface TopicWithPosts {
  topic: Topic;
  posts: Post[];
}

