-- Multitenancy: Create tenants table
CREATE TABLE tenants (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  slug text UNIQUE NOT NULL,           -- e.g., 'kaznmu', 'knu'
  name text NOT NULL,                   -- 'Kazakh National Medical University'
  domain text,                          -- Optional custom domain
  logo_url text,                        -- Tenant logo
  settings jsonb DEFAULT '{}',          -- Tenant-specific config
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Add super-admin flag to users (global admin across all tenants)
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_superadmin boolean DEFAULT false;

-- Create default tenant (KazNMU) for existing data
INSERT INTO tenants (id, slug, name, is_active) 
VALUES ('00000000-0000-0000-0000-000000000001', 'kaznmu', 'Kazakh National Medical University', true);

-- Create indexes
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_domain ON tenants(domain) WHERE domain IS NOT NULL;
CREATE INDEX idx_tenants_active ON tenants(is_active);
