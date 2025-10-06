import React, { useEffect, useMemo, useRef, useState } from 'react'
import { api } from '../api/client'
import { PresignedUploader } from '../components/PresignedUploader'
import { Card, CardHeader, CardTitle, CardContent } from '../components/ui/card'

import { Document as PdfDoc, Page as PdfPage, pdfjs } from 'react-pdf'
pdfjs.GlobalWorkerOptions.workerSrc = new URL('pdfjs-dist/build/pdf.worker.min.mjs', import.meta.url).toString()

import { Input } from '../components/ui/input'
import { Textarea } from '../components/ui/textarea'
import { Button } from '../components/ui/button'
import { MentionChips, MentionResults, type Mention } from '../components/mentions'
import { useTranslation } from 'react-i18next'

type Version = { id:string; storage_path?:string; mime_type?:string; size_bytes?:number }
type Doc = { id:string; title:string; kind:string }
type Comment = { id:string; body:string; author_id:string; resolved:boolean; parent_id?:string|null; mentions?:string[] }

export function DocumentDetail({ docId }: { docId: string }) {
  const { t: T } = useTranslation('common')
  const [doc, setDoc] = useState<Doc | null>(null)
  const [versions, setVersions] = useState<Version[]>([])
  const [comments, setComments] = useState<Comment[]>([])
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)
  const [supportS3, setSupportS3] = useState<boolean | null>(null)
  const [newComment, setNewComment] = useState('')
  const [replyTo, setReplyTo] = useState<string | null>(null)
  const [mentionQ, setMentionQ] = useState('')
  const [mentions, setMentions] = useState<{id:string,name:string}[]>([])
  const [chosenMentions, setChosenMentions] = useState<Mention[]>([])

  async function load() {
    const d = await api(`/documents/${docId}`)
    setDoc(d.doc)
    setVersions(d.versions)
    const cs = await api(`/documents/${docId}/comments`)
    setComments(cs)
  }
  useEffect(()=>{ load() }, [docId])

  // detect S3: try presign without sending file; if 400 -> local only
  useEffect(()=>{
    (async()=>{
      try {
        const res = await api(`/documents/${docId}/presign`, { method:'POST', body: JSON.stringify({ filename: 'ping.txt', content_type: 'text/plain' }) })
        setSupportS3(!!res.url)
      } catch { setSupportS3(false) }
    })()
  },[docId])

  // mentions search
  useEffect(()=>{
    const t = setTimeout(async ()=>{
      if (mentionQ.length<1) { setMentions([]); return }
      const res = await api(`/admin/users?q=${encodeURIComponent(mentionQ)}`)
      setMentions(res.map((u:any)=>({id:u.id,name:u.name})))
    }, 250)
    return ()=>clearTimeout(t)
  },[mentionQ])

  const thread = useMemo(()=>{
    const byParent: Record<string, Comment[]> = {}
    comments.forEach(c => {
      const p = c.parent_id || 'root'
      byParent[p] = byParent[p] || []
      byParent[p].push(c)
    })
    return byParent
  },[comments])

  async function addComment() {
    await api(`/documents/${docId}/comments`, {
      method:'POST',
      body: JSON.stringify({ body: newComment, parent_id: replyTo, mentions: chosenMentions.map(m=>m.id) })
    })
    setNewComment(''); setReplyTo(null); setChosenMentions([])
    await load()
  }

  async function uploadLocal(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const form = e.target as HTMLFormElement
    const input = form.querySelector('input[type=file]') as HTMLInputElement
    const f = input?.files?.[0]; if (!f) return
    const data = new FormData()
    data.append('file', f)
    await fetch(`${import.meta.env.VITE_API_URL}/documents/${docId}/versions`, { method:'POST', body:data, headers: {} })
    await load()
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader><CardTitle>{T('breadcrumbs.document')}</CardTitle></CardHeader>
        <CardContent>
          <div className="text-sm">Title: {doc?.title}</div>
          <div className="text-sm">Kind: {doc?.kind}</div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>Preview</CardTitle></CardHeader>
        <CardContent>
          {previewUrl ? (
            <div className="border rounded p-2 bg-gray-50">
              <PdfDoc file={previewUrl} onLoadError={()=>{}}>
                <PdfPage pageNumber={1} width={600} />
              </PdfDoc>
              <div className="text-xs text-gray-600 mt-2">Showing first page. Click below to open the full file.</div>
            </div>
          ) : <div className="text-sm text-gray-600">No preview available.</div>}
          <div className="mt-2">
            <button className="underline text-sm" onClick={()=>{ if(previewUrl) window.open(previewUrl, '_blank') }}>Open full PDF</button>
          </div>
        </CardContent>
      </Card>


      <Card>
        <CardHeader><CardTitle>Upload</CardTitle></CardHeader>
        <CardContent className="space-y-2">
          {supportS3 ? (
            <PresignedUploader docId={docId} />
          ) : (
            <form className="space-y-2" onSubmit={uploadLocal}>
              <input type="file" accept=".pdf,.docx" />
              <Button type="submit">Upload</Button>
            </form>
          )}
          <div className="text-xs text-gray-500">Auto-detected: {supportS3 ? 'S3 pre-signed' : 'Local multipart'}</div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>Preview</CardTitle></CardHeader>
        <CardContent>
          {previewUrl ? (
            <div className="border rounded p-2 bg-gray-50">
              <PdfDoc file={previewUrl} onLoadError={()=>{}}>
                <PdfPage pageNumber={1} width={600} />
              </PdfDoc>
              <div className="text-xs text-gray-600 mt-2">Showing first page. Click below to open the full file.</div>
            </div>
          ) : <div className="text-sm text-gray-600">No preview available.</div>}
          <div className="mt-2">
            <button className="underline text-sm" onClick={()=>{ if(previewUrl) window.open(previewUrl, '_blank') }}>Open full PDF</button>
          </div>
        </CardContent>
      </Card>


      <Card>
        <CardHeader><CardTitle>Versions</CardTitle></CardHeader>
        <CardContent>
          <ul className="list-disc pl-5 text-sm">
            {versions.map(v => <li key={v.id}>{v.id} • {v.mime_type} • {v.size_bytes} bytes</li>)}
          </ul>
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>Preview</CardTitle></CardHeader>
        <CardContent>
          {previewUrl ? (
            <div className="border rounded p-2 bg-gray-50">
              <PdfDoc file={previewUrl} onLoadError={()=>{}}>
                <PdfPage pageNumber={1} width={600} />
              </PdfDoc>
              <div className="text-xs text-gray-600 mt-2">Showing first page. Click below to open the full file.</div>
            </div>
          ) : <div className="text-sm text-gray-600">No preview available.</div>}
          <div className="mt-2">
            <button className="underline text-sm" onClick={()=>{ if(previewUrl) window.open(previewUrl, '_blank') }}>Open full PDF</button>
          </div>
        </CardContent>
      </Card>


      <Card>
        <CardHeader><CardTitle>Comments</CardTitle></CardHeader>
        <CardContent>
          <div className="space-y-3">
            {(thread['root']||[]).map(c => (
              <div key={c.id} className="border rounded p-2">
                <div className="text-sm">{c.body}</div>
                <div className="text-xs text-gray-500">Mentions: {(c.mentions||[]).length}</div>
                {(thread[c.id]||[]).map(rc => (
                  <div key={rc.id} className="ml-4 mt-2 border-l pl-3 text-sm">{rc.body}</div>
                ))}
                <div className="mt-2">
                  <Button className="text-xs" onClick={()=>setReplyTo(c.id)}>Reply</Button>
                </div>
              </div>
            ))}
          </div>

          <div className="mt-4 space-y-2">
            {replyTo && <div className="text-xs text-gray-600">Replying to: {replyTo}</div>}
            <Textarea placeholder="Write a comment… Use mentions below." value={newComment} onChange={e=>setNewComment(e.target.value)} />
            <div className="space-y-2">
              <Input placeholder="Search @mentions" value={mentionQ} onChange={e=>setMentionQ(e.target.value)} />
              <MentionResults list={mentions} onPick={(m)=>{ if(!chosenMentions.find(x=>x.id===m.id)) setChosenMentions([...chosenMentions, m]); setMentionQ(''); }} />
              <MentionChips list={chosenMentions} onRemove={(id)=>setChosenMentions(chosenMentions.filter(m=>m.id!==id))} />
            </div>
            <div className="flex gap-2">
              <Button onClick={addComment}>Post</Button>
              {replyTo && <Button className="border-rose-400" onClick={()=>setReplyTo(null)}>Cancel reply</Button>}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
