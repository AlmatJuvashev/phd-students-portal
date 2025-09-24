import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../api/client'
import { Card, CardHeader, CardTitle, CardContent } from '../components/ui/card'

type Row = { id:string; title:string; kind:string; current_version_id?:string }

export function DocumentsList({ studentId }:{ studentId: string }) {
  const { data } = useQuery({ queryKey:['docs', studentId], queryFn: ()=> api(`/students/${studentId}/documents`) })
  return (
    <div className="space-y-3">
      {(data||[]).map((d:Row)=> (
        <Card key={d.id}>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <span>{d.title}</span>
              <a className="underline text-sm" href={`/documents/${d.id}`}>Open</a>
            </CardTitle>
          </CardHeader>
          <CardContent className="text-xs text-gray-600">
            Kind: {d.kind} • Version: {d.current_version_id ? 'yes' : 'no'}
          </CardContent>
        </Card>
      ))}
      {(!data || data.length===0) && <div className="text-sm text-gray-600">No documents yet.</div>}
    </div>
  )
}
