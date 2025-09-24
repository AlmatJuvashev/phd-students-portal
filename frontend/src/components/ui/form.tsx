import * as React from 'react'
import { FormProvider, type UseFormReturn } from 'react-hook-form'

export function Form({ children, ...props }: React.ComponentProps<typeof FormProvider>) {
  return <FormProvider {...props}>{children}</FormProvider>
}
export function FormField({ children }: { children: React.ReactNode }) { return <div className="space-y-1">{children}</div> }
export function FormItem({ children }: { children: React.ReactNode }) { return <div className="space-y-1">{children}</div> }
export function FormLabel(props: React.HTMLAttributes<HTMLLabelElement>) { return <label className="text-sm font-medium" {...props} /> }
export function FormMessage({ children }: { children?: React.ReactNode }) { return <p className="text-xs text-rose-600">{children}</p> }
