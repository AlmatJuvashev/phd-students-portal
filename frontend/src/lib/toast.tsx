import React from 'react'

type Toast = { id: number, text: string }
const Ctx = React.createContext<{toasts:Toast[], push:(t:string)=>void}>({toasts:[], push: ()=>{} })

export function ToastProvider({ children }:{children:React.ReactNode}) {
  const [toasts, set] = React.useState<Toast[]>([])
  const push = (text:string) => set(t => [...t, { id: Date.now()+Math.random(), text }])
  return <Ctx.Provider value={{toasts, push}}>
    {children}
    <div className="fixed bottom-4 right-4 space-y-2">
      {toasts.map(t => <div key={t.id} className="rounded-md border bg-white px-3 py-2 shadow">{t.text}</div>)}
    </div>
  </Ctx.Provider>
}

export function useToast(){ return React.useContext(Ctx) }
