-- Migration: Sync node_instances to journey_states
-- This ensures all existing node progress is visible in the frontend journey map

-- Sync all node_instances to journey_states for Demo University tenant
INSERT INTO journey_states (tenant_id, user_id, node_id, state, updated_at)
SELECT 
    ni.tenant_id,
    ni.user_id,
    ni.node_id,
    ni.state,
    COALESCE(ni.updated_at, ni.submitted_at, ni.opened_at, NOW())
FROM node_instances ni
WHERE ni.tenant_id = 'dd000000-0000-0000-0000-d00000000001'
ON CONFLICT (user_id, node_id) 
DO UPDATE SET 
    state = EXCLUDED.state,
    tenant_id = EXCLUDED.tenant_id,
    updated_at = NOW();

-- Also sync for any other tenants that might have data
INSERT INTO journey_states (tenant_id, user_id, node_id, state, updated_at)
SELECT DISTINCT ON (ni.user_id, ni.node_id)
    ni.tenant_id,
    ni.user_id,
    ni.node_id,
    ni.state,
    COALESCE(ni.updated_at, ni.submitted_at, ni.opened_at, NOW())
FROM node_instances ni
WHERE ni.tenant_id != 'dd000000-0000-0000-0000-d00000000001'
ORDER BY ni.user_id, ni.node_id, ni.updated_at DESC NULLS LAST
ON CONFLICT (user_id, node_id) 
DO UPDATE SET 
    state = EXCLUDED.state,
    tenant_id = EXCLUDED.tenant_id,
    updated_at = NOW();
