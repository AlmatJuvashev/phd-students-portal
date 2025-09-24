import * as React from 'react'
export function DropdownMenu({ trigger, children }: { trigger: React.ReactNode, children: React.ReactNode }) {
  const [open, setOpen] = React.useState(false)
  return (
    <div className="relative inline-block">
      <div onClick={()=>setOpen(o=>!o)}>{trigger}</div>
      {open && <div className="absolute right-0 mt-2 min-w-[12rem] rounded-xl border bg-white p-2 shadow">{children}</div>}
    </div>
  )
}
export function DropdownItem({ onClick, children }: { onClick?: ()=>void, children: React.ReactNode }) {
  return <button onClick={onClick} className="block w-full text-left rounded-md px-3 py-1.5 text-sm hover:bg-muted">{children}</button>
}
