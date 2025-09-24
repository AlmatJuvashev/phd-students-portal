import * as React from 'react'
export function Accordion({ children }: { children: React.ReactNode }) { return <div className="divide-y rounded-xl border">{children}</div> }
export function AccordionItem({ header, children }: { header: React.ReactNode, children: React.ReactNode }) {
  const [open, setOpen] = React.useState(false)
  return (
    <div>
      <button onClick={()=>setOpen(o=>!o)} className="w-full text-left p-3 font-medium bg-muted">{header}</button>
      <div className={`${open ? 'block' : 'hidden'} p-3`}>{children}</div>
    </div>
  )
}
