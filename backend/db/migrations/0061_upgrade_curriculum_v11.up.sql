-- Upgrade Programs
ALTER TABLE programs ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE programs SET tenant_id = '00000000-0000-0000-0000-000000000001' WHERE tenant_id IS NULL;
ALTER TABLE programs ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE programs ADD COLUMN IF NOT EXISTS title jsonb;
ALTER TABLE programs ADD COLUMN IF NOT EXISTS description jsonb;
ALTER TABLE programs ADD COLUMN IF NOT EXISTS credits int DEFAULT 0;
ALTER TABLE programs ADD COLUMN IF NOT EXISTS duration_months int DEFAULT 36;

-- Create Courses
CREATE TABLE IF NOT EXISTS courses (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id uuid NOT NULL REFERENCES tenants(id),
  program_id uuid REFERENCES programs(id) ON DELETE SET NULL,
  code text,
  title jsonb NOT NULL,
  description jsonb,
  credits int DEFAULT 0,
  workload_hours int DEFAULT 0,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Journey Maps (Playbooks linked to Programs)
CREATE TABLE IF NOT EXISTS journey_maps (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  program_id uuid NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
  title jsonb,
  version text,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now()
);

-- Journey Node Definitions (The "Playbook Payload" dismantled)
CREATE TABLE IF NOT EXISTS journey_node_definitions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  journey_map_id uuid NOT NULL REFERENCES journey_maps(id) ON DELETE CASCADE,
  parent_node_id uuid REFERENCES journey_node_definitions(id),
  slug text NOT NULL, -- e.g. "VI_attestation_file"
  type text NOT NULL, -- "form", "task", "gateway"
  title jsonb,
  description jsonb,
  module_key text, -- "I", "II"
  coordinates jsonb DEFAULT '{"x": 0, "y": 0}'::jsonb,
  config jsonb DEFAULT '{}'::jsonb, -- dynamic config
  prerequisites text[], -- slugs
  created_at timestamptz DEFAULT now(),
  UNIQUE(journey_map_id, slug)
);

-- Cohorts
CREATE TABLE IF NOT EXISTS cohorts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  program_id uuid NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
  name text NOT NULL,
  start_date date,
  end_date date,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now()
);

-- Fix for existing cohorts table from migration 0030
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='cohorts' AND column_name='program_id') THEN
        ALTER TABLE cohorts ADD COLUMN program_id uuid REFERENCES programs(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_programs_tenant ON programs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_courses_tenant ON courses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_courses_program ON courses(program_id);
CREATE INDEX IF NOT EXISTS idx_journey_maps_program ON journey_maps(program_id);
CREATE INDEX IF NOT EXISTS idx_journey_nodes_map ON journey_node_definitions(journey_map_id);
CREATE INDEX IF NOT EXISTS idx_cohorts_program ON cohorts(program_id);
