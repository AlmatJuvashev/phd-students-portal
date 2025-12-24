#!/usr/bin/env python3
"""
Script to fix handler tests by adding tenant_id to node_instances INSERT statements
and Gin context.
"""

import re
import sys

def fix_node_instances_insert(content):
    """Add tenant_id to INSERT INTO node_instances statements."""
    
    # Pattern 1: INSERT without tenant_id
    pattern1 = r'INSERT INTO node_instances \(user_id,'
    replacement1 = r'INSERT INTO node_instances (tenant_id, user_id,'
    content = re.sub(pattern1, replacement1, content)
    
    # Pattern 2: VALUES without tenant_id (need to add $1 and shift other params)
    # This is trickier - we need to find VALUES clauses and add tenantID as first param
    
    return content

def add_tenant_id_variable(content):
    """Add tenantID variable declaration if not present."""
    
    # Check if tenantID is already declared
    if 'tenantID := "00000000-0000-0000-0000-000000000001"' in content:
        return content
    
    # Find a good place to add it - after SetupTestDB() typically
    pattern = r'(db, teardown := testutils\.SetupTestDB\(\)\n\tdefer teardown\(\))\n'
    replacement = r'\1\n\n\ttenantID := "00000000-0000-0000-0000-000000000001"\n'
    
    content = re.sub(pattern, replacement, content, count=1)
    
    return content

def add_tenant_id_to_context(content):
    """Add tenant_id to Gin context if not present."""
    
    # Pattern: c.Set("claims", ...) followed by c.Next() without tenant_id
    pattern = r'(c\.Set\("claims", [^)]+\))\n(\s+c\.Next\(\))'
    
    def replacement(match):
        claims_line = match.group(1)
        next_line = match.group(2)
        indent = next_line.split('c.Next')[0]
        
        # Check if tenant_id is already there
        if 'tenant_id' in claims_line:
            return match.group(0)
        
        return f'{claims_line}\n{indent}c.Set("tenant_id", tenantID)\n{next_line}'
    
    content = re.sub(pattern, replacement, content)
    
    return content

def main():
    if len(sys.argv) != 2:
        print("Usage: fix_test.py <test_file.go>")
        sys.exit(1)
    
    filename = sys.argv[1]
    
    with open(filename, 'r') as f:
        content = f.read()
    
    # Apply fixes
    content = add_tenant_id_variable(content)
    content = fix_node_instances_insert(content)
    content = add_tenant_id_to_context(content)
    
    with open(filename, 'w') as f:
        f.write(content)
    
    print(f"Fixed {filename}")

if __name__ == '__main__':
    main()
