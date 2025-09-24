import React from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { api } from '../api/client'
import { Input } from '../components/ui/input'
import { Button } from '../components/ui/button'
import { Card, CardHeader, CardTitle, CardContent } from '../components/ui/card'
import { useToast } from '../components/toast'

const Schema = z.object({
  first_name: z.string().min(1),
  last_name: z.string().min(1),
  email: z.string().email(),
  role: z.enum(['student','advisor','chair','admin']),
})
type Form = z.infer<typeof Schema>

export function AdminUsers() {
  const { register, handleSubmit, reset, formState:{errors, isSubmitting} } = useForm<Form>({ resolver: zodResolver(Schema), defaultValues: { role:'student' } })
  const [created, setCreated] = React.useState<{username:string,temp_password:string}|null>(null)
  const { push } = useToast()
  const onSubmit = async (data: Form) => {
    try {
      const res = await api('/admin/users', { method:'POST', body: JSON.stringify(data) })
      setCreated(res); reset()
      push({ title:'User created', description:'Credentials generated' })
    } catch (e:any) {
      push({ title:'Error', description:e.message })
    }
  }
  return (
    <div className="max-w-lg mx-auto mt-8 space-y-4">
      <h2 className="text-xl font-semibold">Create user</h2>
      <form className="space-y-3" onSubmit={handleSubmit(onSubmit)}>
        <Input placeholder="First name" {...register('first_name')} />
        {errors.first_name && <div className="text-xs text-rose-600">{errors.first_name.message}</div>}
        <Input placeholder="Last name" {...register('last_name')} />
        {errors.last_name && <div className="text-xs text-rose-600">{errors.last_name.message}</div>}
        <Input placeholder="Email" {...register('email')} />
        {errors.email && <div className="text-xs text-rose-600">{errors.email.message}</div>}
        <select className="w-full border p-2 rounded" {...register('role')}>
          <option value="student">student</option>
          <option value="advisor">advisor</option>
          <option value="chair">chair</option>
          <option value="admin">admin</option>
        </select>
        <Button className="w-full" disabled={isSubmitting}>Create</Button>
      </form>
      {created && (
        <Card>
          <CardHeader><CardTitle>Credentials</CardTitle></CardHeader>
          <CardContent>
            <div className="text-sm"><strong>Username:</strong> {created.username}</div>
            <div className="text-sm"><strong>Password:</strong> {created.temp_password}</div>
            <Button variant="outline" className="mt-2" onClick={()=>navigator.clipboard.writeText(`${created.username} • ${created.temp_password}`)}>Copy</Button>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
