import React, { useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { DropdownMenu, DropdownItem } from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { ChevronDown, RefreshCw } from 'lucide-react'
import { cn } from '@/lib/utils'

export const RoleSwitcher: React.FC = () => {
  const { user, switchRole } = useAuth()
  const [switching, setSwitching] = useState(false)

  if (!user || !user.available_roles || user.available_roles.length <= 1) {
    return null
  }

  const handleRoleSwitch = async (role: string) => {
    if (role === user.active_role) return
    
    setSwitching(true)
    try {
      await switchRole(role)
    } finally {
      setSwitching(false)
    }
  }

  // Helper to format role name nicely
  const formatRole = (r: string) => r.replace('_', ' ').replace(/\b\w/g, c => c.toUpperCase())

  return (
    <DropdownMenu
        trigger={
            <Button variant="outline" size="sm" className="ml-2 gap-2" disabled={switching}>
                {switching ? <RefreshCw className="h-4 w-4 animate-spin" /> : null}
                {user.active_role ? formatRole(user.active_role) : 'Select Role'}
                <ChevronDown className="h-4 w-4 opacity-50" />
            </Button>
        }
    >
        {user.available_roles.map((role) => (
            <DropdownItem 
                key={role} 
                onClick={() => handleRoleSwitch(role)}
            >
                <span className={cn(role === user.active_role && "bg-accent font-medium block w-full")}>
                  {formatRole(role)}
                  {role === user.active_role && " (Active)"}
                </span>
            </DropdownItem>
        ))}
    </DropdownMenu>
  )
}
