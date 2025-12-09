-- Reset E1_apply_omid node for test user tu6260
-- This script resets the node state to allow re-completion with the new upload requirement

-- Find the user
DO $$
DECLARE
    v_user_id uuid;
    v_playbook_version_id uuid;
    v_node_instance_id uuid;
BEGIN
    -- Get user ID for tu6260
    SELECT id INTO v_user_id FROM users WHERE username = 'tu6260';
    
    IF v_user_id IS NULL THEN
        RAISE NOTICE 'User tu6260 not found';
        RETURN;
    END IF;
    
    RAISE NOTICE 'Found user: %', v_user_id;
    
    -- Get the current playbook version
    SELECT id INTO v_playbook_version_id 
    FROM playbook_versions 
    ORDER BY created_at DESC 
    LIMIT 1;
    
    IF v_playbook_version_id IS NULL THEN
        RAISE NOTICE 'No playbook version found';
        RETURN;
    END IF;
    
    RAISE NOTICE 'Using playbook version: %', v_playbook_version_id;
    
    -- Get the node instance for E1_apply_omid
    SELECT id INTO v_node_instance_id
    FROM node_instances
    WHERE user_id = v_user_id
      AND playbook_version_id = v_playbook_version_id
      AND node_id = 'E1_apply_omid';
    
    IF v_node_instance_id IS NULL THEN
        RAISE NOTICE 'Node instance E1_apply_omid not found for user';
        RETURN;
    END IF;
    
    RAISE NOTICE 'Found node instance: %', v_node_instance_id;
    
    -- Delete attachments for this node (if any)
    DELETE FROM node_instance_slot_attachments
    WHERE slot_id IN (
        SELECT id FROM node_instance_slots WHERE node_instance_id = v_node_instance_id
    );
    
    RAISE NOTICE 'Deleted slot attachments';
    
    -- Delete slots for this node
    DELETE FROM node_instance_slots WHERE node_instance_id = v_node_instance_id;
    
    RAISE NOTICE 'Deleted slots';
    
    -- Delete form revisions
    DELETE FROM node_instance_form_revisions WHERE node_instance_id = v_node_instance_id;
    
    RAISE NOTICE 'Deleted form revisions';
    
    -- Delete outcomes
    DELETE FROM node_outcomes WHERE node_instance_id = v_node_instance_id;
    
    RAISE NOTICE 'Deleted outcomes';
    
    -- Delete events
    DELETE FROM node_events WHERE node_instance_id = v_node_instance_id;
    
    RAISE NOTICE 'Deleted events';
    
    -- Reset the node instance state
    UPDATE node_instances
    SET state = 'active',
        submitted_at = NULL,
        updated_at = now(),
        current_rev = 0
    WHERE id = v_node_instance_id;
    
    RAISE NOTICE 'Reset node instance state to active';
    RAISE NOTICE 'Successfully reset E1_apply_omid node for user tu6260';
    
END $$;
