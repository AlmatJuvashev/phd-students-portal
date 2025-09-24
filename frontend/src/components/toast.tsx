import React, { createContext, useContext, useState } from 'react'

type Toast = { id:number, title?:string, description?:string }
type Ctx = { toasts: Toast[], push: (t: Omit<Toast,'id'>)=>void, remove:(id:number)=>void }
const Ctx = createContext<Ctx>({ toasts:[], push: ()=>{}, remove: ()=>{} })

export function ToastProvider({ children }:{children:React.ReactNode}) {
  const [toasts, setToasts] = useState<Toast[]>([])
  const push = (t: Omit<Toast,'id'>) => setToasts(ts=>[...ts, { id: Date.now(), ...t }])
  const remove = (id:number) => setToasts(ts=>ts.filter(x=>x.id!==id))
  return <Ctx.Provider value={{ toasts, push, remove }}>
    {children}
    <div className="fixed bottom-4 right-4 space-y-2">
      {toasts.map(t => (
        <div key={t.id} className="rounded-xl border bg-white p-3 shadow">
          <div className="font-medium">{t.title}</div>
          <div className="text-sm text-gray-600">{t.description}</div>
          <button className="text-xs underline mt-1" onClick={()=>remove(t.id)}>Close</button>
        </div>
      ))}
    </div>
  </Ctx.Provider>
}
export function useToast(){ return useContext(Ctx) }
