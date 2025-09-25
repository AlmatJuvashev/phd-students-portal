CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('superadmin','admin','student','advisor','chair');
CREATE TYPE step_status AS ENUM ('todo','in_progress','submitted','needs_changes','done');
CREATE TYPE doc_kind AS ENUM ('dissertation','publication_list','review','bioethics','ncste_norm','ncste_pub','rector_app','council_packet','video','order','other');

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  username text UNIQUE NOT NULL,
  email text UNIQUE NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  role user_role NOT NULL,
  password_hash text NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Password reset tokens (one active per email)
CREATE TABLE password_reset_tokens (
  email text PRIMARY KEY,
  token text NOT NULL,
  expires_at timestamptz NOT NULL
);

-- Minimal checklist scaffolding (extend later)
CREATE TABLE checklist_modules (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  code text UNIQUE NOT NULL,
  title text NOT NULL,
  description text,
  sort_order int NOT NULL
);

CREATE TABLE checklist_steps (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  module_id uuid NOT NULL REFERENCES checklist_modules(id) ON DELETE CASCADE,
  code text UNIQUE NOT NULL,
  title text NOT NULL,
  description text,
  requires_upload boolean NOT NULL DEFAULT false,
  form_schema jsonb,
  sort_order int NOT NULL
);

CREATE TABLE student_steps (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  step_id uuid NOT NULL REFERENCES checklist_steps(id) ON DELETE CASCADE,
  status step_status NOT NULL DEFAULT 'todo',
  data jsonb,
  due_at timestamptz,
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (user_id, step_id)
);

-- Documents & versions (simplified for starter)
CREATE TABLE documents (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind doc_kind NOT NULL,
  title text NOT NULL,
  current_version_id uuid,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE document_versions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  document_id uuid NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  storage_path text NOT NULL,
  mime_type text NOT NULL,
  size_bytes int NOT NULL,
  uploaded_by uuid NOT NULL REFERENCES users(id),
  note text,
  created_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE documents
  ADD CONSTRAINT fk_current_version
  FOREIGN KEY (current_version_id) REFERENCES document_versions(id);

-- Comments system
CREATE TABLE comments (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  document_id uuid NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  content text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
