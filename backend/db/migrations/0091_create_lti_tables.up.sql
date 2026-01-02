CREATE TABLE lti_tools (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    
    -- LTI Configuration
    client_id VARCHAR(255) NOT NULL, -- The Client ID we provide to the Tool
    initiate_login_url TEXT NOT NULL,
    redirection_uris TEXT[] NOT NULL,
    public_jwks_url TEXT, -- Or public_key text. For LTI 1.3 usually JWKS URL.
    
    -- Deployment Info (LTI 1.3)
    deployment_id UUID DEFAULT uuid_generate_v4(),
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- We might want deployments per course if the tool is installed multiple times with diff settings
-- but usually LTI tools are installed at Tenant level and available to courses.
-- We can add lti_context_deployments later if we heavily use Deep Linking per course.
-- For now, Tenant-level installation is sufficient.

CREATE INDEX idx_lti_tools_tenant_id ON lti_tools(tenant_id);
CREATE INDEX idx_lti_tools_client_id ON lti_tools(client_id);
