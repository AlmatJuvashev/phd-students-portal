"use client"

import type React from "react"

import { useState, useEffect } from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Users, GraduationCap, MessageSquare } from "lucide-react"
import type { ChatRoom, ChatRoomType, RoomFormValues } from "@/lib/types/chat"

interface ChatRoomFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  room?: ChatRoom | null
  onSubmit: (data: RoomFormValues) => void
}

export function ChatRoomFormModal({ open, onOpenChange, room, onSubmit }: ChatRoomFormModalProps) {
  const [formData, setFormData] = useState<RoomFormValues>({
    name: "",
    type: "cohort",
    description: "",
  })
  const [errors, setErrors] = useState<Partial<Record<keyof RoomFormValues, string>>>({})

  useEffect(() => {
    if (open) {
      setFormData({
        name: room?.name || "",
        type: room?.type || "cohort",
        description: room?.description || "",
      })
      setErrors({})
    }
  }, [open, room])

  const validate = (): boolean => {
    const newErrors: Partial<Record<keyof RoomFormValues, string>> = {}
    if (!formData.name.trim()) {
      newErrors.name = "Название обязательно"
    } else if (formData.name.length < 3) {
      newErrors.name = "Минимум 3 символа"
    }
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (validate()) {
      onSubmit(formData)
      onOpenChange(false)
    }
  }

  const isEditing = !!room

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] max-h-[90dvh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEditing ? "Редактировать группу" : "Создать группу"}</DialogTitle>
          <DialogDescription>
            {isEditing ? "Измените параметры группы чата" : "Заполните информацию о новой группе чата"}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4 py-4">
            {/* Name field */}
            <div className="space-y-2">
              <Label htmlFor="name">
                Название <span className="text-destructive">*</span>
              </Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Например: PhD 2025 Биомедицина"
                className={errors.name ? "border-destructive" : ""}
              />
              {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="type">Тип группы</Label>
              <Select
                value={formData.type}
                onValueChange={(value: ChatRoomType) => setFormData({ ...formData, type: value })}
              >
                <SelectTrigger id="type">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="cohort">
                    <div className="flex items-center gap-2">
                      <Users className="h-4 w-4 text-blue-600" />
                      <span>Когорта</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="advisory">
                    <div className="flex items-center gap-2">
                      <GraduationCap className="h-4 w-4 text-emerald-600" />
                      <span>Научное руководство</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="other">
                    <div className="flex items-center gap-2">
                      <MessageSquare className="h-4 w-4 text-slate-600" />
                      <span>Другое</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Description field */}
            <div className="space-y-2">
              <Label htmlFor="description">Описание</Label>
              <Textarea
                id="description"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Краткое описание группы..."
                rows={3}
                className="resize-none"
              />
            </div>
          </div>
          <DialogFooter className="flex-col-reverse sm:flex-row gap-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} className="w-full sm:w-auto">
              Отмена
            </Button>
            <Button type="submit" className="w-full sm:w-auto">
              {isEditing ? "Сохранить" : "Создать"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
